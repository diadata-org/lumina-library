package scrapers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

// ---- fake wsConn ----

type fakeConn struct {
	readCount  int
	writeCount int
	closed     bool

	// control ReadMessage behavior
	readMsgType int
	readMsgData []byte
	readErr     error
}

func (c *fakeConn) ReadMessage() (int, []byte, error) {
	c.readCount++
	if c.readErr != nil {
		return 0, nil, c.readErr
	}
	return c.readMsgType, c.readMsgData, nil
}

func (c *fakeConn) WriteMessage(messageType int, data []byte) error {
	c.writeCount++
	return nil
}

func (c *fakeConn) ReadJSON(v interface{}) error {
	return nil
}

func (c *fakeConn) WriteJSON(v interface{}) error {
	c.writeCount++
	return nil
}

func (c *fakeConn) Close() error {
	c.closed = true
	return nil
}

// ---- fake hooks ----

type testHooks struct {
	exKey string
	url   string

	onOpenCalled    bool
	subscribeCalls  []string
	onMessageCalls  int
	readLoopHandled bool

	// Dial related
	dialConn  wsConn
	dialErr   error
	dialCount int
}

// implement ScraperHooks
func (h *testHooks) ExchangeKey() string { return h.exKey }
func (h *testHooks) WSURL() string       { return h.url }

func (h *testHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	h.onOpenCalled = true
}

func (h *testHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	if subscribe {
		h.subscribeCalls = append(h.subscribeCalls, pair.ForeignName)
	}
	return nil
}

func (h *testHooks) OnMessage(bs *BaseCEXScraper, messageType int, data []byte, lock *sync.RWMutex) {
	h.onMessageCalls++
}

func (h *testHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	// default false, use Base's default ReadMessage loop
	return h.readLoopHandled
}

func (h *testHooks) TickerKeyFromForeign(foreign string) string {
	// in test, simply return original string
	return foreign
}

func (h *testHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// implement CustomDialer
func (h *testHooks) Dial(ctx context.Context, url string) (wsConn, error) {
	h.dialCount++
	return h.dialConn, h.dialErr
}

// ---- tests ----

// confirm: connectWithRetry uses CustomDialer and exits after first Dial success
func TestBaseCEXScraper_ConnectWithRetry_UsesCustomDialer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := &testHooks{
		exKey:    "test_exchange",
		url:      "wss://example.test/ws",
		dialConn: &fakeConn{},
		dialErr:  nil,
	}
	bs := &BaseCEXScraper{
		hooks:           h,
		maxErrCount:     20,
		restartWaitTime: 1,
	}

	ok := bs.connectWithRetry(ctx)
	if !ok {
		t.Fatalf("expected connectWithRetry to succeed")
	}

	if h.dialCount != 1 {
		t.Fatalf("expected Dial to be called once, got %d", h.dialCount)
	}

	if bs.wsClient == nil {
		t.Fatalf("expected wsClient to be set after successful connect")
	}
}

// confirm: dialOnce uses custom Dial when hooks implements CustomDialer
func TestBaseCEXScraper_DialOnce_UsesCustomDialer(t *testing.T) {
	ctx := context.Background()

	fc := &fakeConn{}
	h := &testHooks{
		exKey:    "test_exchange",
		url:      "wss://example.test/ws",
		dialConn: fc,
		dialErr:  nil,
	}
	bs := &BaseCEXScraper{
		hooks: h,
	}

	conn, err := bs.dialOnce(ctx)
	if err != nil {
		t.Fatalf("dialOnce returned error: %v", err)
	}
	if conn != fc {
		t.Fatalf("expected dialOnce to return fakeConn from CustomDialer")
	}
	if h.dialCount != 1 {
		t.Fatalf("expected CustomDialer.Dial to be called once, got %d", h.dialCount)
	}
}

// confirm: diffPairMap's added / removed / changed logic
func TestDiffPairMap(t *testing.T) {
	last := map[string]int64{
		"BTC-USDT": 60,
		"ETH-USDT": 60,
	}
	current := map[string]int64{
		"BTC-USDT": 120, // delay changed
		"XRP-USDT": 30,  // new
	}

	added, removed, changed := diffPairMap(last, current)

	// helper
	contains := func(slice []string, v string) bool {
		for _, s := range slice {
			if s == v {
				return true
			}
		}
		return false
	}

	if !contains(added, "XRP-USDT") || len(added) != 1 {
		t.Fatalf("expected added to contain only XRP-USDT, got %+v", added)
	}
	if !contains(removed, "ETH-USDT") || len(removed) != 1 {
		t.Fatalf("expected removed to contain only ETH-USDT, got %+v", removed)
	}
	if !contains(changed, "BTC-USDT") || len(changed) != 1 {
		t.Fatalf("expected changed to contain only BTC-USDT, got %+v", changed)
	}
}

// confirm: runReadLoop exits when ReadMessage returns "closed network connection"
// (with handleErrorReadJSON logic)
func TestBaseCEXScraper_RunReadLoop_StopsOnClosedNetworkConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fc := &fakeConn{
		readErr: errors.New("use of closed network connection"),
	}
	h := &testHooks{
		exKey: "test_exchange",
		url:   "wss://example.test/ws",
	}
	bs := &BaseCEXScraper{
		hooks:           h,
		wsClient:        fc,
		maxErrCount:     20,
		restartWaitTime: 1,
	}

	var lock sync.RWMutex

	done := make(chan struct{})
	go func() {
		bs.runReadLoop(ctx, &lock)
		close(done)
	}()

	select {
	case <-done:
		// ok, loop exited
	case <-time.After(2 * time.Second):
		t.Fatalf("runReadLoop did not exit on closed network connection error")
	}
}
