package scrapers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

func TestCryptodotcomParseTradeMessage(t *testing.T) {
	cases := []struct {
		name   string
		input  cryptodotcomWSResponse
		expect []models.Trade
	}{
		{
			name: "single valid trade",
			input: cryptodotcomWSResponse{
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:   "abc123",
							Timestamp: 1721923200000,
							Price:     "12345.67",
							Volume:    "0.01",
							Side:      "BUY",
						},
					},
				},
			},
			expect: []models.Trade{
				{
					Price:          12345.67,
					Volume:         0.01,
					Time:           time.Unix(0, 1721923200000*1e6),
					Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
					ForeignTradeID: "abc123",
				},
			},
		},
		{
			name: "single valid sell trade",
			input: cryptodotcomWSResponse{
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:   "abc123",
							Timestamp: 1721923200000,
							Price:     "12345.67",
							Volume:    "0.01",
							Side:      "SELL",
						},
					},
				},
			},
			expect: []models.Trade{
				{
					Price:          12345.67,
					Volume:         -0.01,
					Time:           time.Unix(0, 1721923200000*1e6),
					Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
					ForeignTradeID: "abc123",
				},
			},
		},
		{
			name: "invalid price (should return no trades)",
			input: cryptodotcomWSResponse{
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:   "badprice",
							Timestamp: 1721923200000,
							Price:     "not-a-number",
							Volume:    "0.01",
							Side:      "BUY",
						},
					},
				},
			},
			expect: nil,
		},
		{
			name: "invalid volume (should return no trades)",
			input: cryptodotcomWSResponse{
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:   "badvol",
							Timestamp: 1721923200000,
							Price:     "222.2",
							Volume:    "not-a-number",
							Side:      "BUY",
						},
					},
				},
			},
			expect: nil,
		},
		{
			name: "multiple valid trades",
			input: cryptodotcomWSResponse{
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:   "one",
							Timestamp: 1000000000000,
							Price:     "10",
							Volume:    "1",
							Side:      "BUY",
						},
						{
							TradeID:   "two",
							Timestamp: 2000000000000,
							Price:     "20",
							Volume:    "2",
							Side:      "BUY",
						},
					},
				},
			},
			expect: []models.Trade{
				{
					Price:          10,
					Volume:         1,
					Time:           time.Unix(0, 1000000000000*1e6),
					Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
					ForeignTradeID: "one",
				},
				{
					Price:          20,
					Volume:         2,
					Time:           time.Unix(0, 2000000000000*1e6),
					Exchange:       Exchanges[CRYPTODOTCOM_EXCHANGE],
					ForeignTradeID: "two",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := cryptodotcomParseTradeMessage(tc.input)
			if len(got) != len(tc.expect) {
				t.Fatalf("expected %d trades, got %d", len(tc.expect), len(got))
			}
			for i := range got {
				if got[i].Exchange != tc.expect[i].Exchange {
					t.Errorf("Exchange: got %v, want %v", got[i].Exchange, tc.expect[i].Exchange)
				}
				if !got[i].Time.Equal(tc.expect[i].Time) {
					t.Errorf("Time: got %v, want %v", got[i].Time, tc.expect[i].Time)
				}
				if got[i].Price != tc.expect[i].Price {
					t.Errorf("Price: got %v, want %v", got[i].Price, tc.expect[i].Price)
				}
				if got[i].Volume != tc.expect[i].Volume {
					t.Errorf("Volume: got %v, want %v", got[i].Volume, tc.expect[i].Volume)
				}
				if got[i].ForeignTradeID != tc.expect[i].ForeignTradeID {
					t.Errorf("ForeignTradeID: got %v, want %v", got[i].ForeignTradeID, tc.expect[i].ForeignTradeID)
				}
			}
		})
	}
}

func TestCryptodotcomFetchTrades(t *testing.T) {
	// Prepare mock websocket connection
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			cryptodotcomWSResponse{
				Method: "not_heartbeat",
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:     "123",
							Timestamp:   time.Now().UnixMilli(),
							Price:       "1234.56",
							Volume:      "10.5",
							Side:        "BUY",
							ForeignName: "BTC_USDT",
						},
					},
				},
			},
		},
	}

	// Setup tickerPairMap to match ForeignName
	tickerPairMap := map[string]models.Pair{
		"BTCUSDT": {
			QuoteToken: models.Asset{Symbol: "BTC"},
			BaseToken:  models.Asset{Symbol: "USDT"},
		},
	}

	tradesCh := make(chan models.Trade, 1)
	lock := &sync.RWMutex{}
	scraper := &cryptodotcomScraper{
		wsClient:            mockWs,
		tradesChannel:       tradesCh,
		subscribeChannel:    make(chan models.ExchangePair),
		tickerPairMap:       tickerPairMap,
		lastTradeTimeMap:    make(map[string]time.Time),
		maxErrCount:         2,
		restartWaitTime:     0,
		tradeTimeoutSeconds: 10000, // ensure trade not filtered by age
	}

	// Run fetchTrades in goroutine
	go scraper.fetchTrades(lock)

	// Wait for result or timeout
	select {
	case trade := <-tradesCh:
		if trade.Price != 1234.56 {
			t.Errorf("unexpected trade price: got %v", trade.Price)
		}
		if trade.Volume != 10.5 {
			t.Errorf("unexpected trade volume: got %v", trade.Volume)
		}
		if trade.ForeignTradeID != "123" {
			t.Errorf("unexpected trade ForeignTradeID: got %v", trade.ForeignTradeID)
		}
		if trade.QuoteToken.Symbol != "BTC" || trade.BaseToken.Symbol != "USDT" {
			t.Errorf("unexpected tokens: %v, %v", trade.QuoteToken, trade.BaseToken)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for trade")
	}
}

func TestCryptoDotComNewAndCloseAndChannel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan string, 1)
	NewCryptodotcomScraper(ctx, nil, ch, &wg)
	// TradesChannel/Close (with nil wsClient)
	s := &cryptodotcomScraper{}
	s.TradesChannel()
	s.Close(func() {})
}

func TestCryptoDotComSubscribe_SendHeartbeat(t *testing.T) {
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &cryptodotcomScraper{wsClient: mockWs}

	pair := models.ExchangePair{ForeignName: "BTC-USDT"}
	err := s.subscribe(pair, true, lock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err = s.sendHeartbeat(42, lock)
	if err != nil {
		t.Errorf("unexpected error from sendHeartbeat: %v", err)
	}
}

func TestCryptoDotComResubscribe(t *testing.T) {
	ch := make(chan models.ExchangePair, 1)
	ch <- models.ExchangePair{ForeignName: "BTC-USDT"}
	lock := &sync.RWMutex{}
	mockWs := &mockWsConn{}
	s := &cryptodotcomScraper{
		wsClient:         mockWs,
		subscribeChannel: ch,
	}
	ctx, cancel := context.WithCancel(context.Background())

	// Cover normal path (subscribe/unsubscribe)
	go s.resubscribe(ctx, lock)
	time.Sleep(10 * time.Millisecond)
	cancel()
	// Let ctx.Done() branch be taken
	time.Sleep(10 * time.Millisecond)

	// Cover error path (simulate WriteJSON error)
	ch2 := make(chan models.ExchangePair, 1)
	ch2 <- models.ExchangePair{ForeignName: "BTC-USDT"}
	mockWs2 := &mockWsConn{
		writeJSONErrs: []error{
			errors.New("write error"),
			nil,
		},
	}
	s2 := &cryptodotcomScraper{
		wsClient:         mockWs2,
		subscribeChannel: ch2,
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	go s2.resubscribe(ctx2, lock)
	time.Sleep(10 * time.Millisecond)
	cancel2()
}

func TestCryptoDotComHandleWSResponse_ErrorPath(t *testing.T) {
	s := &cryptodotcomScraper{
		tradesChannel:       make(chan models.Trade, 2),
		tickerPairMap:       map[string]models.Pair{"BTCUSDT": {}},
		lastTradeTimeMap:    make(map[string]time.Time),
		tradeTimeoutSeconds: 1,
	}
	lock := &sync.RWMutex{}

	// Error 1: Bad Message
	badMsg := cryptodotcomWSResponse{
		Result: cryptodotcomWSResponseResult{
			Data: []cryptodotcomWSResponseData{
				{
					TradeID:     "bad",
					Timestamp:   time.Now().UnixMilli(),
					Price:       "notafloat",
					Volume:      "0.2",
					ForeignName: "BTC_USDT",
				},
			},
		},
	}
	s.handleWSResponse(badMsg, lock)
	select {
	case tr := <-s.tradesChannel:
		t.Errorf("did not expect trade for parse error, got %v", tr)
	default:
		// success
	}

	// Error 2: Old Trade
	msg := cryptodotcomWSResponse{
		Result: cryptodotcomWSResponseResult{
			Data: []cryptodotcomWSResponseData{
				{
					TradeID:     "old",
					Timestamp:   time.Now().Add(-2*time.Second).UnixNano() / 1e6, // too old
					Price:       "10",
					Volume:      "1",
					ForeignName: "BTC_USDT",
				},
			},
		},
	}
	s.handleWSResponse(msg, lock)
	select {
	case tr := <-s.tradesChannel:
		t.Errorf("did not expect trade for too old, got %v", tr)
	default:
		// success
	}
}

func TestCryptoDotComFetchTrades(t *testing.T) {
	mockWs := &mockWsConn{
		readJSONQueue: []interface{}{
			cryptodotcomWSResponse{ // heartbeat
				ID:     9,
				Method: "public/heartbeat",
			},
			cryptodotcomWSResponse{ // valid trade
				Method: "not_heartbeat",
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:     "999",
							Timestamp:   time.Now().UnixMilli(),
							Price:       "50",
							Volume:      "0.01",
							Side:        "BUY",
							ForeignName: "BTC_USDT",
						},
					},
				},
			},
			cryptodotcomWSResponse{ // invalid trade (parse error)
				Method: "not_heartbeat",
				Result: cryptodotcomWSResponseResult{
					Data: []cryptodotcomWSResponseData{
						{
							TradeID:     "bad",
							Timestamp:   time.Now().UnixMilli(),
							Price:       "nope",
							Volume:      "nope",
							Side:        "BUY",
							ForeignName: "BTC_USDT",
						},
					},
				},
			},
		},
		readJSONErrs: []error{
			nil, nil, nil,
		},
	}
	tradesCh := make(chan models.Trade, 1)
	lock := &sync.RWMutex{}
	scraper := &cryptodotcomScraper{
		wsClient:            mockWs,
		tradesChannel:       tradesCh,
		subscribeChannel:    make(chan models.ExchangePair),
		tickerPairMap:       map[string]models.Pair{"BTCUSDT": {}},
		lastTradeTimeMap:    make(map[string]time.Time),
		tradeTimeoutSeconds: 9999,
	}
	go scraper.fetchTrades(lock)

	// Heartbeat should trigger sendHeartbeat (not produce trade)
	time.Sleep(10 * time.Millisecond) // let goroutine run
	// Should get the valid trade
	select {
	case trade := <-tradesCh:
		if trade.Price != 50 {
			t.Errorf("unexpected trade: %v", trade)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for trade")
	}
}
