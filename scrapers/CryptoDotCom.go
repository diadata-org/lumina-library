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

type cryptodotcomWSSubscribeMessage struct {
	ID     int                  `json:"id"`
	Method string               `json:"method"`
	Params cryptodotcomChannels `json:"params"`
}

type cryptodotcomChannels struct {
	Channels []string `json:"channels"`
}

type cryptodotcomWSResponse struct {
	ID     int                          `json:"id"`
	Method string                       `json:"method"`
	Code   int                          `json:"code"`
	Result cryptodotcomWSResponseResult `json:"result"`
}

type cryptodotcomWSResponseResult struct {
	InstrumentName string                       `json:"instrument_name"`
	Subscription   string                       `json:"subscription"`
	Channel        string                       `json:"channel"`
	Data           []cryptodotcomWSResponseData `json:"data"`
}

type cryptodotcomWSResponseData struct {
	TradeID     string `json:"d"`
	Timestamp   int64  `json:"t"`
	Price       string `json:"p"`
	Volume      string `json:"q"`
	Side        string `json:"s"`
	ForeignName string `json:"i"`
}

type cryptodotcomScraper struct {
	wsClient            wsConn
	tradesChannel       chan models.Trade
	subscribeChannel    chan models.ExchangePair
	unsubscribeChannel  chan models.ExchangePair
	watchdogCancel      map[string]context.CancelFunc
	tickerPairMap       map[string]models.Pair
	lastTradeTimeMap    map[string]time.Time
	maxErrCount         int
	restartWaitTime     int
	genesis             time.Time
	tradeTimeoutSeconds int
}

var (
	cryptodotcomWSBaseString = "wss://stream.crypto.com/v2/market"
)

func NewCryptodotcomScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("Crypto.com - Started scraper.")

	scraper := cryptodotcomScraper{
		tradesChannel:       make(chan models.Trade),
		subscribeChannel:    make(chan models.ExchangePair),
		unsubscribeChannel:  make(chan models.ExchangePair),
		watchdogCancel:      make(map[string]context.CancelFunc),
		tickerPairMap:       models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:    make(map[string]time.Time),
		maxErrCount:         20,
		restartWaitTime:     5,
		genesis:             time.Now(),
		tradeTimeoutSeconds: 120,
	}

	// Dial websocket API.
	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(cryptodotcomWSBaseString, nil)
	if err != nil {
		log.Errorf("Crypto.com - Dial ws base string: %v.", err)
		failoverChannel <- string(CRYPTODOTCOM_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs and initialize cryptodotcomLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("Crypto.com - Subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("Crypto.com - Subscribed to pair %s.", pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @cryptodotcomWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper
}

func (scraper *cryptodotcomScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("Crypto.com - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("Crypto.com - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("Crypto.com - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *cryptodotcomScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	envKey := strings.ToUpper(CRYPTODOTCOM_EXCHANGE) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "30"))
	if err != nil {
		log.Errorf("Crypto.com - Failed to parse %s: %v.", envKey, err)
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Get the initial config.
	last, err := models.GetExchangePairMap(CRYPTODOTCOM_EXCHANGE)
	if err != nil {
		log.Errorf("Crypto.com - GetExchangePairMap: %v.", err)
		return
	}

	// Watch for config changes.
	for {
		select {
		case <-ticker.C:
			current, err := models.GetExchangePairMap(CRYPTODOTCOM_EXCHANGE)
			if err != nil {
				log.Errorf("Crypto.com - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, current)
			// Update the last config.
			last = current
		case <-ctx.Done():
			log.Debugf("Crypto.com - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *cryptodotcomScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("Crypto.com - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("Crypto.com - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("Crypto.com - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		err = scraper.subscribe(ep, true, lock)
		if err != nil {
			log.Errorf("Crypto.com - Failed to subscribe to %s: %v", ep.ForeignName, err)
			continue // Don't start watchdog if subscription failed
		}
		// Start watchdog for this pair.
		scraper.startWatchdogForPair(ctx, lock, ep)
		key := strings.ReplaceAll(ep.ForeignName, "-", "")
		// Add the pair to the ticker pair map.
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
		log.Infof("Crypto.com - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *cryptodotcomScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(CRYPTODOTCOM_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", CRYPTODOTCOM_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(CRYPTODOTCOM_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", CRYPTODOTCOM_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *cryptodotcomScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *cryptodotcomScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *cryptodotcomScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("Crypto.com - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *cryptodotcomScraper) Close(cancel context.CancelFunc) error {
	log.Warn("Crypto.com - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *cryptodotcomScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *cryptodotcomScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int

	for {

		var message cryptodotcomWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, CRYPTODOTCOM_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}
		if message.Method == "public/heartbeat" {
			scraper.sendHeartbeat(message.ID, lock)
			continue
		}

		scraper.handleWSResponse(message, lock)

	}

}

func (scraper *cryptodotcomScraper) handleWSResponse(message cryptodotcomWSResponse, lock *sync.RWMutex) {
	trades, err := cryptodotcomParseTradeMessage(message)
	if err != nil {
		log.Errorf("Crypto.com - parseCryptodotcomTradeMessage: %s.", err.Error())
		// continue
		return
	}

	// Identify ticker symbols with underlying assets.
	for _, trade := range trades {

		// The websocket API returns very old trades when first subscribing. Hence, discard if too old.
		if trade.Time.Before(time.Now().Add(-time.Duration(scraper.tradeTimeoutSeconds) * time.Second)) {
			continue
		}

		pair := strings.Split(message.Result.Data[0].ForeignName, "_")
		if len(pair) > 1 {
			lock.RLock()
			trade.QuoteToken = scraper.tickerPairMap[pair[0]+pair[1]].QuoteToken
			trade.BaseToken = scraper.tickerPairMap[pair[0]+pair[1]].BaseToken
			lock.RUnlock()
		}

		log.Tracef("Crypto.com - got trade: %v -- %s -- %v -- %v -- %s.", trade.Time, trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
		lock.Lock()
		scraper.lastTradeTimeMap[pair[0]+"-"+pair[1]] = trade.Time
		lock.Unlock()

		scraper.tradesChannel <- trade
	}

}

func (scraper *cryptodotcomScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("Crypto.com - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Crypto.com - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("Crypto.com - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Crypto.com - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("Crypto.com - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *cryptodotcomScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	channel := []string{"trade." + strings.Split(pair.ForeignName, "-")[0] + "_" + strings.Split(pair.ForeignName, "-")[1]}
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	a := cryptodotcomWSSubscribeMessage{
		ID:     1,
		Method: subscribeType,
		Params: cryptodotcomChannels{
			Channels: channel,
		},
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func (scraper *cryptodotcomScraper) sendHeartbeat(id int, lock *sync.RWMutex) error {
	defer lock.Unlock()
	a := cryptodotcomWSSubscribeMessage{
		ID:     id,
		Method: "public/respond-heartbeat",
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func cryptodotcomParseTradeMessage(message cryptodotcomWSResponse) (trades []models.Trade, err error) {

	for _, data := range message.Result.Data {
		var (
			price  float64
			volume float64
		)
		price, err = strconv.ParseFloat(data.Price, 64)
		if err != nil {
			return
		}
		volume, err = strconv.ParseFloat(data.Volume, 64)
		if err != nil {
			return
		}
		if data.Side == "SELL" {
			volume *= -1
		}
		timestamp := time.Unix(0, data.Timestamp*1e6)
		foreignTradeID := data.TradeID

		trade := models.Trade{
			Price:          price,
			Volume:         volume,
			Time:           timestamp,
			Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
			ForeignTradeID: foreignTradeID,
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
