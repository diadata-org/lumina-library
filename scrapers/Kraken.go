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

type krakenWSSubscribeMessage struct {
	Method string       `json:"method"`
	Params krakenParams `json:"params"`
}

type krakenParams struct {
	Channel string   `json:"channel"`
	Symbol  []string `json:"symbol"`
}

type krakenWSResponse struct {
	Channel string                 `json:"channel"`
	Type    string                 `json:"type"`
	Data    []krakenWSResponseData `json:"data"`
}

type krakenWSResponseData struct {
	Symbol    string  `json:"symbol"` // e.g. "BTC/USDT"
	Side      string  `json:"side"`   // "buy" / "sell"
	Price     float64 `json:"price"`
	Size      float64 `json:"qty"`
	OrderType string  `json:"order_type"`
	TradeID   int     `json:"trade_id"`
	Time      string  `json:"timestamp"` // "2006-01-02T15:04:05.000000Z"
}

var (
	krakenWSBaseString = "wss://ws.kraken.com/v2"
)

type krakenHooks struct{}

func (krakenHooks) ExchangeKey() string {
	return KRAKEN_EXCHANGE
}

func (krakenHooks) WSURL() string {
	return krakenWSBaseString
}

// Kraken does not need extra ping
func (krakenHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {}

func (krakenHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	method := "unsubscribe"
	if subscribe {
		method = "subscribe"
	}

	msg := &krakenWSSubscribeMessage{
		Method: method,
		Params: krakenParams{
			Channel: "trade",
			Symbol:  []string{strings.ReplaceAll(pair.ForeignName, "-", "/")}, // "BTC-USDT" -> "BTC/USDT"
		},
	}

	return bs.SafeWriteJSON(msg)
}

func (krakenHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// handle each ws text message
func (krakenHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var msg krakenWSResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	if msg.Channel != "trade" {
		return
	}

	for _, d := range msg.Data {
		price, volume, ts, foreignTradeID, err := parseKrakenTradeMessage(d)
		if err != nil {
			log.Errorf("Kraken - parseTradeMessage: %v.", err)
			continue
		}

		// map to pair
		parts := strings.Split(d.Symbol, "/")
		if len(parts) < 2 {
			continue
		}
		key := parts[0] + parts[1] // "BTC/USDT" -> "BTCUSDT"

		lock.RLock()
		pair, ok := bs.tickerPairMap[key]
		lock.RUnlock()
		if !ok {
			continue
		}

		trade := models.Trade{
			QuoteToken:     pair.QuoteToken,
			BaseToken:      pair.BaseToken,
			Price:          price,
			Volume:         volume,
			Time:           ts,
			Exchange:       Exchanges[KRAKEN_EXCHANGE],
			ForeignTradeID: foreignTradeID,
		}

		log.Tracef(
			"Kraken - got trade: %s -- %v -- %v -- %s.",
			trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
			trade.Price,
			trade.Volume,
			trade.ForeignTradeID,
		)

		// lastTradeTimeMap key: use "QUOTE-BASE" format (e.g. BTC-USDT)
		lastKey := trade.QuoteToken.Symbol + "-" + trade.BaseToken.Symbol
		bs.setLastTradeTime(lock, lastKey, trade.Time)

		bs.tradesChannel <- trade
	}
}

// tickerPairMap key: remove '-' (e.g. "BTC-USDT" -> "BTCUSDT")
func (krakenHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: use ForeignName (e.g. "BTC-USDT")
func (krakenHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

func parseKrakenTradeMessage(message krakenWSResponseData) (price float64, volume float64, timestamp time.Time, foreignTradeID string, err error) {
	price = message.Price
	volume = message.Size
	if message.Side == "sell" {
		volume *= -1
	}
	timestamp, err = time.Parse("2006-01-02T15:04:05.000000Z", message.Time)
	if err != nil {
		return
	}

	foreignTradeID = strconv.Itoa(message.TradeID)
	return
}

func NewKrakenScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	wg *sync.WaitGroup,
) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, krakenHooks{})
}
