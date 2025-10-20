package scrapers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// A coinBaseWSSubscribeMessage represents a message to subscribe the public/private channel.
type coinBaseWSSubscribeMessage struct {
	Type     string            `json:"type"`
	Channels []coinBaseChannel `json:"channels"`
}

type coinBaseChannel struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids"`
}

type coinBaseWSResponse struct {
	Type         string `json:"type"`
	TradeID      int64  `json:"trade_id"`
	Sequence     int64  `json:"sequence"`
	MakerOrderID string `json:"maker_order_id"`
	TakerOrderID string `json:"taker_order_id"`
	Time         string `json:"time"`
	ProductID    string `json:"product_id"`
	Size         string `json:"size"`
	Price        string `json:"price"`
	Side         string `json:"side"`
}

type coinbaseScraper struct {
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
	coinbaseWSBaseString = "wss://ws-feed.exchange.coinbase.com"
)

func NewCoinBaseScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("CoinBase - Started scraper.")

	scraper := coinbaseScraper{
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

	// Dial websocket API.
	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(coinbaseWSBaseString, nil)
	if err != nil {
		log.Errorf("CoinBase - Dial ws base string: %v.", err)
		failoverChannel <- string(COINBASE_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs and initialize coinbaseLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("CoinBase - subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("CoinBase - Subscribed to pair %s:%s.", COINBASE_EXCHANGE, pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @coinbaseWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper
}

func (scraper *coinbaseScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("CoinBase - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("CoinBase - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("CoinBase - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *coinbaseScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Keep track of the last config.
	var last map[string]int64

	// Get the initial config.
	cfg, err := models.GetExchangePairMap(COINBASE_EXCHANGE)
	if err != nil {
		log.Errorf("CoinBase - GetExchangePairMap: %v.", err)
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
			cfg, err := models.GetExchangePairMap(COINBASE_EXCHANGE)
			if err != nil {
				log.Errorf("CoinBase - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, cfg)
			// Update the last config.
			last = cfg
		case <-ctx.Done():
			log.Debugf("CoinBase - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *coinbaseScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

	added := make([]string, 0)
	removed := make([]string, 0)

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
	}

	// Unsubscribe from removed pairs.
	for _, p := range removed {
		log.Infof("CoinBase - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("CoinBase - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("CoinBase - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		scraper.subscribeChannel <- ep
		// Start watchdog for this pair.
		scraper.startWatchdogForPair(ctx, lock, ep)
		// Add the pair to the ticker pair map.
		scraper.tickerPairMap[strings.Split(ep.ForeignName, "-")[0]+strings.Split(ep.ForeignName, "-")[1]] = ep.UnderlyingPair
		lock.Lock()
		// Set the last trade time for this pair.
		if _, exists := scraper.lastTradeTimeMap[ep.ForeignName]; !exists {
			scraper.lastTradeTimeMap[ep.ForeignName] = time.Now()
		}
		lock.Unlock()
	}
}

func (scraper *coinbaseScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(COINBASE_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", COINBASE_EXCHANGE, err)
	}
	ep := models.ConstructExchangePair(COINBASE_EXCHANGE, foreignName, delay, idMap)
	return ep, nil
}

func (scraper *coinbaseScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *coinbaseScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *coinbaseScraper) Close(cancel context.CancelFunc) error {
	log.Warn("CoinBase - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *coinbaseScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *coinbaseScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int

	for {

		var message coinBaseWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, COINBASE_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		if message.Type == "match" {
			scraper.handleWSResponse(message, lock)
		}

	}

}

func (scraper *coinbaseScraper) handleWSResponse(message coinBaseWSResponse, lock *sync.RWMutex) {
	trade, err := coinbaseParseTradeMessage(message)
	if err != nil {
		log.Errorf("CoinBase - parseCoinBaseTradeMessage: %s.", err.Error())
		return
	}

	// Identify ticker symbols with underlying assets.
	pair := strings.Split(message.ProductID, "-")
	if len(pair) > 1 {
		trade.QuoteToken = scraper.tickerPairMap[pair[0]+pair[1]].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[pair[0]+pair[1]].BaseToken

		log.Tracef("CoinBase - got trade: %s -- %v -- %v -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
		lock.Lock()
		scraper.lastTradeTimeMap[pair[0]+"-"+pair[1]] = trade.Time
		lock.Unlock()
		scraper.tradesChannel <- trade
	}

}

func (scraper *coinbaseScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("CoinBase - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("CoinBase - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("CoinBase - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("CoinBase - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("CoinBase - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *coinbaseScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}
	a := &coinBaseWSSubscribeMessage{
		Type: subscribeType,
		Channels: []coinBaseChannel{
			{
				Name:       "matches",
				ProductIDs: []string{pair.ForeignName},
			},
		},
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func coinbaseParseTradeMessage(message coinBaseWSResponse) (models.Trade, error) {
	price, err := strconv.ParseFloat(message.Price, 64)
	if err != nil {
		return models.Trade{}, nil
	}
	volume, err := strconv.ParseFloat(message.Size, 64)
	if err != nil {
		return models.Trade{}, nil
	}
	if message.Side == "sell" {
		volume *= -1
	}
	timestamp, err := time.Parse("2006-01-02T15:04:05.000000Z", message.Time)
	if err != nil {
		return models.Trade{}, nil
	}

	foreignTradeID := strconv.Itoa(int(message.TradeID))

	trade := models.Trade{
		Price:          price,
		Volume:         volume,
		Time:           timestamp,
		Exchange:       Exchanges[COINBASE_EXCHANGE],
		ForeignTradeID: foreignTradeID,
	}

	return trade, nil
}
