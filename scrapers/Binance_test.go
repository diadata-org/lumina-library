package scrapers

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// ------- fake wsConn implementation, for intercepting WriteJSON -------

type fakeWSConn struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConn) ReadMessage() (int, []byte, error) {
	return 0, nil, nil
}

func (f *fakeWSConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

func (f *fakeWSConn) ReadJSON(v interface{}) error {
	return nil
}

func (f *fakeWSConn) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConn) Close() error {
	return nil
}

// ------- Test: TickerKey / LastTradeKey -------

func TestBinanceHooks_TickerAndLastTradeKey(t *testing.T) {
	h := &binanceHooks{}
	got := h.TickerKeyFromForeign("BTC-USDT")
	if got != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", got)
	}

	got = h.LastTradeTimeKeyFromForeign("ETH-USDT")
	if got != "ETHUSDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETHUSDT, got %s", got)
	}
}

// ------- Test: Subscribe / Unsubscribe -------

func TestBinanceHooks_SubscribeAndUnsubscribe(t *testing.T) {
	oldInterval := binanceWriteInterval
	binanceWriteInterval = 0
	defer func() { binanceWriteInterval = oldInterval }()

	fc := &fakeWSConn{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := &binanceHooks{}

	var lock sync.RWMutex
	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(*binanceWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *binanceWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "SUBSCRIBE" {
		t.Fatalf("expected Method=SUBSCRIBE, got %s", msg.Method)
	}
	if len(msg.Params) != 1 || msg.Params[0] != "btcusdt@trade" {
		t.Fatalf("expected Params=[btcusdt@trade], got %+v", msg.Params)
	}

	// subscribe = false (UNSUBSCRIBE)
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(*binanceWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *binanceWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "UNSUBSCRIBE" {
		t.Fatalf("expected Method=UNSUBSCRIBE, got %s", msg.Method)
	}
	if len(msg.Params) != 1 || msg.Params[0] != "btcusdt@trade" {
		t.Fatalf("expected Params=[btcusdt@trade], got %+v", msg.Params)
	}
}

// ------- Test: OnMessage parse trade and write to channel -------

func TestBinanceHooks_OnMessage_ValidTrade(t *testing.T) {
	h := &binanceHooks{}

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}

	// prepare tickerPairMap: key use Binance symbol (e.g. "BTCUSDT")
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken: models.Asset{Symbol: "BTC"},
		QuoteToken: models.Asset{
			Symbol: "USDT",
		},
	}

	// construct a WebSocket message (Binance trade)
	msg := binanceWSResponse{
		Timestamp:      1700000000000, // 只是个示例
		Price:          "100.5",
		Volume:         "1.23",
		ForeignTradeID: 123456,
		ForeignName:    "BTCUSDT",
		Type:           "trade",
		Buy:            true,
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var lock sync.RWMutex

	// call OnMessage, simulate TextMessage
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// should write a trade to tradesChannel
	select {
	case trade := <-bs.tradesChannel:
		if trade.Price != 100.5 {
			t.Fatalf("expected trade.Price=100.5, got %v", trade.Price)
		}
		if trade.Volume != 1.23 {
			t.Fatalf("expected trade.Volume=1.23, got %v", trade.Volume)
		}
		if trade.BaseToken.Symbol != "BTC" || trade.QuoteToken.Symbol != "USDT" {
			t.Fatalf("unexpected tokens: base=%s quote=%s",
				trade.BaseToken.Symbol, trade.QuoteToken.Symbol)
		}

		// lastTradeTimeMap should be updated, key use LastTradeTimeKeyFromForeign
		key := h.LastTradeTimeKeyFromForeign("BTC-USDT")
		lock.RLock()
		_, ok := bs.lastTradeTimeMap[key]
		lock.RUnlock()
		if !ok {
			t.Fatalf("expected lastTradeTimeMap[%s] to be set", key)
		}
	default:
		t.Fatalf("expected one trade in tradesChannel, but channel is empty")
	}
}

// ------- Test: OnMessage ignores non-TextMessage / invalid JSON / Type == nil -------

func TestBinanceHooks_OnMessage_IgnoresNonTextOrInvalid(t *testing.T) {
	h := &binanceHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	// 1) non-TextMessage
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)

	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-text message")
	default:
	}

	// 2) invalid JSON
	h.OnMessage(bs, ws.TextMessage, []byte(`not-json`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for invalid JSON")
	default:
	}

	// 3) Type == nil
	msg := binanceWSResponse{
		ForeignName: "BTCUSDT",
		Type:        nil,
	}
	raw, _ := json.Marshal(msg)
	h.OnMessage(bs, ws.TextMessage, raw, &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade when Type is nil")
	default:
	}
}
