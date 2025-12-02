package scrapers

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DEXHooks interface {
	ExchangeName() string

	EnsurePair(
		ctx context.Context,
		base *BaseDEXScraper,
		pool models.Pool,
		lock *sync.RWMutex,
	) (addr common.Address, watchdogDelay int64, err error)

	StartStream(
		ctx context.Context,
		base *BaseDEXScraper,
		addr common.Address,
		tradesChan chan models.Trade,
		lock *sync.RWMutex,
	)

	OnPoolRemoved(
		base *BaseDEXScraper,
		addr common.Address,
		lock *sync.RWMutex,
	)

	// OnOrderChanged is called when pool.Order changed in the config diff.
	// The hook should update any DEX-specific order fields stored in its pool map.
	OnOrderChanged(
		base *BaseDEXScraper,
		addr common.Address,
		oldPool models.Pool,
		newPool models.Pool,
		lock *sync.RWMutex,
	)
}

type BaseDEXScraper struct {
	exchange models.Exchange

	restClient *ethclient.Client
	wsClient   *ethclient.Client

	// Internal control channels shared by all DEX scrapers.
	subscribeChannel   chan common.Address // triggered by watchdog / subscription error to resub a pool
	unsubscribeChannel chan common.Address // triggered by config diff to stop a pool

	// Per-pool runtime state.
	lastTradeTimeMap map[common.Address]time.Time
	streamCancel     map[common.Address]context.CancelFunc
	watchdogCancel   map[string]context.CancelFunc // key: addr.Hex()

	// Optional: used if you need backoff tuning in hooks; not enforced by base.
	waitTime int
}

func NewBaseDEXScraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	hooks DEXHooks,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	wg *sync.WaitGroup,
) *BaseDEXScraper {
	if wg != nil {
		wg.Done()
	}

	var err error
	base := &BaseDEXScraper{
		exchange: models.Exchange{
			Name:       exchangeName,
			Blockchain: blockchain,
		},
		subscribeChannel:   make(chan common.Address),
		unsubscribeChannel: make(chan common.Address),
		lastTradeTimeMap:   make(map[common.Address]time.Time),
		streamCancel:       make(map[common.Address]context.CancelFunc),
		watchdogCancel:     make(map[string]context.CancelFunc),
	}

	// Init waitTime
	base.waitTime, err = strconv.Atoi(utils.Getenv(strings.ToUpper(exchangeName)+"_WAIT_TIME", "500"))
	if err != nil {
		log.Errorf("BaseDEX - Failed to parse waitTime for exchange %s: %v", exchangeName, err)
	}

	// Init REST / WS eth clients
	restURI := utils.Getenv(strings.ToUpper(exchangeName)+"_URI_REST", "")
	wsURI := utils.Getenv(strings.ToUpper(exchangeName)+"_URI_WS", "")

	if restURI != "" {
		base.restClient, err = ethclient.Dial(restURI)
		if err != nil {
			log.Errorf("%s - init rest client: %v", exchangeName, err)
		}
	} else {
		log.Warnf("%s - no REST URI configured", exchangeName)
	}

	if wsURI != "" {
		base.wsClient, err = ethclient.Dial(wsURI)
		if err != nil {
			log.Errorf("%s - init ws client: %v", exchangeName, err)
		}
	} else {
		log.Warnf("%s - no WS URI configured", exchangeName)
	}

	var lock sync.RWMutex

	// Start config-based pools (initial set)
	base.startInitialPools(ctx, hooks, pools, tradesChannel, &lock)

	// Start generic handlers
	base.startResubHandler(ctx, hooks, tradesChannel, &lock)
	base.startUnsubHandler(ctx, hooks, &lock)
	go base.watchConfig(ctx, hooks, tradesChannel, &lock, pools)

	return base
}

// DEX hooks typically need this to read metadata (token0/token1, decimals, etc.).
func (b *BaseDEXScraper) RESTClient() *ethclient.Client {
	return b.restClient
}

// Hooks must use this in StartStream.
func (b *BaseDEXScraper) WSClient() *ethclient.Client {
	return b.wsClient
}

// SubscribeChannel returns the internal subscribe channel used for resubscribing pools.
func (b *BaseDEXScraper) SubscribeChannel() chan common.Address {
	return b.subscribeChannel
}

// UnsubscribeChannel returns the internal unsubscribe channel used for stopping pools.
func (b *BaseDEXScraper) UnsubscribeChannel() chan common.Address {
	return b.unsubscribeChannel
}

// Hooks should call this whenever they emit a trade for addr.
func (b *BaseDEXScraper) UpdateLastTradeTime(addr common.Address, t time.Time, lock *sync.RWMutex) {
	lock.Lock()
	b.lastTradeTimeMap[addr] = t
	lock.Unlock()
}

func (b *BaseDEXScraper) startInitialPools(
	ctx context.Context,
	hooks DEXHooks,
	pools []models.Pool,
	trades chan models.Trade,
	lock *sync.RWMutex,
) {
	for _, p := range pools {
		if err := b.startPool(ctx, hooks, p, trades, lock); err != nil {
			log.Errorf("%s - startInitialPool %s: %v", hooks.ExchangeName(), p.Address, err)
		}
	}
}

func (b *BaseDEXScraper) startResubHandler(
	ctx context.Context,
	hooks DEXHooks,
	trades chan models.Trade,
	lock *sync.RWMutex,
) {
	go func() {
		for {
			select {
			case addr := <-b.subscribeChannel:
				log.Infof("%s - Resubscribing to pool: %s", hooks.ExchangeName(), addr.Hex())
				// Reset lastTradeTime to avoid immediate watchdog retrigger.
				lock.Lock()
				b.lastTradeTimeMap[addr] = time.Now()
				// Stop old stream if any.
				if c, ok := b.streamCancel[addr]; ok && c != nil {
					c()
					delete(b.streamCancel, addr)
				}
				lock.Unlock()
				// Restart stream.
				pctx, cancel := context.WithCancel(ctx)
				b.streamCancel[addr] = cancel
				hooks.StartStream(pctx, b, addr, trades, lock)

			case <-ctx.Done():
				log.Infof("%s - Stopping resubscription handler.", hooks.ExchangeName())
				return
			}
		}
	}()
}

func (b *BaseDEXScraper) startUnsubHandler(
	ctx context.Context,
	hooks DEXHooks,
	lock *sync.RWMutex,
) {
	go func() {
		for {
			select {
			case addr := <-b.unsubscribeChannel:
				log.Infof("%s - Unsubscribing pool: %s", hooks.ExchangeName(), addr.Hex())
				b.stopPool(addr, hooks, lock)
			case <-ctx.Done():
				log.Infof("%s - unsubscribe handler stopped", hooks.ExchangeName())
				return
			}
		}
	}()
}

// startPool:
//  1. hooks.EnsurePair(...) to create/update DEX-specific pool entry and get (addr, delay)
//  2. init lastTradeTimeMap
//  3. restart watchdog for this addr
//  4. start stream if not already running
func (b *BaseDEXScraper) startPool(
	ctx context.Context,
	hooks DEXHooks,
	pool models.Pool,
	trades chan models.Trade,
	lock *sync.RWMutex,
) error {
	addr, delay, err := hooks.EnsurePair(ctx, b, pool, lock)
	if err != nil {
		return err
	}
	if delay <= 0 {
		delay = 300
	}

	// Init lastTradeTime, avoid triggering watchdog immediately.
	lock.Lock()
	if _, ok := b.lastTradeTimeMap[addr]; !ok {
		b.lastTradeTimeMap[addr] = time.Now()
	}
	lock.Unlock()

	// Restart watchdog
	b.restartWatchdogForAddr(ctx, addr, delay, lock)

	// Start listening (if already running, don't repeat).
	if _, ok := b.streamCancel[addr]; ok {
		return nil
	}
	pctx, cancel := context.WithCancel(ctx)
	b.streamCancel[addr] = cancel
	hooks.StartStream(pctx, b, addr, trades, lock)

	return nil
}

// stopPool cancels stream + watchdog + lastTradeTimeMap, and let hooks clean its own metadata.
func (b *BaseDEXScraper) stopPool(
	addr common.Address,
	hooks DEXHooks,
	lock *sync.RWMutex,
) {
	if c, ok := b.streamCancel[addr]; ok && c != nil {
		c()
		delete(b.streamCancel, addr)
	}
	key := addr.Hex()
	if c, ok := b.watchdogCancel[key]; ok && c != nil {
		c()
		delete(b.watchdogCancel, key)
	}
	lock.Lock()
	delete(b.lastTradeTimeMap, addr)
	lock.Unlock()

	// Let DEX-specific scraper clean its own poolMap etc.
	hooks.OnPoolRemoved(b, addr, lock)
}

func (b *BaseDEXScraper) restartWatchdogForAddr(
	ctx context.Context,
	addr common.Address,
	delay int64,
	lock *sync.RWMutex,
) {
	key := addr.Hex()
	if c, ok := b.watchdogCancel[key]; ok && c != nil {
		c()
		delete(b.watchdogCancel, key)
	}
	wdCtx, cancel := context.WithCancel(ctx)
	b.watchdogCancel[key] = cancel

	t := time.NewTicker(time.Duration(delay) * time.Second)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				lock.RLock()
				last, ok := b.lastTradeTimeMap[addr]
				lock.RUnlock()
				if !ok {
					// Pool might have been removed.
					continue
				}
				if time.Since(last) >= time.Duration(delay)*time.Second {
					// Trigger resubscribe.
					b.subscribeChannel <- addr
				}
			case <-wdCtx.Done():
				return
			}
		}
	}()
}

func (b *BaseDEXScraper) watchConfig(
	ctx context.Context,
	hooks DEXHooks,
	trades chan models.Trade,
	lock *sync.RWMutex,
	initialPools []models.Pool,
) {
	ex := hooks.ExchangeName()

	envKey := strings.ToUpper(ex) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "3600"))
	if err != nil {
		log.Errorf("%s - Failed to parse %s: %v.", ex, envKey, err)
		return
	}
	tk := time.NewTicker(time.Duration(interval) * time.Second)
	defer tk.Stop()

	last := mapPoolsByAddrLower(initialPools)

	for {
		select {
		case <-tk.C:
			nowList, err := models.PoolsFromConfigFile(ex)
			if err != nil {
				log.Errorf("%s - reload pools: %v", ex, err)
				continue
			}
			now := mapPoolsByAddrLower(nowList)
			b.applyConfigDiff(ctx, hooks, last, now, trades, lock)
			last = now

		case <-ctx.Done():
			log.Infof("%s - watchConfig exit", ex)
			return
		}
	}
}

func mapPoolsByAddrLower(pools []models.Pool) map[string]models.Pool {
	m := make(map[string]models.Pool, len(pools))
	for _, p := range pools {
		m[strings.ToLower(p.Address)] = p
	}
	return m
}

// applyConfigDiff computes:
//   - removed pools: send to unsubscribeChannel;
//   - added pools: startPool;
//   - changed pools: update Order via hooks, and restart watchdog if WatchDogDelay changed.
func (b *BaseDEXScraper) applyConfigDiff(
	ctx context.Context,
	hooks DEXHooks,
	last map[string]models.Pool,
	curr map[string]models.Pool,
	trades chan models.Trade,
	lock *sync.RWMutex,
) {
	ex := hooks.ExchangeName()

	// Remove
	for k := range last {
		if _, ok := curr[k]; !ok {
			addr := common.HexToAddress(k)
			log.Infof("%s - remove pool %s", ex, addr.Hex())
			b.unsubscribeChannel <- addr
		}
	}

	// Add + update
	for k, p := range curr {
		addr := common.HexToAddress(k)
		old, existed := last[k]
		if !existed {
			// New pool
			log.Infof("%s - add pool %s (wd=%d, order=%d)", ex, p.Address, p.WatchDogDelay, p.Order)
			if err := b.startPool(ctx, hooks, p, trades, lock); err != nil {
				log.Errorf("%s - startPool %s: %v", ex, p.Address, err)
			}
			continue
		}

		// Order changed
		if old.Order != p.Order {
			log.Infof("%s - update order %s: %d -> %d", ex, addr.Hex(), old.Order, p.Order)
			hooks.OnOrderChanged(b, addr, old, p, lock)
		}

		// Watchdog delay changed
		newWD := p.WatchDogDelay
		oldWD := old.WatchDogDelay
		if newWD != oldWD {
			log.Infof("%s - update watchdog %s: %d -> %d", ex, addr.Hex(), oldWD, newWD)
			if newWD <= 0 {
				newWD = 300
			}
			b.restartWatchdogForAddr(ctx, addr, newWD, lock)
		}
	}
}
