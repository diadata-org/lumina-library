package scrapers

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// ------- helper: deflate a payload the same way Bitmart's compressed endpoint does -------

func bitmartDeflate(t *testing.T, b []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		t.Fatalf("flate.NewWriter: %v", err)
	}
	if _, err := w.Write(b); err != nil {
		t.Fatalf("flate write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("flate close: %v", err)
	}
	return buf.Bytes()
}

// ------- Test: TickerKey / LastTradeKey -------

func TestBitmartHooks_TickerAndLastTradeKey(t *testing.T) {
	h := bitmartHooks{}
	if got := h.TickerKeyFromForeign("BTC-USDT"); got != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", got)
	}
	if got := h.LastTradeTimeKeyFromForeign("ETH-USDT"); got != "ETH-USDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USDT, got %s", got)
	}
}

// ------- Test: Subscribe / Unsubscribe -------

func TestBitmartHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConn{}
	bs := &BaseCEXScraper{wsClient: fc}
	h := bitmartHooks{}

	var lock sync.RWMutex
	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(bitmartWSRequest)
	if !ok {
		t.Fatalf("expected lastWritten to be bitmartWSRequest, got %T", fc.lastWritten)
	}
	if msg.Op != "subscribe" {
		t.Fatalf("expected Op=subscribe, got %s", msg.Op)
	}
	if len(msg.Args) != 1 || msg.Args[0] != "spot/trade:BTC_USDT" {
		t.Fatalf("expected Args=[spot/trade:BTC_USDT], got %+v", msg.Args)
	}

	// subscribe = false (unsubscribe)
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(bitmartWSRequest)
	if !ok {
		t.Fatalf("expected lastWritten to be bitmartWSRequest, got %T", fc.lastWritten)
	}
	if msg.Op != "unsubscribe" {
		t.Fatalf("expected Op=unsubscribe, got %s", msg.Op)
	}
	if len(msg.Args) != 1 || msg.Args[0] != "spot/trade:BTC_USDT" {
		t.Fatalf("expected Args=[spot/trade:BTC_USDT], got %+v", msg.Args)
	}
}

// ------- Test: OnMessage parses a valid trade (text frame) -------

func TestBitmartHooks_OnMessage_ValidTrade(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "BTC"},
		QuoteToken: models.Asset{Symbol: "USDT"},
	}

	resp := bitmartWSTradeResponse{
		Table: bitmartTradeChannel,
		Data: []bitmartTradeData{
			{Symbol: "BTC_USDT", Price: "100.5", Side: "buy", Size: "1.23", TimestampSec: 1700000000},
		},
	}
	raw, _ := json.Marshal(resp)

	var lock sync.RWMutex
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case trade := <-bs.tradesChannel:
		if trade.Price != 100.5 {
			t.Fatalf("expected Price=100.5, got %v", trade.Price)
		}
		if trade.Volume != 1.23 {
			t.Fatalf("expected Volume=1.23 (buy), got %v", trade.Volume)
		}
		if trade.BaseToken.Symbol != "BTC" || trade.QuoteToken.Symbol != "USDT" {
			t.Fatalf("unexpected tokens: base=%s quote=%s", trade.BaseToken.Symbol, trade.QuoteToken.Symbol)
		}
		if !trade.Time.Equal(time.Unix(1700000000, 0)) {
			t.Fatalf("expected Time from s_t, got %v", trade.Time)
		}
		lock.RLock()
		_, ok := bs.lastTradeTimeMap["BTC-USDT"]
		lock.RUnlock()
		if !ok {
			t.Fatalf("expected lastTradeTimeMap[BTC-USDT] to be set")
		}
	default:
		t.Fatalf("expected one trade in tradesChannel, channel empty")
	}
}

// ------- Test: sell side flips volume sign -------

func TestBitmartHooks_OnMessage_SellSideNegative(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel: make(chan models.Trade, 1),
		tickerPairMap: map[string]models.Pair{
			"BTCUSDT": {BaseToken: models.Asset{Symbol: "BTC"}, QuoteToken: models.Asset{Symbol: "USDT"}},
		},
		lastTradeTimeMap: make(map[string]time.Time),
	}
	resp := bitmartWSTradeResponse{
		Table: bitmartTradeChannel,
		Data: []bitmartTradeData{
			{Symbol: "BTC_USDT", Price: "100", Side: "sell", Size: "2", TimestampSec: 1700000000},
		},
	}
	raw, _ := json.Marshal(resp)

	var lock sync.RWMutex
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case trade := <-bs.tradesChannel:
		if trade.Volume != -2 {
			t.Fatalf("expected Volume=-2 for sell, got %v", trade.Volume)
		}
	default:
		t.Fatalf("expected one trade, channel empty")
	}
}

// ------- Test: unknown ticker key is skipped -------

func TestBitmartHooks_OnMessage_UnknownTickerSkipped(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair), // empty
		lastTradeTimeMap: make(map[string]time.Time),
	}
	resp := bitmartWSTradeResponse{
		Table: bitmartTradeChannel,
		Data: []bitmartTradeData{
			{Symbol: "DOGE_USDT", Price: "1", Side: "buy", Size: "1", TimestampSec: 1700000000},
		},
	}
	raw, _ := json.Marshal(resp)

	var lock sync.RWMutex
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for unknown ticker key")
	default:
	}
}

// ------- Test: malformed price/size entries are skipped, valid ones still emitted -------

func TestBitmartHooks_OnMessage_SkipsUnparseableEntries(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel: make(chan models.Trade, 2),
		tickerPairMap: map[string]models.Pair{
			"BTCUSDT": {BaseToken: models.Asset{Symbol: "BTC"}, QuoteToken: models.Asset{Symbol: "USDT"}},
		},
		lastTradeTimeMap: make(map[string]time.Time),
	}
	resp := bitmartWSTradeResponse{
		Table: bitmartTradeChannel,
		Data: []bitmartTradeData{
			{Symbol: "BTC_USDT", Price: "not-a-number", Side: "buy", Size: "1", TimestampSec: 1700000000},
			{Symbol: "BTC_USDT", Price: "100", Side: "buy", Size: "1", TimestampSec: 1700000001},
		},
	}
	raw, _ := json.Marshal(resp)

	var lock sync.RWMutex
	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case trade := <-bs.tradesChannel:
		if trade.Price != 100 {
			t.Fatalf("expected the valid entry (Price=100), got %v", trade.Price)
		}
	default:
		t.Fatalf("expected the valid entry to be emitted")
	}
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected only one trade; the unparseable entry should be skipped")
	default:
	}
}

// ------- Test: error envelope, pong, and invalid JSON are ignored -------

func TestBitmartHooks_OnMessage_ErrorAndPongIgnored(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	// error envelope
	errResp, _ := json.Marshal(bitmartWSTradeResponse{ErrorCode: "90006", Event: "subscribe", ErrorMessage: "topic limit"})
	h.OnMessage(bs, ws.TextMessage, errResp, &lock)

	// text pong
	h.OnMessage(bs, ws.TextMessage, []byte(bitmartPongMessage), &lock)

	// invalid JSON
	h.OnMessage(bs, ws.TextMessage, []byte("not-json"), &lock)

	// wrong table
	wrong, _ := json.Marshal(bitmartWSTradeResponse{Table: "spot/depth"})
	h.OnMessage(bs, ws.TextMessage, wrong, &lock)

	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for error/pong/invalid/wrong-table frames")
	default:
	}
}

// ------- Test: OnMessage inflates and parses a compressed binary frame -------

func TestBitmartHooks_OnMessage_BinaryInflate(t *testing.T) {
	h := bitmartHooks{}
	bs := &BaseCEXScraper{
		tradesChannel: make(chan models.Trade, 1),
		tickerPairMap: map[string]models.Pair{
			"BTCUSDT": {BaseToken: models.Asset{Symbol: "BTC"}, QuoteToken: models.Asset{Symbol: "USDT"}},
		},
		lastTradeTimeMap: make(map[string]time.Time),
	}
	resp := bitmartWSTradeResponse{
		Table: bitmartTradeChannel,
		Data: []bitmartTradeData{
			{Symbol: "BTC_USDT", Price: "42", Side: "buy", Size: "1", TimestampSec: 1700000000},
		},
	}
	raw, _ := json.Marshal(resp)
	compressed := bitmartDeflate(t, raw)

	var lock sync.RWMutex
	h.OnMessage(bs, ws.BinaryMessage, compressed, &lock)

	select {
	case trade := <-bs.tradesChannel:
		if trade.Price != 42 {
			t.Fatalf("expected Price=42 from inflated frame, got %v", trade.Price)
		}
	default:
		t.Fatalf("expected one trade from compressed frame, channel empty")
	}
}

// ------- Test: bitmartInflate round-trip, malformed input, and size ceiling -------

func TestBitmartInflate(t *testing.T) {
	original := []byte(`{"table":"spot/trade","data":[]}`)
	compressed := bitmartDeflate(t, original)

	got, err := bitmartInflate(compressed)
	if err != nil {
		t.Fatalf("bitmartInflate round-trip error: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("round-trip mismatch: got %q want %q", got, original)
	}

	// malformed input should error, not panic.
	if _, err := bitmartInflate([]byte{0xff, 0xff, 0xff, 0xff}); err == nil {
		t.Fatalf("expected error for malformed deflate input")
	}

	// a frame that decompresses past the ceiling should error.
	big := bytes.Repeat([]byte("A"), bitmartMaxDecompressed+1024)
	if _, err := bitmartInflate(bitmartDeflate(t, big)); err == nil {
		t.Fatalf("expected error for oversized decompressed frame")
	}
}

// ------- Test: chunkPairs boundaries -------

func TestChunkPairs(t *testing.T) {
	mk := func(n int) []models.ExchangePair {
		out := make([]models.ExchangePair, n)
		for i := range out {
			out[i] = models.ExchangePair{ForeignName: "P"}
		}
		return out
	}

	cases := []struct {
		name       string
		pairs      int
		size       int
		wantChunks int
		wantFirst  int
		wantLast   int
	}{
		{"empty", 0, 70, 0, 0, 0},
		{"size zero single chunk", 5, 0, 1, 5, 5},
		{"size negative single chunk", 5, -3, 1, 5, 5},
		{"exact multiple", 140, 70, 2, 70, 70},
		{"remainder", 150, 70, 3, 70, 10},
		{"size larger than len", 3, 70, 1, 3, 3},
		{"empty with size zero", 0, 0, 0, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := chunkPairs(mk(c.pairs), c.size)
			if len(got) != c.wantChunks {
				t.Fatalf("chunks: got %d want %d", len(got), c.wantChunks)
			}
			if c.wantChunks == 0 {
				return
			}
			if len(got[0]) != c.wantFirst {
				t.Fatalf("first chunk: got %d want %d", len(got[0]), c.wantFirst)
			}
			if len(got[len(got)-1]) != c.wantLast {
				t.Fatalf("last chunk: got %d want %d", len(got[len(got)-1]), c.wantLast)
			}
		})
	}
}