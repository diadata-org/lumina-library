package scrapers

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	curvefitwocryptooptimized "github.com/diadata-org/lumina-library/contracts/curve/curvefitwocrypto"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	models "github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

const NativeETHSentinel = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

type CurvePair struct {
	OutAsset models.Asset
	InAsset  models.Asset
	OutIndex int
	InIndex  int
	Address  common.Address
}

type CurveSwap struct {
	ID        string
	Timestamp int64
	Pair      CurvePair
	Amount0   float64 // sold amount
	Amount1   float64 // bought amount
}

// CurveScraper: attached to BaseDEXScraper, only maintains poolMap + some metadata
type CurveScraper struct {
	base       *BaseDEXScraper
	poolMap    map[common.Address][]CurvePair
	exchange   string
	blockchain string
}

type curveHooks struct {
	s *CurveScraper
}

// ExchangeName: used by Base to read config
func (h *curveHooks) ExchangeName() string {
	return h.s.exchange
}

// EnsurePair:
//   - Base will call this once when starting a pool.
//   - We ensure all (Address,Order) for this address, and return the maximum watchdogDelay.
func (h *curveHooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {
	addr = common.HexToAddress(pool.Address)

	// 1) ensure the direction corresponding to pool.Order is already in poolMap
	if err := h.s.ensurePairForPool(ctx, base, pool, lock); err != nil {
		return addr, 0, err
	}

	// 2) watchdogDelay is the WatchDogDelay of the pool itself
	watchdogDelay = pool.WatchDogDelay
	if watchdogDelay <= 0 {
		watchdogDelay = 300
	}
	return addr, watchdogDelay, nil
}

func (s *CurveScraper) ensurePairForPool(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) error {
	restClient := base.RESTClient()
	if restClient == nil {
		return fmt.Errorf("ensurePairForPool: REST client is nil")
	}

	outIdx, inIdx, err := parseIndexCode(pool.Order)
	if err != nil {
		return err
	}

	addr := common.HexToAddress(pool.Address)

	// if the direction already exists, return
	lock.RLock()
	existing := s.poolMap[addr]
	lock.RUnlock()
	for _, p := range existing {
		if p.OutIndex == outIdx && p.InIndex == inIdx {
			return nil
		}
	}

	// need to supplement this direction: get coin address + Asset from chain
	tokenOutAddr, err := coinAddressFromPool(ctx, restClient, addr, outIdx)
	if err != nil || tokenOutAddr == (common.Address{}) {
		return fmt.Errorf("get out-coin addr: %w", err)
	}

	tokenInAddr, err := coinAddressFromPool(ctx, restClient, addr, inIdx)
	if err != nil || tokenInAddr == (common.Address{}) {
		return fmt.Errorf("get in-coin addr: %w", err)
	}

	getAsset := func(a common.Address) (models.Asset, error) {
		if isNative(a) {
			return nativeAsset(s.blockchain), nil
		}
		return models.GetAsset(a, s.blockchain, restClient)
	}

	outAsset, err := getAsset(tokenOutAddr)
	if err != nil {
		return err
	}
	inAsset, err := getAsset(tokenInAddr)
	if err != nil {
		return err
	}

	pair := CurvePair{
		OutAsset: outAsset,
		InAsset:  inAsset,
		OutIndex: outIdx,
		InIndex:  inIdx,
		Address:  addr,
	}

	lock.Lock()
	defer lock.Unlock()
	s.poolMap[addr] = append(s.poolMap[addr], pair)
	return nil
}

// StartStream: Base will call this when starting a stream for an address
func (h *curveHooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	h.s.watchSwaps(ctx, base, addr, tradesChan, lock)
}

// OnPoolRemoved: clear poolMap when an address is removed
func (h *curveHooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addr common.Address,
	lock *sync.RWMutex,
) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()
}

// OnOrderChanged:
// Base thinks "the Order of this address has changed", we do it simply and roughly:
//   - clear the poolMap for this address
//   - re-ensure all Orders from the latest config
func (h *curveHooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
	lock *sync.RWMutex,
) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()

	_, err := h.s.ensureAllOrdersForAddress(context.Background(), base, addr, branchMarketConfig, lock)
	if err != nil {
		log.Errorf("Curve - OnOrderChanged ensureAllOrdersForAddress %s: %v", addr.Hex(), err)
	}
}

// NewCurveScraper: Curve scraper entry, attached to BaseDEXScraper
func NewCurveScraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) {
	if wg != nil {
		wg.Done()
	}

	scraper := &CurveScraper{
		poolMap:    make(map[common.Address][]CurvePair),
		exchange:   exchangeName,
		blockchain: blockchain,
	}

	hooks := &curveHooks{s: scraper}

	//let BaseDEXScraper manage: rest/ws client, watchdog, resub/unsub, config diff by address level
	scraper.base = NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, branchMarketConfig, nil)
}

// ensureAllOrdersForAddress:
//   - find all (Address,Order) for the same address from config
//   - call ensurePairInPoolMap for each
//   - return the maximum WatchDogDelay for this address
func (s *CurveScraper) ensureAllOrdersForAddress(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	branchMarketConfig string,
	lock *sync.RWMutex,
) (maxDelay int64, err error) {
	ex := s.exchange

	pools, err := models.PoolsFromConfigFile(ex, branchMarketConfig)
	if err != nil {
		return 0, fmt.Errorf("load pools for %s: %w", ex, err)
	}

	target := strings.ToLower(addr.Hex())
	found := false

	for _, p := range pools {
		if strings.ToLower(p.Address) != target {
			continue
		}
		found = true

		if e := s.ensurePairInPoolMap(ctx, base, p, lock); e != nil {
			log.Errorf("Curve - ensurePairInPoolMap %s#%d: %v", p.Address, p.Order, e)
			continue
		}

		d := p.WatchDogDelay
		if d <= 0 {
			d = 300
		}
		if d > maxDelay {
			maxDelay = d
		}
	}

	if !found {
		return 0, fmt.Errorf("no config entries for pool %s", addr.Hex())
	}
	if maxDelay <= 0 {
		maxDelay = 300
	}
	return maxDelay, nil
}

// ensurePairInPoolMap: ensure the CurvePair corresponding to a (Address,Order) is already in poolMap
func (s *CurveScraper) ensurePairInPoolMap(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) error {
	outIdx, inIdx, err := parseIndexCode(pool.Order)
	if err != nil {
		return err
	}
	addr := common.HexToAddress(pool.Address)

	lock.RLock()
	existing := s.poolMap[addr]
	lock.RUnlock()
	for _, p := range existing {
		if p.InIndex == inIdx && p.OutIndex == outIdx {
			return nil
		}
	}

	client := base.RESTClient()
	if client == nil {
		return fmt.Errorf("nil REST client for Curve")
	}

	tokenOutAddr, err := coinAddressFromPool(ctx, client, addr, outIdx)
	if err != nil || tokenOutAddr == (common.Address{}) {
		return fmt.Errorf("get out-coin addr: %w", err)
	}
	tokenInAddr, err := coinAddressFromPool(ctx, client, addr, inIdx)
	if err != nil || tokenInAddr == (common.Address{}) {
		return fmt.Errorf("get in-coin addr: %w", err)
	}

	getAsset := func(a common.Address) (models.Asset, error) {
		if isNative(a) {
			return nativeAsset(s.blockchain), nil
		}
		return models.GetAsset(a, s.blockchain, client)
	}

	outAsset, err := getAsset(tokenOutAddr)
	if err != nil {
		return err
	}
	inAsset, err := getAsset(tokenInAddr)
	if err != nil {
		return err
	}

	pair := CurvePair{
		OutAsset: outAsset,
		InAsset:  inAsset,
		OutIndex: outIdx,
		InIndex:  inIdx,
		Address:  addr,
	}

	lock.Lock()
	s.poolMap[addr] = append(s.poolMap[addr], pair)
	lock.Unlock()

	return nil
}

// Swap subscription & processing logic

type rawSwapData struct {
	addr     common.Address
	soldID   int
	boughtID int
	sold     *big.Int
	bought   *big.Int
	swapID   string
}

func (s *CurveScraper) extractSwapData(ev interface{}) (rawSwapData, error) {
	switch e := ev.(type) {
	case curvepool.CurvepoolTokenExchange:
		return rawSwapData{
			addr:     e.Raw.Address,
			soldID:   int(e.SoldId.Int64()),
			boughtID: int(e.BoughtId.Int64()),
			sold:     e.TokensSold,
			bought:   e.TokensBought,
			swapID:   e.Raw.TxHash.Hex(),
		}, nil
	case curvefifactory.CurvefifactoryTokenExchange:
		return rawSwapData{
			addr:     e.Raw.Address,
			soldID:   int(e.SoldId.Int64()),
			boughtID: int(e.BoughtId.Int64()),
			sold:     e.TokensSold,
			bought:   e.TokensBought,
			swapID:   e.Raw.TxHash.Hex(),
		}, nil
	case curvefitwocryptooptimized.CurvefitwocryptooptimizedTokenExchange:
		return rawSwapData{
			addr:     e.Raw.Address,
			soldID:   int(e.SoldId.Int64()),
			boughtID: int(e.BoughtId.Int64()),
			sold:     e.TokensSold,
			bought:   e.TokensBought,
			swapID:   e.Raw.TxHash.Hex(),
		}, nil
	case curvepool.CurvepoolTokenExchangeUnderlying:
		return rawSwapData{
			addr:     e.Raw.Address,
			soldID:   int(e.SoldId.Int64()),
			boughtID: int(e.BoughtId.Int64()),
			sold:     e.TokensSold,
			bought:   e.TokensBought,
			swapID:   e.Raw.TxHash.Hex(),
		}, nil
	default:
		return rawSwapData{}, fmt.Errorf("unknown swap type")
	}
}

type curveFeeds struct {
	Sink           chan *curvepool.CurvepoolTokenExchange
	Sub            event.Subscription
	factorySink    chan *curvefifactory.CurvefifactoryTokenExchange
	factorySub     event.Subscription
	twoSink        chan *curvefitwocryptooptimized.CurvefitwocryptooptimizedTokenExchange
	twoSub         event.Subscription
	underlyingSink chan *curvepool.CurvepoolTokenExchangeUnderlying
	underlyingSub  event.Subscription
}

func (s *CurveScraper) GetSwapsChannel(
	base *BaseDEXScraper,
	address common.Address,
) (*curveFeeds, error) {
	wsClient := base.WSClient()
	if wsClient == nil {
		return nil, fmt.Errorf("nil WS client for Curve")
	}

	filterer, err := curvepool.NewCurvepoolFilterer(address, wsClient)
	if err != nil {
		return nil, fmt.Errorf("curvepool filterer: %v", err)
	}

	curvefiFactoryFilterer, err := curvefifactory.NewCurvefifactoryFilterer(address, wsClient)
	if err != nil {
		return nil, fmt.Errorf("curvefifactory filterer: %v", err)
	}

	filtererTwoCrypto, err := curvefitwocryptooptimized.NewCurvefitwocryptooptimizedFilterer(address, wsClient)
	if err != nil {
		return nil, fmt.Errorf("curvefitwocryptooptimized filterer: %v", err)
	}

	sinks := &curveFeeds{
		Sink:           make(chan *curvepool.CurvepoolTokenExchange),
		factorySink:    make(chan *curvefifactory.CurvefifactoryTokenExchange),
		twoSink:        make(chan *curvefitwocryptooptimized.CurvefitwocryptooptimizedTokenExchange),
		underlyingSink: make(chan *curvepool.CurvepoolTokenExchangeUnderlying),
	}

	var start *uint64
	// Backfill 20 blocks (~5 minutes on Ethereum mainnet) to capture swaps that occurred during websocket reconnection
	const blocksToBackfill = uint64(20)

	restClient := base.RESTClient()
	if restClient == nil {
		log.Warn("Curve - nil REST client, subscribe from latest (no backfill).")
	} else if hdr, err := restClient.HeaderByNumber(context.Background(), nil); err != nil {
		log.Warnf("Curve - header fetch failed (%v), subscribe from latest (no backfill).", err)
	} else if n := hdr.Number.Uint64(); n > blocksToBackfill {
		sb := n - blocksToBackfill
		start = &sb
	}

	okCount := 0

	if sub, err := filterer.WatchTokenExchange(&bind.WatchOpts{Start: start}, sinks.Sink, nil); err == nil {
		sinks.Sub = sub
		okCount++
	}
	if curvefiFactorySub, err := curvefiFactoryFilterer.WatchTokenExchange(&bind.WatchOpts{Start: start}, sinks.factorySink, nil); err == nil {
		sinks.factorySub = curvefiFactorySub
		okCount++
	}
	if subTwoCrypto, err := filtererTwoCrypto.WatchTokenExchange(&bind.WatchOpts{Start: start}, sinks.twoSink, nil); err == nil {
		sinks.twoSub = subTwoCrypto
		okCount++
	}
	if subUnderlying, err := filterer.WatchTokenExchangeUnderlying(&bind.WatchOpts{Start: start}, sinks.underlyingSink, nil); err == nil {
		sinks.underlyingSub = subUnderlying
		okCount++
	}

	if okCount == 0 {
		return nil, fmt.Errorf("failed to establish any subscriptions")
	}

	return sinks, nil
}

// watch swaps for a given address
func (s *CurveScraper) watchSwaps(
	ctx context.Context,
	base *BaseDEXScraper,
	address common.Address,
	tradesChannel chan models.Trade,
	lock *sync.RWMutex,
) {
	feeds, err := s.GetSwapsChannel(base, address)
	if err != nil {
		log.Error("Curve - failed to get swaps channel: ", err)
		return
	}

	var (
		subErr, factorySubErr, twoSubErr, underlyingSubErr <-chan error
	)

	if feeds.Sub != nil {
		subErr = feeds.Sub.Err()
	}
	if feeds.factorySub != nil {
		factorySubErr = feeds.factorySub.Err()
	}
	if feeds.twoSub != nil {
		twoSubErr = feeds.twoSub.Err()
	}
	if feeds.underlyingSub != nil {
		underlyingSubErr = feeds.underlyingSub.Err()
	}

	go func() {
		defer func() {
			if feeds.Sub != nil {
				feeds.Sub.Unsubscribe()
			}
			if feeds.factorySub != nil {
				feeds.factorySub.Unsubscribe()
			}
			if feeds.twoSub != nil {
				feeds.twoSub.Unsubscribe()
			}
			if feeds.underlyingSub != nil {
				feeds.underlyingSub.Unsubscribe()
			}
		}()

		for {
			select {
			case rawSwap, ok := <-feeds.Sink:
				if ok {
					s.handleRawSwap(*rawSwap, base, address, lock, tradesChannel)
				} else {
					base.SubscribeChannel() <- address
					return
				}

			case rawSwapCurvefiFactory, ok := <-feeds.factorySink:
				if ok {
					s.handleRawSwap(*rawSwapCurvefiFactory, base, address, lock, tradesChannel)
				} else {
					base.SubscribeChannel() <- address
					return
				}

			case rawSwapTwoCrypto, ok := <-feeds.twoSink:
				if ok {
					s.handleRawSwap(*rawSwapTwoCrypto, base, address, lock, tradesChannel)
				} else {
					base.SubscribeChannel() <- address
					return
				}

			case rawSwapUnderlying, ok := <-feeds.underlyingSink:
				if ok {
					s.handleRawSwap(*rawSwapUnderlying, base, address, lock, tradesChannel)
				} else {
					base.SubscribeChannel() <- address
					return
				}

			case err := <-subErr:
				log.Errorf("Curve - subscription error for pool %s: %v", address.Hex(), err)
				base.SubscribeChannel() <- address
				return
			case err := <-factorySubErr:
				log.Errorf("Curve - factory subscription error for pool %s: %v", address.Hex(), err)
				base.SubscribeChannel() <- address
				return
			case err := <-twoSubErr:
				log.Errorf("Curve - twoCrypto subscription error for pool %s: %v", address.Hex(), err)
				base.SubscribeChannel() <- address
				return
			case err := <-underlyingSubErr:
				log.Errorf("Curve - underlying subscription error for pool %s: %v", address.Hex(), err)
				base.SubscribeChannel() <- address
				return

			case <-ctx.Done():
				log.Infof("Curve - Shutting down watchSwaps for %s", address.Hex())
				return
			}
		}
	}()
}

// handle different event types uniformly
func (s *CurveScraper) handleRawSwap(
	ev interface{},
	base *BaseDEXScraper,
	address common.Address,
	lock *sync.RWMutex,
	tradesChannel chan models.Trade,
) {
	swapData, err := s.extractSwapData(ev)
	if err != nil {
		log.Error("Curve - error normalizing swap: ", err)
		return
	}

	lock.RLock()
	pairs := s.poolMap[swapData.addr]
	lock.RUnlock()

	if len(pairs) == 0 {
		return
	}

	pair, ok := findPairByIdx(pairs, swapData.soldID, swapData.boughtID)
	if !ok {
		return
	}

	log.Infof("Curve - subscribe to %s with pair %s",
		address.Hex(), pair.OutAsset.Symbol+"-"+pair.InAsset.Symbol,
	)

	var decSold, decBought int
	switch swapData.soldID {
	case pair.InIndex:
		decSold = int(pair.InAsset.Decimals)
	case pair.OutIndex:
		decSold = int(pair.OutAsset.Decimals)
	}
	switch swapData.boughtID {
	case pair.OutIndex:
		decBought = int(pair.OutAsset.Decimals)
	case pair.InIndex:
		decBought = int(pair.InAsset.Decimals)
	}

	soldAmt, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swapData.sold),
		new(big.Float).SetFloat64(math.Pow10(decSold)),
	).Float64()
	boughtAmt, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swapData.bought),
		new(big.Float).SetFloat64(math.Pow10(decBought)),
	).Float64()

	normalizedSwap := CurveSwap{
		ID:        swapData.swapID,
		Timestamp: time.Now().Unix(),
		Pair:      pair,
		Amount0:   soldAmt,
		Amount1:   boughtAmt,
	}

	s.processSwap(normalizedSwap, base, address, lock, tradesChannel)
}

// only process swaps in the InIndex -> OutIndex direction
func (s *CurveScraper) processSwap(
	swap CurveSwap,
	base *BaseDEXScraper,
	address common.Address,
	lock *sync.RWMutex,
	tradesChannel chan models.Trade,
) {
	price, volume := getSwapDataCurve(swap)
	if volume == 0 || math.IsNaN(volume) || math.IsInf(volume, 0) ||
		price == 0 || math.IsNaN(price) || math.IsInf(price, 0) {
		log.Debugf("Curve - skip invalid price/volume, tx=%s", swap.ID)
		return
	}

	t := makeCurveTrade(
		swap.Pair,
		price,
		volume,
		time.Unix(swap.Timestamp, 0),
		address,
		swap.ID,
		s.exchange,
		s.blockchain,
	)

	// update lastTradeTime using Base (one of the key points attached to base)
	base.UpdateLastTradeTime(address, t.Time, lock)

	tradesChannel <- t
	logTrade(t)
}

func getSwapDataCurve(swap CurveSwap) (price float64, volume float64) {
	volume = swap.Amount1
	if volume == 0 || math.IsNaN(volume) || math.IsInf(volume, 0) {
		return math.NaN(), volume
	}
	price = swap.Amount0 / swap.Amount1
	if price == 0 || math.IsNaN(price) || math.IsInf(price, 0) {
		return math.NaN(), volume
	}
	return
}

func makeCurveTrade(
	pair CurvePair,
	price float64,
	volume float64,
	timestamp time.Time,
	address common.Address,
	foreignTradeID string,
	exchangeName string,
	blockchain string,
) models.Trade {
	token0 := pair.OutAsset
	token1 := pair.InAsset
	return models.Trade{
		Price:          price,
		Volume:         volume,
		BaseToken:      token1,
		QuoteToken:     token0,
		Time:           timestamp,
		PoolAddress:    address.Hex(),
		ForeignTradeID: foreignTradeID,
		Exchange:       models.Exchange{Name: exchangeName, Blockchain: blockchain},
	}
}

// Order is two digits, e.g. 01 / 10 / 12 / 21
func parseIndexCode(code int) (outIdx, inIdx int, err error) {
	s := fmt.Sprintf("%02d", code)
	if len(s) != 2 {
		return 0, 0, fmt.Errorf("index code must be 2 digits, got: %s", s)
	}

	outIdx = int(s[0] - '0')
	inIdx = int(s[1] - '0')
	if outIdx < 0 || inIdx < 0 || outIdx == inIdx {
		return 0, 0, fmt.Errorf("invalid index code: %s", s)
	}
	return
}

// find the pair where InIndex->OutIndex == soldID->boughtID
func findPairByIdx(
	pairs []CurvePair,
	soldID int,
	boughtID int,
) (CurvePair, bool) {
	for _, pair := range pairs {
		if pair.InIndex == soldID && pair.OutIndex == boughtID {
			return pair, true
		}
	}
	return CurvePair{}, false
}

// try to resolve coin address from pool+idx
func coinAddressFromPool(
	ctx context.Context,
	client *ethclient.Client,
	pool common.Address,
	idx int,
) (common.Address, error) {
	bi := big.NewInt(int64(idx))

	// 1) common curvepool
	if c, err := curvepool.NewCurvepool(pool, client); err == nil {
		if addr, err := c.Coins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
		if addr, err := c.UnderlyingCoins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
	}

	// 2) TwoCrypto optimized pool
	if c2, err := curvefitwocryptooptimized.NewCurvefitwocryptooptimized(pool, client); err == nil {
		if addr, err := c2.Coins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
	}

	return common.Address{}, fmt.Errorf("cannot resolve coin at index %d for pool %s", idx, pool.Hex())
}

func isNative(addr common.Address) bool {
	return strings.EqualFold(addr.Hex(), NativeETHSentinel)
}

func nativeAsset(chain string) models.Asset {
	return models.Asset{
		Symbol:     "ETH",
		Name:       "Ether",
		Address:    NativeETHSentinel,
		Decimals:   18,
		Blockchain: chain,
	}
}
