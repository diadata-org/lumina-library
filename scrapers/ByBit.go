package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

type byBitWSSubscribeMessage struct {
	OP   string   `json:"op"`
	Args []string `json:"args"`
}

type byBitWSResponse struct {
	Topic     string                   `json:"topic"`
	Timestamp int64                    `json:"ts"`
	Type      string                   `json:"type"`
	Data      []ByBitTradeResponseData `json:"data"`
}

type ByBitTradeResponseData struct {
	TradeID   string `json:"i"`
	Timestamp int64  `json:"T"`
	Price     string `json:"p"`
	Size      string `json:"v"`
	Side      string `json:"S"`
	Symbol    string `json:"s"`
}

type byBitScraper struct {
	wsClient           *ws.Conn
	tradesChannel      chan models.Trade
	subscribeChannel   chan models.ExchangePair
	unsubscribeChannel chan models.ExchangePair
	watchdogCancel     map[string]context.CancelFunc
	tickerPairMap      map[string]models.Pair
	lastTradeTimeMap   map[string]time.Time
	maxErrCount        int
	restartWaitTime    int
	genesis            time.Time
	writeLock          sync.Mutex
}

const byBitWSBaseURL = "wss://stream.bybit.com/v5/public/spot"

func NewByBitScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("Bybit - Started scraper.")

	scraper := byBitScraper{
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

	var dialer ws.Dialer
	conn, _, err := dialer.Dial(byBitWSBaseURL, nil)
	if err != nil {
		log.Errorf("ByBit - WebSocket connection failed: %v", err)
		failoverChannel <- string(BYBIT_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = conn

	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("ByBit - Failed to subscribe to %v: %v", pair.ForeignName, err)
		} else {
			log.Infof("ByBit - Subscribed to %v", pair)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.startByBitPing(ctx)

	go scraper.fetchTrades(ctx, &lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper
}

func (scraper *byBitScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("ByBit - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("ByBit - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("ByBit - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *byBitScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Keep track of the last config.
	var last map[string]int64

	// Get the initial config.
	cfg, err := models.GetExchangePairMap(BYBIT_EXCHANGE)
	if err != nil {
		log.Errorf("ByBit - GetExchangePairMap: %v.", err)
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
			cfg, err := models.GetExchangePairMap(BYBIT_EXCHANGE)
			if err != nil {
				log.Errorf("ByBit - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, cfg)
			// Update the last config.
			last = cfg
		case <-ctx.Done():
			log.Debugf("ByBit - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *byBitScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("ByBit - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("ByBit - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("ByBit - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
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

func (scraper *byBitScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(BYBIT_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", BYBIT_EXCHANGE, err)
	}
	ep := models.ConstructExchangePair(BYBIT_EXCHANGE, foreignName, delay, idMap)
	return ep, nil
}

func (scraper *byBitScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *byBitScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *byBitScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *byBitScraper) Close(cancel context.CancelFunc) error {
	log.Warn("ByBit - call scraper.Close().")
	cancel()
	if scraper.wsClient != nil {
		return scraper.wsClient.Close()
	}
	return nil
}

func (scraper *byBitScraper) startByBitPing(ctx context.Context) {
	log.Info("ByBit - Sent Ping.")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("ByBit - Ping routine stopped.")
			return
		case <-ticker.C:
			err := scraper.safeWriteMessage(ws.PingMessage, []byte{})
			if err != nil {
				log.Errorf("ByBit - Ping error: %v", err)
				return
			}

		}
	}
}

func (scraper *byBitScraper) fetchTrades(ctx context.Context, lock *sync.RWMutex) {
	var errCount int

	for {
		select {
		case <-ctx.Done():
			log.Infof("ByBit - Stopping WebSocket reader")
			return
		default:
			messageType, message, err := scraper.wsClient.ReadMessage()
			if err != nil {
				log.Errorf("ByBit - ReadMessage error: %v", err)
				if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, BYBIT_EXCHANGE, scraper.restartWaitTime) {
					return
				}
				continue
			}
			if messageType != ws.TextMessage {
				log.Warnf("ByBit - Non-text WebSocket message received, type: %d", messageType)
				continue
			}
			scraper.handleMessage(message, lock)
		}
	}
}

func (scraper *byBitScraper) handleMessage(message []byte, lock *sync.RWMutex) {
	if strings.Contains(string(message), "\"success\"") {
		log.Infof("ByBit - Subscription success ack: %s", string(message))
		return
	}

	var resp byBitWSResponse
	if err := json.Unmarshal(message, &resp); err != nil {
		log.Errorf("ByBit - Failed to unmarshal message: %v", err)
		return
	}

	if resp.Type != "snapshot" {
		return
	}

	for _, data := range resp.Data {
		price, err := strconv.ParseFloat(data.Price, 64)
		if err != nil {
			log.Errorf("ByBit - Invalid price: %v", err)
			return
		}

		volume, err := strconv.ParseFloat(data.Size, 64)
		side := data.Side
		if side == "Sell" {
			volume = -1 * volume
		}

		if err != nil {
			log.Errorf("ByBit - Invalid volume: %v", err)
			return
		}

		pairName := data.Symbol
		lock.Lock()
		pair, exists := scraper.tickerPairMap[pairName]
		lock.Unlock()
		if !exists {
			log.Warnf("ByBit - Unknown pair: %s", pairName)
			return
		}

		trade := models.Trade{
			Price:      price,
			Volume:     volume,
			Time:       time.Now(),
			Exchange:   Exchanges[BYBIT_EXCHANGE],
			BaseToken:  pair.BaseToken,
			QuoteToken: pair.QuoteToken,
		}

		lock.Lock()
		scraper.lastTradeTimeMap[pair.QuoteToken.Symbol+"-"+pair.BaseToken.Symbol] = time.Now()
		lock.Unlock()

		scraper.tradesChannel <- trade
		log.Tracef("ByBit - Trade: %s-%s | Side: %v | Price: %f | Volume: %f", pair.QuoteToken.Symbol, pair.BaseToken.Symbol, side, price, volume)
	}
}

func (scraper *byBitScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("ByBit - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("ByBit - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("ByBit - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("ByBit - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("ByBit - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *byBitScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}
	topic := "publicTrade." + strings.ReplaceAll(pair.ForeignName, "-", "")
	subscribeMsg := byBitWSSubscribeMessage{
		OP:   subscribeType,
		Args: []string{topic},
	}
	lock.Lock()
	return scraper.safeWriteJSON(subscribeMsg)
}

func (scraper *byBitScraper) safeWriteJSON(v interface{}) error {
	scraper.writeLock.Lock()
	defer scraper.writeLock.Unlock()
	return scraper.wsClient.WriteJSON(v)
}

func (scraper *byBitScraper) safeWriteMessage(messageType int, data []byte) error {
	scraper.writeLock.Lock()
	defer scraper.writeLock.Unlock()
	return scraper.wsClient.WriteMessage(messageType, data)
}
