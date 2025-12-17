package scrapers

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// --------- Fake wsConn implementation for GateIO tests ---------

// fakeWSConnGateIO is a fake implementation of wsConn.
// It records WriteJSON calls so we can assert on subscription messages.
type fakeWSConnGateIO struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConnGateIO) ReadMessage() (int, []byte, error) {
	// Not needed in these tests.
	return 0, nil, nil
}

func (f *fakeWSConnGateIO) WriteMessage(messageType int, data []byte) error {
	// Not needed in these tests.
	return nil
}

func (f *fakeWSConnGateIO) ReadJSON(v interface{}) error {
	// Not needed in these tests.
	return nil
}

func (f *fakeWSConnGateIO) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConnGateIO) Close() error {
	return nil
}

// --------- Helpers ---------

// ensureGateExchangeMap makes sure the global Exchanges map
// has an entry for GATEIO_EXCHANGE so parseGateIOTrade can set Exchange.
func ensureGateExchangeMap() {
	if Exchanges == nil {
		Exchanges = make(map[string]models.Exchange)
	}
	if _, ok := Exchanges[GATEIO_EXCHANGE]; !ok {
		Exchanges[GATEIO_EXCHANGE] = models.Exchange{}
	}
}

// --------- Key mapping tests ---------

func TestGateIOHooks_TickerAndLastTradeKeys(t *testing.T) {
	h := gateIOHooks{}

	// TickerKeyFromForeign should remove '-' (e.g. BTC-USDT -> BTCUSDT)
	if got := h.TickerKeyFromForeign("BTC-USDT"); got != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", got)
	}

	// LastTradeTimeKeyFromForeign should keep the original foreign name
	if got := h.LastTradeTimeKeyFromForeign("ETH-USDT"); got != "ETH-USDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USDT, got %s", got)
	}
}

// --------- Subscribe / Unsubscribe tests ---------

func TestGateIOHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnGateIO{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := gateIOHooks{}
	var lock sync.RWMutex

	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true -> Event = "subscribe", Channel = "spot.trades", Payload = ["BTC_USDT"]
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call for subscribe, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(*SubscribeGate)
	if !ok {
		t.Fatalf("expected lastWritten to be *SubscribeGate, got %T", fc.lastWritten)
	}
	if msg.Event != "subscribe" {
		t.Fatalf("expected Event=subscribe, got %s", msg.Event)
	}
	if msg.Channel != "spot.trades" {
		t.Fatalf("expected Channel=spot.trades, got %s", msg.Channel)
	}
	if len(msg.Payload) != 1 || msg.Payload[0] != "BTC_USDT" {
		t.Fatalf("expected Payload=[BTC_USDT], got %+v", msg.Payload)
	}

	// subscribe = false -> Event = "unsubscribe"
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(*SubscribeGate)
	if !ok {
		t.Fatalf("expected lastWritten to be *SubscribeGate, got %T", fc.lastWritten)
	}
	if msg.Event != "unsubscribe" {
		t.Fatalf("expected Event=unsubscribe, got %s", msg.Event)
	}
	if msg.Channel != "spot.trades" {
		t.Fatalf("expected Channel=spot.trades, got %s", msg.Channel)
	}
	if len(msg.Payload) != 1 || msg.Payload[0] != "BTC_USDT" {
		t.Fatalf("expected Payload=[BTC_USDT], got %+v", msg.Payload)
	}
}

// --------- parseGateIOTrade tests ---------

func TestParseGateIOTrade_Basic(t *testing.T) {
	ensureGateExchangeMap()

	now := time.Now().Unix()
	msg := GateIOResponseTrade{
		Result: struct {
			ID           int    `json:"id"`
			CreateTime   int    `json:"create_time"`
			CreateTimeMs string `json:"create_time_ms"`
			Side         string `json:"side"`
			CurrencyPair string `json:"currency_pair"`
			Amount       string `json:"amount"`
			Price        string `json:"price"`
		}{
			ID:           123,
			CreateTime:   int(now),
			CreateTimeMs: strconv.FormatInt(now*1000, 10),
			Side:         "buy",
			CurrencyPair: "BTC_USDT",
			Amount:       "0.5",
			Price:        "20000.5",
		},
	}

	bs := &BaseCEXScraper{
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "BTC"},
		QuoteToken: models.Asset{Symbol: "USDT"},
	}

	var lock sync.RWMutex

	tr, ok := parseGateIOTrade(msg, bs, &lock)
	if !ok {
		t.Fatalf("parseGateIOTrade returned ok=false")
	}
	if tr.Price != 20000.5 {
		t.Fatalf("expected Price=20000.5, got %v", tr.Price)
	}
	if tr.Volume != 0.5 {
		t.Fatalf("expected Volume=0.5, got %v", tr.Volume)
	}
	if tr.BaseToken.Symbol != "BTC" || tr.QuoteToken.Symbol != "USDT" {
		t.Fatalf("unexpected tokens: base=%s quote=%s",
			tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
	}
	// ForeignTradeID should be hex of Result.ID
	if tr.ForeignTradeID != strconv.FormatInt(123, 16) {
		t.Fatalf("expected ForeignTradeID=%s, got %s",
			strconv.FormatInt(123, 16), tr.ForeignTradeID)
	}
	if tr.Exchange != Exchanges[GATEIO_EXCHANGE] {
		t.Fatalf("expected Exchange to be GateIO exchange entry")
	}
}

func TestParseGateIOTrade_SellNegativeVolume(t *testing.T) {
	ensureGateExchangeMap()

	msg := GateIOResponseTrade{
		Result: struct {
			ID           int    `json:"id"`
			CreateTime   int    `json:"create_time"`
			CreateTimeMs string `json:"create_time_ms"`
			Side         string `json:"side"`
			CurrencyPair string `json:"currency_pair"`
			Amount       string `json:"amount"`
			Price        string `json:"price"`
		}{
			ID:           1,
			CreateTime:   int(time.Now().Unix()),
			CreateTimeMs: "0",
			Side:         "sell",
			CurrencyPair: "ETH_USDT",
			Amount:       "3.0",
			Price:        "100.0",
		},
	}

	bs := &BaseCEXScraper{
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	bs.tickerPairMap["ETHUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "ETH"},
		QuoteToken: models.Asset{Symbol: "USDT"},
	}
	var lock sync.RWMutex

	tr, ok := parseGateIOTrade(msg, bs, &lock)
	if !ok {
		t.Fatalf("parseGateIOTrade returned ok=false for sell trade")
	}
	if tr.Price != 100.0 {
		t.Fatalf("expected Price=100.0, got %v", tr.Price)
	}
	if tr.Volume != -3.0 {
		t.Fatalf("expected Volume=-3.0 for sell side, got %v", tr.Volume)
	}
}

// --------- OnMessage tests ---------

func TestGateIOHooks_OnMessage_ValidTrade(t *testing.T) {
	ensureGateExchangeMap()

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 10),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	// matching key: BTCUSDT
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "USDT"},
		QuoteToken: models.Asset{Symbol: "BTC"},
	}
	h := gateIOHooks{}
	var lock sync.RWMutex

	now := int(time.Now().Unix())
	msg := GateIOResponseTrade{
		Time:    now,
		Channel: "spot.trades",
		Event:   "update",
		Result: struct {
			ID           int    `json:"id"`
			CreateTime   int    `json:"create_time"`
			CreateTimeMs string `json:"create_time_ms"`
			Side         string `json:"side"`
			CurrencyPair string `json:"currency_pair"`
			Amount       string `json:"amount"`
			Price        string `json:"price"`
		}{
			ID:           42,
			CreateTime:   now,
			CreateTimeMs: strconv.FormatInt(int64(now*1000), 10),
			Side:         "buy",
			CurrencyPair: "BTC_USDT",
			Amount:       "1.25",
			Price:        "25000.0",
		},
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case tr := <-bs.tradesChannel:
		if tr.Price != 25000.0 {
			t.Fatalf("expected Price=25000.0, got %v", tr.Price)
		}
		if tr.Volume != 1.25 {
			t.Fatalf("expected Volume=1.25, got %v", tr.Volume)
		}
		if tr.BaseToken.Symbol != "USDT" || tr.QuoteToken.Symbol != "BTC" {
			t.Fatalf("unexpected tokens: base=%s quote=%s",
				tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
		}
		// lastTradeTimeMap key should be "BTC-USDT" (Quote-Base)
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

func TestGateIOHooks_OnMessage_IgnoresNonTextOrInvalid(t *testing.T) {
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	h := gateIOHooks{}
	var lock sync.RWMutex

	// Non-text message type should be ignored
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

// --------- ReadLoop handled flag ---------

func TestGateIOHooks_ReadLoopHandledFalse(t *testing.T) {
	h := gateIOHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex

	handled := h.ReadLoop(context.Background(), bs, &lock)
	if handled {
		t.Fatalf("expected gateIOHooks.ReadLoop to return false")
	}
}
