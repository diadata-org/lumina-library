package scrapers

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
)

// ---- fake hooks implement DEXHooks, for observing BaseDEXScraper behavior ----

type fakeDEXHooks struct {
	exName string

	mu sync.Mutex

	ensureCalls      int
	startStreamCalls int
	removed          []common.Address
	orderChanges     []struct {
		addr common.Address
		old  models.Pool
		new  models.Pool
	}

	// control the return of EnsurePair
	ensureDelay int64
}

func (f *fakeDEXHooks) ExchangeName() string {
	return f.exName
}

func (f *fakeDEXHooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {
	f.mu.Lock()
	f.ensureCalls++
	f.mu.Unlock()

	addr = common.HexToAddress(pool.Address)
	delay := pool.WatchDogDelay
	if f.ensureDelay > 0 {
		delay = f.ensureDelay
	}
	return addr, delay, nil
}

func (f *fakeDEXHooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	f.mu.Lock()
	f.startStreamCalls++
	f.mu.Unlock()
}

func (f *fakeDEXHooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addr common.Address,
	lock *sync.RWMutex,
) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.removed = append(f.removed, addr)
}

func (f *fakeDEXHooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
	lock *sync.RWMutex,
) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.orderChanges = append(f.orderChanges, struct {
		addr common.Address
		old  models.Pool
		new  models.Pool
	}{addr, oldPool, newPool})
}

// ---- Tests ----

// construct a BaseDEXScraper without calling NewBaseDEXScraper (to avoid real ethclient.Dial)
func newTestBase(exchangeName, chain string) *BaseDEXScraper {
	return &BaseDEXScraper{
		exchange: models.Exchange{
			Name:       exchangeName,
			Blockchain: chain,
		},
		subscribeChannel:   make(chan common.Address, 10),
		unsubscribeChannel: make(chan common.Address, 10),
		lastTradeTimeMap:   make(map[common.Address]time.Time),
		streamCancel:       make(map[common.Address]context.CancelFunc),
		watchdogCancel:     make(map[string]context.CancelFunc),
	}
}

// Test: startPool will call EnsurePair + StartStream when a new pool is created, and write to streamCancel and lastTradeTimeMap
func TestBaseDEXScraper_StartPool_NewPoolStartsStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base := newTestBase("UniswapV3", "ethereum")
	hooks := &fakeDEXHooks{exName: "UniswapV3"}

	var lock sync.RWMutex
	trades := make(chan models.Trade, 10)

	pool := models.Pool{
		Address:       "0x0000000000000000000000000000000000000123",
		WatchDogDelay: 5,
		Order:         0,
	}

	if err := base.startPool(ctx, hooks, pool, trades, &lock); err != nil {
		t.Fatalf("startPool returned error: %v", err)
	}

	addr := common.HexToAddress(pool.Address)

	// confirm lastTradeTimeMap contains this address
	lock.RLock()
	_, ok := base.lastTradeTimeMap[addr]
	lock.RUnlock()
	if !ok {
		t.Fatalf("expected lastTradeTimeMap to contain addr %s", addr.Hex())
	}

	// confirm streamCancel contains a cancel function
	if _, ok := base.streamCancel[addr]; !ok {
		t.Fatalf("expected streamCancel to contain addr %s", addr.Hex())
	}

	// fake hooks calls count
	hooks.mu.Lock()
	defer hooks.mu.Unlock()
	if hooks.ensureCalls != 1 {
		t.Fatalf("expected EnsurePair to be called once, got %d", hooks.ensureCalls)
	}
	if hooks.startStreamCalls != 1 {
		t.Fatalf("expected StartStream to be called once, got %d", hooks.startStreamCalls)
	}
}

// Test: if stream already exists, startPool will not restart StartStream
func TestBaseDEXScraper_StartPool_DoesNotRestartExistingStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base := newTestBase("UniswapV3", "ethereum")
	hooks := &fakeDEXHooks{exName: "UniswapV3"}

	var lock sync.RWMutex
	trades := make(chan models.Trade, 10)

	pool := models.Pool{
		Address:       "0x0000000000000000000000000000000000000ABC",
		WatchDogDelay: 5,
		Order:         0,
	}
	addr := common.HexToAddress(pool.Address)

	// pre-place a streamCancel, simulate an existing stream is running
	base.streamCancel[addr] = func() {}

	if err := base.startPool(ctx, hooks, pool, trades, &lock); err != nil {
		t.Fatalf("startPool returned error: %v", err)
	}

	hooks.mu.Lock()
	defer hooks.mu.Unlock()
	// EnsurePair will be called again (because startPool calls hooks.EnsurePair first)
	if hooks.ensureCalls != 1 {
		t.Fatalf("expected EnsurePair to be called once, got %d", hooks.ensureCalls)
	}
	// but StartStream should not be called again
	if hooks.startStreamCalls != 0 {
		t.Fatalf("expected StartStream NOT to be called, got %d", hooks.startStreamCalls)
	}
}

// Test: startResubHandler will cancel the old stream and start a new one when it receives an address from subscribeChannel
func TestBaseDEXScraper_StartResubHandler_Resubscribes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base := newTestBase("UniswapV3", "ethereum")
	hooks := &fakeDEXHooks{exName: "UniswapV3"}

	var lock sync.RWMutex
	trades := make(chan models.Trade, 10)

	addr := common.HexToAddress("0x0000000000000000000000000000000000000DEF")

	var canceled int32
	base.streamCancel[addr] = func() {
		atomic.AddInt32(&canceled, 1)
	}

	// start resub handler
	base.startResubHandler(ctx, hooks, trades, &lock)

	// trigger watchdog / subscription error behavior: write address to subscribeChannel
	base.subscribeChannel <- addr

	// wait for a short time, let the goroutine handle the resubscription
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for resubscription handling")
		default:
			// expect cancel to be called once, StartStream to be called once
			hooks.mu.Lock()
			startCalls := hooks.startStreamCalls
			hooks.mu.Unlock()

			if atomic.LoadInt32(&canceled) == 1 && startCalls == 1 {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Test: applyConfigDiff will write the address to unsubscribeChannel when a pool is removed
func TestBaseDEXScraper_ApplyConfigDiff_RemovedPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base := newTestBase("UniswapV3", "ethereum")
	hooks := &fakeDEXHooks{exName: "UniswapV3"}

	var lock sync.RWMutex
	trades := make(chan models.Trade, 10)

	// last has one pool, curr is empty => should trigger remove
	lastAddr := "0x0000000000000000000000000000000000000AAA"
	last := map[string]models.Pool{
		strings.ToLower(lastAddr): {
			Address:       lastAddr,
			WatchDogDelay: 5,
			Order:         0,
		},
	}
	curr := map[string]models.Pool{}

	go base.applyConfigDiff(ctx, hooks, last, curr, trades, "", &lock)

	// read an address from unsubscribeChannel, check if it is lastAddr
	select {
	case addr := <-base.unsubscribeChannel:
		if addr.Hex() != common.HexToAddress(lastAddr).Hex() {
			t.Fatalf("expected unsubscribe addr %s, got %s", lastAddr, addr.Hex())
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for unsubscribeChannel")
	}
}
