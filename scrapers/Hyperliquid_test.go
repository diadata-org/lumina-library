package scrapers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

const hyperliquidSpotMetaFixture = `{
	"tokens": [
		{"name": "USDC", "index": 0},
		{"name": "PURR", "index": 1},
		{"name": "HYPE", "index": 150},
		{"name": "UBTC", "index": 197}
	],
	"universe": [
		{"name": "PURR/USDC", "tokens": [1, 0], "index": 0},
		{"name": "@107", "tokens": [150, 0], "index": 107},
		{"name": "@142", "tokens": [197, 0], "index": 142}
	]
}`

func collectHyperliquidTrades(h *hyperliquidHooks, bs *BaseCEXScraper, mt int, data string) []models.Trade {
	var lock sync.RWMutex
	done := make(chan struct{})
	go func() {
		h.OnMessage(bs, mt, []byte(data), &lock)
		close(done)
	}()

	var got []models.Trade
	for {
		select {
		case tr := <-bs.tradesChannel:
			got = append(got, tr)
		case <-done:
			return got
		}
	}
}

func TestHyperliquidOnMessage(t *testing.T) {
	pair := models.Pair{
		QuoteToken: models.Asset{Symbol: "BTC"},
		BaseToken:  models.Asset{Symbol: "USDC"},
	}

	tests := []struct {
		name       string
		mt         int
		data       string
		wantTrades []models.Trade
	}{
		{
			name: "non-text message ignored",
			mt:   ws.BinaryMessage,
			data: `{"channel":"trades","data":[{"coin":"@142","side":"B","px":"1","sz":"1","time":1700000000000,"tid":1}]}`,
		},
		{
			name: "pong ignored",
			mt:   ws.TextMessage,
			data: `{"channel":"pong"}`,
		},
		{
			name: "subscription ack ignored",
			mt:   ws.TextMessage,
			data: `{"channel":"subscriptionResponse","data":{"method":"subscribe","subscription":{"type":"trades","coin":"@142"}}}`,
		},
		{
			name: "invalid json ignored",
			mt:   ws.TextMessage,
			data: `{not-json`,
		},
		{
			name: "unknown coin ignored",
			mt:   ws.TextMessage,
			data: `{"channel":"trades","data":[{"coin":"@151","side":"B","px":"2700","sz":"1","time":1700000000000,"tid":1}]}`,
		},
		{
			name: "buy keeps positive volume",
			mt:   ws.TextMessage,
			data: `{"channel":"trades","data":[{"coin":"@142","side":"B","px":"85746.5","sz":"0.25","time":1700000000000,"tid":11}]}`,
			wantTrades: []models.Trade{
				{Price: 85746.5, Volume: 0.25, Time: time.UnixMilli(1700000000000), ForeignTradeID: "11"},
			},
		},
		{
			name: "sell negates volume",
			mt:   ws.TextMessage,
			data: `{"channel":"trades","data":[{"coin":"@142","side":"A","px":"85746.5","sz":"0.25","time":1700000000000,"tid":12}]}`,
			wantTrades: []models.Trade{
				{Price: 85746.5, Volume: -0.25, Time: time.UnixMilli(1700000000000), ForeignTradeID: "12"},
			},
		},
		{
			name: "unparseable price skips that trade",
			mt:   ws.TextMessage,
			data: `{"channel":"trades","data":[{"coin":"@142","side":"B","px":"abc","sz":"1","time":1700000000000,"tid":13},{"coin":"@142","side":"B","px":"85000","sz":"2","time":1700000000001,"tid":14}]}`,
			wantTrades: []models.Trade{
				{Price: 85000, Volume: 2, Time: time.UnixMilli(1700000000001), ForeignTradeID: "14"},
			},
		},
		{
			name: "multiple trades in one frame all emitted",
			mt:   ws.TextMessage,
			data: `{"channel":"trades","data":[{"coin":"@142","side":"B","px":"10","sz":"1","time":1700000000000,"tid":15},{"coin":"@142","side":"A","px":"11","sz":"2","time":1700000000001,"tid":16}]}`,
			wantTrades: []models.Trade{
				{Price: 10, Volume: 1, Time: time.UnixMilli(1700000000000), ForeignTradeID: "15"},
				{Price: 11, Volume: -2, Time: time.UnixMilli(1700000000001), ForeignTradeID: "16"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &hyperliquidHooks{
				coins:         make(map[string]string),
				foreignByCoin: map[string]string{"@142": "BTC-USDC"},
			}
			bs := &BaseCEXScraper{
				hooks:            h,
				tradesChannel:    make(chan models.Trade),
				tickerPairMap:    models.MakeTickerPairMap([]models.ExchangePair{{ForeignName: "BTC-USDC", UnderlyingPair: pair}}),
				lastTradeTimeMap: make(map[string]time.Time),
			}

			got := collectHyperliquidTrades(h, bs, tc.mt, tc.data)

			if len(got) != len(tc.wantTrades) {
				t.Fatalf("trade count = %d, want %d (got %+v)", len(got), len(tc.wantTrades), got)
			}
			for i, want := range tc.wantTrades {
				g := got[i]
				if g.Price != want.Price || g.Volume != want.Volume || !g.Time.Equal(want.Time) || g.ForeignTradeID != want.ForeignTradeID {
					t.Errorf("trade[%d] = %v %v %v %s, want %v %v %v %s", i, g.Price, g.Volume, g.Time, g.ForeignTradeID, want.Price, want.Volume, want.Time, want.ForeignTradeID)
				}
				if g.QuoteToken.Symbol != "BTC" || g.BaseToken.Symbol != "USDC" {
					t.Errorf("trade[%d] pair = %s-%s, want BTC-USDC", i, g.QuoteToken.Symbol, g.BaseToken.Symbol)
				}
				if g.Exchange.Name != HYPERLIQUID_EXCHANGE {
					t.Errorf("trade[%d].Exchange = %q, want %q", i, g.Exchange.Name, HYPERLIQUID_EXCHANGE)
				}
			}
			if len(got) > 0 && !bs.lastTradeTimeMap["BTC-USDC"].Equal(got[len(got)-1].Time) {
				t.Errorf("lastTradeTime = %v, want %v", bs.lastTradeTimeMap["BTC-USDC"], got[len(got)-1].Time)
			}
		})
	}
}

func TestHyperliquidSubscribe(t *testing.T) {
	var metaRequests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metaRequests.Add(1)
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req["type"] != "spotMeta" {
			t.Errorf("unexpected info request %v: %v", req, err)
		}
		w.Write([]byte(hyperliquidSpotMetaFixture))
	}))
	defer srv.Close()

	h := &hyperliquidHooks{
		infoURL:       srv.URL,
		coins:         make(map[string]string),
		foreignByCoin: make(map[string]string),
	}
	fc := &fakeWSConn{}
	bs := &BaseCEXScraper{wsClient: fc}
	var lock sync.RWMutex

	tests := []struct {
		foreignName string
		subscribe   bool
		wantMethod  string
		wantCoin    string
		wantErr     bool
	}{
		{foreignName: "BTC-USDC", subscribe: true, wantMethod: "subscribe", wantCoin: "@142"},
		{foreignName: "PURR-USDC", subscribe: true, wantMethod: "subscribe", wantCoin: "PURR/USDC"},
		{foreignName: "HYPE-USDC", subscribe: false, wantMethod: "unsubscribe", wantCoin: "@107"},
		{foreignName: "ETH-USDC", subscribe: true, wantErr: true},
		{foreignName: "BTCUSDC", subscribe: true, wantErr: true},
	}

	for _, tc := range tests {
		writes := fc.writeJSONCount
		err := h.Subscribe(bs, models.ExchangePair{ForeignName: tc.foreignName}, tc.subscribe, &lock)
		if tc.wantErr {
			if err == nil {
				t.Errorf("%s: expected error", tc.foreignName)
			}
			if fc.writeJSONCount != writes {
				t.Errorf("%s: unresolved pair must not be sent", tc.foreignName)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: %v", tc.foreignName, err)
		}
		req, ok := fc.lastWritten.(hyperliquidWSRequest)
		if !ok || req.Subscription == nil {
			t.Fatalf("%s: written %T %+v", tc.foreignName, fc.lastWritten, fc.lastWritten)
		}
		if req.Method != tc.wantMethod || req.Subscription.Type != "trades" || req.Subscription.Coin != tc.wantCoin {
			t.Errorf("%s: sent %s %s %s, want %s trades %s", tc.foreignName, req.Method, req.Subscription.Type, req.Subscription.Coin, tc.wantMethod, tc.wantCoin)
		}
		if h.foreignByCoin[tc.wantCoin] != tc.foreignName {
			t.Errorf("foreignByCoin[%s] = %q, want %q", tc.wantCoin, h.foreignByCoin[tc.wantCoin], tc.foreignName)
		}
	}

	// BTC-USDC loads spotMeta, PURR and HYPE hit the cache, ETH-USDC reloads it once.
	if n := metaRequests.Load(); n != 2 {
		t.Errorf("spotMeta requests = %d, want 2", n)
	}
}
