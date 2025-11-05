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

type SubscribeGate struct {
	Time    int64    `json:"time"`
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Payload []string `json:"payload"`
}

type GateIOResponseTrade struct {
	Time    int    `json:"time"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Result  struct {
		ID           int    `json:"id"`
		CreateTime   int    `json:"create_time"`
		CreateTimeMs string `json:"create_time_ms"`
		Side         string `json:"side"`
		CurrencyPair string `json:"currency_pair"`
		Amount       string `json:"amount"`
		Price        string `json:"price"`
	} `json:"result"`
}

type gateIOScraper struct {
	wsClient           wsConn
	tradesChannel      chan models.Trade
	subscribeChannel   chan models.ExchangePair
	unsubscribeChannel chan models.ExchangePair
	tickerPairMap      map[string]models.Pair
	lastTradeTimeMap   map[string]time.Time
	maxErrCount        int
	restartWaitTime    int
	genesis            time.Time
	watchdogCancel     map[string]context.CancelFunc
}

var (
	_GateIOsocketurl string = "wss://api.gateio.ws/ws/v4/"
)

func NewGateIOScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("GateIO - Started scraper.")

	scraper := gateIOScraper{
		tradesChannel:      make(chan models.Trade),
		subscribeChannel:   make(chan models.ExchangePair),
		unsubscribeChannel: make(chan models.ExchangePair),
		tickerPairMap:      models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:   make(map[string]time.Time),
		maxErrCount:        20,
		restartWaitTime:    5,
		genesis:            time.Now(),
		watchdogCancel:     make(map[string]context.CancelFunc),
	}

	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(_GateIOsocketurl, nil)
	if err != nil {
		log.Errorf("GateIO - Dial ws base string: %v.", err)
		failoverChannel <- string(GATEIO_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("GateIO - subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("GateIO - Subscribed to pair %s.", pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Start watchdog for each pair.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}
	return &scraper
}

func (scraper *gateIOScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("GateIO - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("GateIO - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("GateIO - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *gateIOScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	envKey := strings.ToUpper(GATEIO_EXCHANGE) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "30"))
	if err != nil {
		log.Errorf("GateIO - Failed to parse %s: %v.", envKey, err)
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Get the initial config.
	last, err := models.GetExchangePairMap(GATEIO_EXCHANGE)
	if err != nil {
		log.Errorf("GateIO - GetExchangePairMap: %v.", err)
		return
	}

	// Watch for config changes.
	for {
		select {
		case <-ticker.C:
			current, err := models.GetExchangePairMap(GATEIO_EXCHANGE)
			if err != nil {
				log.Errorf("GateIO - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, current)
			// Update the last config.
			last = current
		case <-ctx.Done():
			log.Debugf("GateIO - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *gateIOScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("GateIO - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("GateIO - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("GateIO - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		err = scraper.subscribe(ep, true, lock)
		if err != nil {
			log.Errorf("GateIO - Failed to subscribe to %s: %v", ep.ForeignName, err)
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
		log.Infof("GateIO - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *gateIOScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("GateIO - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *gateIOScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(GATEIO_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", GATEIO_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(GATEIO_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", GATEIO_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *gateIOScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *gateIOScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *gateIOScraper) Close(cancel context.CancelFunc) error {
	log.Warn("GateIO - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *gateIOScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *gateIOScraper) fetchTrades(lock *sync.RWMutex) {
	var errCount int
	for {

		var message GateIOResponseTrade
		if err := scraper.wsClient.ReadJSON(&message); err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, GATEIO_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		trade := scraper.handleWSResponse(message, lock)

		lock.Lock()
		scraper.lastTradeTimeMap[trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol] = trade.Time
		lock.Unlock()

		log.Tracef("GateIO - got trade: %s -- %v -- %v -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)

		scraper.tradesChannel <- trade

	}
}

func (scraper *gateIOScraper) handleWSResponse(message GateIOResponseTrade, lock *sync.RWMutex) models.Trade {
	var (
		f64Price     float64
		f64Volume    float64
		exchangepair models.Pair
		err          error
	)

	f64Price, err = strconv.ParseFloat(message.Result.Price, 64)
	if err != nil {
		log.Errorf("GateIO - error parsing float Price %v: %v.", message.Result.Price, err)
		return models.Trade{}
	}

	f64Volume, err = strconv.ParseFloat(message.Result.Amount, 64)
	if err != nil {
		log.Errorf("GateIO - error parsing float Volume %v: %v.", message.Result.Amount, err)
		return models.Trade{}
	}

	if message.Result.Side == "sell" {
		f64Volume = -f64Volume
	}
	lock.RLock()
	exchangepair = scraper.tickerPairMap[strings.Split(message.Result.CurrencyPair, "_")[0]+strings.Split(message.Result.CurrencyPair, "_")[1]]
	lock.RUnlock()
	t := models.Trade{
		QuoteToken:     exchangepair.QuoteToken,
		BaseToken:      exchangepair.BaseToken,
		Price:          f64Price,
		Volume:         f64Volume,
		Time:           time.Unix(int64(message.Result.CreateTime), 0),
		Exchange:       Exchanges[GATEIO_EXCHANGE],
		ForeignTradeID: strconv.FormatInt(int64(message.Result.ID), 16),
	}

	return t
}

func (scraper *gateIOScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("GateIO - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("GateIO - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("GateIO - Resubscribe pair %s: %v", pair.ForeignName, err)
			} else {
				log.Debugf("GateIO - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Warnf("GateIO - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *gateIOScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	gateioPairTicker := strings.Split(pair.ForeignName, "-")[0] + "_" + strings.Split(pair.ForeignName, "-")[1]
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}
	a := &SubscribeGate{
		Event:   subscribeType,
		Time:    time.Now().Unix(),
		Channel: "spot.trades",
		Payload: []string{gateioPairTicker},
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}
