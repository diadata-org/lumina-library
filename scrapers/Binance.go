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
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
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
	wsClient           wsConn
	tradesChannel      chan models.Trade
	subscribeChannel   chan models.ExchangePair
	unsubscribeChannel chan models.ExchangePair
	watchdogCancel     map[string]context.CancelFunc
	tickerPairMap      map[string]models.Pair
	lastTradeTimeMap   map[string]time.Time
	maxErrCount        int
	restartWaitTime    int
	genesis            time.Time
	apiConnectRetries  int
	proxyIndex         int
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
		tradesChannel:      make(chan models.Trade),
		subscribeChannel:   make(chan models.ExchangePair),
		unsubscribeChannel: make(chan models.ExchangePair),
		watchdogCancel:     make(map[string]context.CancelFunc),
		tickerPairMap:      models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:   make(map[string]time.Time),
		maxErrCount:        20,
		restartWaitTime:    5,
		genesis:            time.Now(),
		proxyIndex:         0,
	}

	err := errors.New("cannot connect to API")
	var errCount int
	for err != nil {

		if errCount > 2*scraper.apiConnectRetries {
			failoverChannel <- BINANCE_EXCHANGE
			return &scraper
		}

		err = scraper.connectToAPI(pairs)
		if err != nil {
			errCount++
			scraper.apiConnectRetries++
			time.Sleep(time.Duration(binanceApiWaitSeconds) * time.Second)
		}
	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @binanceWatchdogDelay[pair].
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper

}

func (scraper *binanceScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("Binance - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("Binance - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("Binance - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
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
	for _, p := range removed {
		log.Infof("Binance - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("Binance - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("Binance - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		scraper.subscribeChannel <- ep
		// Start watchdog for this pair.
		scraper.startWatchdogForPair(ctx, lock, ep)
		// Add the pair to the ticker pair map.
		key := strings.ReplaceAll(ep.ForeignName, "-", "")
		lock.Lock()
		scraper.tickerPairMap[key] = ep.UnderlyingPair
		// Set the last trade time for this pair.
		if _, exists := scraper.lastTradeTimeMap[ep.ForeignName]; !exists {
			scraper.lastTradeTimeMap[ep.ForeignName] = time.Now()
		}
		lock.Unlock()
	}
	// Resubscribe to changed pairs.
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
	ep := models.ConstructExchangePair(BINANCE_EXCHANGE, foreignName, delay, idMap)
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
		trade.QuoteToken = scraper.tickerPairMap[message.ForeignName].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[message.ForeignName].BaseToken

		log.Tracef("Binance - got trade %s -- %v -- %v -- %v.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
		lock.Lock()
		scraper.lastTradeTimeMap[trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol] = trade.Time
		lock.Unlock()

		scraper.tradesChannel <- trade
	}

}

func (scraper *binanceScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			log.Debugf("Binance - scraper with genesis %v: Resubscribe pair %s.", scraper.genesis, pair.ForeignName)
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("Binance - scraper with genesis %v: Unsubscribe pair %s: %v.", scraper.genesis, pair.ForeignName, err)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("Binance - Resubscribe pair %s: %v.", pair.ForeignName, err)
			}
		case <-ctx.Done():
			log.Debugf("Binance - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *binanceScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "UNSUBSCRIBE"
	if subscribe {
		subscribeType = "SUBSCRIBE"
	}
	pairTicker := strings.ToLower(strings.Split(pair.ForeignName, "-")[0] + strings.Split(pair.ForeignName, "-")[1])
	subscribeMessage := &binanceWSSubscribeMessage{
		Method: subscribeType,
		Params: []string{pairTicker + "@trade"},
		ID:     1,
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(subscribeMessage)
}

func (scraper *binanceScraper) connectToAPI(pairs []models.ExchangePair) error {

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
	for _, pair := range pairs {
		wsAssetsString += "/" + strings.ToLower(strings.Split(pair.ForeignName, "-")[0]) + strings.ToLower(strings.Split(pair.ForeignName, "-")[1]) + "@trade"
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
