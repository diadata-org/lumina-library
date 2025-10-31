package scrapers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type binanceWSSubscribeMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type binanceWSResponse struct {
	Timestamp      int64       `json:"T"`
	Price          string      `json:"p"`
	Volume         string      `json:"q"`
	ForeignTradeID int         `json:"t"`
	ForeignName    string      `json:"s"`
	Type           interface{} `json:"e"`
	Buy            bool        `json:"m"`
}

type binanceScraper struct {
	wsClient          wsConn
	tradesChannel     chan models.Trade
	subscribeChannel  chan models.ExchangePair
	watchdogCancel    map[string]context.CancelFunc
	tickerPairMap     map[string]models.Pair
	lastTradeTimeMap  map[string]time.Time
	maxErrCount       int
	restartWaitTime   int
	genesis           time.Time
	apiConnectRetries int
	proxyIndex        int
	writeLimiter      *rate.Limiter
	writeQueue        chan interface{}
	idCounter         int64
	activePairs       map[string]models.ExchangePair // foreignName -> exchangePair
}

const (
	BINANCE_API_MAX_RETRIES = 5
)

var (
	binanceWSBaseString   = "wss://stream.binance.com:9443/ws"
	binanceApiWaitSeconds = 5
)

func NewBinanceScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Infof("Binance - Started scraper at %v.", time.Now())

	scraper := binanceScraper{
		tradesChannel:    make(chan models.Trade),
		subscribeChannel: make(chan models.ExchangePair),
		watchdogCancel:   make(map[string]context.CancelFunc),
		tickerPairMap:    models.MakeTickerPairMap(pairs),
		lastTradeTimeMap: make(map[string]time.Time),
		maxErrCount:      20,
		restartWaitTime:  5,
		genesis:          time.Now(),
		proxyIndex:       0,
		writeLimiter:     rate.NewLimiter(rate.Every(300*time.Millisecond), 1),
		writeQueue:       make(chan interface{}, 1024),
		idCounter:        0,
		activePairs:      make(map[string]models.ExchangePair),
	}

	// seed active pairs
	for _, pair := range pairs {
		scraper.activePairs[pair.ForeignName] = pair
	}

	err := errors.New("cannot connect to API")
	var errCount int
	for err != nil {

		if errCount > 2*scraper.apiConnectRetries {
			failoverChannel <- BINANCE_EXCHANGE
			return &scraper
		}

		err = scraper.connectToAPI()
		if err != nil {
			errCount++
			scraper.apiConnectRetries++
			time.Sleep(time.Duration(binanceApiWaitSeconds) * time.Second)
		}
	}

	go scraper.wsWriter(&lock)
	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @binanceWatchdogDelay[pair].
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper

}

func (scraper *binanceScraper) wsWriter(lock *sync.RWMutex) {
	for msg := range scraper.writeQueue {
		_ = scraper.writeLimiter.Wait(context.Background())
		lock.Lock()
		err := scraper.wsClient.WriteJSON(msg)
		lock.Unlock()
		if err != nil {
			log.Errorf("Binance - WriteJSON failed: %v", err)
			time.Sleep(200 * time.Microsecond)
			select {
			case scraper.writeQueue <- msg:
			default:
				log.Errorf("Binance - WriteQueue is full, dropping message.")
			}
		}
	}
}

func (scraper *binanceScraper) nextID() int {
	return int(atomic.AddInt64(&scraper.idCounter, 1))

}

func (s *binanceScraper) sendBatch(subscribe bool, topics []string) {
	if len(topics) == 0 {
		return
	}

	// keep batches small to be safe (e.g. 10 per message)
	const batchSize = 10
	m := "UNSUBSCRIBE"
	if subscribe {
		m = "SUBSCRIBE"
	}

	for i := 0; i < len(topics); i += batchSize {
		end := i + batchSize
		if end > len(topics) {
			end = len(topics)
		}
		payload := binanceWSSubscribeMessage{
			Method: m,
			Params: topics[i:end],
			ID:     s.nextID(),
		}
		s.writeQueue <- payload
	}
}

func (scraper *binanceScraper) subscribe(ep models.ExchangePair, doSub bool) error {
	topic := getTopic(ep.ForeignName)
	scraper.sendBatch(doSub, []string{topic})
	return nil
}

func getTopic(foreignName string) string {
	return strings.ToLower(strings.ReplaceAll(foreignName, "-", "")) + "@trade"
}

func (scraper *binanceScraper) processUnsubscribe(removed []string, lock *sync.RWMutex) {
	topics := make([]string, 0, len(removed))
	for _, p := range removed {
		log.Infof("Binance - Removed pair %s.", p)
		// stop watchdog for this pair
		scraper.stopWatchdogForPair(lock, p)
		lock.Lock()
		// delete last trade time for this pair
		delete(scraper.lastTradeTimeMap, p)
		// delete active pair from the map
		delete(scraper.activePairs, p) // keep the authoritative set in sync
		lock.Unlock()
		topics = append(topics, getTopic(p))
	}
	// unsubscribe from these pairs
	scraper.sendBatch(false, topics)
}

func (scraper *binanceScraper) processSubscribe(ctx context.Context, added []models.ExchangePair, lock *sync.RWMutex) {
	topics := make([]string, 0, len(added))
	for _, ep := range added {
		log.Infof("Binance - Added pair %s with delay %v.", ep.ForeignName, ep.WatchDogDelay)
		// maps
		key := strings.ReplaceAll(ep.ForeignName, "-", "")
		lock.Lock()
		scraper.tickerPairMap[key] = ep.UnderlyingPair
		if _, ok := scraper.lastTradeTimeMap[ep.ForeignName]; !ok {
			scraper.lastTradeTimeMap[ep.ForeignName] = time.Now()
		}
		scraper.activePairs[ep.ForeignName] = ep
		lock.Unlock()

		// start watchdog for this pair
		scraper.startWatchdogForPair(ctx, lock, ep)

		// add topic to the list of topics to subscribe to
		topics = append(topics, getTopic(ep.ForeignName))
	}
	scraper.sendBatch(true, topics)
}

func (scraper *binanceScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Keep track of the last config.
	var last map[string]int64

	// Get the initial config.
	cfg, err := models.GetExchangePairMap(BINANCE_EXCHANGE)
	if err != nil {
		log.Errorf("Binance - GetExchangePairMap: %v.", err)
		return
	} else {
		// Apply the initial config.
		last = cfg
		scraper.applyConfigDiff(ctx, lock, nil, cfg)
	}

	// Watch for config changes.
	for {
		select {
		case <-ticker.C:
			cfg, err := models.GetExchangePairMap(BINANCE_EXCHANGE)
			if err != nil {
				log.Errorf("Binance - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, cfg)
			// Update the last config.
			last = cfg
		case <-ctx.Done():
			log.Debugf("Binance - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *binanceScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

	added := make([]string, 0)
	removed := make([]string, 0)
	changed := make([]string, 0)

	// If last is nil, add all pairs from current.
	if last == nil {
		for p := range current {
			added = append(added, p)
		}
	} else {
		// If last is not nil, check for added and removed pairs.
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
	}

	// Unsubscribe from removed pairs.
	scraper.processUnsubscribe(removed, lock)

	// added -> subscribe to pairs and start watchdog
	if len(added) > 0 {
		eps := make([]models.ExchangePair, 0, len(added))
		for _, name := range added {
			ep, err := scraper.getExchangePairInfo(name, current[name])
			if err != nil {
				log.Errorf("Binance - Failed to GetExchangePairInfo for new pair %s: %v.", name, err)
				continue
			}
			eps = append(eps, ep)
		}
		scraper.processSubscribe(ctx, eps, lock)
	}

	// changed -> restart watchdog for pairs
	for _, p := range changed {
		newDelay := current[p]
		log.Infof("Binance - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *binanceScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("Binance - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *binanceScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(BINANCE_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", BINANCE_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(BINANCE_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", BINANCE_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *binanceScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
	// Check if watchdog is already running for this pair.
	lock.Lock()
	if cancel, exists := scraper.watchdogCancel[pair.ForeignName]; exists && cancel != nil {
		lock.Unlock()
		return
	}
	lock.Unlock()

	wdCtx, cancel := context.WithCancel(ctx)
	lock.Lock()
	scraper.watchdogCancel[pair.ForeignName] = cancel
	lock.Unlock()

	// Start watchdog for this pair.
	watchdogTicker := time.NewTicker(time.Duration(pair.WatchDogDelay) * time.Second)
	go watchdog(wdCtx, pair, watchdogTicker, scraper.lastTradeTimeMap, pair.WatchDogDelay, scraper.subscribeChannel, lock)
}

func (scraper *binanceScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *binanceScraper) Close(cancel context.CancelFunc) error {
	log.Warn("Binance - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *binanceScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *binanceScraper) fetchTrades(lock *sync.RWMutex) {
	var errCount int

	for {

		var message binanceWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, BINANCE_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		if message.Type == nil {
			continue
		}

		trade := binanceParseWSResponse(message)
		lock.RLock()
		trade.QuoteToken = scraper.tickerPairMap[message.ForeignName].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[message.ForeignName].BaseToken
		lock.RUnlock()
		log.Tracef("Binance - got trade %s -- %v -- %v -- %v.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
		lock.Lock()
		scraper.lastTradeTimeMap[trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol] = trade.Time
		lock.Unlock()

		scraper.tradesChannel <- trade
	}

}

func (scraper *binanceScraper) resubscribe(ctx context.Context) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			log.Debugf("Binance - scraper with genesis %v: Resubscribe pair %s.", scraper.genesis, pair.ForeignName)
			err := scraper.subscribe(pair, false)
			if err != nil {
				log.Errorf("Binance - scraper with genesis %v: Unsubscribe pair %s: %v.", scraper.genesis, pair.ForeignName, err)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true)
			if err != nil {
				log.Errorf("Binance - Resubscribe pair %s: %v.", pair.ForeignName, err)
			}
		case <-ctx.Done():
			log.Debugf("Binance - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *binanceScraper) connectToAPI() error {

	// Switch to alternative Proxy whenever too many retries on the first.
	if scraper.apiConnectRetries > BINANCE_API_MAX_RETRIES {
		log.Errorf("too many timeouts for Binance api connection with proxy %v. Switch to alternative proxy.", scraper.proxyIndex)
		scraper.apiConnectRetries = 0
		scraper.proxyIndex = (scraper.proxyIndex + 1) % 2
	}

	username := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(scraper.proxyIndex)+"_USERNAME", "")
	password := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(scraper.proxyIndex)+"_PASSWORD", "")
	user := url.UserPassword(username, password)
	host := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(scraper.proxyIndex)+"_HOST", "")
	var d ws.Dialer
	if host != "" {
		d = ws.Dialer{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http", // or "https" depending on your proxy
				User:   user,
				Host:   host,
				Path:   "/",
			}),
		}
	}

	wsAssetsString := ""
	for _, pair := range scraper.activePairs {
		wsAssetsString += "/" + getTopic(pair.ForeignName)
	}

	// Connect to Binance API.
	conn, _, err := d.Dial(binanceWSBaseString+wsAssetsString, nil)
	if err != nil {
		log.Errorf("Binance - Connect to API: %s.", err.Error())
		return err
	}
	scraper.wsClient = conn
	return nil

}

func binanceParseWSResponse(message binanceWSResponse) (trade models.Trade) {
	var err error
	trade.Exchange = Exchanges[BINANCE_EXCHANGE]
	trade.Time = time.Unix(0, message.Timestamp*1000000)
	trade.Price, err = strconv.ParseFloat(message.Price, 64)
	if err != nil {
		log.Errorf("Binance - Parse price: %v.", err)
	}
	trade.Volume, err = strconv.ParseFloat(message.Volume, 64)
	if err != nil {
		log.Errorf("Binance - Parse volume: %v.", err)
	}
	if !message.Buy {
		trade.Volume *= -1
	}
	trade.ForeignTradeID = strconv.Itoa(int(message.ForeignTradeID))
	return
}
