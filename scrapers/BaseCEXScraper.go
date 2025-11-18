package scrapers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

type ScraperHooks interface {
	// exchange key (e.g. BYBIT_EXCHANGE) for logging and config fetching
	ExchangeKey() string
	// WebSocket URL
	WSURL() string

	// after open connection (or reconnect) optional initialization (like start ping)
	OnOpen(ctx context.Context, bs *BaseCEXScraper)

	// construct and send subscribe/unsubscribe (usually use bs.SafeWriteJSON/Message)
	Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error

	// parse single ws raw message (default read loop will pass bytes to here)
	OnMessage(bs *BaseCEXScraper, messageType int, data []byte, lock *sync.RWMutex)

	// optional: completely take over read loop (return true means handled, Base will not run default loop)
	// for example Coinbase prefers ReadJSON(&struct{})
	ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool)

	// convert external pair.ForeignName to tickerPairMap key (like remove hyphen/case)
	TickerKeyFromForeign(foreign string) string

	// optional: customize last trade mapping key for different exchanges (default use foreignName itself)
	LastTradeTimeKeyFromForeign(foreign string) string
}

type DialerHooks interface {
	Dial(ctx context.Context, url string) (wsConn, error)
}

type BaseCEXScraper struct {
	wsClient wsConn
	hooks    ScraperHooks

	tradesChannel      chan models.Trade
	subscribeChannel   chan models.ExchangePair
	unsubscribeChannel chan models.ExchangePair

	watchdogCancel   map[string]context.CancelFunc
	tickerPairMap    map[string]models.Pair // key: hooks.TickerKeyFromForeign()
	lastTradeTimeMap map[string]time.Time   // key: hooks.LastTradeTimeKeyFromForeign()

	maxErrCount     int
	restartWaitTime int
	genesis         time.Time

	writeLock sync.Mutex
}

// construct common Scraper, called by NewXxxScraper of specific exchanges
func NewBaseCEXScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	failoverChannel chan string,
	wg *sync.WaitGroup,
	hooks ScraperHooks,
) *BaseCEXScraper {
	defer wg.Done()

	var lock sync.RWMutex
	key := strings.ToUpper(hooks.ExchangeKey())
	log.Infof("%s - Started scraper.", key)

	bs := &BaseCEXScraper{
		hooks:              hooks,
		tradesChannel:      make(chan models.Trade),
		subscribeChannel:   make(chan models.ExchangePair),
		unsubscribeChannel: make(chan models.ExchangePair),
		watchdogCancel:     make(map[string]context.CancelFunc),
		tickerPairMap:      models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:   make(map[string]time.Time),
		maxErrCount:        20,
		restartWaitTime:    5,
		genesis:            time.Now(),
	}

	// 1) connect: if implements DialerHooks, use custom Dial, otherwise use default Dial
	var (
		conn wsConn
		err  error
	)

	if dh, ok := hooks.(DialerHooks); ok {
		// custom Dial (e.g. Binance uses proxy)
		conn, err = dh.Dial(ctx, hooks.WSURL())
	} else {
		// default: no proxy, direct dial
		var d ws.Dialer
		c, _, e := d.Dial(hooks.WSURL(), nil)
		conn, err = c, e
	}

	if err != nil {
		log.Errorf("%s - WebSocket connection failed: %v", key, err)
		failoverChannel <- key
		return bs
	}
	bs.wsClient = conn

	// 2) exchange custom initialization (like ping)
	hooks.OnOpen(ctx, bs)

	// 3) initial subscribe + initialize lastTradeTimeMap
	for _, p := range pairs {
		if err := hooks.Subscribe(bs, p, true, &lock); err != nil {
			log.Errorf("%s - Failed to subscribe to %v: %v", key, p.ForeignName, err)
		} else {
			log.Infof("%s - Subscribed to %v", key, p.ForeignName)
			bs.setLastTradeTime(&lock, hooks.LastTradeTimeKeyFromForeign(p.ForeignName), time.Now())
		}
	}

	// 4) go routines: read loop, resubscribe, process unsubscribe, watch config, watchdog
	go bs.runReadLoop(ctx, &lock)
	go bs.runResubscribe(ctx, &lock)
	go bs.runProcessUnsubscribe(ctx, &lock)
	go bs.runWatchConfig(ctx, &lock)

	for _, p := range pairs {
		bs.startWatchdogForPair(ctx, &lock, p)
	}

	return bs
}

// export trades channel
func (bs *BaseCEXScraper) TradesChannel() chan models.Trade { return bs.tradesChannel }

func (bs *BaseCEXScraper) Close(cancel context.CancelFunc) error {
	log.Warnf("%s - call scraper.Close().", strings.ToUpper(bs.hooks.ExchangeKey()))
	cancel()
	if bs.wsClient != nil {
		return bs.wsClient.Close()
	}
	return nil
}

// thread safe write: JSON & Message
func (bs *BaseCEXScraper) SafeWriteJSON(v interface{}) error {
	bs.writeLock.Lock()
	defer bs.writeLock.Unlock()
	return bs.wsClient.WriteJSON(v)
}
func (bs *BaseCEXScraper) SafeWriteMessage(messageType int, data []byte) error {
	bs.writeLock.Lock()
	defer bs.writeLock.Unlock()
	return bs.wsClient.WriteMessage(messageType, data)
}

// ---------------- public loop ----------------

func (bs *BaseCEXScraper) runReadLoop(ctx context.Context, lock *sync.RWMutex) {
	// give to hooks to completely take over
	if handled := bs.hooks.ReadLoop(ctx, bs, lock); handled {
		return
	}

	// default: ReadMessage -> give to OnMessage
	var errCount int
	for {
		select {
		case <-ctx.Done():
			log.Infof("%s - Stopping WebSocket reader", strings.ToUpper(bs.hooks.ExchangeKey()))
			return
		default:
			messageType, msg, err := bs.wsClient.ReadMessage()
			if err != nil {
				if handleErrorReadJSON(err, &errCount, bs.maxErrCount, bs.hooks.ExchangeKey(), bs.restartWaitTime) {
					return
				}
				continue
			}
			bs.hooks.OnMessage(bs, messageType, msg, lock)
		}
	}
}

func (bs *BaseCEXScraper) runResubscribe(ctx context.Context, lock *sync.RWMutex) {
	exKey := strings.ToUpper(bs.hooks.ExchangeKey())
	for {
		select {
		case pair := <-bs.subscribeChannel:
			// unsubscribe first, then subscribe
			if err := bs.hooks.Subscribe(bs, pair, false, lock); err != nil {
				log.Errorf("%s - Unsubscribe pair %s: %v.", exKey, pair.ForeignName, err)
			} else {
				log.Debugf("%s - Unsubscribed pair %s.", exKey, pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			if err := bs.hooks.Subscribe(bs, pair, true, lock); err != nil {
				log.Errorf("%s - Resubscribe pair %s: %v.", exKey, pair.ForeignName, err)
			} else {
				log.Debugf("%s - Subscribed to pair %s.", exKey, pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("%s - Close resubscribe routine of scraper with genesis: %v.", exKey, bs.genesis)
			return
		}
	}
}

func (bs *BaseCEXScraper) runProcessUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	exKey := strings.ToUpper(bs.hooks.ExchangeKey())
	for {
		select {
		case pair := <-bs.unsubscribeChannel:
			if err := bs.hooks.Subscribe(bs, pair, false, lock); err != nil {
				log.Errorf("%s - Unsubscribe pair %s: %v.", exKey, pair.ForeignName, err)
			} else {
				log.Infof("%s - Unsubscribed pair %s.", exKey, pair.ForeignName)
			}
			// clean lastTradeTime, watchdog
			bs.deleteLastTradeTime(lock, bs.hooks.LastTradeTimeKeyFromForeign(pair.ForeignName))
			bs.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("%s - Close processUnsubscribe routine of scraper with genesis: %v.", exKey, bs.genesis)
			return
		}
	}
}

func (bs *BaseCEXScraper) runWatchConfig(ctx context.Context, lock *sync.RWMutex) {
	ex := strings.ToUpper(bs.hooks.ExchangeKey())
	envKey := ex + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "30"))
	if err != nil {
		log.Errorf("%s - Failed to parse %s: %v.", ex, envKey, err)
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// initial config
	last, err := models.GetExchangePairMap(bs.hooks.ExchangeKey())
	if err != nil {
		log.Errorf("%s - GetExchangePairMap: %v.", ex, err)
		return
	}

	for {
		select {
		case <-ticker.C:
			current, err := models.GetExchangePairMap(bs.hooks.ExchangeKey())
			if err != nil {
				log.Errorf("%s - GetExchangePairMap: %v.", ex, err)
				continue
			}
			bs.applyConfigDiff(ctx, lock, last, current)
			last = current
		case <-ctx.Done():
			log.Debugf("%s - Close watchConfig routine of scraper with genesis: %v.", ex, bs.genesis)
			return
		}
	}
}

func (bs *BaseCEXScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last, current map[string]int64) {
	added, removed, changed := diffPairMap(last, current)

	// removed -> unsubscribe
	for _, p := range removed {
		log.Infof("%s - Removed pair %s.", strings.ToUpper(bs.hooks.ExchangeKey()), p)
		bs.unsubscribeChannel <- models.ExchangePair{ForeignName: p}
	}

	// added -> subscribe + watchdog + map initialize
	for _, p := range added {
		delay := current[p]
		log.Infof("%s - Added pair %s with delay %v.", strings.ToUpper(bs.hooks.ExchangeKey()), p, delay)

		ep, err := bs.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("%s - GetExchangePairInfo(%s) err: %v.", strings.ToUpper(bs.hooks.ExchangeKey()), p, err)
			continue
		}
		if err := bs.hooks.Subscribe(bs, ep, true, lock); err != nil {
			log.Errorf("%s - Subscribe %s err: %v", strings.ToUpper(bs.hooks.ExchangeKey()), ep.ForeignName, err)
			continue
		}
		bs.startWatchdogForPair(ctx, lock, ep)
		key := bs.hooks.TickerKeyFromForeign(ep.ForeignName)
		lock.Lock()
		bs.tickerPairMap[key] = ep.UnderlyingPair
		if _, ok := bs.lastTradeTimeMap[bs.hooks.LastTradeTimeKeyFromForeign(ep.ForeignName)]; !ok {
			bs.lastTradeTimeMap[bs.hooks.LastTradeTimeKeyFromForeign(ep.ForeignName)] = time.Now()
		}
		lock.Unlock()
	}

	// changed -> restart watchdog (only delay)
	for _, p := range changed {
		newDelay := current[p]
		bs.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func diffPairMap(last, current map[string]int64) (added, removed, changed []string) {
	added, removed, changed = []string{}, []string{}, []string{}
	if last == nil {
		for p := range current {
			added = append(added, p)
		}
		return
	}
	for p := range current {
		if _, ok := last[p]; !ok {
			added = append(added, p)
		}
	}
	for p := range last {
		if _, ok := current[p]; !ok {
			removed = append(removed, p)
		}
	}
	for p, newDelay := range current {
		if oldDelay, ok := last[p]; ok && oldDelay != newDelay {
			changed = append(changed, p)
		}
	}
	return
}

// ---------------- watchdog & config helper ----------------

func (bs *BaseCEXScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(bs.hooks.ExchangeKey())
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", bs.hooks.ExchangeKey(), err)
	}
	ep, err := models.ConstructExchangePair(bs.hooks.ExchangeKey(), foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", bs.hooks.ExchangeKey(), foreignName, delay, err)
	}
	return ep, nil
}

func (bs *BaseCEXScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
	lock.Lock()
	if cancel, exists := bs.watchdogCancel[pair.ForeignName]; exists && cancel != nil {
		lock.Unlock()
		return
	}
	lock.Unlock()

	wdCtx, cancel := context.WithCancel(ctx)
	lock.Lock()
	bs.watchdogCancel[pair.ForeignName] = cancel
	lock.Unlock()

	ticker := time.NewTicker(time.Duration(pair.WatchDogDelay) * time.Second)
	go watchdog(wdCtx, pair, ticker, bs.lastTradeTimeMap, pair.WatchDogDelay, bs.subscribeChannel, lock)
}

func (bs *BaseCEXScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	if cancel, ok := bs.watchdogCancel[foreignName]; ok && cancel != nil {
		cancel()
		delete(bs.watchdogCancel, foreignName)
		log.Debugf("%s - Stopped watchdog for pair %s.", strings.ToUpper(bs.hooks.ExchangeKey()), foreignName)
	}
	lock.Unlock()
}

func (bs *BaseCEXScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	bs.stopWatchdogForPair(lock, foreignName)
	ep, err := bs.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("%s - GetExchangePairInfo(%s) err: %v.", strings.ToUpper(bs.hooks.ExchangeKey()), foreignName, err)
		return
	}
	bs.startWatchdogForPair(ctx, lock, ep)
}

func (bs *BaseCEXScraper) setLastTradeTime(lock *sync.RWMutex, key string, t time.Time) {
	lock.Lock()
	bs.lastTradeTimeMap[key] = t
	lock.Unlock()
}
func (bs *BaseCEXScraper) deleteLastTradeTime(lock *sync.RWMutex, key string) {
	lock.Lock()
	delete(bs.lastTradeTimeMap, key)
	lock.Unlock()
}
