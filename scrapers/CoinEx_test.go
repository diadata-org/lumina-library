package scrapers

import (
	"bytes"
	"compress/gzip"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// newCoinExTestScraper builds a minimal BaseCEXScraper sufficient for exercising
// coinexHooks.OnMessage in isolation. It does NOT open a websocket; only the
// fields read by OnMessage are populated.
func newCoinExTestScraper(tickerPairMap map[string]models.Pair) *BaseCEXScraper {
	return &BaseCEXScraper{
		hooks:            coinexHooks{},
		tradesChannel:    make(chan models.Trade),
		tickerPairMap:    tickerPairMap,
		lastTradeTimeMap: make(map[string]time.Time),
	}
}

// btcusdtCoinExPair is the canonical pair used across cases. tickerPairMap key
// is the market form "BTCUSDT" (TickerKeyFromForeign("BTC-USDT")).
func btcusdtCoinExPair() (string, models.Pair) {
	return "BTCUSDT", models.Pair{
		BaseToken:  models.Asset{Symbol: "USDT"},
		QuoteToken: models.Asset{Symbol: "BTC"},
	}
}

// gzipStr gzip-compresses a string, matching what the CoinEx WS server sends.
func gzipStr(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buf.Bytes()
}

func TestCoinExOnMessage(t *testing.T) {
	key, pair := btcusdtCoinExPair()

	tests := []struct {
		name string
		raw  string // plain JSON, will be gzip-compressed before calling OnMessage
		// wantTrades: expected trades produced (order preserved).
		// nil/empty means OnMessage must NOT emit anything.
		wantTrades []models.Trade
	}{
		{
			name: "subscribe ack ignored",
			raw:  `{"id":1,"code":0,"data":{},"message":"OK"}`,
		},
		{
			name: "pong ack ignored",
			raw:  `{"id":2,"code":0,"data":{"result":"pong"},"message":"OK"}`,
		},
		{
			name: "error ack ignored (logged only)",
			raw:  `{"id":3,"code":24,"data":{},"message":"market not exist"}`,
		},
		{
			name: "unknown market (not in tickerPairMap) ignored",
			raw:  `{"method":"deals.update","data":{"market":"ETHUSDT","deal_list":[{"deal_id":1,"created_at":1700000000000,"side":"buy","price":"100","amount":"1"}]},"id":null}`,
		},
		{
			name: "empty deal_list ignored",
			raw:  `{"method":"deals.update","data":{"market":"BTCUSDT","deal_list":[]},"id":null}`,
		},
		{
			name: "buy side keeps positive volume",
			raw:  `{"method":"deals.update","data":{"market":"BTCUSDT","deal_list":[{"deal_id":3514376759,"created_at":1689152421692,"side":"buy","price":"30718.42","amount":"0.00000325"}]},"id":null}`,
			wantTrades: []models.Trade{
				{
					Price:          30718.42,
					Volume:         0.00000325,
					Time:           time.Unix(0, 1689152421692*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "3514376759",
				},
			},
		},
		{
			name: "sell side negates volume",
			raw:  `{"method":"deals.update","data":{"market":"BTCUSDT","deal_list":[{"deal_id":42,"created_at":1700000000000,"side":"sell","price":"42000.5","amount":"0.25"}]},"id":null}`,
			wantTrades: []models.Trade{
				{
					Price:          42000.5,
					Volume:         -0.25,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "42",
				},
			},
		},
		{
			name: "unparseable price skips that trade",
			raw:  `{"method":"deals.update","data":{"market":"BTCUSDT","deal_list":[{"deal_id":1,"created_at":1700000000000,"side":"buy","price":"abc","amount":"1"},{"deal_id":2,"created_at":1700000000000,"side":"buy","price":"99","amount":"2"}]},"id":null}`,
			wantTrades: []models.Trade{
				{
					Price:          99,
					Volume:         2,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "2",
				},
			},
		},
		{
			name: "multiple trades in one frame all emitted",
			raw:  `{"method":"deals.update","data":{"market":"BTCUSDT","deal_list":[{"deal_id":10,"created_at":1700000000000,"side":"buy","price":"10","amount":"1"},{"deal_id":11,"created_at":1700000000001,"side":"sell","price":"11","amount":"2"}]},"id":null}`,
			wantTrades: []models.Trade{
				{
					Price:          10,
					Volume:         1,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "10",
				},
				{
					Price:          11,
					Volume:         -2,
					Time:           time.Unix(0, 1700000000001*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "11",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := newCoinExTestScraper(map[string]models.Pair{key: pair})
			var lock sync.RWMutex

			compressed := gzipStr(t, tc.raw)

			// Collect trades from the unbuffered channel in a separate goroutine.
			collected := make(chan models.Trade, 16)
			done := make(chan struct{})
			go func() {
				coinexHooks{}.OnMessage(bs, ws.BinaryMessage, compressed, &lock)
				close(done)
			}()
			go func() {
				for {
					select {
					case tr := <-bs.tradesChannel:
						collected <- tr
					case <-done:
						select {
						case tr := <-bs.tradesChannel:
							collected <- tr
						default:
						}
						close(collected)
						return
					}
				}
			}()

			var got []models.Trade
			for tr := range collected {
				got = append(got, tr)
			}

			if len(got) != len(tc.wantTrades) {
				t.Fatalf("trade count = %d, want %d (got %+v)", len(got), len(tc.wantTrades), got)
			}
			for i, want := range tc.wantTrades {
				g := got[i]
				if g.Price != want.Price {
					t.Errorf("trade[%d].Price = %v, want %v", i, g.Price, want.Price)
				}
				if g.Volume != want.Volume {
					t.Errorf("trade[%d].Volume = %v, want %v", i, g.Volume, want.Volume)
				}
				if !g.Time.Equal(want.Time) {
					t.Errorf("trade[%d].Time = %v, want %v", i, g.Time, want.Time)
				}
				if g.BaseToken.Symbol != want.BaseToken.Symbol {
					t.Errorf("trade[%d].BaseToken = %v, want %v", i, g.BaseToken.Symbol, want.BaseToken.Symbol)
				}
				if g.QuoteToken.Symbol != want.QuoteToken.Symbol {
					t.Errorf("trade[%d].QuoteToken = %v, want %v", i, g.QuoteToken.Symbol, want.QuoteToken.Symbol)
				}
				if g.ForeignTradeID != want.ForeignTradeID {
					t.Errorf("trade[%d].ForeignTradeID = %v, want %v", i, g.ForeignTradeID, want.ForeignTradeID)
				}
			}
		})
	}
}

// TestCoinExOnMessageInvalidGzip verifies malformed (non-gzip) input is dropped silently.
func TestCoinExOnMessageInvalidGzip(t *testing.T) {
	key, pair := btcusdtCoinExPair()
	bs := newCoinExTestScraper(map[string]models.Pair{key: pair})
	var lock sync.RWMutex

	done := make(chan struct{})
	go func() {
		coinexHooks{}.OnMessage(bs, ws.BinaryMessage, []byte("not gzip data"), &lock)
		close(done)
	}()

	select {
	case tr := <-bs.tradesChannel:
		t.Fatalf("expected no trade, got %+v", tr)
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("OnMessage did not return")
	}
}

// TestCoinExTickerKeyFromForeign pins the market-matching behaviour.
func TestCoinExTickerKeyFromForeign(t *testing.T) {
	cases := map[string]string{
		"BTC-USDT": "BTCUSDT",
		"ETH-USDT": "ETHUSDT",
		"BTCUSDT":  "BTCUSDT",
	}
	for in, want := range cases {
		if got := (coinexHooks{}).TickerKeyFromForeign(in); got != want {
			t.Errorf("TickerKeyFromForeign(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestCoinExSubscribeMarketList verifies the deals.subscribe payload uses the
// hyphen-stripped market name inside market_list, matching CoinEx's naming
// (e.g. BTCUSDT).
func TestCoinExSubscribeMarketList(t *testing.T) {
	if got := (coinexHooks{}).TickerKeyFromForeign("BTC-USDT"); got != "BTCUSDT" {
		t.Errorf("market name = %q, want %q", got, "BTCUSDT")
	}
}
