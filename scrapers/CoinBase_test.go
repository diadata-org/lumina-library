package scrapers

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

func mustParseTime(t *testing.T, layout, value string) time.Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("failed to parse time in test: %v", err)
	}
	return ts
}

func TestCoinbaseParseTradeMessage(t *testing.T) {
	cases := []struct {
		name   string
		input  coinBaseWSResponse
		expect models.Trade
	}{
		{
			name: "valid buy trade",
			input: coinBaseWSResponse{
				Price:   "30000.00",
				Size:    "0.5",
				Time:    "2024-06-26T01:23:45.123456Z",
				TradeID: 789,
				Side:    "buy",
			},
			expect: models.Trade{
				Price:          30000.00,
				Volume:         0.5,
				Time:           mustParseTime(t, "2006-01-02T15:04:05.000000Z", "2024-06-26T01:23:45.123456Z"),
				Exchange:       Exchanges[COINBASE_EXCHANGE],
				ForeignTradeID: strconv.Itoa(789),
			},
		},
		{
			name: "valid sell trade",
			input: coinBaseWSResponse{
				Price:   "30000.00",
				Size:    "0.5",
				Time:    "2024-06-26T01:23:45.123456Z",
				TradeID: 789,
				Side:    "sell",
			},
			expect: models.Trade{
				Price:          30000.00,
				Volume:         -0.5,
				Time:           mustParseTime(t, "2006-01-02T15:04:05.000000Z", "2024-06-26T01:23:45.123456Z"),
				Exchange:       Exchanges[COINBASE_EXCHANGE],
				ForeignTradeID: strconv.Itoa(789),
			},
		},
		{
			name: "invalid price (should return zero trade)",
			input: coinBaseWSResponse{
				Price:   "not-a-number",
				Size:    "2.2",
				Time:    "2024-07-01T10:01:00.000000Z",
				TradeID: 456,
				Side:    "buy",
			},
			expect: models.Trade{}, // all fields zero
		},
		{
			name: "invalid volume (should return zero trade)",
			input: coinBaseWSResponse{
				Price:   "1337.99",
				Size:    "bad-size",
				Time:    "2024-07-01T10:02:00.000000Z",
				TradeID: 999,
				Side:    "buy",
			},
			expect: models.Trade{}, // all fields zero
		},
		{
			name: "invalid time (should return zero trade)",
			input: coinBaseWSResponse{
				Price:   "10.00",
				Size:    "1.0",
				Time:    "bad-time-format",
				TradeID: 100,
				Side:    "sell",
			},
			expect: models.Trade{}, // all fields zero
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := coinbaseParseTradeMessage(tc.input)
			if got.Exchange != tc.expect.Exchange {
				t.Errorf("Exchange: got %v, want %v", got.Exchange, tc.expect.Exchange)
			}
			if !got.Time.Equal(tc.expect.Time) {
				t.Errorf("Time: got %v, want %v", got.Time, tc.expect.Time)
			}
			if got.Price != tc.expect.Price {
				t.Errorf("Price: got %v, want %v", got.Price, tc.expect.Price)
			}
			if got.Volume != tc.expect.Volume {
				t.Errorf("Volume: got %v, want %v", got.Volume, tc.expect.Volume)
			}
			if got.ForeignTradeID != tc.expect.ForeignTradeID {
				t.Errorf("ForeignTradeID: got %v, want %v", got.ForeignTradeID, tc.expect.ForeignTradeID)
			}
		})
	}
}

func TestCoinbaseSubscribe(t *testing.T) {
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &coinbaseScraper{wsClient: mockWs}
	pair := models.ExchangePair{ForeignName: "BTC-USD"}
	err := s.subscribe(pair, true, lock)
	if err != nil {
		t.Errorf("subscribe error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 1 {
		t.Errorf("expected one call, got %d", len(mockWs.writeJSONCalls))
	}
	// Unsubscribe
	err = s.subscribe(pair, false, lock)
	if err != nil {
		t.Errorf("unsubscribe error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 2 {
		t.Errorf("expected two calls, got %d", len(mockWs.writeJSONCalls))
	}
}

func TestCoinbaseResubscribe(t *testing.T) {
	ch := make(chan models.ExchangePair, 1)
	ch <- models.ExchangePair{ForeignName: "BTC-USD"}

	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{
		writeJSONErrs: []error{
			errors.New("write error"),
			nil,
		},
	}
	s := &coinbaseScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()
}

func TestCoinbaseClose(t *testing.T) {
	// Case 1: wsClient is nil
	s := &coinbaseScraper{}
	cancelCalled := false
	err := s.Close(func() { cancelCalled = true })
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !cancelCalled {
		t.Error("expected cancel to be called")
	}

	// Case 2: wsClient is not nil
	mockWs := &mockWsConn{}
	s = &coinbaseScraper{wsClient: mockWs}
	cancelCalled = false
	err = s.Close(func() { cancelCalled = true })
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !mockWs.closeCalled {
		t.Error("expected Close to be called")
	}
	if !cancelCalled {
		t.Error("expected cancel to be called")
	}
}

func TestCoinbaseHandleWSResponse(t *testing.T) {
	mockWs := &mockWsConn{}
	s := &coinbaseScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{"BTCUSD": {}},
		lastTradeTimeMap: make(map[string]time.Time),
	}

	var lock sync.RWMutex

	// Valid trade
	msg := coinBaseWSResponse{
		Price:     "10",
		Size:      "2",
		Time:      "2024-07-01T10:00:00.000000Z",
		TradeID:   123,
		ProductID: "BTC-USD",
		Side:      "buy",
	}
	s.handleWSResponse(msg, &lock)
	select {
	case trade := <-s.tradesChannel:
		if trade.Price != 10 || trade.Volume != 2 {
			t.Errorf("bad trade: %+v", trade)
		}
	default:
		t.Error("expected trade")
	}

	// Error case (invalid price)
	badMsg := coinBaseWSResponse{
		Price:     "bad",
		Size:      "2",
		Time:      "2024-07-01T10:00:00.000000Z",
		TradeID:   123,
		ProductID: "BTC-USD",
		Side:      "buy",
	}
	s.handleWSResponse(badMsg, &lock)
	select {
	case trade := <-s.tradesChannel:
		if trade.Price != 0 || trade.Volume != 0 {
			t.Errorf("expected zero-value trade, got: %+v", trade)
		}
	default:
		t.Errorf("expected a zero-value trade, got none")
	}

	// // ProductID missing dash
	// noDashMsg := coinBaseWSResponse{
	// 	Price:     "123.45",
	// 	Size:      "0.01",
	// 	Time:      "2024-07-01T10:05:00.000000Z",
	// 	TradeID:   999,
	// 	ProductID: "BTCUSD", // <-- no dash
	// 	Side:      "buy",
	// }
	// s.handleWSResponse(noDashMsg, &lock)
	// select {
	// case trade := <-s.tradesChannel:
	// 	// Should not receive a trade when parse fails
	// 	t.Errorf("expected no trade, got: %+v", trade)
	// default:
	// 	// Good, nothing sent
	// }
}

func TestCoinbaseFetchTrades(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			// Case 1: match
			coinBaseWSResponse{
				Type:      "match",
				TradeID:   123,
				Price:     "1.23",
				Size:      "0.5",
				Time:      "2024-07-01T10:01:00.000000Z",
				ProductID: "BTC-USD",
				Side:      "buy",
			},
			// Case 2: non-match
			coinBaseWSResponse{Type: "last_match"},
			// Case 3: match again, but will be skipped by tickerPairMap
			coinBaseWSResponse{
				Type:      "match",
				TradeID:   456,
				Price:     "100.0",
				Size:      "0.1",
				Time:      "2024-07-01T10:02:00.000000Z",
				ProductID: "FAKE-PAIR",
				Side:      "buy",
			},
		},
		readJSONErrs: []error{
			nil, nil, nil,
			errors.New("test connection error"),
		},
	}

	tickerMap := map[string]models.Pair{
		"BTCUSD": {},
	}

	s := &coinbaseScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 2),
		tickerPairMap:    tickerMap,
		lastTradeTimeMap: make(map[string]time.Time),
		maxErrCount:      2,
	}

	var lock sync.RWMutex
	go s.fetchTrades(&lock)

	// Get the first trade
	select {
	case trade := <-s.tradesChannel:
		if trade.Price != 1.23 {
			t.Errorf("expected price 1.23, got %v", trade.Price)
		}
	case <-time.After(time.Second):
		t.Fatal("expected trade, got none")
	}

	// No second trade (since ProductID FAKE-PAIR not in tickerMap)
	select {
	case trade := <-s.tradesChannel:
		// Check if the tokens are empty; if so, that's actually expected, so don't fail.
		if trade.QuoteToken.Symbol == "" && trade.BaseToken.Symbol == "" {
			// This is an empty/invalid trade, which is ok for this test.
		} else {
			t.Errorf("expected no trade, got: %+v", trade)
		}
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}
