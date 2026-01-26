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

// subscribe message
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

// WebSocket response
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

// parse single trade
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
	// okex timestamp is in ms
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

type okexHooks struct{}

func (okexHooks) ExchangeKey() string {
	return OKEX_EXCHANGE
}

func (okexHooks) WSURL() string {
	return OKExWSBaseString
}

// OKEx does not need extra ping
func (okexHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {}

// subscribe / unsubscribe
func (okexHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	op := "unsubscribe"
	if subscribe {
		op = "subscribe"
	}

	msg := &OKExWSSubscribeMessage{
		OP: op,
		Args: []OKExArgs{
			{
				Channel: "trades",
				InstIDs: pair.ForeignName, // e.g. "BTC-USDT"
			},
		},
	}

	return bs.SafeWriteJSON(msg)
}

func (okexHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (okexHooks) OnMessage(bs *BaseCEXScraper, messageType int, data []byte, lock *sync.RWMutex) {
	if messageType != ws.TextMessage {
		return
	}

	var msg OKEXWSResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Errorf("OKEx - unmarshal ws message error: %v", err)
		return
	}

	// each Data is one trade
	for _, d := range msg.Data {
		handleOKExData(&d, bs, lock)
	}
}

// tickerPairMap key: "BTCUSDT"
func (okexHooks) TickerKeyFromForeign(foreign string) string {
	parts := strings.Split(foreign, "-")
	if len(parts) < 2 {
		return foreign
	}
	return parts[0] + parts[1]
}

// lastTradeTimeMap key: "BTC-USDT"
func (okexHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// actually process single OKEXDATA, construct Trade and put into channel
func handleOKExData(d *OKEXDATA, bs *BaseCEXScraper, lock *sync.RWMutex) {
	trade, err := OKExParseTradeMessage(*d)
	if err != nil {
		log.Errorf("OKEx - parseOKExTradeMessage: %s.", err.Error())
		return
	}

	// InstID e.g. "BTC-USDT"
	pairParts := strings.Split(d.InstID, "-")
	if len(pairParts) < 2 {
		log.Warnf("OKEx - Unexpected instId format: %q", d.InstID)
		return
	}

	// tickerPairMap key: "BTCUSDT"
	tickerKey := pairParts[0] + pairParts[1]

	lock.RLock()
	pair, ok := bs.tickerPairMap[tickerKey]
	lock.RUnlock()
	if !ok {
		log.Tracef("OKEx - unknown ticker key in tickerPairMap: %s", tickerKey)
		return
	}

	trade.QuoteToken = pair.QuoteToken
	trade.BaseToken = pair.BaseToken

	log.Tracef("OKEx - got trade: %s -- %v -- %v -- %s -- %s.",
		trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
		trade.Price, trade.Volume, trade.ForeignTradeID, trade.Time,
	)

	// lastTradeTimeMap key: "BTC-USDT"
	lastKey := d.InstID
	bs.setLastTradeTime(lock, lastKey, trade.Time)

	bs.tradesChannel <- trade
}

func NewOKExScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, okexHooks{}, branchMarketConfig)
}
