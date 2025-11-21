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

var coinbaseWSBaseString = "wss://ws-feed.exchange.coinbase.com"

type coinBaseWSSubscribeMessage struct {
	Type     string            `json:"type"`
	Channels []coinBaseChannel `json:"channels"`
}

type coinBaseChannel struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids"`
}

type coinBaseWSResponse struct {
	Type      string `json:"type"` // "match"
	TradeID   int64  `json:"trade_id"`
	Time      string `json:"time"`       // "2006-01-02T15:04:05.000000Z"
	ProductID string `json:"product_id"` // "BTC-USD"
	Size      string `json:"size"`
	Price     string `json:"price"`
	Side      string `json:"side"` // "buy"/"sell"
}

type coinbaseHooks struct{}

func (coinbaseHooks) ExchangeKey() string { return COINBASE_EXCHANGE }
func (coinbaseHooks) WSURL() string       { return coinbaseWSBaseString }

func (coinbaseHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) { /* no ping */ }

func (coinbaseHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	typ := "unsubscribe"
	if subscribe {
		typ = "subscribe"
	}
	msg := &coinBaseWSSubscribeMessage{
		Type: typ,
		Channels: []coinBaseChannel{{
			Name:       "matches",
			ProductIDs: []string{pair.ForeignName},
		}},
	}
	return bs.SafeWriteJSON(msg)
}

func (coinbaseHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (coinbaseHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var m coinBaseWSResponse
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}
	if m.Type != "match" {
		return
	}

	price, err := strconv.ParseFloat(m.Price, 64)
	if err != nil {
		return
	}
	vol, err := strconv.ParseFloat(m.Size, 64)
	if err != nil {
		return
	}
	if m.Side == "sell" {
		vol = -vol
	}
	ts, err := time.Parse("2006-01-02T15:04:05.000000Z", m.Time)
	if err != nil {
		return
	}

	parts := strings.Split(m.ProductID, "-")
	if len(parts) < 2 {
		return
	}
	key := parts[0] + parts[1] // key for tickerPairMap: remove '-'

	lock.RLock()
	pair, ok := bs.tickerPairMap[key]
	lock.RUnlock()
	if !ok {
		return
	}

	trade := models.Trade{
		Price:          price,
		Volume:         vol,
		Time:           ts,
		Exchange:       Exchanges[COINBASE_EXCHANGE],
		ForeignTradeID: strconv.Itoa(int(m.TradeID)),
		BaseToken:      pair.BaseToken,
		QuoteToken:     pair.QuoteToken,
	}

	bs.setLastTradeTime(lock, parts[0]+"-"+parts[1], trade.Time)

	bs.tradesChannel <- trade
}

func (coinbaseHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}
func (coinbaseHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign // "BTC-USDT"
}

func NewCoinBaseScraper(ctx context.Context, pairs []models.ExchangePair, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, coinbaseHooks{})
}
