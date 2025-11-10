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

// A krakenWSSubscribeMessage represents a message to subscribe the public/private channel.
type krakenWSSubscribeMessage struct {
	Method string       `json:"method"`
	Params krakenParams `json:"params"`
}

type krakenParams struct {
	Channel string   `json:"channel"`
	Symbol  []string `json:"symbol"`
}

type krakenWSResponse struct {
	Channel string                 `json:"channel"`
	Type    string                 `json:"type"`
	Data    []krakenWSResponseData `json:"data"`
}

type krakenWSResponseData struct {
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Size      float64 `json:"qty"`
	OrderType string  `json:"order_type"`
	TradeID   int     `json:"trade_id"`
	Time      string  `json:"timestamp"`
}

type krakenScraper struct {
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
}

var (
	krakenWSBaseString = "wss://ws.kraken.com/v2"
)

func NewKrakenScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("Kraken - Started scraper.")

	scraper := krakenScraper{
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

	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(krakenWSBaseString, nil)
	if err != nil {
		log.Errorf("Kraken - Dial ws base string: %v.", err)
		failoverChannel <- string(KRAKEN_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("Kraken - Subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("Kraken - Subscribed to pair %s.", pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}

	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @krakenWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper

}

func (scraper *krakenScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("Kraken - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("Kraken - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("Kraken - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *krakenScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 60 minutes.
	envKey := strings.ToUpper(KRAKEN_EXCHANGE) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "3600"))
	if err != nil {
		log.Errorf("Kraken - Failed to parse %s: %v.", envKey, err)
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Get the initial config.
	last, err := models.GetExchangePairMap(KRAKEN_EXCHANGE)
	if err != nil {
		log.Errorf("Kraken - GetExchangePairMap: %v.", err)
		return
	}

	// Watch for config changes.
	for {
		select {
		case <-ticker.C:
			current, err := models.GetExchangePairMap(KRAKEN_EXCHANGE)
			if err != nil {
				log.Errorf("Kraken - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, current)
			// Update the last config.
			last = current
		case <-ctx.Done():
			log.Debugf("Kraken - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *krakenScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("Kraken - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("Kraken - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("Kraken - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		err = scraper.subscribe(ep, true, lock)
		if err != nil {
			log.Errorf("Kraken - Failed to subscribe to %s: %v", ep.ForeignName, err)
			continue // Don't start watchdog if subscription failed
		}
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
		log.Infof("Kraken - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *krakenScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("Kraken - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *krakenScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(KRAKEN_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", KRAKEN_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(KRAKEN_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", KRAKEN_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *krakenScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *krakenScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *krakenScraper) Close(cancel context.CancelFunc) error {
	log.Warn("Kraken - Call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *krakenScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *krakenScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int
	for {

		var message krakenWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, KRAKEN_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		if message.Channel == "trade" {
			for _, data := range message.Data {

				// Parse trade quantities.
				price, volume, timestamp, foreignTradeID, err := parseKrakenTradeMessage(data)
				if err != nil {
					log.Errorf("Kraken - parseTradeMessage: %v.", err)
					continue
				}

				// Identify ticker symbols with underlying assets.
				pair := strings.Split(data.Symbol, "/")
				var exchangepair models.Pair
				if len(pair) > 1 {
					lock.RLock()
					exchangepair = scraper.tickerPairMap[pair[0]+pair[1]]
					lock.RUnlock()
				}

				trade := models.Trade{
					QuoteToken:     exchangepair.QuoteToken,
					BaseToken:      exchangepair.BaseToken,
					Price:          price,
					Volume:         volume,
					Time:           timestamp,
					Exchange:       Exchanges[KRAKEN_EXCHANGE],
					ForeignTradeID: foreignTradeID,
				}
				log.Tracef("Kraken - got trade: %s -- %v -- %v -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
				lock.Lock()
				scraper.lastTradeTimeMap[exchangepair.QuoteToken.Symbol+"-"+exchangepair.BaseToken.Symbol] = trade.Time
				lock.Unlock()
				scraper.tradesChannel <- trade
			}
		}
	}
}

func (scraper *krakenScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("Kraken - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Kraken - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("Kraken - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Kraken - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("Kraken - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *krakenScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}
	a := &krakenWSSubscribeMessage{
		Method: subscribeType,
		Params: krakenParams{
			Channel: "trade",
			Symbol:  []string{strings.ReplaceAll(pair.ForeignName, "-", "/")},
		},
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func parseKrakenTradeMessage(message krakenWSResponseData) (price float64, volume float64, timestamp time.Time, foreignTradeID string, err error) {
	price = message.Price
	volume = message.Size
	if message.Side == "sell" {
		volume *= -1
	}
	timestamp, err = time.Parse("2006-01-02T15:04:05.000000Z", message.Time)
	if err != nil {
		return
	}

	foreignTradeID = strconv.Itoa(message.TradeID)
	return
}
