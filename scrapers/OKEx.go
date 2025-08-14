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
	wsClient         *ws.Conn
	tradesChannel    chan models.Trade
	subscribeChannel chan models.ExchangePair
	tickerPairMap    map[string]models.Pair
	lastTradeTimeMap map[string]time.Time
	maxErrCount      int
	restartWaitTime  int
	genesis          time.Time
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
		tradesChannel:    make(chan models.Trade),
		subscribeChannel: make(chan models.ExchangePair),
		tickerPairMap:    models.MakeTickerPairMap(pairs),
		lastTradeTimeMap: make(map[string]time.Time),
		maxErrCount:      20,
		restartWaitTime:  5,
		genesis:          time.Now(),
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

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @OKExWatchdogDelayMap.
	for _, pair := range pairs {
		envVar := strings.ToUpper(OKEX_EXCHANGE) + "_WATCHDOG_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[0] + "_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[1]
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "300"), 10, 64)
		if err != nil {
			log.Errorf("OKEx - Parse OKExWatchdogDelay: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdog(ctx, pair, watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, &lock)
		go scraper.resubscribe(ctx, &lock)
	}

	return &scraper
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

		log.Infof("OKEx - got trade: %s -- %v -- %v -- %s -- %s.", trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID, trade.Time)
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
