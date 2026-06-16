package scrapers

import (
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// newBitgetTestScraper builds a minimal BaseCEXScraper sufficient for exercising
// bitgetHooks.OnMessage in isolation. It does NOT open a websocket; only the
// fields read by OnMessage are populated.
func newBitgetTestScraper(tickerPairMap map[string]models.Pair) *BaseCEXScraper {
	return &BaseCEXScraper{
		hooks:            bitgetHooks{},
		tradesChannel:    make(chan models.Trade),
		tickerPairMap:    tickerPairMap,
		lastTradeTimeMap: make(map[string]time.Time),
	}
}

// btcusdtPair is the canonical pair used across cases. tickerPairMap key is the
// instId form "BTCUSDT" (TickerKeyFromForeign("BTC-USDT")).
func btcusdtPair() (string, models.Pair) {
	return "BTCUSDT", models.Pair{
		BaseToken:  models.Asset{Symbol: "BTC"},
		QuoteToken: models.Asset{Symbol: "USDT"},
	}
}

func TestBitgetOnMessage(t *testing.T) {
	key, pair := btcusdtPair()

	tests := []struct {
		name string
		mt   int
		data string
		// wantTrades: expected trades produced (order preserved).
		// nil/empty means OnMessage must NOT emit anything.
		wantTrades []models.Trade
	}{
		{
			name: "non-text message ignored",
			mt:   ws.BinaryMessage,
			data: `{"action":"snapshot"}`,
		},
		{
			name: "pong heartbeat short-circuits",
			mt:   ws.TextMessage,
			data: "pong",
		},
		{
			name: "invalid json ignored",
			mt:   ws.TextMessage,
			data: `{not-json`,
		},
		{
			name: "subscribe ack (event set) ignored",
			mt:   ws.TextMessage,
			data: `{"event":"subscribe","arg":{"instType":"SPOT","channel":"trade","instId":"BTCUSDT"}}`,
		},
		{
			name: "error frame ignored (logged only)",
			mt:   ws.TextMessage,
			data: `{"event":"error","code":"30001","msg":"channel not exist"}`,
		},
		{
			name: "non-trade channel ignored",
			mt:   ws.TextMessage,
			data: `{"action":"snapshot","arg":{"channel":"ticker","instId":"BTCUSDT"},"data":[{"ts":"1700000000000","price":"100","size":"1","side":"buy"}]}`,
		},
		{
			name: "unknown instId (not in tickerPairMap) ignored",
			mt:   ws.TextMessage,
			data: `{"action":"snapshot","arg":{"channel":"trade","instId":"ETHUSDT"},"data":[{"ts":"1700000000000","price":"100","size":"1","side":"buy"}]}`,
		},
		{
			name: "buy side keeps positive volume",
			mt:   ws.TextMessage,
			data: `{"action":"update","arg":{"channel":"trade","instId":"BTCUSDT"},"data":[{"ts":"1700000000000","price":"42000.5","size":"0.25","side":"buy","tradeId":"t1"}]}`,
			wantTrades: []models.Trade{
				{
					Price:          42000.5,
					Volume:         0.25,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "t1",
				},
			},
		},
		{
			name: "sell side negates volume",
			mt:   ws.TextMessage,
			data: `{"action":"update","arg":{"channel":"trade","instId":"BTCUSDT"},"data":[{"ts":"1700000000000","price":"42000.5","size":"0.25","side":"sell","tradeId":"t2"}]}`,
			wantTrades: []models.Trade{
				{
					Price:          42000.5,
					Volume:         -0.25,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "t2",
				},
			},
		},
		{
			name: "millisecond timestamp parsed correctly",
			mt:   ws.TextMessage,
			data: `{"action":"snapshot","arg":{"channel":"trade","instId":"BTCUSDT"},"data":[{"ts":"1712345678901","price":"1","size":"1","side":"buy","tradeId":"t3"}]}`,
			wantTrades: []models.Trade{
				{
					Price:          1,
					Volume:         1,
					Time:           time.Unix(0, 1712345678901*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "t3",
				},
			},
		},
		{
			name: "unparseable price skips that trade",
			mt:   ws.TextMessage,
			data: `{"action":"update","arg":{"channel":"trade","instId":"BTCUSDT"},"data":[{"ts":"1700000000000","price":"abc","size":"1","side":"buy","tradeId":"bad"},{"ts":"1700000000000","price":"99","size":"2","side":"buy","tradeId":"good"}]}`,
			wantTrades: []models.Trade{
				{
					Price:          99,
					Volume:         2,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "good",
				},
			},
		},
		{
			name: "multiple trades in one frame all emitted",
			mt:   ws.TextMessage,
			data: `{"action":"snapshot","arg":{"channel":"trade","instId":"BTCUSDT"},"data":[{"ts":"1700000000000","price":"10","size":"1","side":"buy","tradeId":"a"},{"ts":"1700000000001","price":"11","size":"2","side":"sell","tradeId":"b"}]}`,
			wantTrades: []models.Trade{
				{
					Price:          10,
					Volume:         1,
					Time:           time.Unix(0, 1700000000000*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "a",
				},
				{
					Price:          11,
					Volume:         -2,
					Time:           time.Unix(0, 1700000000001*int64(time.Millisecond)),
					BaseToken:      pair.BaseToken,
					QuoteToken:     pair.QuoteToken,
					ForeignTradeID: "b",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := newBitgetTestScraper(map[string]models.Pair{key: pair})
			var lock sync.RWMutex

			// Collect trades from the unbuffered channel in a separate goroutine.
			collected := make(chan models.Trade, 16)
			done := make(chan struct{})
			go func() {
				bitgetHooks{}.OnMessage(bs, tc.mt, []byte(tc.data), &lock)
				close(done)
			}()
			go func() {
				for {
					select {
					case tr := <-bs.tradesChannel:
						collected <- tr
					case <-done:
						// Drain anything already buffered, then stop.
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

// TestBitgetTickerKeyFromForeign pins the instId-matching behaviour.
func TestBitgetTickerKeyFromForeign(t *testing.T) {
	cases := map[string]string{
		"BTC-USDT": "BTCUSDT",
		"ETH-USDT": "ETHUSDT",
		"BTCUSDT":  "BTCUSDT",
	}
	for in, want := range cases {
		if got := (bitgetHooks{}).TickerKeyFromForeign(in); got != want {
			t.Errorf("TickerKeyFromForeign(%q) = %q, want %q", in, got, want)
		}
	}
}