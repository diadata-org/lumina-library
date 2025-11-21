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

// CustomDialer is an optional interface that exchanges (like Binance) can implement to customize Dial (e.g., using proxy)
type CustomDialer interface {
	Dial(ctx context.Context, url string) (wsConn, error)
}

type BaseCEXScraper struct {
	wsClient wsConn
	hooks    ScraperHooks

	cancel             context.CancelFunc
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
	wg *sync.WaitGroup,
	hooks ScraperHooks,
) *BaseCEXScraper {
	defer wg.Done()

	var lock sync.RWMutex
	key := strings.ToUpper(hooks.ExchangeKey())
	log.Infof("%s - Started scraper.", key)

	ctx, cancel := context.WithCancel(ctx)

	bs := &BaseCEXScraper{
		hooks:              hooks,
		cancel:             cancel,
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

	// 1) connect + retry
	if ok := bs.connectWithRetry(ctx); !ok {
		log.Errorf("%s - initial WebSocket connection failed.", key)
		// ctx is already canceled, return directly (wsClient is nil)
		bs.cancel()
		return bs
	}

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

// first call hooks' custom Dial, if not implemented, use default Dialer
func (bs *BaseCEXScraper) dialOnce(ctx context.Context) (wsConn, error) {
	if d, ok := bs.hooks.(CustomDialer); ok {
		// exchange implemented Dial (can use proxy)
		return d.Dial(ctx, bs.hooks.WSURL())
	}

	// default: no proxy, direct Dial
	var dialer ws.Dialer
	c, _, err := dialer.Dial(bs.hooks.WSURL(), nil)
	return c, err
}

// connection logic with retry (using restartWaitTime and incrementing multiplier)
func (bs *BaseCEXScraper) connectWithRetry(ctx context.Context) bool {
	key := strings.ToUpper(bs.hooks.ExchangeKey())
	maxRetries, err := strconv.Atoi(utils.Getenv("MAX_RETRIES", "10"))
	if err != nil {
		log.Errorf("%s - Failed to parse MAX_RETRIES: %v.", key, err)
		maxRetries = 10
	}
	numRetries := 1

	for {
		select {
		case <-ctx.Done():
			log.Warnf("%s - context canceled during connect, stop retry", key)
			return false
		default:
		}

		conn, err := bs.dialOnce(ctx)
		if err == nil {
			bs.wsClient = conn
			log.Infof("%s - WebSocket connected.", key)
			return true
		}

		log.Errorf("%s - WebSocket connection failed (try %d): %v", key, numRetries, err)
		sleepSec := numRetries * bs.restartWaitTime
		if numRetries >= maxRetries {
			log.Errorf("%s - Max retries reached, giving up.", key)
			return false
		}
		time.Sleep(time.Duration(sleepSec) * time.Second)
		numRetries++
	}
}

// export trades channel
func (bs *BaseCEXScraper) TradesChannel() chan models.Trade { return bs.tradesChannel }

func (bs *BaseCEXScraper) Close(cancel context.CancelFunc) error {
	log.Warnf("%s - call scraper.Close().", strings.ToUpper(bs.hooks.ExchangeKey()))
	cancel()
	if cancel != nil {
		cancel()
	}
	if bs.cancel != nil {
		bs.cancel()
	}
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

	if bs.wsClient == nil {
		log.Warnf("%s - wsClient is nil, exiting read loop", strings.ToUpper(bs.hooks.ExchangeKey()))
		return
	}

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
	ex := bs.hooks.ExchangeKey()
	go WatchConfigLoop(ctx, ex, 3600, func(ctx context.Context, last, current map[string]int64) {
		bs.applyConfigDiff(ctx, lock, last, current)
	})
}

func WatchConfigLoop(
	ctx context.Context,
	exKey string, // e.g. "MEXC" / hooks.ExchangeKey()
	defaultIntervalSec int, // Base uses 3600
	apply func(ctx context.Context, last, current map[string]int64),
) {
	envKey := strings.ToUpper(exKey) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, strconv.Itoa(defaultIntervalSec)))
	if err != nil {
		log.Errorf("%s - Failed to parse %s: %v.", strings.ToUpper(exKey), envKey, err)
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	last, err := models.GetExchangePairMap(exKey)
	if err != nil {
		log.Errorf("%s - GetExchangePairMap: %v.", strings.ToUpper(exKey), err)
		return
	}

	for {
		select {
		case <-ticker.C:
			current, err := models.GetExchangePairMap(exKey)
			if err != nil {
				log.Errorf("%s - GetExchangePairMap: %v.", strings.ToUpper(exKey), err)
				continue
			}
			apply(ctx, last, current)
			last = current
		case <-ctx.Done():
			log.Debugf("%s - Close watchConfig routine.", strings.ToUpper(exKey))
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
	StartWatchdogForPair(
		ctx, lock, pair,
		bs.watchdogCancel,
		bs.lastTradeTimeMap,
		bs.subscribeChannel,
	)
}

func (bs *BaseCEXScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	StopWatchdogForPair(lock, foreignName, bs.watchdogCancel)
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
