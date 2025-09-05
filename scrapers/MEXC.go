package scrapers

import (
	"context"
	"encoding/json"
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
	connections      []WSConnection
	tradesChannel    chan models.Trade
	subscribeChannel chan models.ExchangePair
	tickerPairMap    map[string]models.Pair
	lastTradeTimeMap map[string]time.Time
	maxErrCount      int
	restartWaitTime  int
	genesis          time.Time
	maxSubscriptions int
	pairConnIndex    map[string]int
	mu               sync.RWMutex
}

var (
	MEXCWSBaseString       = "wss://wbs-api.mexc.com/ws"
	maxSubscriptionPerConn = 20
)

func NewMEXCScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	log.Info("MEXC - Started scraper.")

	scraper := MEXCScraper{
		tradesChannel:    make(chan models.Trade),
		subscribeChannel: make(chan models.ExchangePair),
		tickerPairMap:    models.MakeTickerPairMap(pairs),
		lastTradeTimeMap: make(map[string]time.Time),
		maxErrCount:      20,
		restartWaitTime:  5,
		genesis:          time.Now(),
		maxSubscriptions: maxSubscriptionPerConn,
		connections:      make([]WSConnection, 0),
		pairConnIndex:    make(map[string]int),
		mu:               sync.RWMutex{},
	}

	if _, err := scraper.newConn(failoverChannel); err != nil {
		log.Errorf("MEXC - newConn failed: %v.", err)
		return &scraper
	}

	go func() {
		pingMsg := map[string]string{"method": "PING"}
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Infof("MEXC - Sent Ping...")
			scraper.mu.Lock()
			conns := make([]*ws.Conn, 0, len(scraper.connections))
			for _, c := range scraper.connections {
				if c.wsConn != nil {
					conns = append(conns, c.wsConn)
				}
			}
			scraper.mu.Unlock()

			for _, c := range conns {
				scraper.mu.Lock()
				if err := c.WriteJSON(pingMsg); err != nil {
					log.Error("ping error: ", err)
				}
				scraper.mu.Unlock()
			}
			// for _, c := range scraper.connections {
			// 	scraper.mu.Lock()
			// 	if err := c.wsConn.WriteJSON(pingMsg); err != nil {
			// 		log.Error("ping error: ", err)
			// 		return
			// 	}
			// 	scraper.mu.Unlock()
			// }
		}
	}()

	// Subscribe to pairs and initialize MEXCLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, failoverChannel); err != nil {
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

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @MEXCWatchdogDelayMap.
	for _, pair := range pairs {
		envVar := strings.ToUpper(MEXC_EXCHANGE) + "_WATCHDOG_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[0] + "_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[1]
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "300"), 10, 64)
		if err != nil {
			log.Errorf("MEXC - Parse MEXCWatchdogDelay: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdog(ctx, pair, watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, &scraper.mu)
		go scraper.resubscribe(ctx, failoverChannel)
	}

	return &scraper
}

func (s *MEXCScraper) newConn(failoverChannel chan string) (*WSConnection, error) {
	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(MEXCWSBaseString, nil)
	if err != nil {
		log.Errorf("MEXC - Failed to open WebSocket connection: %v", err)
		failoverChannel <- string(MEXC_EXCHANGE)
		return nil, err
	}
	conn := WSConnection{
		wsConn:           wsClient,
		numSubscriptions: 0,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections = append(s.connections, conn)
	log.Debugf("MEXC - New WS connection established. Total connections: %d", len(s.connections))
	return &s.connections[len(s.connections)-1], nil
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
		trade.QuoteToken = scraper.tickerPairMap[pair].QuoteToken
		trade.BaseToken = scraper.tickerPairMap[pair].BaseToken

		log.Tracef("MEXC - got trade: %s -- %v -- %v -- %s -- %v.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID, trade.Time)
		scraper.mu.Lock()
		scraper.lastTradeTimeMap[trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol] = trade.Time
		scraper.mu.Unlock()
		scraper.tradesChannel <- trade
	}

}

func (s *MEXCScraper) resubscribe(ctx context.Context, failoverChannel chan string) {
	for {
		select {
		case pair := <-s.subscribeChannel:
			err := s.subscribe(pair, false, failoverChannel) // unsubscribe first
			if err != nil {
				log.Errorf("MEXC - Unsubscribe failed for %s: %v", pair.ForeignName, err)
			}
			time.Sleep(2 * time.Second)

			err = s.subscribe(pair, true, failoverChannel)
			if err != nil {
				log.Errorf("MEXC - Resubscribe failed for %s: %v", pair.ForeignName, err)
			}

		case <-ctx.Done():
			log.Debugf("MEXC - Close resubscribe routine of scraper with genesis: %v.", s.genesis)
			return
		}
	}
}

func (s *MEXCScraper) subscribe(pair models.ExchangePair, subscribe bool, failoverChannel chan string) error {
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
			if _, err := s.newConn(failoverChannel); err != nil {
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
