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

type SubscribeGate struct {
	Time    int64    `json:"time"`
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Payload []string `json:"payload"`
}

type GateIOResponseTrade struct {
	Time    int    `json:"time"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Result  struct {
		ID           int    `json:"id"`
		CreateTime   int    `json:"create_time"`
		CreateTimeMs string `json:"create_time_ms"`
		Side         string `json:"side"`
		CurrencyPair string `json:"currency_pair"` // e.g. "BTC_USDT"
		Amount       string `json:"amount"`
		Price        string `json:"price"`
	} `json:"result"`
}

const gateIOWSURL = "wss://api.gateio.ws/ws/v4/"

type gateIOHooks struct{}

func (gateIOHooks) ExchangeKey() string {
	return GATEIO_EXCHANGE
}

func (gateIOHooks) WSURL() string {
	return gateIOWSURL
}

// GateIO does not need extra ping
func (gateIOHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {}

func (gateIOHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	parts := strings.Split(pair.ForeignName, "-")
	if len(parts) != 2 {
		log.Errorf("GateIO - invalid ForeignName: %s", pair.ForeignName)
		return nil
	}
	gateioPairTicker := parts[0] + "_" + parts[1]

	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	msg := &SubscribeGate{
		Event:   subscribeType,
		Time:    time.Now().Unix(),
		Channel: "spot.trades",
		Payload: []string{gateioPairTicker},
	}

	return bs.SafeWriteJSON(msg)
}

func (gateIOHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (gateIOHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var msg GateIOResponseTrade
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	if msg.Channel != "spot.trades" || msg.Event != "update" && msg.Event != "" {
		log.Debugf("GateIO - ignore non trade msg: %+v", msg.Event)
	}

	trade, ok := parseGateIOTrade(msg, bs, lock)
	if !ok {
		return
	}

	lastKey := trade.QuoteToken.Symbol + "-" + trade.BaseToken.Symbol
	bs.setLastTradeTime(lock, lastKey, trade.Time)

	log.Tracef("GateIO - got trade: %s -- %v -- %v -- %s.",
		trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
		trade.Price, trade.Volume, trade.ForeignTradeID,
	)

	bs.tradesChannel <- trade
}

// tickerPairMap key: remove '-' (e.g. BTC-USDT -> BTCUSDT)
func (gateIOHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: use ForeignName (e.g. BTC-USDT)
func (gateIOHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

func parseGateIOTrade(
	message GateIOResponseTrade,
	bs *BaseCEXScraper,
	lock *sync.RWMutex,
) (models.Trade, bool) {
	var (
		price  float64
		volume float64
		err    error
	)

	if message.Result.Price == "" || message.Result.Amount == "" {
		return models.Trade{}, false
	}

	price, err = strconv.ParseFloat(message.Result.Price, 64)
	if err != nil {
		log.Errorf("GateIO - error parsing float Price %v: %v.", message.Result.Price, err)
		return models.Trade{}, false
	}

	volume, err = strconv.ParseFloat(message.Result.Amount, 64)
	if err != nil {
		log.Errorf("GateIO - error parsing float Volume %v: %v.", message.Result.Amount, err)
		return models.Trade{}, false
	}

	if message.Result.Side == "sell" {
		volume = -volume
	}

	// map to pair: the original map key is "BTCUSDT"
	cpParts := strings.Split(message.Result.CurrencyPair, "_")
	if len(cpParts) != 2 {
		return models.Trade{}, false
	}
	tickerKey := cpParts[0] + cpParts[1]

	lock.RLock()
	exchangePair, ok := bs.tickerPairMap[tickerKey]
	lock.RUnlock()
	if !ok {
		return models.Trade{}, false
	}

	trade := models.Trade{
		QuoteToken:     exchangePair.QuoteToken,
		BaseToken:      exchangePair.BaseToken,
		Price:          price,
		Volume:         volume,
		Time:           time.Unix(int64(message.Result.CreateTime), 0),
		Exchange:       Exchanges[GATEIO_EXCHANGE],
		ForeignTradeID: strconv.FormatInt(int64(message.Result.ID), 16),
	}

	return trade, true
}

func NewGateIOScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	wg *sync.WaitGroup,
) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, gateIOHooks{})
}
