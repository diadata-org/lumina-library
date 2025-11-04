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

var dummyPair = models.Pair{
	QuoteToken: models.Asset{Symbol: "USDT", Blockchain: "Ethereum", Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7"},
	BaseToken:  models.Asset{Symbol: "BTC", Blockchain: "Bitcoin", Address: "0x0000000000000000000000000000000000000000"},
}

func TestHandleWSResponse_GateIO(t *testing.T) {
	scraper := &gateIOScraper{
		tickerPairMap: map[string]models.Pair{
			"BTCUSDT": dummyPair,
		},
	}

	cases := []struct {
		name   string
		input  GateIOResponseTrade
		expect models.Trade
	}{
		{
			name: "valid buy trade",
			input: GateIOResponseTrade{
				Result: struct {
					ID           int    `json:"id"`
					CreateTime   int    `json:"create_time"`
					CreateTimeMs string `json:"create_time_ms"`
					Side         string `json:"side"`
					CurrencyPair string `json:"currency_pair"`
					Amount       string `json:"amount"`
					Price        string `json:"price"`
				}{
					ID:           123456,
					CreateTime:   1721923200,
					CreateTimeMs: "1721923200000",
					Side:         "buy",
					CurrencyPair: "BTC_USDT",
					Amount:       "0.123",
					Price:        "42000.99",
				},
			},
			expect: models.Trade{
				QuoteToken:     dummyPair.QuoteToken,
				BaseToken:      dummyPair.BaseToken,
				Price:          42000.99,
				Volume:         0.123,
				Time:           time.Unix(1721923200, 0),
				Exchange:       Exchanges[GATEIO_EXCHANGE],
				ForeignTradeID: strconv.FormatInt(123456, 16),
			},
		},
		{
			name: "valid sell trade",
			input: GateIOResponseTrade{
				Result: struct {
					ID           int    `json:"id"`
					CreateTime   int    `json:"create_time"`
					CreateTimeMs string `json:"create_time_ms"`
					Side         string `json:"side"`
					CurrencyPair string `json:"currency_pair"`
					Amount       string `json:"amount"`
					Price        string `json:"price"`
				}{
					ID:           123456,
					CreateTime:   1721923200,
					CreateTimeMs: "1721923200000",
					Side:         "sell",
					CurrencyPair: "BTC_USDT",
					Amount:       "0.123",
					Price:        "42000.99",
				},
			},
			expect: models.Trade{
				QuoteToken:     dummyPair.QuoteToken,
				BaseToken:      dummyPair.BaseToken,
				Price:          42000.99,
				Volume:         -0.123,
				Time:           time.Unix(1721923200, 0),
				Exchange:       Exchanges[GATEIO_EXCHANGE],
				ForeignTradeID: strconv.FormatInt(123456, 16),
			},
		},
		{
			name: "invalid price (should return zero trade)",
			input: GateIOResponseTrade{
				Result: struct {
					ID           int    `json:"id"`
					CreateTime   int    `json:"create_time"`
					CreateTimeMs string `json:"create_time_ms"`
					Side         string `json:"side"`
					CurrencyPair string `json:"currency_pair"`
					Amount       string `json:"amount"`
					Price        string `json:"price"`
				}{
					ID:           77,
					CreateTime:   1000000000,
					CreateTimeMs: "1000000000000",
					Side:         "buy",
					CurrencyPair: "BTC_USDT",
					Amount:       "0.5",
					Price:        "bad-price",
				},
			},
			expect: models.Trade{}, // zero trade on error
		},
		{
			name: "invalid amount (should return zero trade)",
			input: GateIOResponseTrade{
				Result: struct {
					ID           int    `json:"id"`
					CreateTime   int    `json:"create_time"`
					CreateTimeMs string `json:"create_time_ms"`
					Side         string `json:"side"`
					CurrencyPair string `json:"currency_pair"`
					Amount       string `json:"amount"`
					Price        string `json:"price"`
				}{
					ID:           88,
					CreateTime:   2000000000,
					CreateTimeMs: "2000000000000",
					Side:         "buy",
					CurrencyPair: "BTC_USDT",
					Amount:       "bad-amount",
					Price:        "1000.99",
				},
			},
			expect: models.Trade{}, // zero trade on error
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := scraper.handleWSResponse(tc.input, &sync.RWMutex{})
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
			if got.QuoteToken != tc.expect.QuoteToken {
				t.Errorf("QuoteToken: got %+v, want %+v", got.QuoteToken, tc.expect.QuoteToken)
			}
			if got.BaseToken != tc.expect.BaseToken {
				t.Errorf("BaseToken: got %+v, want %+v", got.BaseToken, tc.expect.BaseToken)
			}
		})
	}
}

func TestGateIOSubscribe(t *testing.T) {
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &gateIOScraper{wsClient: mockWs}
	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	err := s.subscribe(pair, true, lock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 1 {
		t.Errorf("expected WriteJSON to be called once")
	}
	msg, _ := mockWs.writeJSONCalls[0].(*SubscribeGate)
	if msg.Event != "subscribe" {
		t.Errorf("expected Event=subscribe, got %v", msg.Event)
	}

	err = s.subscribe(pair, false, lock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockWs.writeJSONCalls) != 2 {
		t.Errorf("expected WriteJSON to be called twice")
	}
	msg2, _ := mockWs.writeJSONCalls[1].(*SubscribeGate)
	if msg2.Event != "unsubscribe" {
		t.Errorf("expected Event=unsubscribe, got %v", msg2.Event)
	}
}

func TestGateIOResubscribe(t *testing.T) {
	ch := make(chan models.ExchangePair, 1)
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{
		writeJSONErrs: []error{
			errors.New("write error"),
			nil,
		},
	}
	s := &gateIOScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	// Test error path
	ch <- models.ExchangePair{ForeignName: "BTC-USDT"}
	ctx, cancel := context.WithCancel(context.Background())
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Test success path
	ch = make(chan models.ExchangePair, 1)
	mockWs = &mockWsConn{}
	s = &gateIOScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	ch <- models.ExchangePair{ForeignName: "BTC-USDT"}
	ctx, cancel = context.WithCancel(context.Background())
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()
}

func TestGateIOClose(t *testing.T) {
	// Case 1: wsClient is nil
	s := &gateIOScraper{}
	err := s.Close(func() {})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// Case 2: wsClient is set
	mockWs := &mockWsConn{}
	s = &gateIOScraper{wsClient: mockWs}
	called := false
	err = s.Close(func() { called = true })
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !mockWs.closeCalled {
		t.Error("expected Close() to call wsClient.Close()")
	}
	if !called {
		t.Error("expected cancel func to be called")
	}
}

func TestGateIOTradesChannel(t *testing.T) {
	ch := make(chan models.Trade)
	s := &gateIOScraper{tradesChannel: ch}
	if s.TradesChannel() != ch {
		t.Error("TradesChannel() did not return the expected channel")
	}
}

func TestGateIOFetchTrades(t *testing.T) {
	mockWs := &mockWsConn{}
	tradesCh := make(chan models.Trade, 1)
	pairMap := map[string]models.Pair{
		"BTCUSDT": {
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USDT"},
		},
	}

	// Simulate a trade message as received from GateIO.
	mockWs.readJSONQueue = []interface{}{
		GateIOResponseTrade{
			Result: struct {
				ID           int    `json:"id"`
				CreateTime   int    `json:"create_time"`
				CreateTimeMs string `json:"create_time_ms"`
				Side         string `json:"side"`
				CurrencyPair string `json:"currency_pair"`
				Amount       string `json:"amount"`
				Price        string `json:"price"`
			}{
				ID:           1001,
				CreateTime:   int(time.Now().Unix()),
				CreateTimeMs: "",
				Side:         "buy",
				CurrencyPair: "BTC_USDT",
				Amount:       "0.5",
				Price:        "65000",
			},
		},
	}

	scraper := &gateIOScraper{
		wsClient:         mockWs,
		tradesChannel:    tradesCh,
		tickerPairMap:    pairMap,
		lastTradeTimeMap: map[string]time.Time{},
	}

	var lock sync.RWMutex

	go scraper.fetchTrades(&lock)

	select {
	case trade := <-tradesCh:
		if trade.Price != 65000 {
			t.Errorf("expected price 65000, got %v", trade.Price)
		}
		if trade.Volume != 0.5 {
			t.Errorf("expected volume 0.5, got %v", trade.Volume)
		}
		if trade.QuoteToken.Symbol != "BTC" || trade.BaseToken.Symbol != "USDT" {
			t.Errorf("unexpected tokens: %+v, %+v", trade.QuoteToken, trade.BaseToken)
		}
	case <-time.After(time.Second):
		t.Fatal("did not receive trade")
	}
}

func TestGateIOFetchTrades_ErrorPaths(t *testing.T) {
	// Case: ReadJSON always errors, should eventually exit
	mockWs := &mockWsConn{
		readJSONErrs: []error{errors.New("fail1"), errors.New("fail2"), errors.New("fail3")},
	}
	scraper := &gateIOScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{},
		lastTradeTimeMap: map[string]time.Time{},
		maxErrCount:      2, // low threshold for test speed
		restartWaitTime:  0,
	}
	var lock sync.RWMutex
	go scraper.fetchTrades(&lock)
	time.Sleep(50 * time.Millisecond)

	// Case: handleWSResponse returns zero trade
	mockWs = &mockWsConn{
		readJSONQueue: []interface{}{
			GateIOResponseTrade{
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
					CreateTime:   1,
					CreateTimeMs: "",
					Side:         "buy",
					CurrencyPair: "BTC_USDT",
					Amount:       "bad", // invalid
					Price:        "1000",
				},
			},
		},
	}
	scraper = &gateIOScraper{
		wsClient:         mockWs,
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    map[string]models.Pair{"BTCUSDT": dummyPair},
		lastTradeTimeMap: map[string]time.Time{},
	}
	go scraper.fetchTrades(&lock)
	time.Sleep(20 * time.Millisecond)
}
