package scrapers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// ----------------- Fake wsConn for Coinbase tests -----------------

// fakeWSConnCB is a minimal fake implementation of wsConn that records
// WriteJSON calls so we can assert on the messages sent by Subscribe.
type fakeWSConnCB struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConnCB) ReadMessage() (int, []byte, error) {
	// Not used in these tests.
	return 0, nil, nil
}

func (f *fakeWSConnCB) WriteMessage(messageType int, data []byte) error {
	// Not used in these tests.
	return nil
}

func (f *fakeWSConnCB) ReadJSON(v interface{}) error {
	// Not used in these tests.
	return nil
}

func (f *fakeWSConnCB) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConnCB) Close() error {
	return nil
}

// ----------------- Tests for key mapping helpers -----------------

func TestCoinbaseHooks_TickerAndLastTradeKey(t *testing.T) {
	h := coinbaseHooks{}

	// TickerKeyFromForeign should remove "-"
	gotTicker := h.TickerKeyFromForeign("BTC-USD")
	if gotTicker != "BTCUSD" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSD, got %s", gotTicker)
	}

	// LastTradeTimeKeyFromForeign should keep the original "BASE-QUOTE" form
	gotLastKey := h.LastTradeTimeKeyFromForeign("ETH-USD")
	if gotLastKey != "ETH-USD" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USD, got %s", gotLastKey)
	}
}

// ----------------- Tests for Subscribe / Unsubscribe -----------------

func TestCoinbaseHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnCB{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := coinbaseHooks{}

	var lock sync.RWMutex
	pair := models.ExchangePair{ForeignName: "BTC-USD"}

	// subscribe = true -> type = "subscribe"
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call, got %d", fc.writeJSONCount)
	}

	msg, ok := fc.lastWritten.(*coinBaseWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *coinBaseWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Type != "subscribe" {
		t.Fatalf("expected Type=subscribe, got %s", msg.Type)
	}
	if len(msg.Channels) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(msg.Channels))
	}
	ch := msg.Channels[0]
	if ch.Name != "matches" {
		t.Fatalf("expected channel name 'matches', got %s", ch.Name)
	}
	if len(ch.ProductIDs) != 1 || ch.ProductIDs[0] != "BTC-USD" {
		t.Fatalf("expected ProductIDs=[BTC-USD], got %+v", ch.ProductIDs)
	}

	// subscribe = false -> type = "unsubscribe"
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}

	msg, ok = fc.lastWritten.(*coinBaseWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *coinBaseWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Type != "unsubscribe" {
		t.Fatalf("expected Type=unsubscribe, got %s", msg.Type)
	}
	if len(msg.Channels) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(msg.Channels))
	}
	ch = msg.Channels[0]
	if len(ch.ProductIDs) != 1 || ch.ProductIDs[0] != "BTC-USD" {
		t.Fatalf("expected ProductIDs=[BTC-USD], got %+v", ch.ProductIDs)
	}
}

// ----------------- Tests for OnMessage -----------------

func TestCoinbaseHooks_OnMessage_ValidMatch(t *testing.T) {
	h := coinbaseHooks{}

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}

	// tickerPairMap key is "BASEQUOTE" (without "-"), e.g. "BTCUSD"
	bs.tickerPairMap["BTCUSD"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "BTC"},
		QuoteToken: models.Asset{Symbol: "USD"},
	}

	// Build a Coinbase "match" message for BTC-USD
	resp := coinBaseWSResponse{
		Type:      "match",
		TradeID:   42,
		Time:      "2024-01-02T15:04:05.000000Z",
		ProductID: "BTC-USD",
		Size:      "0.5",
		Price:     "30000.12",
		Side:      "buy",
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var lock sync.RWMutex

	// Call OnMessage as if it came from WebSocket as TextMessage.
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// We expect one trade in tradesChannel.
	select {
	case tr := <-bs.tradesChannel:
		// Price and volume
		if tr.Price != 30000.12 {
			t.Fatalf("expected Price=30000.12, got %v", tr.Price)
		}
		if tr.Volume != 0.5 {
			t.Fatalf("expected Volume=0.5, got %v", tr.Volume)
		}
		// Tokens
		if tr.BaseToken.Symbol != "BTC" || tr.QuoteToken.Symbol != "USD" {
			t.Fatalf("unexpected tokens: base=%s quote=%s",
				tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
		}
		// ForeignTradeID
		if tr.ForeignTradeID != "42" {
			t.Fatalf("expected ForeignTradeID=42, got %s", tr.ForeignTradeID)
		}
		// lastTradeTimeMap key should be "BTC-USD"
		lock.RLock()
		_, ok := bs.lastTradeTimeMap["BTC-USD"]
		lock.RUnlock()
		if !ok {
			t.Fatalf("expected lastTradeTimeMap[BTC-USD] to be set")
		}
	default:
		t.Fatalf("expected one trade in tradesChannel, but channel is empty")
	}
}

func TestCoinbaseHooks_OnMessage_IgnoresNonTextOrNonMatch(t *testing.T) {
	h := coinbaseHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	// Non-text message should be ignored.
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-text message")
	default:
	}

	// Invalid JSON should be ignored.
	h.OnMessage(bs, ws.TextMessage, []byte(`not-json`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for invalid JSON")
	default:
	}

	// Type != "match" should be ignored.
	resp := coinBaseWSResponse{
		Type:      "subscriptions",
		ProductID: "BTC-USD",
	}
	raw, _ := json.Marshal(resp)
	h.OnMessage(bs, ws.TextMessage, raw, &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-match type")
	default:
	}
}

// ----------------- Optional: ReadLoop handled flag -----------------

func TestCoinbaseHooks_ReadLoopHandledFalse(t *testing.T) {
	h := coinbaseHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex

	// Coinbase ReadLoop always returns false, so BaseCEXScraper.runReadLoop
	// will use the default ReadMessage path.
	handled := h.ReadLoop(context.Background(), bs, &lock)
	if handled {
		t.Fatalf("expected coinbaseHooks.ReadLoop to return false")
	}
}
