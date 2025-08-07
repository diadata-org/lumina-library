package scrapers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
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
	wsClient         *ws.Conn
	tradesChannel    chan models.Trade
	subscribeChannel chan models.ExchangePair
	tickerPairMap    map[string]models.Pair
	lastTradeTimeMap map[string]time.Time
	maxErrCount      int
	restartWaitTime  int
	genesis          time.Time
	writeLock        sync.Mutex
}

const byBitWSBaseURL = "wss://stream.bybit.com/v5/public/spot"

func NewByBitScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("Bybit - Started scraper.")

	scraper := byBitScraper{
		tradesChannel:    make(chan models.Trade),
		subscribeChannel: make(chan models.ExchangePair),
		tickerPairMap:    models.MakeTickerPairMap(pairs),
		lastTradeTimeMap: make(map[string]time.Time),
		maxErrCount:      20,
		restartWaitTime:  5,
		genesis:          time.Now(),
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

	for _, pair := range pairs {
		envVar := strings.ToUpper(BYBIT_EXCHANGE) + "_WATCHDOG_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[0] + "_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[1]
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "300"), 10, 64)
		if err != nil {
			log.Errorf("ByBit - Parse bybitWatchdogDelayMap: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdog(ctx, pair, watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, &lock)
		go scraper.resubscribe(ctx, &lock)
	}

	return &scraper
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
		log.Infof("ByBit - Trade: %s-%s | Side: %v | Price: %f | Volume: %f", pair.BaseToken.Symbol, pair.QuoteToken.Symbol, side, price, volume)
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
