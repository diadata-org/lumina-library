package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// A WebSocketSubscribeMessage represents a message to subscribe the public/private channel.
type kuCoinWSSubscribeMessage struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	PrivateChannel bool   `json:"privateChannel"`
	Response       bool   `json:"response"`
}

type kuCoinWSResponse struct {
	Type    string       `json:"type"`
	Topic   string       `json:"topic"`
	Subject string       `json:"subject"`
	Data    kuCoinWSData `json:"data"`
}

type kuCoinWSData struct {
	Sequence string `json:"sequence"`
	Type     string `json:"type"`
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
	Price    string `json:"price"`
	Size     string `json:"size"`
	TradeID  string `json:"tradeId"`
	Time     string `json:"time"`
}

type kucoinScraper struct {
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
	kucoinWSBaseString    = "wss://ws-api-spot.kucoin.com/"
	kucoinTokenURL        = "https://api.kucoin.com/api/v1/bullet-public"
	kucoinPingIntervalFix = int64(10)
)

func NewKuCoinScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("KuCoin - Started scraper.")

	token, pingInterval, err := getPublicKuCoinToken(kucoinTokenURL)
	if err != nil {
		log.Errorf("KuCoin - getPublicKuCoinToken: %v.", err)
	}

	scraper := kucoinScraper{
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
	wsClient, _, err := wsDialer.Dial(kucoinWSBaseString+"?token="+token, nil)
	if err != nil {
		log.Errorf("KuCoin - Dial ws base string: %v.", err)
		failoverChannel <- string(KUCOIN_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("KuCoin - Subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("KuCoin - Subscribe to pair %s.", pair.ForeignName)
		}
	}

	go scraper.ping(ctx, pingInterval, time.Now(), &lock)
	go scraper.fetchTrades(&lock)
	go scraper.resubscribe(ctx, &lock)
	go scraper.processUnsubscribe(ctx, &lock)
	go scraper.watchConfig(ctx, &lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @kucoinWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &lock, pair)
	}

	return &scraper

}

func (scraper *kucoinScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int
	for {

		var message kuCoinWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, KUCOIN_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		if message.Type == "pong" {
			log.Debug("KuCoin - Successful ping: received pong.")
		} else if message.Type == "message" {
			scraper.handleWSResponse(message, lock)
		}

	}
}

func (scraper *kucoinScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false, lock); err != nil {
				log.Errorf("KuCoin - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("KuCoin - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("KuCoin - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *kucoinScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	// Check for config changes every 30 seconds.
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Keep track of the last config.
	var last map[string]int64

	// Get the initial config.
	cfg, err := models.GetExchangePairMap(KUCOIN_EXCHANGE)
	if err != nil {
		log.Errorf("KuCoin - GetExchangePairMap: %v.", err)
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
			cfg, err := models.GetExchangePairMap(KUCOIN_EXCHANGE)
			if err != nil {
				log.Errorf("KuCoin - GetExchangePairMap: %v.", err)
				continue
			}
			// Apply the config changes.
			scraper.applyConfigDiff(ctx, lock, last, cfg)
			// Update the last config.
			last = cfg
		case <-ctx.Done():
			log.Debugf("KuCoin - Close watchConfig routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *kucoinScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

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
		log.Infof("KuCoin - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("KuCoin - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("KuCoin - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		scraper.subscribeChannel <- ep
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
		log.Infof("KuCoin - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *kucoinScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("KuCoin - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *kucoinScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(KUCOIN_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", KUCOIN_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(KUCOIN_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", KUCOIN_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *kucoinScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
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

func (scraper *kucoinScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	lock.Lock()
	cancel, ok := scraper.watchdogCancel[foreignName]
	if ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, foreignName)
	}
	lock.Unlock()
}

func (scraper *kucoinScraper) handleWSResponse(message kuCoinWSResponse, lock *sync.RWMutex) {
	// Parse trade quantities.
	price, volume, timestamp, foreignTradeID, err := parseKuCoinTradeMessage(message)
	if err != nil {
		log.Errorf("KuCoin - parseTradeMessage: %v.", err)
		return
	}

	// Identify ticker symbols with underlying assets.
	pair := strings.Split(message.Data.Symbol, "-")
	if len(pair) < 2 {
		log.Warnf("KuCoin - Unexpected symbol format: %q", message.Data.Symbol)
		return
	}

	lock.RLock()
	exchangepair := scraper.tickerPairMap[pair[0]+pair[1]]
	lock.RUnlock()

	trade := models.Trade{
		QuoteToken:     exchangepair.QuoteToken,
		BaseToken:      exchangepair.BaseToken,
		Price:          price,
		Volume:         volume,
		Time:           timestamp,
		Exchange:       Exchanges[KUCOIN_EXCHANGE],
		ForeignTradeID: foreignTradeID,
	}

	lock.Lock()
	scraper.lastTradeTimeMap[pair[0]+"-"+pair[1]] = trade.Time
	lock.Unlock()

	log.Tracef("KuCoin - got trade: %s -- %v -- %v -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
	scraper.tradesChannel <- trade
}

func (scraper *kucoinScraper) Close(cancel context.CancelFunc) error {
	log.Warn("KuCoin - Call scraper.Close()")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *kucoinScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *kucoinScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("KuCoin - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("KuCoin - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("KuCoin - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("KuCoin - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("KuCoin - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *kucoinScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	a := &kuCoinWSSubscribeMessage{
		Type:  subscribeType,
		Topic: "/market/match:" + pair.ForeignName,
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func parseKuCoinTradeMessage(message kuCoinWSResponse) (price float64, volume float64, timestamp time.Time, foreignTradeID string, err error) {
	price, err = strconv.ParseFloat(message.Data.Price, 64)
	if err != nil {
		return
	}
	volume, err = strconv.ParseFloat(message.Data.Size, 64)
	if err != nil {
		return
	}
	if message.Data.Side == "sell" {
		volume *= -1
	}
	timeMilliseconds, err := strconv.Atoi(message.Data.Time)
	if err != nil {
		return
	}
	timestamp = time.Unix(0, int64(timeMilliseconds))
	foreignTradeID = message.Data.TradeID
	return
}

// A WebSocketMessage represents a message between the WebSocket client and server.
type kuCoinWSMessage struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type kuCoinPostResponse struct {
	Code string `json:"code"`
	Data struct {
		Token           string            `json:"token"`
		InstanceServers []instanceServers `json:"instanceServers"`
	} `json:"data"`
}

type instanceServers struct {
	PingInterval int64 `json:"pingInterval"`
}

// Send ping to server.
func (scraper *kucoinScraper) ping(ctx context.Context, pingInterval int64, starttime time.Time, lock *sync.RWMutex) {
	var ping kuCoinWSMessage
	ping.Type = "ping"
	tick := time.NewTicker(time.Duration(kucoinPingIntervalFix) * time.Second)

	for {
		select {
		case <-tick.C:
			lock.Lock()
			if err := scraper.wsClient.WriteJSON(ping); err != nil {
				log.Errorf("KuCoin - Send ping: %s.", err.Error())
				lock.Unlock()
				return
			}
			lock.Unlock()
		case <-ctx.Done():
			log.Warn("KuCoin - Close ping.")
			return
		}
	}
}

// getPublicKuCoinToken returns a token for public market data along with the pingInterval in seconds.
func getPublicKuCoinToken(url string) (token string, pingInterval int64, err error) {
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)
	data, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return
	}
	defer data.Body.Close()
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return
	}

	var postResp kuCoinPostResponse
	err = json.Unmarshal(body, &postResp)
	if err != nil {
		return
	}
	if len(postResp.Data.InstanceServers) > 0 {
		pingInterval = postResp.Data.InstanceServers[0].PingInterval
	}
	token = postResp.Data.Token
	return
}
