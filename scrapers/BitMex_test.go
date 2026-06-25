package scrapers

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// ---------------------------------------------------------------------------
// Helpers
//
// NOTE (verify against your local BaseCEXScraper definition):
//   - These tests assume `tickerPairMap` and `tradesChannel` are package-visible
//     fields on *BaseCEXScraper that can be populated directly in-package.
//   - They assume `setLastTradeTime(lock, foreignName, t)` is safe to call with
//     a real (possibly zero-value) scraper and does not panic on a nil/empty
//     internal map. If it dereferences an uninitialised map, initialise that
//     field here too.
// If BaseCEXScraper hides these or needs more setup, adjust newTestScraper.
// ---------------------------------------------------------------------------

func mustAsset(symbol string) models.Asset {
	// Minimal asset; only Symbol is exercised by the trade path / logging.
	return models.Asset{Symbol: symbol}
}

// newTestScraper builds a BaseCEXScraper wired with a buffered trades channel
// and a ticker pair map keyed exactly the way MakeTickerPairMap keys it for
// BitMEX (separator stripped, e.g. "ETHUSDT").
func newTestScraper(t *testing.T, bufferedTrades int) (*BaseCEXScraper, chan models.Trade) {
	t.Helper()

	tradesCh := make(chan models.Trade, bufferedTrades)

	ethPair := models.Pair{
		QuoteToken: mustAsset("ETH"),
		BaseToken:  mustAsset("USDT"),
	}
	btcPair := models.Pair{
		QuoteToken: mustAsset("XBT"),
		BaseToken:  mustAsset("USDT"),
	}

	bs := &BaseCEXScraper{
		tickerPairMap: map[string]models.Pair{
			"ETHUSDT": ethPair,
			"XBTUSDT": btcPair,
		},
		tradesChannel: tradesCh,
		// setLastTradeTime writes into this map; it must be non-nil or the
		// write panics.
		lastTradeTimeMap: make(map[string]time.Time),
	}
	return bs, tradesCh
}

// drain returns all trades currently buffered in ch without blocking.
func drain(ch chan models.Trade) []models.Trade {
	var out []models.Trade
	for {
		select {
		case tr := <-ch:
			out = append(out, tr)
		default:
			return out
		}
	}
}

// ---------------------------------------------------------------------------
// bitMexNativeSymbol — pins the symbol contract (no separator on the wire)
// ---------------------------------------------------------------------------

func TestBitMexNativeSymbol(t *testing.T) {
	cases := []struct {
		name    string
		foreign string
		want    string
	}{
		{"eth_usdt", "ETH-USDT", "ETHUSDT"},
		{"xbt_usdt", "XBT-USDT", "XBTUSDT"},
		{"already_stripped", "ETHUSDT", "ETHUSDT"},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := bitMexNativeSymbol(tc.foreign); got != tc.want {
				t.Fatalf("bitMexNativeSymbol(%q) = %q, want %q", tc.foreign, got, tc.want)
			}
		})
	}
}

func TestBitMexTickerKeyFromForeign(t *testing.T) {
	// The lookup key must match MakeTickerPairMap's stripped key.
	if got := (bitMexHooks{}).TickerKeyFromForeign("ETH-USDT"); got != "ETHUSDT" {
		t.Fatalf("TickerKeyFromForeign(ETH-USDT) = %q, want ETHUSDT", got)
	}
}

// ---------------------------------------------------------------------------
// processBitMexTrades — sign convention, unknown-symbol skip, fields
// ---------------------------------------------------------------------------

func TestProcessBitMexTrades(t *testing.T) {
	ts := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name        string
		in          []bitMexWSTrade
		wantCount   int
		wantVolumes []float64 // expected Volume for each emitted trade, in order
		wantPrices  []float64
	}{
		{
			name: "buy_is_positive",
			in: []bitMexWSTrade{
				{Symbol: "ETHUSDT", Side: "Buy", Size: 2.5, Price: 3000, Timestamp: ts, TrdMatchID: "a"},
			},
			wantCount:   1,
			wantVolumes: []float64{2.5},
			wantPrices:  []float64{3000},
		},
		{
			name: "sell_is_negative",
			in: []bitMexWSTrade{
				{Symbol: "ETHUSDT", Side: "Sell", Size: 1.0, Price: 3001, Timestamp: ts, TrdMatchID: "b"},
			},
			wantCount:   1,
			wantVolumes: []float64{-1.0},
			wantPrices:  []float64{3001},
		},
		{
			name: "unknown_symbol_skipped",
			in: []bitMexWSTrade{
				{Symbol: "DOGEUSDT", Side: "Buy", Size: 5, Price: 0.1, Timestamp: ts, TrdMatchID: "c"},
			},
			wantCount: 0,
		},
		{
			name: "mixed_known_and_unknown",
			in: []bitMexWSTrade{
				{Symbol: "ETHUSDT", Side: "Buy", Size: 1, Price: 3000, Timestamp: ts, TrdMatchID: "d"},
				{Symbol: "NOPEUSDT", Side: "Sell", Size: 9, Price: 1, Timestamp: ts, TrdMatchID: "e"},
				{Symbol: "XBTUSDT", Side: "Sell", Size: 0.5, Price: 60000, Timestamp: ts, TrdMatchID: "f"},
			},
			wantCount:   2,
			wantVolumes: []float64{1, -0.5},
			wantPrices:  []float64{3000, 60000},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bs, ch := newTestScraper(t, len(tc.in))
			var lock sync.RWMutex

			processBitMexTrades(bs, &lock, tc.in)

			got := drain(ch)
			if len(got) != tc.wantCount {
				t.Fatalf("emitted %d trades, want %d", len(got), tc.wantCount)
			}
			for i := range got {
				if got[i].Volume != tc.wantVolumes[i] {
					t.Errorf("trade[%d].Volume = %v, want %v", i, got[i].Volume, tc.wantVolumes[i])
				}
				if got[i].Price != tc.wantPrices[i] {
					t.Errorf("trade[%d].Price = %v, want %v", i, got[i].Price, tc.wantPrices[i])
				}
			}
		})
	}
}

func TestProcessBitMexTradesMapsTokens(t *testing.T) {
	bs, ch := newTestScraper(t, 1)
	var lock sync.RWMutex

	processBitMexTrades(bs, &lock, []bitMexWSTrade{
		{Symbol: "ETHUSDT", Side: "Buy", Size: 1, Price: 3000, Timestamp: time.Now(), TrdMatchID: "x"},
	})

	got := drain(ch)
	if len(got) != 1 {
		t.Fatalf("emitted %d trades, want 1", len(got))
	}
	if got[0].QuoteToken.Symbol != "ETH" || got[0].BaseToken.Symbol != "USDT" {
		t.Fatalf("token mapping wrong: quote=%q base=%q, want ETH/USDT",
			got[0].QuoteToken.Symbol, got[0].BaseToken.Symbol)
	}
	if got[0].ForeignTradeID != "x" {
		t.Errorf("ForeignTradeID = %q, want x", got[0].ForeignTradeID)
	}
}

// ---------------------------------------------------------------------------
// OnMessage — partial vs insert filtering (the staleness guard)
// ---------------------------------------------------------------------------

func TestBitMexOnMessageDropsPartial(t *testing.T) {
	bs, ch := newTestScraper(t, 4)
	var lock sync.RWMutex
	h := bitMexHooks{}

	partial := mustJSON(t, bitMexWSResponse{
		Table:  "trade",
		Action: "partial",
		Data: []bitMexWSTrade{
			{Symbol: "ETHUSDT", Side: "Buy", Size: 1, Price: 3000, Timestamp: time.Now(), TrdMatchID: "p1"},
			{Symbol: "ETHUSDT", Side: "Sell", Size: 2, Price: 3001, Timestamp: time.Now(), TrdMatchID: "p2"},
		},
	})

	h.OnMessage(bs, ws.TextMessage, partial, &lock)

	if got := drain(ch); len(got) != 0 {
		t.Fatalf("partial snapshot emitted %d trades, want 0 (must be dropped)", len(got))
	}
}

func TestBitMexOnMessageProcessesInsert(t *testing.T) {
	bs, ch := newTestScraper(t, 4)
	var lock sync.RWMutex
	h := bitMexHooks{}

	insert := mustJSON(t, bitMexWSResponse{
		Table:  "trade",
		Action: "insert",
		Data: []bitMexWSTrade{
			{Symbol: "ETHUSDT", Side: "Buy", Size: 1.5, Price: 3000, Timestamp: time.Now(), TrdMatchID: "i1"},
		},
	})

	h.OnMessage(bs, ws.TextMessage, insert, &lock)

	got := drain(ch)
	if len(got) != 1 {
		t.Fatalf("insert emitted %d trades, want 1", len(got))
	}
	if got[0].Volume != 1.5 {
		t.Errorf("Volume = %v, want 1.5", got[0].Volume)
	}
}

func TestBitMexOnMessageIgnoresNonTradeFrames(t *testing.T) {
	bs, ch := newTestScraper(t, 4)
	var lock sync.RWMutex
	h := bitMexHooks{}

	// pong, ack, and error frames must all produce no trades and not panic.
	frames := [][]byte{
		[]byte("pong"),
		mustJSON(t, bitMexWSResponse{Success: true, Subscribe: "trade:ETHUSDT"}),
		mustJSON(t, bitMexWSResponse{Error: "Unknown table", Status: 400}),
	}
	for _, f := range frames {
		h.OnMessage(bs, ws.TextMessage, f, &lock)
	}

	if got := drain(ch); len(got) != 0 {
		t.Fatalf("non-trade frames emitted %d trades, want 0", len(got))
	}
}

func mustJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}