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
	"github.com/diadata-org/lumina-library/scrapers/mexcproto"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type WSConnection struct {
	wsConn           *ws.Conn
	numSubscriptions int
}

type MEXCScraper struct {
	connections        []WSConnection
	tradesChannel      chan models.Trade
	subscribeChannel   chan models.ExchangePair
	unsubscribeChannel chan models.ExchangePair
	watchdogCancel     map[string]context.CancelFunc
	tickerPairMap      map[string]models.Pair
	lastTradeTimeMap   map[string]time.Time
	maxErrCount        int
	restartWaitTime    int
	pingPeriod         int
	genesis            time.Time
	maxSubscriptions   int
	pairConnIndex      map[string]int
	mu                 sync.RWMutex
}

var (
	MEXCWSBaseString       = "wss://wbs-api.mexc.com/ws"
	maxSubscriptionPerConn = 20
)

func NewMEXCScraper(ctx context.Context, pairs []models.ExchangePair, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	log.Info("MEXC - Started scraper.")

	pingPeriod, err := strconv.Atoi(utils.Getenv("MEXC_PING_PERIOD_SECONDS", "15"))
	if err != nil {
		log.Errorf("parse MEXC_PING_PERIOD_SECONDS: %v. Set to default 15", err)
		pingPeriod = 15
	}

	scraper := MEXCScraper{
		tradesChannel:      make(chan models.Trade),
		subscribeChannel:   make(chan models.ExchangePair),
		unsubscribeChannel: make(chan models.ExchangePair),
		watchdogCancel:     make(map[string]context.CancelFunc),
		tickerPairMap:      models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:   make(map[string]time.Time),
		maxErrCount:        20,
		restartWaitTime:    5,
		pingPeriod:         pingPeriod,
		genesis:            time.Now(),
		maxSubscriptions:   maxSubscriptionPerConn,
		connections:        make([]WSConnection, 0),
		pairConnIndex:      make(map[string]int),
		mu:                 sync.RWMutex{},
	}

	if _, err := scraper.newConn(ctx); err != nil {
		log.Errorf("MEXC - newConn failed: %v.", err)
		return &scraper
	}

	go scraper.pingServer()

	// Subscribe to pairs and initialize MEXCLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true); err != nil {
			log.Errorf("MEXC - subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("MEXC - Subscribed to pair %s:%s.", MEXC_EXCHANGE, pair.ForeignName)
			scraper.mu.Lock()
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
			scraper.mu.Unlock()
		}
	}

	for _, conn := range scraper.connections {
		go scraper.fetchTrades(conn)
	}
	go scraper.resubscribe(ctx)
	go scraper.processUnsubscribe(ctx, &scraper.mu)
	go scraper.watchConfig(ctx, &scraper.mu)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @MEXCWatchdogDelayMap.
	for _, pair := range pairs {
		scraper.startWatchdogForPair(ctx, &scraper.mu, pair)
	}

	return &scraper
}

func (scraper *MEXCScraper) pingServer() {

	pingMsg := map[string]string{"method": "PING"}
	ticker := time.NewTicker(time.Duration(scraper.pingPeriod) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		scraper.mu.Lock()
		for i, c := range scraper.connections {
			if err := c.wsConn.WriteJSON(pingMsg); err != nil {
				log.Errorf("MEXC - ping error for connection %v: %v", i, err)
				continue
			}
			log.Infof("MEXC - Sent Ping to connection %v.", i)
		}
		scraper.mu.Unlock()
		log.Infof("MEXC - Sent Pings")
	}

}

func (s *MEXCScraper) newConn(ctx context.Context) (*WSConnection, error) {
	var wsDialer ws.Dialer
	numRetries := 1

	for {
		// allow cancellation
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("MEXC - context canceled while dialing: %w", ctx.Err())
		default:
		}

		wsClient, _, err := wsDialer.Dial(MEXCWSBaseString, nil)
		if err == nil {
			conn := WSConnection{
				wsConn:           wsClient,
				numSubscriptions: 0,
			}
			s.mu.Lock()
			s.connections = append(s.connections, conn)
			idx := len(s.connections) - 1
			s.mu.Unlock()

			log.Debugf("MEXC - New WS connection established. Total connections: %d", len(s.connections))
			return &s.connections[idx], nil
		}

		log.Errorf("MEXC - Dial ws base string failed (try %d): %v", numRetries, err)
		// backoff: numRetries * restartWaitTime
		time.Sleep(time.Duration(numRetries*s.restartWaitTime) * time.Second)
		numRetries++
	}
}

func (scraper *MEXCScraper) processUnsubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.unsubscribeChannel:
			// Unsubscribe from this pair.
			if err := scraper.subscribe(pair, false); err != nil {
				log.Errorf("MEXC - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Infof("MEXC - Unsubscribed pair %s.", pair.ForeignName)
			}
			// Delete last trade time for this pair.
			lock.Lock()
			delete(scraper.lastTradeTimeMap, pair.ForeignName)
			lock.Unlock()
			scraper.stopWatchdogForPair(lock, pair.ForeignName)
		case <-ctx.Done():
			log.Debugf("MEXC - Close processUnsubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *MEXCScraper) watchConfig(ctx context.Context, lock *sync.RWMutex) {
	go WatchConfigLoop(ctx, MEXC_EXCHANGE, 30, func(ctx context.Context, last, current map[string]int64) {
		scraper.applyConfigDiff(ctx, lock, last, current)
	})
}

func (scraper *MEXCScraper) applyConfigDiff(ctx context.Context, lock *sync.RWMutex, last map[string]int64, current map[string]int64) {

	added, removed, changed := diffPairMap(last, current)

	// Unsubscribe from removed pairs.
	for _, p := range removed {
		log.Infof("MEXC - Removed pair %s.", p)
		scraper.unsubscribeChannel <- models.ExchangePair{
			ForeignName: p,
		}
	}
	// Subscribe to added pairs.
	for _, p := range added {
		// Get the delay for this pair.
		delay := current[p]
		log.Infof("MEXC - Added pair %s with delay %v.", p, delay)

		ep, err := scraper.getExchangePairInfo(p, delay)
		if err != nil {
			log.Errorf("MEXC - Failed to GetExchangePairInfo for new pair %s: %v.", p, err)
			continue
		}
		err = scraper.subscribe(ep, true)
		if err != nil {
			log.Errorf("MEXC - Failed to subscribe to %s: %v", ep.ForeignName, err)
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
		log.Infof("MEXC - Changed pair %s with delay %v.", p, newDelay)
		scraper.restartWatchdogForPair(ctx, lock, p, newDelay)
	}
}

func (scraper *MEXCScraper) restartWatchdogForPair(ctx context.Context, lock *sync.RWMutex, foreignName string, newDelay int64) {
	// 1. Stop the watchdog for the pair.
	scraper.stopWatchdogForPair(lock, foreignName)
	// 2. Get the new exchange pair info (only for watchdog, no effect on subscription).
	ep, err := scraper.getExchangePairInfo(foreignName, newDelay)
	if err != nil {
		log.Errorf("MEXC - Failed to GetExchangePairInfo for changed pair %s: %v.", foreignName, err)
		return
	}
	// 3. Start the watchdog for the pair with the new delay.
	scraper.startWatchdogForPair(ctx, lock, ep)
}

func (scraper *MEXCScraper) getExchangePairInfo(foreignName string, delay int64) (models.ExchangePair, error) {
	idMap, err := models.GetSymbolIdentificationMap(MEXC_EXCHANGE)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", MEXC_EXCHANGE, err)
	}
	ep, err := models.ConstructExchangePair(MEXC_EXCHANGE, foreignName, delay, idMap)
	if err != nil {
		return models.ExchangePair{}, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", MEXC_EXCHANGE, foreignName, delay, err)
	}
	return ep, nil
}

func (scraper *MEXCScraper) startWatchdogForPair(ctx context.Context, lock *sync.RWMutex, pair models.ExchangePair) {
	StartWatchdogForPair(
		ctx, lock, pair,
		scraper.watchdogCancel,
		scraper.lastTradeTimeMap,
		scraper.subscribeChannel,
	)
}

func (scraper *MEXCScraper) stopWatchdogForPair(lock *sync.RWMutex, foreignName string) {
	StopWatchdogForPair(lock, foreignName, scraper.watchdogCancel)
}

func (scraper *MEXCScraper) Close(cancel context.CancelFunc) error {
	log.Warn("MEXC - call scraper.Close().")
	cancel()

	for i, conn := range scraper.connections {
		if conn.wsConn != nil {
			err := conn.wsConn.Close()
			if err != nil {
				log.Errorf("MEXC - failed to close connection %d: %v", i, err)
			} else {
				log.Infof("MEXC - closed connection %d", i)
			}
		}
	}

	return nil
}

func (scraper *MEXCScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *MEXCScraper) fetchTrades(conn WSConnection) {
	// Read trades stream.
	var errCount int

	for {
		_, payload, err := conn.wsConn.ReadMessage()
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, MEXC_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}

		switch {
		case len(payload) > 0 && (payload[0] == '{' || payload[0] == '['):
			var msg map[string]any
			if err := json.Unmarshal(payload, &msg); err != nil {
				log.Errorf("failed to parse JSON: %v", err)
				return
			}
			log.Debugf("Received JSON message: %+v", msg["msg"])

		default:
			decodedMessage := &mexcproto.PushDataV3ApiWrapper{}
			if err := proto.Unmarshal(payload, decodedMessage); err != nil {
				log.Println("protobuf unmarshal error:", err)
				continue
			}

			// Received Message: channel:"spot@public.aggre.deals.v3.api.pb@100ms@BTCUSDC"
			// publicAggreDeals:{deals:{price:"115069.04"  quantity:"0.000042"  tradeType:1
			// time:1755519868593}  deals:{price:"115069"  quantity:"0.000011"  tradeType:2
			// time:1755519868593}  deals:{price:"115069.08"  quantity:"0.00003"  tradeType:1
			// time:1755519868593}  deals:{price:"115069.02"  quantity:"0.000032"  tradeType:1
			// time:1755519868593}  deals:{price:"115069.01"  quantity:"0.000035"  tradeType:2
			// time:1755519868594}  deals:{price:"115069.06"  quantity:"0.000034"  tradeType:1
			// time:1755519868594}  eventType:"spot@public.aggre.deals.v3.api.pb@100ms"}
			// symbol:"BTCUSDC"  sendTime:1755519868617

			ch := strings.ToLower(decodedMessage.GetChannel())
			sym := decodedMessage.GetSymbol()
			tradeID := decodedMessage.GetSendTime()

			switch {
			case strings.Contains(ch, "public.aggre.deals.v3.api.pb"):
				dealsMsg := decodedMessage.GetPublicAggreDeals()
				if dealsMsg == nil {
					log.Debug("aggre.deals wrapper has no PublicAggreDeals payload")
					break
				}
				for idx, trade := range dealsMsg.GetDeals() {
					scraper.handleWSResponse(trade, sym, tradeID+int64(idx))
				}
			default:
				// handle other channels if you subscribe to them later
				log.Debugf("unhandled channel: %s", decodedMessage.GetChannel())
			}
		}
	}
}

func (scraper *MEXCScraper) handleWSResponse(message *mexcproto.PublicAggreDealsV3ApiItem, pair string, tradeID int64) {
	trade, err := MEXCParseTradeMessage(message, tradeID)
	if err != nil {
		log.Errorf("MEXC - parseMEXCTradeMessage: %s.", err.Error())
		return
	}

	// Identify ticker symbols with underlying assets.
	if pair != "" {
		scraper.mu.RLock()
		trade.QuoteToken = scraper.tickerPairMap[pair].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[pair].BaseToken
		scraper.mu.RUnlock()

		log.Tracef("MEXC - got trade: %s -- %v -- %v -- %s -- %v.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID, trade.Time)
		scraper.mu.Lock()
		scraper.lastTradeTimeMap[trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol] = trade.Time
		scraper.mu.Unlock()
		scraper.tradesChannel <- trade
	}

}

func (s *MEXCScraper) resubscribe(ctx context.Context) {
	for {
		select {
		case pair := <-s.subscribeChannel:
			err := s.subscribe(pair, false) // unsubscribe first
			if err != nil {
				log.Errorf("MEXC - Unsubscribe failed for %s: %v", pair.ForeignName, err)
			}
			time.Sleep(2 * time.Second)

			err = s.subscribe(pair, true)
			if err != nil {
				log.Errorf("MEXC - Resubscribe failed for %s: %v", pair.ForeignName, err)
			}

		case <-ctx.Done():
			log.Debugf("MEXC - Close resubscribe routine of scraper with genesis: %v.", s.genesis)
			return
		}
	}
}

func (s *MEXCScraper) subscribe(pair models.ExchangePair, subscribe bool) error {
	foreignName := strings.ReplaceAll(pair.ForeignName, "-", "")
	topic := "spot@public.aggre.deals.v3.api.pb@100ms@" + foreignName

	subscriptionMessage := map[string]interface{}{
		"method": "UNSUBSCRIPTION",
		"params": []string{topic},
	}
	if subscribe {
		subscriptionMessage["method"] = "SUBSCRIPTION"
	}

	subMsg, _ := json.Marshal(subscriptionMessage)

	if subscribe {
		// Try to find an existing connection with available slots
		var targetConnID int = -1
		for i, conn := range s.connections {
			if conn.numSubscriptions < s.maxSubscriptions {
				targetConnID = i
				break
			}
		}

		// If all are full, create a new connection
		if targetConnID == -1 {
			if _, err := s.newConn(context.Background()); err != nil {
				log.Errorf("MEXC - Failed to create new connection for %s: %v", pair.ForeignName, err)
				return err
			}
			targetConnID = len(s.connections) - 1
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		err := s.connections[targetConnID].wsConn.WriteMessage(ws.TextMessage, subMsg)
		if err != nil {
			log.Errorf("MEXC - Failed to send SUBSCRIPTION message for %s: %v", pair.ForeignName, err)
			return err
		}

		s.connections[targetConnID].numSubscriptions++
		s.pairConnIndex[pair.ForeignName] = targetConnID
		log.Debugf("MEXC - Subscribed to %s on connection %d", pair.ForeignName, targetConnID)
	} else {
		// Unsubscribe logic
		s.mu.Lock()
		connID, ok := s.pairConnIndex[pair.ForeignName]
		s.mu.Unlock()
		if !ok {
			log.Warnf("MEXC - No connection found for pair %s to unsubscribe", pair.ForeignName)
			return nil
		}

		s.mu.Lock()
		err := s.connections[connID].wsConn.WriteMessage(ws.TextMessage, subMsg)
		if err != nil {
			log.Errorf("MEXC - Failed to send UNSUBSCRIPTION message for %s: %v", pair.ForeignName, err)
			return err
		}

		s.connections[connID].numSubscriptions--
		delete(s.pairConnIndex, pair.ForeignName)
		s.mu.Unlock()
		log.Infof("MEXC - Unsubscribed from %s on connection %d", pair.ForeignName, connID)
	}

	return nil
}

func MEXCParseTradeMessage(message *mexcproto.PublicAggreDealsV3ApiItem, tradeID int64) (models.Trade, error) {
	price, err := strconv.ParseFloat(message.GetPrice(), 64)
	if err != nil {
		return models.Trade{}, nil
	}
	volume, err := strconv.ParseFloat(message.GetQuantity(), 64)
	if err != nil {
		return models.Trade{}, nil
	}
	if message.GetTradeType() == 2 {
		volume *= -1
	}
	timestamp := time.Unix(0, message.GetTime()*int64(time.Millisecond))
	timestamp.Format(time.RFC3339)

	foreignTradeID := strconv.Itoa(int(tradeID))

	trade := models.Trade{
		Price:          price,
		Volume:         volume,
		Time:           timestamp,
		Exchange:       Exchanges[MEXC_EXCHANGE],
		ForeignTradeID: foreignTradeID,
	}

	return trade, nil
}
