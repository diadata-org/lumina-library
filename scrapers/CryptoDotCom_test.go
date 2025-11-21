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

// ---------- Fake wsConn for Crypto.com tests ----------

// fakeWSConnCrypto is a fake implementation of wsConn.
// It records WriteJSON calls so we can assert on the messages sent.
type fakeWSConnCrypto struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConnCrypto) ReadMessage() (int, []byte, error) {
	// Not used in these tests.
	return 0, nil, nil
}

func (f *fakeWSConnCrypto) WriteMessage(messageType int, data []byte) error {
	// Not used in these tests.
	return nil
}

func (f *fakeWSConnCrypto) ReadJSON(v interface{}) error {
	// Not used in these tests.
	return nil
}

func (f *fakeWSConnCrypto) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConnCrypto) Close() error {
	return nil
}

// ---------- Key mapping helpers ----------

func TestCryptodotcomHooks_TickerAndLastTradeKeys(t *testing.T) {
	h := cryptodotcomHooks{}

	// TickerKeyFromForeign should remove "-"
	gotTicker := h.TickerKeyFromForeign("BTC-USDT")
	if gotTicker != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", gotTicker)
	}

	// LastTradeTimeKeyFromForeign should keep "BASE-QUOTE"
	gotLastKey := h.LastTradeTimeKeyFromForeign("ETH-USDT")
	if gotLastKey != "ETH-USDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USDT, got %s", gotLastKey)
	}
}

// ---------- Subscribe / Unsubscribe ----------

func TestCryptodotcomHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnCrypto{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := cryptodotcomHooks{}
	var lock sync.RWMutex

	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true -> method = "subscribe", channel = "trade.BTC_USDT"
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(cryptodotcomWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be cryptodotcomWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "subscribe" {
		t.Fatalf("expected Method=subscribe, got %s", msg.Method)
	}
	if len(msg.Params.Channels) != 1 || msg.Params.Channels[0] != "trade.BTC_USDT" {
		t.Fatalf("expected Channels=[trade.BTC_USDT], got %+v", msg.Params.Channels)
	}

	// subscribe = false -> method = "unsubscribe"
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(cryptodotcomWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be cryptodotcomWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "unsubscribe" {
		t.Fatalf("expected Method=unsubscribe, got %s", msg.Method)
	}
	if len(msg.Params.Channels) != 1 || msg.Params.Channels[0] != "trade.BTC_USDT" {
		t.Fatalf("expected Channels=[trade.BTC_USDT], got %+v", msg.Params.Channels)
	}
}

// ---------- Heartbeat handling ----------

func TestCryptodotcomHooks_OnMessage_Heartbeat(t *testing.T) {
	fc := &fakeWSConnCrypto{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := cryptodotcomHooks{}
	var lock sync.RWMutex

	// Simulate a heartbeat message from server
	msg := cryptodotcomWSResponse{
		ID:     123,
		Method: "public/heartbeat",
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// We expect one heartbeat response to be written
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call for heartbeat, got %d", fc.writeJSONCount)
	}
	resp, ok := fc.lastWritten.(cryptodotcomWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected heartbeat response to be cryptodotcomWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if resp.Method != "public/respond-heartbeat" {
		t.Fatalf("expected Method=public/respond-heartbeat, got %s", resp.Method)
	}
	if resp.ID != 123 {
		t.Fatalf("expected ID=123, got %d", resp.ID)
	}
}

// ---------- Trade handling ----------

func ensureCryptoExchangeMap() {
	// Ensure global Exchanges map has a value for CRYPTODOTCOM_EXCHANGE
	if Exchanges == nil {
		Exchanges = make(map[string]models.Exchange)
	}
	if _, ok := Exchanges[CRYPTODOTCOM_EXCHANGE]; !ok {
		Exchanges[CRYPTODOTCOM_EXCHANGE] = models.Exchange{}
	}
}

func TestCryptodotcomHooks_OnMessage_Trades(t *testing.T) {
	ensureCryptoExchangeMap()

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 10),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	// Ticker key is "BTCUSDT" (no dash)
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "USDT"},
		QuoteToken: models.Asset{Symbol: "BTC"},
	}

	h := cryptodotcomHooks{}
	var lock sync.RWMutex

	now := time.Now()
	// Crypto.com timestamp in microseconds: use recent time to avoid timeout filter
	tsMicro := now.UnixNano() / 1e6

	msg := cryptodotcomWSResponse{
		Method: "subscribe", // any non-heartbeat method
		Result: cryptodotcomWSResponseResult{
			Channel: "trade.BTC_USDT",
			Data: []cryptodotcomWSResponseData{
				{
					TradeID:     "abc123",
					Timestamp:   tsMicro,
					Price:       "30000.5",
					Volume:      "0.1",
					Side:        "BUY",
					ForeignName: "BTC_USDT",
				},
			},
		},
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	// Call OnMessage as text
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// We expect exactly one trade in tradesChannel
	select {
	case tr := <-bs.tradesChannel:
		if tr.Price != 30000.5 {
			t.Fatalf("expected Price=30000.5, got %v", tr.Price)
		}
		if tr.Volume != 0.1 {
			t.Fatalf("expected Volume=0.1, got %v", tr.Volume)
		}
		if tr.BaseToken.Symbol != "USDT" || tr.QuoteToken.Symbol != "BTC" {
			t.Fatalf("unexpected tokens: base=%s quote=%s",
				tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
		}
		if tr.ForeignTradeID != "abc123" {
			t.Fatalf("expected ForeignTradeID=abc123, got %s", tr.ForeignTradeID)
		}
		// lastTradeTimeMap key should be "BTC-USDT"
		lock.RLock()
		_, ok := bs.lastTradeTimeMap["BTC-USDT"]
		lock.RUnlock()
		if !ok {
			t.Fatalf("expected lastTradeTimeMap[BTC-USDT] to be set")
		}
	default:
		t.Fatalf("expected one trade in tradesChannel, but channel is empty")
	}
}

func TestCryptodotcomParseTradeMessage_SellVolumeNegative(t *testing.T) {
	ensureCryptoExchangeMap()

	msg := cryptodotcomWSResponse{
		Result: cryptodotcomWSResponseResult{
			Data: []cryptodotcomWSResponseData{
				{
					TradeID:   "t1",
					Timestamp: time.Now().UnixNano() / 1e3,
					Price:     "100.5",
					Volume:    "2.0",
					Side:      "SELL",
				},
			},
		},
	}

	trades, err := cryptodotcomParseTradeMessage(msg)
	if err != nil {
		t.Fatalf("cryptodotcomParseTradeMessage returned error: %v", err)
	}
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
	if trades[0].Price != 100.5 {
		t.Fatalf("expected Price=100.5, got %v", trades[0].Price)
	}
	if trades[0].Volume != -2.0 {
		t.Fatalf("expected Volume=-2.0 for SELL side, got %v", trades[0].Volume)
	}
}

// ---------- OnMessage ignores cases ----------

func TestCryptodotcomHooks_OnMessage_IgnoresNonTextOrInvalid(t *testing.T) {
	h := cryptodotcomHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	// Non-text message should be ignored
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-text message")
	default:
	}

	// Invalid JSON should be ignored
	h.OnMessage(bs, ws.TextMessage, []byte(`not-json`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for invalid JSON")
	default:
	}
}

// ---------- ReadLoop handled flag ----------

func TestCryptodotcomHooks_ReadLoopHandledFalse(t *testing.T) {
	h := cryptodotcomHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex

	handled := h.ReadLoop(context.Background(), bs, &lock)
	if handled {
		t.Fatalf("expected cryptodotcomHooks.ReadLoop to return false")
	}
}
