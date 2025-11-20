package scrapers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
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
	cryptodotcomWSBaseString        = "wss://stream.crypto.com/v2/market"
	cryptodotcomTradeTimeoutSeconds = 120
)

type cryptodotcomHooks struct{}

func (cryptodotcomHooks) ExchangeKey() string {
	return CRYPTODOTCOM_EXCHANGE
}

func (cryptodotcomHooks) WSURL() string {
	return cryptodotcomWSBaseString
}

func (cryptodotcomHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	// The code to handle heartbeat must be placed on the path after "receiving messages" (i.e. OnMessage), not in OnOpen.
	// Crypto.com's heartbeat is "server sends first, we respond".
	// There are two characteristics:
	// 1. must get the id from the server before responding, and the id is different each time;
	// 2. only when you "receive the heartbeat message", you know whether to respond and which id to respond with.
}

// subscribe/unsubscribe
func (cryptodotcomHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	baseQuote := strings.Split(pair.ForeignName, "-")
	if len(baseQuote) != 2 {
		log.Errorf("Crypto.com - invalid ForeignName: %s", pair.ForeignName)
		return nil
	}
	channel := "trade." + baseQuote[0] + "_" + baseQuote[1] // e.g. BTC-USDT -> trade.BTC_USDT

	method := "unsubscribe"
	if subscribe {
		method = "subscribe"
	}

	msg := cryptodotcomWSSubscribeMessage{
		ID:     1,
		Method: method,
		Params: cryptodotcomChannels{
			Channels: []string{channel},
		},
	}

	return bs.SafeWriteJSON(msg)
}

func (cryptodotcomHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// handle each ws text message
func (cryptodotcomHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var msg cryptodotcomWSResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	// heartbeat: public/heartbeat -> public/respond-heartbeat
	if msg.Method == "public/heartbeat" {
		sendCryptodotcomHeartbeat(bs, msg.ID)
		return
	}

	// other trade messages
	trades, err := cryptodotcomParseTradeMessage(msg)
	if err != nil {
		log.Errorf("Crypto.com - parseCryptodotcomTradeMessage: %s", err.Error())
		return
	}
	if len(trades) == 0 || len(msg.Result.Data) == 0 {
		return
	}

	// BTC_USDT
	pairParts := strings.Split(msg.Result.Data[0].ForeignName, "_")
	if len(pairParts) < 2 {
		return
	}
	tickerKey := pairParts[0] + pairParts[1] // "BTCUSDT"

	for _, trade := range trades {
		// discard too old trade
		if trade.Time.Before(time.Now().Add(-time.Duration(cryptodotcomTradeTimeoutSeconds) * time.Second)) {
			continue
		}

		// map to token information
		lock.RLock()
		pair, ok := bs.tickerPairMap[tickerKey]
		lock.RUnlock()
		if !ok {
			continue
		}

		trade.QuoteToken = pair.QuoteToken
		trade.BaseToken = pair.BaseToken

		log.Tracef(
			"Crypto.com - got trade: %v -- %s -- %v -- %v -- %s.",
			trade.Time,
			trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
			trade.Price,
			trade.Volume,
			trade.ForeignTradeID,
		)

		// lastTradeTimeMap key use "QUOTE-BASE" format (e.g. BTC-USDT)
		lastKey := pairParts[0] + "-" + pairParts[1]
		bs.setLastTradeTime(lock, lastKey, trade.Time)

		bs.tradesChannel <- trade
	}
}

// tickerPairMap key: remove '-' (e.g. BTC-USDT -> BTCUSDT)
func (cryptodotcomHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: use the original "BASE-QUOTE" format (e.g. BTC-USDT)
func (cryptodotcomHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

func sendCryptodotcomHeartbeat(bs *BaseCEXScraper, id int) {
	msg := cryptodotcomWSSubscribeMessage{
		ID:     id,
		Method: "public/respond-heartbeat",
	}
	if err := bs.SafeWriteJSON(msg); err != nil {
		log.Errorf("Crypto.com - send heartbeat error: %v", err)
	}
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

func NewCryptodotcomScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	wg *sync.WaitGroup,
) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, cryptodotcomHooks{})
}
