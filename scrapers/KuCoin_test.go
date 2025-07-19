package scrapers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

func TestParseKuCoinTradeMessage(t *testing.T) {
	cases := []struct {
		name        string
		input       kuCoinWSResponse
		wantPrice   float64
		wantVolume  float64
		wantTime    time.Time
		wantTradeID string
		expectErr   bool
	}{
		{
			name: "valid buy trade",
			input: kuCoinWSResponse{
				Data: kuCoinWSData{
					Price:   "33225.1",
					Size:    "0.23",
					Side:    "buy",
					TradeID: "abc123",
					Time:    "1721923200000", // ms since epoch
				},
			},
			wantPrice:   33225.1,
			wantVolume:  0.23,
			wantTime:    time.Unix(0, 1721923200000),
			wantTradeID: "abc123",
			expectErr:   false,
		},
		{
			name: "valid sell trade",
			input: kuCoinWSResponse{
				Data: kuCoinWSData{
					Price:   "33225.1",
					Size:    "0.23",
					Side:    "sell",
					TradeID: "abc123",
					Time:    "1721923200000", // ms since epoch
				},
			},
			wantPrice:   33225.1,
			wantVolume:  -0.23,
			wantTime:    time.Unix(0, 1721923200000),
			wantTradeID: "abc123",
			expectErr:   false,
		},
		{
			name: "invalid price",
			input: kuCoinWSResponse{
				Data: kuCoinWSData{
					Price:   "bad-price",
					Size:    "0.2",
					Side:    "buy",
					TradeID: "abc456",
					Time:    "1721923200000",
				},
			},
			expectErr: true,
		},
		{
			name: "invalid volume",
			input: kuCoinWSResponse{
				Data: kuCoinWSData{
					Price:   "33225.1",
					Size:    "bad-size",
					Side:    "buy",
					TradeID: "abc789",
					Time:    "1721923200000",
				},
			},
			expectErr: true,
		},
		{
			name: "invalid timestamp",
			input: kuCoinWSResponse{
				Data: kuCoinWSData{
					Price:   "10000.1",
					Size:    "1.0",
					Side:    "buy",
					TradeID: "badts",
					Time:    "bad-timestamp",
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			price, volume, ts, tradeID, err := parseKuCoinTradeMessage(tc.input)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectErr {
				if price != tc.wantPrice {
					t.Errorf("Price: got %v, want %v", price, tc.wantPrice)
				}
				if volume != tc.wantVolume {
					t.Errorf("Volume: got %v, want %v", volume, tc.wantVolume)
				}
				if !ts.Equal(tc.wantTime) {
					t.Errorf("Time: got %v, want %v", ts, tc.wantTime)
				}
				if tradeID != tc.wantTradeID {
					t.Errorf("TradeID: got %v, want %v", tradeID, tc.wantTradeID)
				}
			}
		})
	}
}

func TestKuCoinTradesChannel(t *testing.T) {
	s := &kucoinScraper{tradesChannel: make(chan models.Trade)}
	if s.TradesChannel() == nil {
		t.Error("expected non-nil tradesChannel")
	}
}

func TestKuCoinClose(t *testing.T) {
	mockWs := &mockWsConn{}
	s := &kucoinScraper{wsClient: mockWs}
	called := false
	cancel := func() { called = true }
	err := s.Close(cancel)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !mockWs.closeCalled {
		t.Errorf("expected Close to be called on wsClient")
	}
	if !called {
		t.Errorf("expected cancel func to be called")
	}
}

func TestKuCoinSubscribe(t *testing.T) {
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &kucoinScraper{wsClient: mockWs}
	pair := models.ExchangePair{ForeignName: "ETH-USDT"}
	err := s.subscribe(pair, true, lock)
	if err != nil {
		t.Errorf("unexpected error on subscribe: %v", err)
	}
	err = s.subscribe(pair, false, lock)
	if err != nil {
		t.Errorf("unexpected error on unsubscribe: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 2 {
		t.Errorf("expected 2 WriteJSON calls, got %d", len(mockWs.writeJSONCalls))
	}
}

func TestKuCoinResubscribe(t *testing.T) {
	ch := make(chan models.ExchangePair, 1)
	ch <- models.ExchangePair{ForeignName: "ETH-USDT"}
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &kucoinScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()
}

func TestKuCoinPing(t *testing.T) {
	mockWs := &mockWsConn{}
	s := &kucoinScraper{wsClient: mockWs}
	lock := &sync.RWMutex{}
	ctx, cancel := context.WithCancel(context.Background())
	go s.ping(ctx, 1, time.Now(), lock)
	time.Sleep(20 * time.Millisecond)
	cancel()
}

func TestGetPublicKuCoinToken(t *testing.T) {
	// Serve a fake KuCoin response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code": "200000",
			"data": map[string]interface{}{
				"token": "dummy-token",
				"instanceServers": []map[string]interface{}{
					{"pingInterval": int64(10)},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	token, ping, err := getPublicKuCoinToken(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "dummy-token" {
		t.Errorf("unexpected token: %v", token)
	}
	if ping != 10 {
		t.Errorf("unexpected pingInterval: %v", ping)
	}
}

func TestKucoinHandleWSResponse(t *testing.T) {
	pair := models.Pair{
		QuoteToken: models.Asset{Symbol: "ETH"},
		BaseToken:  models.Asset{Symbol: "USDT"},
	}
	s := &kucoinScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{"ETHUSDT": pair},
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	msg := kuCoinWSResponse{
		Data: kuCoinWSData{
			Symbol:  "ETH-USDT",
			Side:    "buy",
			Price:   "4200.99",
			Size:    "0.12",
			TradeID: "123abc",
			Time:    "1721923200000",
		},
	}

	s.handleWSResponse(msg, &lock)

	select {
	case trade := <-s.tradesChannel:
		if trade.Price != 4200.99 {
			t.Errorf("expected price 4200.99, got %v", trade.Price)
		}
		if trade.Volume != 0.12 {
			t.Errorf("expected volume 0.12, got %v", trade.Volume)
		}
		if trade.ForeignTradeID != "123abc" {
			t.Errorf("expected ForeignTradeID '123abc', got %v", trade.ForeignTradeID)
		}
		if trade.QuoteToken.Symbol != "ETH" || trade.BaseToken.Symbol != "USDT" {
			t.Errorf("unexpected tokens: %v, %v", trade.QuoteToken, trade.BaseToken)
		}
	case <-time.After(time.Second):
		t.Fatal("expected trade, got none")
	}
}

func TestKucoinHandleWSResponse_ParseError(t *testing.T) {
	pair := models.Pair{
		QuoteToken: models.Asset{Symbol: "ETH"},
		BaseToken:  models.Asset{Symbol: "USDT"},
	}
	s := &kucoinScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{"ETHUSDT": pair},
		lastTradeTimeMap: make(map[string]time.Time),
	}
	var lock sync.RWMutex

	msg := kuCoinWSResponse{
		Data: kuCoinWSData{
			Symbol:  "ETH-USDT",
			Side:    "buy",
			Price:   "badprice",
			Size:    "0.12",
			TradeID: "123abc",
			Time:    "1721923200000",
		},
	}

	s.handleWSResponse(msg, &lock)

	select {
	case trade := <-s.tradesChannel:
		t.Fatalf("expected no trade, got %+v", trade)
	default:
	}
}

func TestKucoinHandleWSResponse_UnmappedPair(t *testing.T) {
	scraper := &kucoinScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{}, // empty, to force missing mapping
		lastTradeTimeMap: map[string]time.Time{},
	}
	var lock sync.RWMutex
	msg := kuCoinWSResponse{
		Type: "message",
		Data: kuCoinWSData{
			Symbol:  "FOO-BAR",
			Side:    "buy",
			Price:   "100.0",
			Size:    "1.0",
			TradeID: "999",
			Time:    strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		},
	}
	go scraper.handleWSResponse(msg, &lock)
	select {
	case trade := <-scraper.tradesChannel:
		// Trade tokens should be zero-valued
		if trade.QuoteToken.Symbol != "" || trade.BaseToken.Symbol != "" {
			t.Errorf("expected zero-value tokens for missing mapping, got %+v", trade)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("did not receive trade")
	}
}

func TestKucoinHandleWSResponse_InvalidPair(t *testing.T) {
	s := &kucoinScraper{
		tradesChannel: make(chan models.Trade, 1),
		tickerPairMap: map[string]models.Pair{},
	}
	var lock sync.RWMutex

	msg := kuCoinWSResponse{
		Data: kuCoinWSData{
			Symbol:  "BTCUSDT", // No dash
			Price:   "100.1",
			Size:    "1",
			Side:    "buy",
			TradeID: "foo",
			Time:    "1721923200000",
		},
	}

	s.handleWSResponse(msg, &lock)

	select {
	case trade := <-s.tradesChannel:
		t.Fatalf("expected no trade, got %+v", trade)
	default:
	}
}

func TestKucoinFetchTrades(t *testing.T) {
	mockWs := &mockWsConn{}
	tradesCh := make(chan models.Trade, 1)
	pairMap := map[string]models.Pair{
		"ETHUSDT": {
			QuoteToken: models.Asset{Symbol: "ETH"},
			BaseToken:  models.Asset{Symbol: "USDT"},
		},
	}

	nowMs := int(time.Now().UnixNano() / 1e6)

	mockWs.readJSONQueue = []interface{}{
		kuCoinWSResponse{
			Type: "pong",
		},
		kuCoinWSResponse{
			Type: "message",
			Data: kuCoinWSData{
				Symbol:  "ETH-USDT",
				Side:    "buy",
				Price:   "3620.42",
				Size:    "0.19",
				TradeID: "12345",
				Time:    strconv.Itoa(nowMs),
			},
		},
	}

	scraper := &kucoinScraper{
		wsClient:         mockWs,
		tradesChannel:    tradesCh,
		tickerPairMap:    pairMap,
		lastTradeTimeMap: map[string]time.Time{},
	}

	var lock sync.RWMutex

	go scraper.fetchTrades(&lock)

	select {
	case trade := <-tradesCh:
		if trade.Price != 3620.42 {
			t.Errorf("expected price 3620.42, got %v", trade.Price)
		}
		if trade.Volume != 0.19 {
			t.Errorf("expected volume 0.19, got %v", trade.Volume)
		}
		if trade.QuoteToken.Symbol != "ETH" || trade.BaseToken.Symbol != "USDT" {
			t.Errorf("unexpected tokens: %+v, %+v", trade.QuoteToken, trade.BaseToken)
		}
		if trade.ForeignTradeID != "12345" {
			t.Errorf("unexpected ForeignTradeID: %v", trade.ForeignTradeID)
		}
	case <-time.After(time.Second):
		t.Fatal("did not receive trade")
	}
}

func TestKucoinFetchTrades_HandleError(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{},
		readJSONErrs:  []error{errors.New("test error"), errors.New("test error"), errors.New("test error")},
	}
	scraper := &kucoinScraper{
		wsClient:      mockWs,
		tradesChannel: make(chan models.Trade, 1),
	}
	var lock sync.RWMutex

	go scraper.fetchTrades(&lock)

	time.Sleep(50 * time.Millisecond)
}
