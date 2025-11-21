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

const byBitWSBaseURL = "wss://stream.bybit.com/v5/public/spot"

type byBitWSSubscribeMessage struct {
	OP   string   `json:"op"`
	Args []string `json:"args"`
}

type byBitWSResponse struct {
	Topic     string                   `json:"topic"`
	Timestamp int64                    `json:"ts"`
	Type      string                   `json:"type"`
	Data      []byBitTradeResponseData `json:"data"`
}

type byBitTradeResponseData struct {
	TradeID   string `json:"i"`
	Timestamp int64  `json:"T"`
	Price     string `json:"p"`
	Size      string `json:"v"`
	Side      string `json:"S"` // Sell/Buy
	Symbol    string `json:"s"` // e.g. BTCUSDT
}

type bybitHooks struct{}

func (bybitHooks) ExchangeKey() string { return BYBIT_EXCHANGE }
func (bybitHooks) WSURL() string       { return byBitWSBaseURL }

func (bybitHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	// 10s ping
	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				_ = bs.SafeWriteMessage(ws.PingMessage, []byte{})
			}
		}
	}()
}

func (bybitHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	op := "unsubscribe"
	if subscribe {
		op = "subscribe"
	}
	msg := byBitWSSubscribeMessage{
		OP:   op,
		Args: []string{"publicTrade." + strings.ReplaceAll(pair.ForeignName, "-", "")},
	}
	return bs.SafeWriteJSON(msg)
}

func (bybitHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}
	// subscription ack
	if strings.Contains(string(data), "\"success\"") {
		return
	}
	var resp byBitWSResponse
	if err := json.Unmarshal(data, &resp); err != nil || resp.Type != "snapshot" {
		return
	}
	for _, d := range resp.Data {
		price, err := strconv.ParseFloat(d.Price, 64)
		if err != nil {
			continue
		}
		vol, err := strconv.ParseFloat(d.Size, 64)
		if err != nil {
			continue
		}
		if d.Side == "Sell" {
			vol = -vol
		}
		pairKey := d.Symbol // already removed '-'
		lock.RLock()
		pair, ok := bs.tickerPairMap[pairKey]
		lock.RUnlock()
		if !ok {
			return
		}
		trade := models.Trade{
			Price:      price,
			Volume:     vol,
			Time:       time.Now(),
			Exchange:   Exchanges[BYBIT_EXCHANGE],
			BaseToken:  pair.BaseToken,
			QuoteToken: pair.QuoteToken,
		}
		bs.setLastTradeTime(lock, pair.QuoteToken.Symbol+"-"+pair.BaseToken.Symbol, time.Now())
		bs.tradesChannel <- trade
	}
}

func (bybitHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (bybitHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}
func (bybitHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign // bybit's lastTradeTimeMap key use "QUOTE-BASE"
}

func NewByBitScraper(ctx context.Context, pairs []models.ExchangePair, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, bybitHooks{})
}
