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

// A OKExWSSubscribeMessage represents a message to subscribe the public/private channel.
type OKExWSSubscribeMessage struct {
	OP   string     `json:"op"`
	Args []OKExArgs `json:"args"`
}

type OKExArgs struct {
	Channel string `json:"channel"`
	InstIDs string `json:"instId"`
}

type Response struct {
	Channel string     `json:"channel"`
	Data    [][]string `json:"data"`
	Binary  int        `json:"binary"`
}

type Responses []Response

type OKExScraper struct {
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
}

type OKEXMarket struct {
	Alias     string `json:"alias"`
	BaseCcy   string `json:"baseCcy"`
	Category  string `json:"category"`
	CtMult    string `json:"ctMult"`
	CtType    string `json:"ctType"`
	CtVal     string `json:"ctVal"`
	CtValCcy  string `json:"ctValCcy"`
	ExpTime   string `json:"expTime"`
	InstID    string `json:"instId"`
	InstType  string `json:"instType"`
	Lever     string `json:"lever"`
	ListTime  string `json:"listTime"`
	LotSz     string `json:"lotSz"`
	MinSz     string `json:"minSz"`
	OptType   string `json:"optType"`
	QuoteCcy  string `json:"quoteCcy"`
	SettleCcy string `json:"settleCcy"`
	State     string `json:"state"`
	Stk       string `json:"stk"`
	TickSz    string `json:"tickSz"`
	Uly       string `json:"uly"`
}

type AllOKEXMarketResponse struct {
	Code string       `json:"code"`
	Data []OKEXMarket `json:"data"`
	Msg  string       `json:"msg"`
}

type OKEXWSResponse struct {
	Arg  []OKExArgs `json:"args"`
	Data []OKEXDATA `json:"data"`
}

type OKEXDATA struct {
	InstID  string `json:"instId"`
	TradeID string `json:"tradeId"`
	Px      string `json:"px"`
	Sz      string `json:"sz"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

var (
	OKExWSBaseString = "wss://ws.okx.com:8443/ws/v5/public"
)

func NewOKExScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("OKEx - Started scraper.")

	scraper := OKExScraper{
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
	wsClient, _, err := wsDialer.Dial(OKExWSBaseString, nil)
	if err != nil {
		log.Errorf("OKEx - Dial ws base string: %v.", err)
		failoverChannel <- string(OKEX_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs and initialize OKExLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("OKEx - subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("OKEx - Subscribed to pair %s:%s.", OKEX_EXCHANGE, pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @OKExWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper
}

func (scraper *OKExScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("OKEx - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("OKEx - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("OKEx - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *OKExScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Keep track of the last config.
	var last map[string]int64

	// Get the initial config.
	cfg, err := models.GetExchangePairMap(OKEX_EXCHANGE)
	if err != nil {
		log.Errorf("OKEx - GetExchangePairMap: %v.", err)
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
			cfg, err := models.GetExchangePairMap(OKEX_EXCHANGE)
			if err != nil {
				log.Errorf("OKEx - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, cfg)
			// Update the last config.
			last = cfg
		case <-ctx.Done():
			log.Debugf("OKEx - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *OKExScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("OKEx - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("OKEx - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("OKEx - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		scraper.subscribeChannel <- ep
		// Start watchdog for this pair.
		scraper.startWatchdogForPair(ctx, lock, ep)
		key := strings.ReplaceAll(ep.ForeignName, "-", "")
		// Add the pair to the ticker pair map.
		scraper.tickerPairMap[key] = ep.UnderlyingPair
		lock.Lock()
		// Set the last trade time for this pair.
		if _, exists := scraper.lastTradeTimeMap[ep.ForeignName]; !exists {
			scraper.lastTradeTimeMap[ep.ForeignName] = time.Now()
		}
		lock.Unlock()
	}
	// Resubscribe to changed pairs.
	for _, p := range changed {
		newDelay := current[p]
		log.Infof("OKEx - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *OKExScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("OKEx - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *OKExScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(OKEX_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", OKEX_EXCHANGE, err)
	}
	ep := models.ConstructExchangePair(OKEX_EXCHANGE, foreignName, delay, idMap)
	return ep, nil
}

func (scraper *OKExScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *OKExScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *OKExScraper) Close(cancel context.CancelFunc) error {
	log.Warn("OKEx - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *OKExScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *OKExScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int

	for {

		var message OKEXWSResponse
		messageType, messageTemp, err := scraper.wsClient.ReadMessage()
		if err != nil {
			if handleErrorReadMessage(err, &errCount, scraper.maxErrCount, OKEX_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		} else {
			switch messageType {
			case ws.TextMessage:
				err := json.Unmarshal(messageTemp, &message)
				if err != nil {
					log.Errorln("Error parsing reponse")
				}
				if len(message.Data) > 0 {
					scraper.handleWSResponse(message.Data[0], lock)
				}
			}
		}
	}
}

func (scraper *OKExScraper) handleWSResponse(data OKEXDATA, lock *sync.RWMutex) {
	trade, err := OKExParseTradeMessage(data)
	if err != nil {
		log.Errorf("OKEx - parseOKExTradeMessage: %s.", err.Error())
		return
	}
	// Identify ticker symbols with underlying assets.
	ep := data.InstID
	pair := strings.Split(ep, "-")
	if len(pair) > 1 {
		trade.QuoteToken = scraper.tickerPairMap[pair[0]+pair[1]].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[pair[0]+pair[1]].BaseToken

		log.Tracef("OKEx - got trade: %s -- %v -- %v -- %s -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID, trade.Time)
		lock.Lock()
		scraper.lastTradeTimeMap[pair[0]+"-"+pair[1]] = trade.Time
		lock.Unlock()
		scraper.tradesChannel <- trade
	}
}

func (scraper *OKExScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("OKEx - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("OKEx - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("OKEx - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("OKEx - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("OKEx - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *OKExScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	var allPairs []OKExArgs

	allPairs = append(allPairs, OKExArgs{Channel: "trades", InstIDs: pair.ForeignName})

	a := &OKExWSSubscribeMessage{
		OP:   subscribeType,
		Args: allPairs,
	}

	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func OKExParseTradeMessage(message OKEXDATA) (models.Trade, error) {
	price, err := strconv.ParseFloat(message.Px, 64)
	if err != nil {
		return models.Trade{}, nil
	}
	volume, err := strconv.ParseFloat(message.Sz, 64)
	if err != nil {
		return models.Trade{}, nil
	}
	if message.Side == "sell" {
		volume *= -1
	}

	tsInt, err := strconv.ParseInt(message.Ts, 10, 64)
	if err != nil {
		return models.Trade{}, nil
	}
	timestamp := time.UnixMilli(tsInt)

	foreignTradeID := message.TradeID

	trade := models.Trade{
		Price:          price,
		Volume:         volume,
		Time:           timestamp,
		Exchange:       Exchanges[OKEX_EXCHANGE],
		ForeignTradeID: foreignTradeID,
	}

	return trade, nil
}

// If @handleErrorReadMessage returns true, the calling function should return. Otherwise continue.
func handleErrorReadMessage(err error, errCount *int, maxErrCount int, exchange string, restartWaitTime int) bool {
	log.Errorf("%s - ReadMessage: %v", exchange, err)
	*errCount++

	if strings.Contains(err.Error(), "closed network connection") {
		return true
	}

	if *errCount > maxErrCount {
		log.Warnf("%s - too many errors. wait for %v seconds and restart scraper.", exchange, restartWaitTime)
		time.Sleep(time.Duration(restartWaitTime) * time.Second)
		return true
	}

	return false
}
