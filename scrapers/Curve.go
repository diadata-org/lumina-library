package scrapers

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	curvefitwocryptooptimized "github.com/diadata-org/lumina-library/contracts/curve/curvefitwocrypto"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

const NativeETHSentinel = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

var (
	restDialCurve = ""
	wsDialCurve   = ""
)

type CurvePair struct {
	OutAsset    models.Asset
	InAsset     models.Asset
	OutIndex    int
	InIndex     int
	ForeignName string
	Address     common.Address
}

type CurveSwap struct {
	ID        string
	Timestamp int64
	Pair      CurvePair
	Amount0   float64
	Amount1   float64
}

type CurveScraper struct {
	exchange           models.Exchange
	poolMap            map[common.Address][]CurvePair
	subscribeChannel   chan common.Address
	unsubscribeChannel chan common.Address
	watchdogCancel     map[string]context.CancelFunc
	swapStreamCancel   map[common.Address]context.CancelFunc
	lastTradeTimeMap   map[common.Address]time.Time
	restClient         *ethclient.Client
	wsClient           *ethclient.Client
	waitTime           int
}

func NewCurveScraper(ctx context.Context, exchangeName string, blockchain string, pools []models.Pool, tradesChannel chan models.Trade, wg *sync.WaitGroup) {
	var err error
	var scraper CurveScraper
	log.Info("Started Curve scraper.")

	scraper.exchange.Name = exchangeName
	scraper.exchange.Blockchain = blockchain
	scraper.subscribeChannel = make(chan common.Address)
	scraper.unsubscribeChannel = make(chan common.Address)
	scraper.watchdogCancel = make(map[string]context.CancelFunc)
	scraper.swapStreamCancel = make(map[common.Address]context.CancelFunc)
	scraper.lastTradeTimeMap = make(map[common.Address]time.Time)
	scraper.waitTime, err = strconv.Atoi(utils.Getenv(strings.ToUpper(exchangeName)+"_WAIT_TIME", "500"))
	if err != nil {
		log.Error("parse waitTime: ", err)
	}

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(CURVE_EXCHANGE)+"_URI_REST", restDialCurve))
	if err != nil {
		log.Error("Curve - init rest client: ", err)
	}
	scraper.wsClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(CURVE_EXCHANGE)+"_URI_WS", wsDialCurve))
	if err != nil {
		log.Error("Curve - init ws client: ", err)
	}

	if err := scraper.makePoolMap(pools); err != nil {
		log.Fatalf("Curve - makePoolMap failed: %v", err)
	}

	var lock sync.RWMutex

	// resubscribe handler
	go scraper.startResubHandler(ctx, tradesChannel, &lock)

	// unsubscribe handler
	go scraper.startUnsubHandler(ctx, &lock)

	// watch config updates
	go scraper.watchConfig(ctx, exchangeName, tradesChannel, &lock)
}

func (scraper *CurveScraper) startResubHandler(ctx context.Context, trades chan models.Trade, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-scraper.subscribeChannel:
				lock.Lock()
				scraper.lastTradeTimeMap[addr] = time.Now()
				// if already running, restart the stream
				if cancel, ok := scraper.swapStreamCancel[addr]; ok && cancel != nil {
					cancel()
					delete(scraper.swapStreamCancel, addr)
				}
				lock.Unlock()
				go scraper.watchSwaps(ctx, addr, trades, lock)
			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				return
			}
		}
	}()
}

func (scraper *CurveScraper) startUnsubHandler(ctx context.Context, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-scraper.unsubscribeChannel:
				scraper.stopPool(addr, lock)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// start the pool, ensure the current order is in the poolMap[addr], initialize lastTradeTime, restart watchdog, start watchSwaps
func (s *CurveScraper) startPool(ctx context.Context, pool models.Pool, trades chan models.Trade, lock *sync.RWMutex) error {
	addr := common.HexToAddress(pool.Address)

	// 1) ensure the current order is in the poolMap[addr], will automatically add token information and append
	if err := s.ensurePairInPoolMap(pool); err != nil {
		return fmt.Errorf("ensurePairInPoolMap(%s): %w", pool.Address, err)
	}

	// 2) initialize lastTradeTime, avoid triggering watchdog immediately
	lock.Lock()
	if _, ok := s.lastTradeTimeMap[addr]; !ok {
		s.lastTradeTimeMap[addr] = time.Now()
	}
	lock.Unlock()

	// 3) restart watchdog (use the maximum WatchDogDelay when there are multiple orders)
	delay := pool.WatchDogDelay

	s.restartWatchdogForAddr(ctx, addr, delay, lock)

	// 4) start watchSwaps (if already running, don't repeat)
	if _, ok := s.swapStreamCancel[addr]; ok {
		return nil
	}
	pctx, cancel := context.WithCancel(ctx)
	s.swapStreamCancel[addr] = cancel
	go s.watchSwaps(pctx, addr, trades, lock)
	return nil
}

// stop the pool, stop the stream and watchdog, clean the state
func (s *CurveScraper) stopPool(addr common.Address, lock *sync.RWMutex) {
	if c, ok := s.swapStreamCancel[addr]; ok && c != nil {
		c()
		delete(s.swapStreamCancel, addr)
	}
	key := addr.Hex()
	if c, ok := s.watchdogCancel[key]; ok && c != nil {
		c()
		delete(s.watchdogCancel, key)
	}
	lock.Lock()
	delete(s.lastTradeTimeMap, addr)
	lock.Unlock()
}

func (s *CurveScraper) restartWatchdogForAddr(ctx context.Context, addr common.Address, delay int64, lock *sync.RWMutex) {
	key := addr.Hex()
	if c, ok := s.watchdogCancel[key]; ok && c != nil {
		c()
		delete(s.watchdogCancel, key)
	}
	wdCtx, cancel := context.WithCancel(ctx)
	s.watchdogCancel[key] = cancel

	t := time.NewTicker(time.Duration(delay) * time.Second)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				lock.RLock()
				last := s.lastTradeTimeMap[addr]
				lock.RUnlock()
				if time.Since(last) >= time.Duration(delay)*time.Second {
					// trigger resubscribe
					s.subscribeChannel <- addr
				}
			case <-wdCtx.Done():
				return
			}
		}
	}()
}

// ensure the current order is in the poolMap[addr], will automatically add token information and append
func (s *CurveScraper) ensurePairInPoolMap(pool models.Pool) error {
	if s.poolMap == nil {
		s.poolMap = make(map[common.Address][]CurvePair)
	}

	outIdx, inIdx, err := parseIndexCode(pool.Order)
	if err != nil {
		return err
	}

	addr := common.HexToAddress(pool.Address)

	// if already exists, don't repeat
	for _, p := range s.poolMap[addr] {
		if p.InIndex == inIdx && p.OutIndex == outIdx {
			// if already exists, it's ready
			return nil
		}
	}

	// add token information and append
	ctx := context.Background()
	tokenOutAddr, err := coinAddressFromPool(ctx, s.restClient, addr, outIdx)
	if err != nil || tokenOutAddr == (common.Address{}) {
		return fmt.Errorf("get out-coin addr: %w", err)
	}
	tokenInAddr, err := coinAddressFromPool(ctx, s.restClient, addr, inIdx)
	if err != nil || tokenInAddr == (common.Address{}) {
		return fmt.Errorf("get in-coin addr: %w", err)
	}

	getAsset := func(a common.Address) (models.Asset, error) {
		if isNative(a) {
			return nativeAsset(s.exchange.Blockchain), nil
		}
		return models.GetAsset(a, s.exchange.Blockchain, s.restClient)
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
		OutAsset: outAsset, InAsset: inAsset,
		OutIndex: outIdx, InIndex: inIdx,
		Address: addr,
	}
	s.poolMap[addr] = append(s.poolMap[addr], pair)
	return nil
}

// flatten the slice to map[key]Pool; key = lower(addr) + "#" + Order(original string)
func flatByAddrOrder(pools []models.Pool) map[string]models.Pool {
	m := make(map[string]models.Pool, len(pools))
	for _, p := range pools {
		k := strings.ToLower(p.Address) + "#" + strconv.Itoa(p.Order)
		m[k] = p
	}
	return m
}

// calculate the effective watchdog for each address based on the flat map (max)
func aggWatchdogByAddr(m map[string]models.Pool) map[string]int64 {
	out := make(map[string]int64)
	for k, p := range m {
		addrLower := strings.SplitN(k, "#", 2)[0]
		d := p.WatchDogDelay
		if d <= 0 {
			d = 300
		}
		if d > out[addrLower] {
			out[addrLower] = d
		}
	}
	return out
}

func (s *CurveScraper) watchConfig(ctx context.Context, exchangeName string, trades chan models.Trade, lock *sync.RWMutex) {
	// Check for config changes every 60 minutes.
	envKey := strings.ToUpper(CURVE_EXCHANGE) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "3600"))
	if err != nil {
		log.Errorf("Curve - Failed to parse %s: %v.", envKey, err)
		return
	}
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	// load initial pools
	curList, err := models.PoolsFromConfigFile(exchangeName)
	if err != nil {
		log.Errorf("Curve - load initial pools: %v", err)
		curList = nil
	}
	lastFlat := flatByAddrOrder(curList)
	lastDelay := aggWatchdogByAddr(lastFlat)

	// start existing addresses (one time for each address; other orders will be added in ensurePairInPoolMap)
	s.startAllFromFlat(ctx, lastFlat, trades, lock)

	for {
		select {
		case <-t.C:
			nowList, err := models.PoolsFromConfigFile(exchangeName)
			if err != nil {
				log.Errorf("Curve - reload pools: %v", err)
				continue
			}
			nowFlat := flatByAddrOrder(nowList)
			nowDelay := aggWatchdogByAddr(nowFlat)

			s.applyConfigDiff(ctx, lastFlat, nowFlat, lastDelay, nowDelay, trades, lock)
			lastFlat = nowFlat
			lastDelay = nowDelay

		case <-ctx.Done():
			log.Info("Curve - watchConfig exit")
			return
		}
	}
}

func (s *CurveScraper) startAllFromFlat(ctx context.Context, flat map[string]models.Pool, trades chan models.Trade, lock *sync.RWMutex) {
	// first merge by address, start each address once (internal will ensurePairInPoolMap append all orders)
	seenAddr := make(map[string]bool)
	for _, p := range flat {
		addrLower := strings.ToLower(p.Address)
		if seenAddr[addrLower] {
			continue
		}
		// start this address (will ensure the current order, and start the stream & watchdog)
		if err := s.startPool(ctx, p, trades, lock); err != nil {
			log.Errorf("Curve - startPool %s: %v", p.Address, err)
		}
		seenAddr[addrLower] = true
	}
	// then ensure all other orders for the same address are appended to the poolMap
	for _, p := range flat {
		if err := s.ensurePairInPoolMap(p); err != nil {
			log.Errorf("Curve - ensurePairInPoolMap %s#%v: %v", p.Address, p.Order, err)
		}
	}
}

func (s *CurveScraper) applyConfigDiff(
	ctx context.Context,
	lastFlat map[string]models.Pool,
	currFlat map[string]models.Pool,
	lastDelay map[string]int64, // address -> last effective watchdog
	currDelay map[string]int64, // address -> current effective watchdog
	trades chan models.Trade,
	lock *sync.RWMutex,
) {
	// ---- 1) handle deletion: find removed (addr,order) keys ----
	removed := make([]string, 0)
	for k := range lastFlat {
		if _, ok := currFlat[k]; !ok {
			removed = append(removed, k)
		}
	}
	// remove these pairs from the poolMap; if an address has no pairs anymore, stop it
	toStop := make(map[string]bool)
	for _, key := range removed {
		parts := strings.SplitN(key, "#", 2)
		addrLower, order := parts[0], parts[1]
		addr := common.HexToAddress(addrLower)

		// remove the pair corresponding to this order from the poolMap[addr]
		outIdx, inIdx, _ := parseIndexCode(mustAtoi(order)) // order is a two-digit string like "10"
		pairs := s.poolMap[addr]
		filtered := make([]CurvePair, 0, len(pairs))
		for _, pr := range pairs {
			if !(pr.InIndex == inIdx && pr.OutIndex == outIdx) {
				filtered = append(filtered, pr)
			}
		}
		s.poolMap[addr] = filtered

		if len(filtered) == 0 {
			toStop[addrLower] = true
		}
	}
	for addrLower := range toStop {
		addr := common.HexToAddress(addrLower)
		log.Infof("Curve - stop pool %s (all orders removed)", addr.Hex())
		s.stopPool(addr, lock)
	}

	// ---- 2) handle addition: find new (addr,order) keys, ensure pair in poolMap, and start the address if necessary ----
	// first record which addresses are "new" (no entries before)
	newlyAddress := make(map[string]bool)
	for k, p := range currFlat {
		if _, ok := lastFlat[k]; ok {
			// not a new entry
			continue
		}
		// new entry: first append the pair to the poolMap
		if err := s.ensurePairInPoolMap(p); err != nil {
			log.Errorf("Curve - ensurePairInPoolMap(add) %s#%v: %v", p.Address, p.Order, err)
			continue
		}
		addrLower := strings.ToLower(p.Address)
		// if this address is completely not in last, need to start the stream/watchdog for it
		if !hasAnyAddr(lastFlat, addrLower) {
			newlyAddress[addrLower] = true
		}
	}
	// for new addresses: start the stream/watchdog
	for addrLower := range newlyAddress {
		// find any order for this address in currFlat, use it to startPool (startPool will restart watchdog and start stream)
		var pick *models.Pool
		for k, p := range currFlat {
			if strings.HasPrefix(k, addrLower+"#") {
				pick = &p
				break
			}
		}
		if pick != nil {
			if err := s.startPool(ctx, *pick, trades, lock); err != nil {
				log.Errorf("Curve - startPool(new addr) %s: %v", pick.Address, err)
			} else {
				log.Infof("Curve - added new pool %s", pick.Address)
			}
		}
	}

	// ---- 3) handle watchdog updates: compare lastDelay vs currDelay for each address, restart watchdog if different ----
	for addrLower, newD := range currDelay {
		oldD := lastDelay[addrLower]
		if newD <= 0 {
			newD = 300
		}
		if oldD <= 0 {
			oldD = 300
		}
		if newD != oldD {
			addr := common.HexToAddress(addrLower)
			log.Infof("Curve - update watchdog %s: %d -> %d", addr.Hex(), oldD, newD)
			s.restartWatchdogForAddr(ctx, addr, newD, lock)
		}
	}
}

func hasAnyAddr(flat map[string]models.Pool, addrLower string) bool {
	prefix := addrLower + "#"
	for k := range flat {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// only used for parseIndexCode (Order is a two-digit string like "10")
func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func getSwapDataCurve(swap CurveSwap) (price float64, volume float64) {
	volume = swap.Amount1
	price = swap.Amount0 / swap.Amount1
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

func findPairByIdx(pairs []CurvePair, soldID int, boughtID int) (CurvePair, bool) {
	for _, pair := range pairs {
		if pair.InIndex == soldID && pair.OutIndex == boughtID {
			return pair, true
		}
	}
	return CurvePair{}, false
}

func (scraper *CurveScraper) watchSwaps(ctx context.Context, address common.Address, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	feeds, err := scraper.GetSwapsChannel(address)
	if err != nil {
		log.Error("Curve - failed to get swaps channel: ", err)
		return
	}

	var (
		subErr, factorySubErr, twoSubErr, underlyingSubErr <-chan error
	)

	// first check if the subscription is nil
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
					swap, err := scraper.extractSwapData(*rawSwap)
					if err != nil {
						log.Error("Curve - error normalizing swap: ", err)
						continue
					}
					pairs := scraper.poolMap[swap.addr]
					if len(pairs) == 0 {
						log.Warnf("Curve - Sink - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Warnf("Curve - Sink - no pair found for address %s", swap.addr.Hex())
						continue
					}
					log.Infof("Curve - subscribe to %s with pair %s", address.Hex(), pair.OutAsset.Symbol+"-"+pair.InAsset.Symbol)
					var decSold, decBought int
					switch swap.soldID {
					case pair.InIndex:
						decSold = int(pair.InAsset.Decimals)
					case pair.OutIndex:
						decSold = int(pair.OutAsset.Decimals)
					}
					switch swap.boughtID {
					case pair.OutIndex:
						decBought = int(pair.OutAsset.Decimals)
					case pair.InIndex:
						decBought = int(pair.InAsset.Decimals)
					}
					soldAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.sold), new(big.Float).SetFloat64(math.Pow10(decSold))).Float64()
					boughtAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.bought), new(big.Float).SetFloat64(math.Pow10(decBought))).Float64()
					normalizedSwap := CurveSwap{
						ID:        swap.swapID,
						Timestamp: time.Now().Unix(),
						Pair:      pair,
						Amount0:   soldAmt,
						Amount1:   boughtAmt,
					}
					scraper.processSwap(normalizedSwap, pairs, address, lock, tradesChannel)
				} else {
					scraper.subscribeChannel <- address
					return
				}
			case rawSwapCurvefiFactory, ok := <-feeds.factorySink:
				if ok {
					swap, err := scraper.extractSwapData(*rawSwapCurvefiFactory)
					if err != nil {
						log.Error("Curve - error normalizing swap: ", err)
						continue
					}
					pairs := scraper.poolMap[swap.addr]
					if len(pairs) == 0 {
						log.Warnf("Curve - factorySink - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Warnf("Curve - factorySink - no pair found for address %s", swap.addr.Hex())
						continue
					}
					log.Infof("Curve - subscribe to %s with pair %s", address.Hex(), pair.OutAsset.Symbol+"-"+pair.InAsset.Symbol)
					var decSold, decBought int
					switch swap.soldID {
					case pair.InIndex:
						decSold = int(pair.InAsset.Decimals)
					case pair.OutIndex:
						decSold = int(pair.OutAsset.Decimals)
					}
					switch swap.boughtID {
					case pair.OutIndex:
						decBought = int(pair.OutAsset.Decimals)
					case pair.InIndex:
						decBought = int(pair.InAsset.Decimals)
					}
					soldAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.sold), new(big.Float).SetFloat64(math.Pow10(decSold))).Float64()
					boughtAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.bought), new(big.Float).SetFloat64(math.Pow10(decBought))).Float64()
					normalizedSwap := CurveSwap{
						ID:        swap.swapID,
						Timestamp: time.Now().Unix(),
						Pair:      pair,
						Amount0:   soldAmt,
						Amount1:   boughtAmt,
					}
					scraper.processSwap(normalizedSwap, pairs, address, lock, tradesChannel)
				} else {
					scraper.subscribeChannel <- address
					return
				}
			case rawSwapTwoCrypto, ok := <-feeds.twoSink:
				if ok {
					swap, err := scraper.extractSwapData(*rawSwapTwoCrypto)
					if err != nil {
						log.Error("Curve - error normalizing swap: ", err)
						continue
					}
					pairs := scraper.poolMap[swap.addr]
					if len(pairs) == 0 {
						log.Warnf("Curve - twoSink - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Warnf("Curve - twoSink - no pair found for address %s", swap.addr.Hex())
						continue
					}
					log.Infof("Curve - subscribe to %s with pair %s", address.Hex(), pair.OutAsset.Symbol+"-"+pair.InAsset.Symbol)
					var decSold, decBought int
					switch swap.soldID {
					case pair.InIndex:
						decSold = int(pair.InAsset.Decimals)
					case pair.OutIndex:
						decSold = int(pair.OutAsset.Decimals)
					}
					switch swap.boughtID {
					case pair.OutIndex:
						decBought = int(pair.OutAsset.Decimals)
					case pair.InIndex:
						decBought = int(pair.InAsset.Decimals)
					}
					soldAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.sold), new(big.Float).SetFloat64(math.Pow10(decSold))).Float64()
					boughtAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.bought), new(big.Float).SetFloat64(math.Pow10(decBought))).Float64()
					normalizedSwap := CurveSwap{
						ID:        swap.swapID,
						Timestamp: time.Now().Unix(),
						Pair:      pair,
						Amount0:   soldAmt,
						Amount1:   boughtAmt,
					}
					scraper.processSwap(normalizedSwap, pairs, address, lock, tradesChannel)
				} else {
					scraper.subscribeChannel <- address
					return
				}
			case rawSwapUnderlying, ok := <-feeds.underlyingSink:
				if ok {
					swap, err := scraper.extractSwapData(*rawSwapUnderlying)
					if err != nil {
						log.Error("Curve - error normalizing swap: ", err)
						continue
					}
					pairs := scraper.poolMap[swap.addr]
					if len(pairs) == 0 {
						log.Warnf("Curve - underlyingSink - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Warnf("Curve - underlyingSink - no pair found for address %s", swap.addr.Hex())
						continue
					}
					log.Infof("Curve - subscribe to %s with pair %s", address.Hex(), pair.OutAsset.Symbol+"-"+pair.InAsset.Symbol)
					var decSold, decBought int
					switch swap.soldID {
					case pair.InIndex:
						decSold = int(pair.InAsset.Decimals)
					case pair.OutIndex:
						decSold = int(pair.OutAsset.Decimals)
					}
					switch swap.boughtID {
					case pair.OutIndex:
						decBought = int(pair.OutAsset.Decimals)
					case pair.InIndex:
						decBought = int(pair.InAsset.Decimals)
					}
					soldAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.sold), new(big.Float).SetFloat64(math.Pow10(decSold))).Float64()
					boughtAmt, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.bought), new(big.Float).SetFloat64(math.Pow10(decBought))).Float64()
					normalizedSwap := CurveSwap{
						ID:        swap.swapID,
						Timestamp: time.Now().Unix(),
						Pair:      pair,
						Amount0:   soldAmt,
						Amount1:   boughtAmt,
					}
					scraper.processSwap(normalizedSwap, pairs, address, lock, tradesChannel)
				} else {
					scraper.subscribeChannel <- address
					return
				}
			case err := <-subErr:
				log.Errorf("Subscription error for pool %s: %v", address.Hex(), err)
				scraper.subscribeChannel <- address
				return
			case err := <-factorySubErr:
				log.Errorf("Subscription error for pool %s: %v", address.Hex(), err)
				scraper.subscribeChannel <- address
				return
			case err := <-twoSubErr:
				log.Errorf("Subscription error for pool %s: %v", address.Hex(), err)
				scraper.subscribeChannel <- address
				return
			case err := <-underlyingSubErr:
				log.Errorf("Subscription error for pool %s: %v", address.Hex(), err)
				scraper.subscribeChannel <- address
				return
			case <-ctx.Done():
				log.Infof("Shutting down watchSwaps for %s", address.Hex())
				return
			}
		}
	}()
}

func (scraper *CurveScraper) processSwap(swap CurveSwap, pairs []CurvePair, address common.Address, lock *sync.RWMutex, tradesChannel chan models.Trade) {
	for _, pair := range pairs {
		// only do InIndex -> OutIndex direction
		if swap.Pair.InIndex == pair.InIndex && swap.Pair.OutIndex == pair.OutIndex {
			price, volume := getSwapDataCurve(swap)
			t := makeCurveTrade(pair, price, volume, time.Unix(swap.Timestamp, 0), address, swap.ID, scraper.exchange.Name, scraper.exchange.Blockchain)

			// Update lastTradeTimeMap
			lock.Lock()
			scraper.lastTradeTimeMap[address] = t.Time
			lock.Unlock()

			tradesChannel <- t
			logTrade(t)
		}
	}
}

type rawSwapData struct {
	addr     common.Address
	soldID   int
	boughtID int
	sold     *big.Int
	bought   *big.Int
	swapID   string
}

func (scraper *CurveScraper) extractSwapData(ev interface{}) (rawSwapData, error) {
	switch s := ev.(type) {
	case curvepool.CurvepoolTokenExchange:
		return rawSwapData{
			addr:     s.Raw.Address,
			soldID:   int(s.SoldId.Int64()),
			boughtID: int(s.BoughtId.Int64()),
			sold:     s.TokensSold,
			bought:   s.TokensBought,
			swapID:   s.Raw.TxHash.Hex(),
		}, nil
	case curvefifactory.CurvefifactoryTokenExchange:
		return rawSwapData{
			addr:     s.Raw.Address,
			soldID:   int(s.SoldId.Int64()),
			boughtID: int(s.BoughtId.Int64()),
			sold:     s.TokensSold,
			bought:   s.TokensBought,
			swapID:   s.Raw.TxHash.Hex(),
		}, nil
	case curvefitwocryptooptimized.CurvefitwocryptooptimizedTokenExchange:
		return rawSwapData{
			addr:     s.Raw.Address,
			soldID:   int(s.SoldId.Int64()),
			boughtID: int(s.BoughtId.Int64()),
			sold:     s.TokensSold,
			bought:   s.TokensBought,
			swapID:   s.Raw.TxHash.Hex(),
		}, nil
	case curvepool.CurvepoolTokenExchangeUnderlying:
		return rawSwapData{
			addr:     s.Raw.Address,
			soldID:   int(s.SoldId.Int64()),
			boughtID: int(s.BoughtId.Int64()),
			sold:     s.TokensSold,
			bought:   s.TokensBought,
			swapID:   s.Raw.TxHash.Hex(),
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

func (scraper *CurveScraper) GetSwapsChannel(address common.Address) (*curveFeeds, error) {
	filterer, err := curvepool.NewCurvepoolFilterer(address, scraper.wsClient)
	if err != nil {
		return nil, fmt.Errorf("curvepool filterer: %v", err)
	}

	curvefiFactoryFilterer, err := curvefifactory.NewCurvefifactoryFilterer(address, scraper.wsClient)
	if err != nil {
		return nil, fmt.Errorf("curvefifactory filterer: %v", err)
	}

	filtererTwoCrypto, err := curvefitwocryptooptimized.NewCurvefitwocryptooptimizedFilterer(address, scraper.wsClient)
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
	blocksToBackfill := uint64(20)
	if hdr, err := scraper.restClient.HeaderByNumber(context.Background(), nil); err != nil {
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

// try to get coin address from pool contract
func coinAddressFromPool(ctx context.Context, client *ethclient.Client, pool common.Address, idx int) (common.Address, error) {
	bi := big.NewInt(int64(idx))

	// 1) start with curvepool which is the most common pool
	if c, err := curvepool.NewCurvepool(pool, client); err == nil {
		// try coins(idx)
		if addr, err := c.Coins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
		// try underlying_coins(idx)
		if addr, err := c.UnderlyingCoins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
	}

	// 2) TwoCrypto optimized pool (like 2pool/cryptoswap)
	if c2, err := curvefitwocryptooptimized.NewCurvefitwocryptooptimized(pool, client); err == nil {
		if addr, err := c2.Coins(&bind.CallOpts{Context: ctx}, bi); err == nil && addr != (common.Address{}) {
			return addr, nil
		}
		// TwoCrypto usually doesn't have underlying_coins
	}

	return common.Address{}, fmt.Errorf("cannot resolve coin at index %d for pool %s", idx, pool.Hex())
}

func (scraper *CurveScraper) makePoolMap(pools []models.Pool) error {
	scraper.poolMap = make(map[common.Address][]CurvePair)
	var assetMap = make(map[common.Address]models.Asset)

	ctx := context.Background()

	for _, pool := range pools {
		outIdx, inIdx, err := parseIndexCode(pool.Order)
		if err != nil {
			log.Error("Curve - failed to parse index code: ", err)
			continue
		}
		addr := common.HexToAddress(pool.Address)

		tokenOutAddr, err := coinAddressFromPool(ctx, scraper.restClient, addr, outIdx)
		if err != nil || tokenOutAddr == (common.Address{}) {
			log.Errorf("Curve - failed to get coin address from pool: %v", err)
			continue
		}
		if _, ok := assetMap[tokenOutAddr]; !ok {
			if isNative(tokenOutAddr) {
				assetMap[tokenOutAddr] = nativeAsset(scraper.exchange.Blockchain)
			} else {
				token0, err := models.GetAsset(tokenOutAddr, scraper.exchange.Blockchain, scraper.restClient)
				if err != nil {
					return err
				}
				assetMap[tokenOutAddr] = token0
			}
		}

		tokenInAddr, err := coinAddressFromPool(ctx, scraper.restClient, addr, inIdx)
		if err != nil || tokenInAddr == (common.Address{}) {
			log.Errorf("Curve - failed to get coin address from pool: %v", err)
			continue
		}
		if _, ok := assetMap[tokenInAddr]; !ok {
			if isNative(tokenInAddr) {
				assetMap[tokenInAddr] = nativeAsset(scraper.exchange.Blockchain)
			} else {
				token1, err := models.GetAsset(tokenInAddr, scraper.exchange.Blockchain, scraper.restClient)
				if err != nil {
					return err
				}
				assetMap[tokenInAddr] = token1
			}
		}

		pair := CurvePair{
			OutAsset: assetMap[tokenOutAddr],
			InAsset:  assetMap[tokenInAddr],
			OutIndex: outIdx,
			InIndex:  inIdx,
			Address:  addr,
		}

		scraper.poolMap[addr] = append(scraper.poolMap[addr], pair)
	}

	return nil
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
