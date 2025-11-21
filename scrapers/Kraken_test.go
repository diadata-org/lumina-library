package scrapers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

// --------- Fake wsConn implementation for Kraken tests ---------

// fakeWSConnKraken is a fake implementation of wsConn.
// It records WriteJSON calls so we can assert on subscription messages.
type fakeWSConnKraken struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConnKraken) ReadMessage() (int, []byte, error) {
	// Not needed for these unit tests.
	return 0, nil, nil
}

func (f *fakeWSConnKraken) WriteMessage(messageType int, data []byte) error {
	// Not needed for these unit tests.
	return nil
}

func (f *fakeWSConnKraken) ReadJSON(v interface{}) error {
	// Not needed for these unit tests.
	return nil
}

func (f *fakeWSConnKraken) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConnKraken) Close() error {
	return nil
}

// --------- Helpers ---------

// ensureKrakenExchangeMap makes sure the global Exchanges map has
// an entry for KRAKEN_EXCHANGE so that Trade.Exchange can be set.
func ensureKrakenExchangeMap() {
	if Exchanges == nil {
		Exchanges = make(map[string]models.Exchange)
	}
	if _, ok := Exchanges[KRAKEN_EXCHANGE]; !ok {
		Exchanges[KRAKEN_EXCHANGE] = models.Exchange{}
	}
}

// --------- Key mapping tests ---------

func TestKrakenHooks_TickerAndLastTradeKeys(t *testing.T) {
	h := krakenHooks{}

	// TickerKeyFromForeign should remove '-' (e.g. BTC-USDT -> BTCUSDT)
	if got := h.TickerKeyFromForeign("BTC-USDT"); got != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", got)
	}

	// LastTradeTimeKeyFromForeign should keep the original foreign name
	if got := h.LastTradeTimeKeyFromForeign("ETH-USDT"); got != "ETH-USDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USDT, got %s", got)
	}
}

// --------- Subscribe / Unsubscribe tests ---------

func TestKrakenHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnKraken{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := krakenHooks{}
	var lock sync.RWMutex

	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true -> Method = "subscribe", Channel = "trade", Symbol = ["BTC/USDT"]
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call for subscribe, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(*krakenWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *krakenWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "subscribe" {
		t.Fatalf("expected Method=subscribe, got %s", msg.Method)
	}
	if msg.Params.Channel != "trade" {
		t.Fatalf("expected Channel=trade, got %s", msg.Params.Channel)
	}
	if len(msg.Params.Symbol) != 1 || msg.Params.Symbol[0] != "BTC/USDT" {
		t.Fatalf("expected Symbol=[BTC/USDT], got %+v", msg.Params.Symbol)
	}

	// subscribe = false -> Method = "unsubscribe"
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(*krakenWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *krakenWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Method != "unsubscribe" {
		t.Fatalf("expected Method=unsubscribe, got %s", msg.Method)
	}
	if msg.Params.Channel != "trade" {
		t.Fatalf("expected Channel=trade, got %s", msg.Params.Channel)
	}
	if len(msg.Params.Symbol) != 1 || msg.Params.Symbol[0] != "BTC/USDT" {
		t.Fatalf("expected Symbol=[BTC/USDT], got %+v", msg.Params.Symbol)
	}
}

// --------- parseKrakenTradeMessage tests ---------

func TestParseKrakenTradeMessage_BasicBuy(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	msg := krakenWSResponseData{
		Symbol:    "BTC/USDT",
		Side:      "buy",
		Price:     25000.5,
		Size:      1.23,
		OrderType: "market",
		TradeID:   123,
		Time:      now.Format("2006-01-02T15:04:05.000000Z"),
	}

	price, volume, ts, ftID, err := parseKrakenTradeMessage(msg)
	if err != nil {
		t.Fatalf("parseKrakenTradeMessage returned error: %v", err)
	}
	if price != 25000.5 {
		t.Fatalf("expected price=25000.5, got %v", price)
	}
	if volume != 1.23 {
		t.Fatalf("expected volume=1.23, got %v", volume)
	}
	if !ts.Equal(now) {
		t.Fatalf("expected timestamp=%v, got %v", now, ts)
	}
	if ftID != "123" {
		t.Fatalf("expected foreignTradeID=123, got %s", ftID)
	}
}

func TestParseKrakenTradeMessage_SellNegativeVolume(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	msg := krakenWSResponseData{
		Symbol:    "ETH/USDT",
		Side:      "sell",
		Price:     100.0,
		Size:      5.0,
		OrderType: "limit",
		TradeID:   99,
		Time:      now.Format("2006-01-02T15:04:05.000000Z"),
	}

	_, volume, _, _, err := parseKrakenTradeMessage(msg)
	if err != nil {
		t.Fatalf("parseKrakenTradeMessage returned error: %v", err)
	}
	if volume != -5.0 {
		t.Fatalf("expected volume=-5.0 for sell trade, got %v", volume)
	}
}

// --------- OnMessage tests ---------

func TestKrakenHooks_OnMessage_ValidTrades(t *testing.T) {
	ensureKrakenExchangeMap()

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 10),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	// tickerPairMap key is "BTCUSDT"
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "USDT"},
		QuoteToken: models.Asset{Symbol: "BTC"},
	}

	h := krakenHooks{}
	var lock sync.RWMutex

	now := time.Now().UTC().Truncate(time.Second)
	msg := krakenWSResponse{
		Channel: "trade",
		Type:    "snapshot",
		Data: []krakenWSResponseData{
			{
				Symbol:    "BTC/USDT",
				Side:      "buy",
				Price:     25000.0,
				Size:      0.5,
				OrderType: "market",
				TradeID:   1,
				Time:      now.Format("2006-01-02T15:04:05.000000Z"),
			},
			{
				Symbol:    "BTC/USDT",
				Side:      "sell",
				Price:     25010.0,
				Size:      1.0,
				OrderType: "limit",
				TradeID:   2,
				Time:      now.Format("2006-01-02T15:04:05.000000Z"),
			},
		},
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	// Expect two trades in the channel.
	trades := make([]models.Trade, 0, 2)
loop:
	for {
		select {
		case tr := <-bs.tradesChannel:
			trades = append(trades, tr)
			if len(trades) == 2 {
				break loop
			}
		default:
			break loop
		}
	}

	if len(trades) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(trades))
	}

	// First trade: buy side, positive volume.
	if trades[0].Price != 25000.0 || trades[0].Volume != 0.5 {
		t.Fatalf("unexpected first trade: price=%v volume=%v", trades[0].Price, trades[0].Volume)
	}
	if trades[0].BaseToken.Symbol != "USDT" || trades[0].QuoteToken.Symbol != "BTC" {
		t.Fatalf("unexpected tokens on first trade: base=%s quote=%s",
			trades[0].BaseToken.Symbol, trades[0].QuoteToken.Symbol)
	}

	// Second trade: sell side, negative volume.
	if trades[1].Price != 25010.0 || trades[1].Volume != -1.0 {
		t.Fatalf("unexpected second trade: price=%v volume=%v", trades[1].Price, trades[1].Volume)
	}

	// lastTradeTimeMap key: Quote-Base format, e.g. "BTC-USDT"
	lock.RLock()
	_, ok := bs.lastTradeTimeMap["BTC-USDT"]
	lock.RUnlock()
	if !ok {
		t.Fatalf("expected lastTradeTimeMap[BTC-USDT] to be set")
	}
}

func TestKrakenHooks_OnMessage_IgnoresNonTradeChannelOrNonText(t *testing.T) {
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	h := krakenHooks{}
	var lock sync.RWMutex

	// Non-text message type should be ignored.
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-text message")
	default:
	}

	// Different channel should be ignored.
	msg := krakenWSResponse{
		Channel: "book", // not "trade"
		Type:    "snapshot",
		Data:    nil,
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	h.OnMessage(bs, ws.TextMessage, raw, &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-trade channel")
	default:
	}
}

// --------- ReadLoop handled flag ---------

func TestKrakenHooks_ReadLoopHandledFalse(t *testing.T) {
	h := krakenHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex

	handled := h.ReadLoop(context.Background(), bs, &lock)
	if handled {
		t.Fatalf("expected krakenHooks.ReadLoop to return false")
	}
}
