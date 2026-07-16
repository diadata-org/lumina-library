package scrapers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// ---------- Fake wsConn for Bitstamp tests ----------

type fakeWSConnBitstamp struct {
	mu             sync.Mutex
	writeJSONCount int
	lastWritten    interface{}
	writtenEvents  []string
}

func (f *fakeWSConnBitstamp) ReadMessage() (int, []byte, error)               { return 0, nil, nil }
func (f *fakeWSConnBitstamp) WriteMessage(messageType int, data []byte) error { return nil }
func (f *fakeWSConnBitstamp) ReadJSON(v interface{}) error                    { return nil }
func (f *fakeWSConnBitstamp) WriteJSON(v interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.writeJSONCount++
	f.lastWritten = v
	if m, ok := v.(bitstampWSSubscribeMessage); ok {
		f.writtenEvents = append(f.writtenEvents, m.Event)
	}
	if m, ok := v.(bitstampWSPingMessage); ok {
		f.writtenEvents = append(f.writtenEvents, m.Event)
	}
	return nil
}
func (f *fakeWSConnBitstamp) Close() error { return nil }

func (f *fakeWSConnBitstamp) snapshot() (int, []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	evs := make([]string, len(f.writtenEvents))
	copy(evs, f.writtenEvents)
	return f.writeJSONCount, evs
}

func ensureBitstampExchangeMap() {
	if Exchanges == nil {
		Exchanges = make(map[string]models.Exchange)
	}
	if _, ok := Exchanges[BITSTAMP_EXCHANGE]; !ok {
		Exchanges[BITSTAMP_EXCHANGE] = models.Exchange{}
	}
}

func newBitstampTestScraper() *BaseCEXScraper {
	return &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 4),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
}

// ---------- Symbol / channel round-tripping ----------

func TestBitstampSymbolRoundTrip(t *testing.T) {
	tests := []struct {
		name        string
		foreign     string
		wantURL     string
		wantChannel string
	}{
		{"btc-usd", "BTC-USD", "btcusd", "live_trades_btcusd"},
		{"eth-eur", "ETH-EUR", "etheur", "live_trades_etheur"},
		{"already-lower", "btc-usd", "btcusd", "live_trades_btcusd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := bitstampForeignToURLSymbol(tt.foreign)
			if url != tt.wantURL {
				t.Fatalf("bitstampForeignToURLSymbol(%q) = %q, want %q", tt.foreign, url, tt.wantURL)
			}
			channel := "live_trades_" + url
			if channel != tt.wantChannel {
				t.Fatalf("channel = %q, want %q", channel, tt.wantChannel)
			}
			back := bitstampURLSymbolFromChannel(channel)
			if back != tt.wantURL {
				t.Fatalf("bitstampURLSymbolFromChannel(%q) = %q, want %q", channel, back, tt.wantURL)
			}
		})
	}
}

func TestBitstampURLSymbolFromChannel_Unexpected(t *testing.T) {
	// A channel without the live_trades_ prefix must yield the empty string,
	// which handleTrade uses as its early-out guard.
	if got := bitstampURLSymbolFromChannel("order_book_btcusd"); got != "" {
		t.Fatalf("expected empty string for non-trade channel, got %q", got)
	}
	if got := bitstampURLSymbolFromChannel(""); got != "" {
		t.Fatalf("expected empty string for empty channel, got %q", got)
	}
}

// ---------- Key mapping helpers ----------

func TestBitstampHooks_TickerAndLastTradeKeys(t *testing.T) {
	h := bitstampHooks{}

	if got := h.TickerKeyFromForeign("BTC-USD"); got != "BTCUSD" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSD, got %s", got)
	}
	// LastTradeTime key must be the raw ForeignName, because the per-pair
	// watchdog reads lastTradeTimeMap[pair.ForeignName] directly.
	if got := h.LastTradeTimeKeyFromForeign("ETH-EUR"); got != "ETH-EUR" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-EUR, got %s", got)
	}
}

// TestBitstampTickerKeyMatchesMakeTickerPairMap locks in the invariant that the
// key handleTrade derives (strings.ToUpper(url_symbol)) is exactly the key
// models.MakeTickerPairMap produces for the same pair.
func TestBitstampTickerKeyMatchesMakeTickerPairMap(t *testing.T) {
	foreign := "BTC-USD"
	ep := models.ExchangePair{
		ForeignName: foreign,
		UnderlyingPair: models.Pair{
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USD"},
		},
	}
	m := models.MakeTickerPairMap([]models.ExchangePair{ep})

	urlSymbol := bitstampForeignToURLSymbol(foreign) // "btcusd"
	handleTradeKey := strings.ToUpper(urlSymbol)     // "BTCUSD"

	if _, ok := m[handleTradeKey]; !ok {
		t.Fatalf("handleTrade key %q not found in MakeTickerPairMap keys %v", handleTradeKey, keysOf(m))
	}
}

// ---------- Subscribe / Unsubscribe ----------

func TestBitstampHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnBitstamp{}
	bs := &BaseCEXScraper{wsClient: fc}
	h := bitstampHooks{}
	var lock sync.RWMutex

	pair := models.ExchangePair{ForeignName: "BTC-USD"}

	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) error: %v", err)
	}
	count, _ := fc.snapshot()
	if count != 1 {
		t.Fatalf("expected 1 WriteJSON call, got %d", count)
	}
	msg, ok := fc.lastWritten.(bitstampWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected bitstampWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Event != "bts:subscribe" {
		t.Fatalf("expected Event=bts:subscribe, got %s", msg.Event)
	}
	if msg.Data.Channel != "live_trades_btcusd" {
		t.Fatalf("expected Channel=live_trades_btcusd, got %s", msg.Data.Channel)
	}

	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) error: %v", err)
	}
	count, _ = fc.snapshot()
	if count != 2 {
		t.Fatalf("expected 2 WriteJSON calls, got %d", count)
	}
	msg, _ = fc.lastWritten.(bitstampWSSubscribeMessage)
	if msg.Event != "bts:unsubscribe" {
		t.Fatalf("expected Event=bts:unsubscribe, got %s", msg.Event)
	}
	if msg.Data.Channel != "live_trades_btcusd" {
		t.Fatalf("expected Channel=live_trades_btcusd, got %s", msg.Data.Channel)
	}
}

// ---------- bitstampParseTrade ----------

func TestBitstampParseTrade(t *testing.T) {
	ensureBitstampExchangeMap()

	micro := int64(1_700_000_000_000_000) // microseconds since epoch
	tests := []struct {
		name       string
		in         bitstampWSTradeData
		wantVolume float64
		wantPrice  float64
		wantID     string
		wantErr    bool
		wantTimeNS int64
	}{
		{
			name: "buy uses price_str/amount_str",
			in: bitstampWSTradeData{
				ID: 42, Type: 0,
				Amount: 0.4, AmountStr: "0.5", // str differs from float to prove str wins
				Price: 30000.4, PriceStr: "30000.5",
				Microtimestamp: strconv.FormatInt(micro, 10),
			},
			wantVolume: 0.5,
			wantPrice:  30000.5,
			wantID:     "42",
			wantTimeNS: micro * 1000,
		},
		{
			name: "sell negative volume from amount_str",
			in: bitstampWSTradeData{
				ID: 7, Type: 1,
				Amount: 0, AmountStr: "2.0",
				Price: 0, PriceStr: "100.5",
				Microtimestamp: strconv.FormatInt(micro, 10),
			},
			wantVolume: -2.0,
			wantPrice:  100.5,
			wantID:     "7",
			wantTimeNS: micro * 1000,
		},
		{
			name: "falls back to float64 when str fields empty",
			in: bitstampWSTradeData{
				ID: 9, Type: 0,
				Amount: 1.25, AmountStr: "",
				Price: 55.5, PriceStr: "",
				Microtimestamp: strconv.FormatInt(micro, 10),
			},
			wantVolume: 1.25,
			wantPrice:  55.5,
			wantID:     "9",
			wantTimeNS: micro * 1000,
		},
		{
			name: "malformed price_str errors",
			in: bitstampWSTradeData{
				ID: 1, Type: 0,
				PriceStr: "not-a-number", AmountStr: "1",
				Microtimestamp: strconv.FormatInt(micro, 10),
			},
			wantErr: true,
		},
		{
			name: "malformed amount_str errors",
			in: bitstampWSTradeData{
				ID: 1, Type: 0,
				PriceStr: "1", AmountStr: "bad",
				Microtimestamp: strconv.FormatInt(micro, 10),
			},
			wantErr: true,
		},
		{
			name: "malformed microtimestamp errors",
			in: bitstampWSTradeData{
				ID: 1, Type: 0,
				PriceStr: "1", AmountStr: "1",
				Microtimestamp: "not-a-number",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := bitstampParseTrade(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tr.Volume != tt.wantVolume {
				t.Fatalf("Volume = %v, want %v", tr.Volume, tt.wantVolume)
			}
			if tr.Price != tt.wantPrice {
				t.Fatalf("Price = %v, want %v", tr.Price, tt.wantPrice)
			}
			if tr.ForeignTradeID != tt.wantID {
				t.Fatalf("ForeignTradeID = %s, want %s", tr.ForeignTradeID, tt.wantID)
			}
			if tr.Time.UnixNano() != tt.wantTimeNS {
				t.Fatalf("Time.UnixNano() = %d, want %d", tr.Time.UnixNano(), tt.wantTimeNS)
			}
		})
	}
}

// ---------- OnMessage: trade routing ----------

func TestBitstampOnMessage_Trade(t *testing.T) {
	ensureBitstampExchangeMap()
	bs := newBitstampTestScraper()
	bs.tickerPairMap["BTCUSD"] = models.Pair{
		QuoteToken: models.Asset{Symbol: "BTC"},
		BaseToken:  models.Asset{Symbol: "USD"},
	}
	h := bitstampHooks{}
	var lock sync.RWMutex

	micro := time.Now().UnixNano() / int64(time.Microsecond)
	td := bitstampWSTradeData{
		ID: 123, Type: 0,
		PriceStr: "30000.5", AmountStr: "0.1",
		Microtimestamp: strconv.FormatInt(micro, 10),
	}
	dataRaw, _ := json.Marshal(td)
	raw, _ := json.Marshal(bitstampWSResponse{Event: "trade", Channel: "live_trades_btcusd", Data: dataRaw})

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case tr := <-bs.tradesChannel:
		if tr.Price != 30000.5 || tr.Volume != 0.1 {
			t.Fatalf("unexpected trade price/volume: %v/%v", tr.Price, tr.Volume)
		}
		if tr.QuoteToken.Symbol != "BTC" || tr.BaseToken.Symbol != "USD" {
			t.Fatalf("unexpected tokens: quote=%s base=%s", tr.QuoteToken.Symbol, tr.BaseToken.Symbol)
		}
		lock.RLock()
		_, ok := bs.lastTradeTimeMap["BTC-USD"]
		lock.RUnlock()
		if !ok {
			t.Fatalf("expected lastTradeTimeMap[BTC-USD] to be set")
		}
	default:
		t.Fatalf("expected a trade, channel empty")
	}
}

// TestBitstampOnMessage_StaleTradeRefreshesLiveness verifies the decoupling of
// liveness tracking from the freshness filter: a stale trade must NOT be
// forwarded to the channel, but MUST still refresh lastTradeTime so the
// watchdog does not falsely fail over during clock skew.
func TestBitstampOnMessage_StaleTradeRefreshesLiveness(t *testing.T) {
	ensureBitstampExchangeMap()
	bs := newBitstampTestScraper()
	bs.tickerPairMap["BTCUSD"] = models.Pair{
		QuoteToken: models.Asset{Symbol: "BTC"},
		BaseToken:  models.Asset{Symbol: "USD"},
	}
	h := bitstampHooks{}
	var lock sync.RWMutex

	before := time.Now()
	staleMicro := time.Now().Add(-time.Duration(bitstampTradeTimeoutSeconds+60)*time.Second).UnixNano() / int64(time.Microsecond)
	td := bitstampWSTradeData{
		ID: 9, Type: 0,
		PriceStr: "1", AmountStr: "1",
		Microtimestamp: strconv.FormatInt(staleMicro, 10),
	}
	dataRaw, _ := json.Marshal(td)
	raw, _ := json.Marshal(bitstampWSResponse{Event: "trade", Channel: "live_trades_btcusd", Data: dataRaw})

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// Not forwarded.
	select {
	case <-bs.tradesChannel:
		t.Fatalf("stale trade must not be forwarded to the channel")
	default:
	}

	// But liveness refreshed with wall-clock now (not the stale trade time).
	lock.RLock()
	ts, ok := bs.lastTradeTimeMap["BTC-USD"]
	lock.RUnlock()
	if !ok {
		t.Fatalf("expected lastTradeTimeMap[BTC-USD] to be set even for stale trade")
	}
	if ts.Before(before) {
		t.Fatalf("expected lastTradeTime to be refreshed to wall-clock now, got stale %v", ts)
	}
}

// ---------- OnMessage: control / non-trade events ----------

func TestBitstampOnMessage_ControlEvents(t *testing.T) {
	bs := newBitstampTestScraper()
	h := bitstampHooks{}
	var lock sync.RWMutex

	controls := []bitstampWSResponse{
		{Event: "bts:subscription_succeeded", Channel: "live_trades_btcusd"},
		{Event: "bts:unsubscription_succeeded", Channel: "live_trades_btcusd"},
		{Event: "bts:request_reconnect", Channel: "live_trades_btcusd"},
		{Event: "bts:pong"},
		{Event: "bts:heartbeat", Data: json.RawMessage(`{"status":"success"}`)},
		{Event: "bts:heartbeat", Data: json.RawMessage(`{"status":"failure"}`)},
		{Event: "some_unknown_event", Channel: "live_trades_btcusd"},
	}
	for _, c := range controls {
		raw, _ := json.Marshal(c)
		h.OnMessage(bs, ws.TextMessage, raw, &lock)
	}
	// None of these should emit a trade.
	select {
	case <-bs.tradesChannel:
		t.Fatalf("control/unknown events must not emit a trade")
	default:
	}
}

func TestBitstampOnMessage_IgnoresNonTextOrInvalid(t *testing.T) {
	bs := newBitstampTestScraper()
	h := bitstampHooks{}
	var lock sync.RWMutex

	// Non-text ignored.
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)
	// Invalid JSON ignored.
	h.OnMessage(bs, ws.TextMessage, []byte(`not-json`), &lock)

	// Trade for an unknown pair is dropped (no tickerPairMap entry).
	micro := time.Now().UnixNano() / int64(time.Microsecond)
	td := bitstampWSTradeData{ID: 1, Type: 0, PriceStr: "1", AmountStr: "1", Microtimestamp: strconv.FormatInt(micro, 10)}
	dataRaw, _ := json.Marshal(td)
	raw, _ := json.Marshal(bitstampWSResponse{Event: "trade", Channel: "live_trades_unknown", Data: dataRaw})
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade")
	default:
	}
}

// ---------- ReadLoop handled flag ----------

func TestBitstampHooks_ReadLoopHandledFalse(t *testing.T) {
	h := bitstampHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex
	if handled := h.ReadLoop(context.Background(), bs, &lock); handled {
		t.Fatalf("expected ReadLoop to return false")
	}
}

// ---------- Client ping loop ----------

// TestBitstampPingLoop_SendsPing verifies the keepalive goroutine writes a
// bts:ping when the interval elapses.
func TestBitstampPingLoop_SendsPing(t *testing.T) {
	fc := &fakeWSConnBitstamp{}
	bs := &BaseCEXScraper{wsClient: fc}

	// Shrink the interval for the test and restore afterwards.
	orig := bitstampPingInterval
	bitstampPingInterval = 10 * time.Millisecond
	defer func() { bitstampPingInterval = orig }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		bitstampPingLoop(ctx, bs)
		close(done)
	}()

	// Wait long enough for at least one tick.
	time.Sleep(40 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("pingLoop did not exit after ctx cancel (possible goroutine leak)")
	}

	count, events := fc.snapshot()
	if count == 0 {
		t.Fatalf("expected at least one ping write")
	}
	foundPing := false
	for _, e := range events {
		if e == "bts:ping" {
			foundPing = true
			break
		}
	}
	if !foundPing {
		t.Fatalf("expected a bts:ping event, got %v", events)
	}
}

// TestBitstampPingLoop_ExitsOnCancel verifies the goroutine returns promptly
// when ctx is canceled before any tick (no leak, no send).
func TestBitstampPingLoop_ExitsOnCancel(t *testing.T) {
	fc := &fakeWSConnBitstamp{}
	bs := &BaseCEXScraper{wsClient: fc}

	orig := bitstampPingInterval
	bitstampPingInterval = time.Hour // ensure no tick fires
	defer func() { bitstampPingInterval = orig }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		bitstampPingLoop(ctx, bs)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("pingLoop did not exit promptly on ctx cancel")
	}

	if count, _ := fc.snapshot(); count != 0 {
		t.Fatalf("expected no ping writes before first tick, got %d", count)
	}
}

// ---------- test-local helper ----------

func keysOf(m map[string]models.Pair) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}