package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

//
// ------- Fake wsConn implementation for KuCoin tests -------
//

// fakeWSConnKuCoin is a fake implementation of wsConn used to
// capture WriteJSON calls and avoid any real network activity.
type fakeWSConnKuCoin struct {
	writeJSONCount int
	lastWritten    interface{}
}

func (f *fakeWSConnKuCoin) ReadMessage() (int, []byte, error) {
	// Not needed for these tests.
	return 0, nil, nil
}

func (f *fakeWSConnKuCoin) WriteMessage(messageType int, data []byte) error {
	// Not needed for these tests.
	return nil
}

func (f *fakeWSConnKuCoin) ReadJSON(v interface{}) error {
	// Not needed for these tests.
	return nil
}

func (f *fakeWSConnKuCoin) WriteJSON(v interface{}) error {
	f.writeJSONCount++
	f.lastWritten = v
	return nil
}

func (f *fakeWSConnKuCoin) Close() error {
	return nil
}

//
// ------- Helpers -------
//

// ensureKuCoinExchangeMap makes sure the global Exchanges map has
// an entry for KUCOIN_EXCHANGE so Trade.Exchange can be set.
func ensureKuCoinExchangeMap() {
	if Exchanges == nil {
		Exchanges = make(map[string]models.Exchange)
	}
	if _, ok := Exchanges[KUCOIN_EXCHANGE]; !ok {
		Exchanges[KUCOIN_EXCHANGE] = models.Exchange{}
	}
}

//
// ------- Key mapping tests -------
//

func TestKuCoinHooks_TickerAndLastTradeKeys(t *testing.T) {
	h := kucoinHooks{}

	// TickerKeyFromForeign should remove '-' (e.g. BTC-USDT -> BTCUSDT)
	if got := h.TickerKeyFromForeign("BTC-USDT"); got != "BTCUSDT" {
		t.Fatalf("TickerKeyFromForeign: expected BTCUSDT, got %s", got)
	}

	// LastTradeTimeKeyFromForeign should keep the original foreign name
	if got := h.LastTradeTimeKeyFromForeign("ETH-USDT"); got != "ETH-USDT" {
		t.Fatalf("LastTradeTimeKeyFromForeign: expected ETH-USDT, got %s", got)
	}
}

//
// ------- Subscribe / Unsubscribe tests -------
//

func TestKuCoinHooks_SubscribeAndUnsubscribe(t *testing.T) {
	fc := &fakeWSConnKuCoin{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := kucoinHooks{}
	var lock sync.RWMutex

	pair := models.ExchangePair{ForeignName: "BTC-USDT"}

	// subscribe = true => Type = "subscribe", Topic = "/market/match:BTC-USDT"
	if err := h.Subscribe(bs, pair, true, &lock); err != nil {
		t.Fatalf("Subscribe(true) returned error: %v", err)
	}
	if fc.writeJSONCount != 1 {
		t.Fatalf("expected 1 WriteJSON call for subscribe, got %d", fc.writeJSONCount)
	}
	msg, ok := fc.lastWritten.(*kuCoinWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *kuCoinWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Type != "subscribe" {
		t.Fatalf("expected Type=subscribe, got %s", msg.Type)
	}
	if msg.Topic != "/market/match:BTC-USDT" {
		t.Fatalf("expected Topic=/market/match:BTC-USDT, got %s", msg.Topic)
	}
	if msg.PrivateChannel != false {
		t.Fatalf("expected PrivateChannel=false, got %v", msg.PrivateChannel)
	}
	if msg.Response != true {
		t.Fatalf("expected Response=true, got %v", msg.Response)
	}

	// subscribe = false => Type = "unsubscribe"
	if err := h.Subscribe(bs, pair, false, &lock); err != nil {
		t.Fatalf("Subscribe(false) returned error: %v", err)
	}
	if fc.writeJSONCount != 2 {
		t.Fatalf("expected 2 WriteJSON calls in total, got %d", fc.writeJSONCount)
	}
	msg, ok = fc.lastWritten.(*kuCoinWSSubscribeMessage)
	if !ok {
		t.Fatalf("expected lastWritten to be *kuCoinWSSubscribeMessage, got %T", fc.lastWritten)
	}
	if msg.Type != "unsubscribe" {
		t.Fatalf("expected Type=unsubscribe, got %s", msg.Type)
	}
	if msg.Topic != "/market/match:BTC-USDT" {
		t.Fatalf("expected Topic=/market/match:BTC-USDT, got %s", msg.Topic)
	}
}

//
// ------- parseKuCoinTradeMessage tests -------
//

func TestParseKuCoinTradeMessage_BasicBuy(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	msg := kuCoinWSResponse{
		Type: "message",
		Data: kuCoinWSData{
			Symbol:  "BTC-USDT",
			Side:    "buy",
			Price:   "25000.5",
			Size:    "1.23",
			TradeID: "trade-123",
			Time:    strconv.FormatInt(nowMs, 10),
		},
	}

	price, volume, ts, ftID, err := parseKuCoinTradeMessage(msg)
	if err != nil {
		t.Fatalf("parseKuCoinTradeMessage returned error: %v", err)
	}
	if price != 25000.5 {
		t.Fatalf("expected price=25000.5, got %v", price)
	}
	if volume != 1.23 {
		t.Fatalf("expected volume=1.23, got %v", volume)
	}
	// Only check that timestamp is non-zero and constructed.
	if ts.IsZero() {
		t.Fatalf("expected non-zero timestamp")
	}
	if ftID != "trade-123" {
		t.Fatalf("expected foreignTradeID=trade-123, got %s", ftID)
	}
}

func TestParseKuCoinTradeMessage_SellNegativeVolume(t *testing.T) {
	msg := kuCoinWSResponse{
		Type: "message",
		Data: kuCoinWSData{
			Symbol:  "ETH-USDT",
			Side:    "sell",
			Price:   "100",
			Size:    "5",
			TradeID: "sell-1",
			Time:    "123456789",
		},
	}
	_, volume, _, _, err := parseKuCoinTradeMessage(msg)
	if err != nil {
		t.Fatalf("parseKuCoinTradeMessage returned error: %v", err)
	}
	if volume != -5.0 {
		t.Fatalf("expected volume=-5.0 for sell side, got %v", volume)
	}
}

//
// ------- OnMessage tests -------
//

func TestKuCoinHooks_OnMessage_PongAndNonMessageIgnored(t *testing.T) {
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	h := kucoinHooks{}
	var lock sync.RWMutex

	// Non-text should be ignored.
	h.OnMessage(bs, ws.BinaryMessage, []byte(`{}`), &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-text message")
	default:
	}

	// "pong" should be logged and ignored as a trade.
	pong := kuCoinWSResponse{
		Type: "pong",
	}
	rawPong, err := json.Marshal(pong)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	h.OnMessage(bs, ws.TextMessage, rawPong, &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for pong message")
	default:
	}

	// Other non-"message" types should be ignored as trades.
	nonMessage := kuCoinWSResponse{
		Type: "welcome",
	}
	rawNonMsg, err := json.Marshal(nonMessage)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	h.OnMessage(bs, ws.TextMessage, rawNonMsg, &lock)
	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for non-message type")
	default:
	}
}

func TestKuCoinHooks_OnMessage_ValidTrade(t *testing.T) {
	ensureKuCoinExchangeMap()

	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair),
		lastTradeTimeMap: make(map[string]time.Time),
	}
	// tickerPairMap key is "BTCUSDT"
	bs.tickerPairMap["BTCUSDT"] = models.Pair{
		BaseToken:  models.Asset{Symbol: "USDT"},
		QuoteToken: models.Asset{Symbol: "BTC"},
	}

	h := kucoinHooks{}
	var lock sync.RWMutex

	nowMs := time.Now().UnixMilli()
	msg := kuCoinWSResponse{
		Type: "message",
		Data: kuCoinWSData{
			Symbol:  "BTC-USDT",
			Side:    "buy",
			Price:   "25000",
			Size:    "0.5",
			TradeID: "t-1",
			Time:    strconv.FormatInt(nowMs, 10),
		},
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case tr := <-bs.tradesChannel:
		if tr.Price != 25000 {
			t.Fatalf("expected trade price=25000, got %v", tr.Price)
		}
		if tr.Volume != 0.5 {
			t.Fatalf("expected volume=0.5, got %v", tr.Volume)
		}
		if tr.BaseToken.Symbol != "USDT" || tr.QuoteToken.Symbol != "BTC" {
			t.Fatalf("unexpected tokens: base=%s quote=%s", tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
		}
	default:
		t.Fatalf("expected one trade in tradesChannel")
	}

	// lastTradeTimeMap key should be "BTC-USDT"
	lock.RLock()
	_, ok := bs.lastTradeTimeMap["BTC-USDT"]
	lock.RUnlock()
	if !ok {
		t.Fatalf("expected lastTradeTimeMap[BTC-USDT] to be set")
	}
}

func TestKuCoinHooks_OnMessage_UnknownPairIgnored(t *testing.T) {
	bs := &BaseCEXScraper{
		tradesChannel:    make(chan models.Trade, 1),
		tickerPairMap:    make(map[string]models.Pair), // empty
		lastTradeTimeMap: make(map[string]time.Time),
	}
	h := kucoinHooks{}
	var lock sync.RWMutex

	msg := kuCoinWSResponse{
		Type: "message",
		Data: kuCoinWSData{
			Symbol:  "DOGE-USDT",
			Side:    "buy",
			Price:   "0.1",
			Size:    "100",
			TradeID: "doge-1",
			Time:    "123",
		},
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	h.OnMessage(bs, ws.TextMessage, raw, &lock)

	select {
	case <-bs.tradesChannel:
		t.Fatalf("expected no trade for unknown pair in tickerPairMap")
	default:
	}
}

//
// ------- OnOpen (ping goroutine) sanity test -------
//

func TestKuCoinHooks_OnOpen_StartsPingRoutine(t *testing.T) {
	// This test only verifies that calling OnOpen does not panic
	// and that a ping goroutine can be started against a fake wsConn.
	fc := &fakeWSConnKuCoin{}
	bs := &BaseCEXScraper{
		wsClient: fc,
	}
	h := kucoinHooks{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h.OnOpen(ctx, bs)

	// Let the goroutine run briefly and then cancel the context.
	time.Sleep(50 * time.Millisecond)
	// No assertions here, just ensuring no panic / deadlock.
}

//
// ------- getPublicKuCoinToken / WSURL tests -------
//

func TestGetPublicKuCoinToken_UsesHTTPResponse(t *testing.T) {
	// Create a fake KuCoin token HTTP endpoint.
	token := "fake-kucoin-token"
	pingInterval := int64(5000)

	respObj := kuCoinPostResponse{
		Code: "200000",
	}
	respObj.Data.Token = token
	respObj.Data.InstanceServers = []instanceServers{
		{PingInterval: pingInterval},
	}

	raw, err := json.Marshal(respObj)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method and path if desired.
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		_, _ = w.Write(raw)
	}))
	defer srv.Close()

	gotToken, gotPing, err := getPublicKuCoinToken(srv.URL)
	if err != nil {
		t.Fatalf("getPublicKuCoinToken returned error: %v", err)
	}
	if gotToken != token {
		t.Fatalf("expected token=%s, got %s", token, gotToken)
	}
	if gotPing != pingInterval {
		t.Fatalf("expected pingInterval=%d, got %d", pingInterval, gotPing)
	}
}

func TestKuCoinHooks_WSURL_UsesTokenFromEndpoint(t *testing.T) {
	// Build a fake token endpoint and override kucoinTokenURL for this test.
	token := "wsurl-token"
	resp := kuCoinPostResponse{
		Code: "200000",
	}
	resp.Data.Token = token
	resp.Data.InstanceServers = []instanceServers{
		{PingInterval: 1000},
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(raw)
	}))
	defer srv.Close()

	// Temporarily override the global kucoinTokenURL for this test.
	origURL := kucoinTokenURL
	kucoinTokenURL = srv.URL
	defer func() { kucoinTokenURL = origURL }()

	h := kucoinHooks{}
	got := h.WSURL()

	// Expected format: kucoinWSBaseString + "?token=" + token
	expectedPrefix := kucoinWSBaseString + "?token="
	if !strings.HasPrefix(got, expectedPrefix) {
		t.Fatalf("expected WSURL prefix %s, got %s", expectedPrefix, got)
	}
	if !bytes.Contains([]byte(got), []byte(token)) {
		t.Fatalf("expected WSURL to contain token=%s, got %s", token, got)
	}
}

//
// ------- ReadLoop handled flag -------
//

func TestKuCoinHooks_ReadLoopHandledFalse(t *testing.T) {
	h := kucoinHooks{}
	bs := &BaseCEXScraper{}
	var lock sync.RWMutex

	handled := h.ReadLoop(context.Background(), bs, &lock)
	if handled {
		t.Fatalf("expected kucoinHooks.ReadLoop to return false")
	}
}
