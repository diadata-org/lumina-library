package scrapers

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/models"
)

func TestParseKrakenTradeMessage(t *testing.T) {
	cases := []struct {
		name    string
		input   krakenWSResponseData
		want    float64
		wantVol float64
		wantTs  string // string for time parsing
		wantID  string
		wantErr bool
	}{
		{
			name: "valid buy trade",
			input: krakenWSResponseData{
				Symbol:    "BTC/USD",
				Side:      "buy",
				Price:     31500.5,
				Size:      0.123,
				OrderType: "market",
				TradeID:   456,
				Time:      "2024-06-26T23:50:55.000000Z",
			},
			want:    31500.5,
			wantVol: 0.123,
			wantTs:  "2024-06-26T23:50:55.000000Z",
			wantID:  strconv.Itoa(456),
			wantErr: false,
		},
		{
			name: "valid sell trade",
			input: krakenWSResponseData{
				Symbol:    "BTC/USD",
				Side:      "sell",
				Price:     31500.5,
				Size:      0.123,
				OrderType: "market",
				TradeID:   456,
				Time:      "2024-06-26T23:50:55.000000Z",
			},
			want:    31500.5,
			wantVol: -0.123,
			wantTs:  "2024-06-26T23:50:55.000000Z",
			wantID:  strconv.Itoa(456),
			wantErr: false,
		},
		{
			name: "invalid timestamp",
			input: krakenWSResponseData{
				Symbol:    "ETH/USD",
				Side:      "buy",
				Price:     2500.0,
				Size:      1.0,
				OrderType: "limit",
				TradeID:   789,
				Time:      "bad-timestamp",
			},
			want:    2500.0,
			wantVol: 1.0,
			wantTs:  "",
			wantID:  "",
			wantErr: true,
		},
	}

	layout := "2006-01-02T15:04:05.000000Z"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			price, volume, timestamp, foreignTradeID, err := parseKrakenTradeMessage(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseKrakenTradeMessage() error = %v, wantErr %v", err, tc.wantErr)
			}
			if price != tc.want {
				t.Errorf("price = %v, want %v", price, tc.want)
			}
			if volume != tc.wantVol {
				t.Errorf("volume = %v, want %v", volume, tc.wantVol)
			}
			if tc.wantTs != "" {
				wantTime, err2 := time.Parse(layout, tc.wantTs)
				if err2 != nil {
					t.Fatalf("bad wantTs in test: %v", err2)
				}
				if !timestamp.Equal(wantTime) {
					t.Errorf("timestamp = %v, want %v", timestamp, wantTime)
				}
			} else if !timestamp.IsZero() {
				t.Errorf("timestamp should be zero, got %v", timestamp)
			}
			if foreignTradeID != tc.wantID {
				t.Errorf("foreignTradeID = %v, want %v", foreignTradeID, tc.wantID)
			}
		})
	}
}

func TestKrakenTradesChannel(t *testing.T) {
	s := &krakenScraper{tradesChannel: make(chan models.Trade)}
	if s.TradesChannel() == nil {
		t.Fatal("expected non-nil channel")
	}
}

func TestKrakenSubscribe(t *testing.T) {
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &krakenScraper{
		wsClient: mockWs,
	}
	pair := models.ExchangePair{
		UnderlyingPair: models.Pair{
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USD"},
		},
	}
	// Test subscribe
	err := s.subscribe(pair, true, lock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 1 {
		t.Errorf("expected WriteJSON to be called once")
	}
	// Test unsubscribe
	err = s.subscribe(pair, false, lock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 2 {
		t.Errorf("expected WriteJSON to be called twice")
	}
}

func TestKrakenResubscribe(t *testing.T) {
	ch := make(chan models.ExchangePair, 1)
	ch <- models.ExchangePair{
		ForeignName: "BTC-USD",
		UnderlyingPair: models.Pair{
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USD"},
		},
	}
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &krakenScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	// You can simulate error on mockWs.WriteJSON if needed
	ctx, cancel := context.WithCancel(context.Background())
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()
}

func TestKrakenClose(t *testing.T) {
	mockWs := &mockWsConn{}
	cancel := func() {}
	s := &krakenScraper{wsClient: mockWs}
	err := s.Close(cancel)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !mockWs.closeCalled {
		t.Errorf("expected Close to be called")
	}
}

func TestKrakenFetchTrades(t *testing.T) {
	mockWs := &mockWsConn{}
	tradesCh := make(chan models.Trade, 1)
	pairMap := map[string]models.Pair{
		"BTCUSD": {
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USD"},
		},
	}

	// Simulate a trade message as received from Kraken.
	mockWs.readJSONQueue = []interface{}{
		krakenWSResponse{
			Channel: "trade",
			Data: []krakenWSResponseData{
				{
					Symbol:    "BTC/USD",
					Side:      "buy",
					Price:     55555.55,
					Size:      1.23,
					TradeID:   987654,
					Time:      time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"),
					OrderType: "market",
				},
			},
		},
	}

	scraper := &krakenScraper{
		wsClient:         mockWs,
		tradesChannel:    tradesCh,
		tickerPairMap:    pairMap,
		lastTradeTimeMap: map[string]time.Time{},
	}

	var lock sync.RWMutex

	go scraper.fetchTrades(&lock)

	select {
	case trade := <-tradesCh:
		if trade.Price != 55555.55 {
			t.Errorf("expected price 55555.55, got %v", trade.Price)
		}
		if trade.Volume != 1.23 {
			t.Errorf("expected volume 1.23, got %v", trade.Volume)
		}
		if trade.QuoteToken.Symbol != "BTC" || trade.BaseToken.Symbol != "USD" {
			t.Errorf("unexpected tokens: %+v, %+v", trade.QuoteToken, trade.BaseToken)
		}
		if trade.ForeignTradeID != "987654" {
			t.Errorf("unexpected ForeignTradeID: %v", trade.ForeignTradeID)
		}
	case <-time.After(time.Second):
		t.Fatal("did not receive trade")
	}
}

func TestKrakenFetchTrades_ReadJSONError(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONErrs: []error{
			errors.New("test error"),
		},
	}
	scraper := &krakenScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{},
		lastTradeTimeMap: map[string]time.Time{},
		maxErrCount:      1,
		restartWaitTime:  0,
	}
	var lock sync.RWMutex
	go scraper.fetchTrades(&lock)
	time.Sleep(10 * time.Millisecond)
}

func TestKrakenFetchTrades_NonTradeChannel(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			krakenWSResponse{
				Channel: "book", // Not "trade"
				Data:    []krakenWSResponseData{},
			},
		},
		readJSONErrs: []error{nil},
	}
	scraper := &krakenScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{},
		lastTradeTimeMap: map[string]time.Time{},
	}
	var lock sync.RWMutex
	go scraper.fetchTrades(&lock)
	time.Sleep(10 * time.Millisecond)
	// No trade should be sent
	select {
	case <-scraper.tradesChannel:
		t.Fatal("should not receive trade for non-trade channel")
	default:
		// pass
	}
}

func TestKrakenFetchTrades_EmptyData(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			krakenWSResponse{
				Channel: "trade",
				Data:    []krakenWSResponseData{},
			},
		},
		readJSONErrs: []error{nil},
	}
	scraper := &krakenScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{},
		lastTradeTimeMap: map[string]time.Time{},
	}
	var lock sync.RWMutex
	go scraper.fetchTrades(&lock)
	time.Sleep(10 * time.Millisecond)
	select {
	case <-scraper.tradesChannel:
		t.Fatal("should not receive trade for empty data")
	default:
		// pass
	}
}

func TestKrakenFetchTrades_ParseTradeError(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			krakenWSResponse{
				Channel: "trade",
				Data: []krakenWSResponseData{
					{
						Symbol:    "BTC/USD",
						Side:      "buy",
						Price:     -1,
						Size:      -1,
						OrderType: "market",
						TradeID:   42,
						Time:      "not-a-time",
					},
				},
			},
		},
		readJSONErrs: []error{nil},
	}
	scraper := &krakenScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{"BTCUSD": {QuoteToken: models.Asset{Symbol: "BTC"}, BaseToken: models.Asset{Symbol: "USD"}}},
		lastTradeTimeMap: map[string]time.Time{},
	}
	var lock sync.RWMutex
	go scraper.fetchTrades(&lock)
	time.Sleep(10 * time.Millisecond)
	select {
	case <-scraper.tradesChannel:
		t.Fatal("should not receive trade for bad timestamp")
	default:
		// pass
	}
}
