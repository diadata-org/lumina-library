package scrapers

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

type cryptodotcomWSSubscribeMessage struct {
	ID     int                  `json:"id"`
	Method string               `json:"method"`
	Params cryptodotcomChannels `json:"params"`
}

type cryptodotcomChannels struct {
	Channels []string `json:"channels"`
}

type cryptodotcomWSResponse struct {
	ID     int                          `json:"id"`
	Method string                       `json:"method"`
	Code   int                          `json:"code"`
	Result cryptodotcomWSResponseResult `json:"result"`
}

type cryptodotcomWSResponseResult struct {
	InstrumentName string                       `json:"instrument_name"`
	Subscription   string                       `json:"subscription"`
	Channel        string                       `json:"channel"`
	Data           []cryptodotcomWSResponseData `json:"data"`
}

type cryptodotcomWSResponseData struct {
	TradeID     string `json:"d"`
	Timestamp   int64  `json:"t"`
	Price       string `json:"p"`
	Volume      string `json:"q"`
	Side        string `json:"s"`
	ForeignName string `json:"i"`
}

type cryptodotcomScraper struct {
	wsClient            *ws.Conn
	tradesChannel       chan models.Trade
	subscribeChannel    chan models.ExchangePair
	tickerPairMap       map[string]models.Pair
	lastTradeTimeMap    map[string]time.Time
	maxErrCount         int
	restartWaitTime     int
	genesis             time.Time
	tradeTimeoutSeconds int
}

var (
	cryptodotcomWSBaseString = "wss://stream.crypto.com/v2/market"
)

func NewCryptodotcomScraper(ctx context.Context, pairs []models.ExchangePair, failoverChannel chan string, wg *sync.WaitGroup) Scraper {
	defer wg.Done()
	var lock sync.RWMutex
	log.Info("Crypto.com - Started scraper.")

	scraper := cryptodotcomScraper{
		tradesChannel:       make(chan models.Trade),
		subscribeChannel:    make(chan models.ExchangePair),
		tickerPairMap:       models.MakeTickerPairMap(pairs),
		lastTradeTimeMap:    make(map[string]time.Time),
		maxErrCount:         20,
		restartWaitTime:     5,
		genesis:             time.Now(),
		tradeTimeoutSeconds: 120,
	}

	// Dial websocket API.
	var wsDialer ws.Dialer
	wsClient, _, err := wsDialer.Dial(cryptodotcomWSBaseString, nil)
	if err != nil {
		log.Errorf("Crypto.com - Dial ws base string: %v.", err)
		failoverChannel <- string(CRYPTODOTCOM_EXCHANGE)
		return &scraper
	}
	scraper.wsClient = wsClient

	// Subscribe to pairs and initialize cryptodotcomLastTradeTimeMap.
	for _, pair := range pairs {
		if err := scraper.subscribe(pair, true, &lock); err != nil {
			log.Errorf("Crypto.com - Subscribe to pair %s: %v.", pair.ForeignName, err)
		} else {
			log.Debugf("Crypto.com - Subscribed to pair %s.", pair.ForeignName)
			scraper.lastTradeTimeMap[pair.ForeignName] = time.Now()
		}
	}

	go scraper.fetchTrades(&lock)

	// Check last trade time for each subscribed pair and resubscribe if no activity for more than @cryptodotcomWatchdogDelayMap.
	for _, pair := range pairs {
		envVar := "CRYPTODOTCOM_WATCHDOG_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[0] + "_" + strings.Split(strings.ToUpper(pair.ForeignName), "-")[1]
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "300"), 10, 64)
		if err != nil {
			log.Errorf("Crypto.com - Parse cryptodotcomWatchdogDelay: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdog(ctx, pair, watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, &lock)
		go scraper.resubscribe(ctx, &lock)
	}

	return &scraper
}

func (scraper *cryptodotcomScraper) Close(cancel context.CancelFunc) error {
	log.Warn("Crypto.com - call scraper.Close().")
	cancel()
	if scraper.wsClient == nil {
		return nil
	}
	return scraper.wsClient.Close()
}

func (scraper *cryptodotcomScraper) TradesChannel() chan models.Trade {
	return scraper.tradesChannel
}

func (scraper *cryptodotcomScraper) fetchTrades(lock *sync.RWMutex) {
	// Read trades stream.
	var errCount int

	for {

		var message cryptodotcomWSResponse
		err := scraper.wsClient.ReadJSON(&message)
		if err != nil {
			if handleErrorReadJSON(err, &errCount, scraper.maxErrCount, CRYPTODOTCOM_EXCHANGE, scraper.restartWaitTime) {
				return
			}
			continue
		}
		if message.Method == "public/heartbeat" {
			scraper.sendHeartbeat(message.ID, lock)
			continue
		}

		scraper.handleWSResponse(message, lock)

	}

}

func (scraper *cryptodotcomScraper) handleWSResponse(message cryptodotcomWSResponse, lock *sync.RWMutex) {
	trades, err := cryptodotcomParseTradeMessage(message)
	if err != nil {
		log.Errorf("Crypto.com - parseCryptodotcomTradeMessage: %s.", err.Error())
		// continue
		return
	}

	// Identify ticker symbols with underlying assets.
	for _, trade := range trades {

		// The websocket API returns very old trades when first subscribing. Hence, discard if too old.
		if trade.Time.Before(time.Now().Add(-time.Duration(scraper.tradeTimeoutSeconds) * time.Second)) {
			continue
		}

		pair := strings.Split(message.Result.Data[0].ForeignName, "_")
		if len(pair) > 1 {
			trade.QuoteToken = scraper.tickerPairMap[pair[0]+pair[1]].QuoteToken
			trade.BaseToken = scraper.tickerPairMap[pair[0]+pair[1]].BaseToken
		}

		log.Tracef("Crypto.com - got trade: %v -- %s -- %v -- %v -- %s.", trade.Time, trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol, trade.Price, trade.Volume, trade.ForeignTradeID)
		lock.Lock()
		scraper.lastTradeTimeMap[pair[0]+"-"+pair[1]] = trade.Time
		lock.Unlock()

		scraper.tradesChannel <- trade
	}

}

func (scraper *cryptodotcomScraper) resubscribe(ctx context.Context, lock *sync.RWMutex) {
	for {
		select {
		case pair := <-scraper.subscribeChannel:
			err := scraper.subscribe(pair, false, lock)
			if err != nil {
				log.Errorf("Crypto.com - Unsubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Crypto.com - Unsubscribed pair %s.", pair.ForeignName)
			}
			time.Sleep(2 * time.Second)
			err = scraper.subscribe(pair, true, lock)
			if err != nil {
				log.Errorf("Crypto.com - Resubscribe pair %s: %v.", pair.ForeignName, err)
			} else {
				log.Debugf("Crypto.com - Subscribed to pair %s.", pair.ForeignName)
			}
		case <-ctx.Done():
			log.Debugf("Crypto.com - Close resubscribe routine of scraper with genesis: %v.", scraper.genesis)
			return
		}
	}
}

func (scraper *cryptodotcomScraper) subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	defer lock.Unlock()
	channel := []string{"trade." + strings.Split(pair.ForeignName, "-")[0] + "_" + strings.Split(pair.ForeignName, "-")[1]}
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	a := cryptodotcomWSSubscribeMessage{
		ID:     1,
		Method: subscribeType,
		Params: cryptodotcomChannels{
			Channels: channel,
		},
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func (scraper *cryptodotcomScraper) sendHeartbeat(id int, lock *sync.RWMutex) error {
	defer lock.Unlock()
	a := cryptodotcomWSSubscribeMessage{
		ID:     id,
		Method: "public/respond-heartbeat",
	}
	lock.Lock()
	return scraper.wsClient.WriteJSON(a)
}

func cryptodotcomParseTradeMessage(message cryptodotcomWSResponse) (trades []models.Trade, err error) {

	for _, data := range message.Result.Data {
		var (
			price  float64
			volume float64
		)
		price, err = strconv.ParseFloat(data.Price, 64)
		if err != nil {
			return
		}
		volume, err = strconv.ParseFloat(data.Volume, 64)
		if err != nil {
			return
		}
		if data.Side == "SELL" {
			volume *= -1
		}
		timestamp := time.Unix(0, data.Timestamp*1e6)
		foreignTradeID := data.TradeID

		trade := models.Trade{
			Price:          price,
			Volume:         volume,
			Time:           timestamp,
			Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
			ForeignTradeID: foreignTradeID,
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
