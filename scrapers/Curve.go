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

type poolSub struct {
	cancel  context.CancelFunc
	lastSub time.Time
}

type CurveScraper struct {
	exchange         models.Exchange
	poolMap          map[common.Address][]CurvePair
	subscribeChannel chan common.Address
	lastTradeTimeMap map[common.Address]time.Time
	restClient       *ethclient.Client
	wsClient         *ethclient.Client
	waitTime         int
}

func NewCurveScraper(ctx context.Context, exchangeName string, blockchain string, pools []models.Pool, tradesChannel chan models.Trade, wg *sync.WaitGroup) {
	var err error
	var scraper CurveScraper
	log.Info("Started Curve scraper.")

	scraper.exchange.Name = exchangeName
	scraper.exchange.Blockchain = blockchain
	scraper.subscribeChannel = make(chan common.Address)
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
	go scraper.mainLoop(ctx, pools, tradesChannel, &lock)
}

func (scraper *CurveScraper) mainLoop(ctx context.Context, pools []models.Pool, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	var wg sync.WaitGroup
	subs := make(map[common.Address]*poolSub)
	seen := map[common.Address]bool{} // to avoid duplicate pools
	for _, pool := range pools {
		addr := common.HexToAddress(pool.Address)
		if seen[addr] {
			continue
		}
		seen[addr] = true

		lock.Lock()
		scraper.lastTradeTimeMap[addr] = time.Now()
		lock.Unlock()

		envVar := strings.ToUpper(scraper.exchange.Name) + "_WATCHDOG_" + pool.Address
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "600"), 10, 64)
		if err != nil {
			log.Errorf("Curve - Parse curveWatchdogDelay: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdogPool(ctx, scraper.exchange.Name, addr, watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, lock)

		time.Sleep(time.Duration(scraper.waitTime) * time.Millisecond)
		subCtx, cancel := context.WithCancel(ctx)
		wg.Add(1)
		go func(address common.Address, cancel context.CancelFunc) {
			defer wg.Done()
			scraper.watchSwaps(subCtx, address, tradesChannel, lock)
		}(addr, cancel)
		subs[addr] = &poolSub{cancel: cancel, lastSub: time.Now()}
	}

	cooldown := 30 * time.Second

	go func() {
		for {
			select {
			case addr := <-scraper.subscribeChannel:
				lock.Lock()
				ps := subs[addr]
				if ps == nil {
					subCtx, cancel := context.WithCancel(ctx)
					wg.Add(1)
					go func(address common.Address, cancel context.CancelFunc) {
						defer wg.Done()
						scraper.watchSwaps(subCtx, address, tradesChannel, lock)
					}(addr, cancel)
					subs[addr] = &poolSub{cancel: cancel, lastSub: time.Now()}
					log.Infof("Initial subscription to pool: %s", addr.Hex())
					lock.Unlock()
					continue
				}

				if time.Since(ps.lastSub) < cooldown {
					// Pool already active, skipping resubscription
					lock.Unlock()
					continue
				}
				// Cancel the old subscription before starting a new one.
				ps.cancel()
				subCtx, cancel := context.WithCancel(ctx)
				wg.Add(1)
				go func(address common.Address, cancel context.CancelFunc) {
					defer wg.Done()
					scraper.watchSwaps(subCtx, address, tradesChannel, lock)
				}(addr, cancel)
				ps.cancel = cancel
				ps.lastSub = time.Now()
				log.Infof("Watchdog re-subscription to pool: %s", addr.Hex())
				lock.Unlock()
			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				lock.Lock()
				for _, ps := range subs {
					ps.cancel()
				}
				lock.Unlock()
				return
			}
		}
	}()
	wg.Wait()
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
						log.Errorf("Curve - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Errorf("Curve - no pair found for address %s", swap.addr.Hex())
						continue
					}
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
						log.Errorf("Curve - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Errorf("Curve - no pair found for address %s", swap.addr.Hex())
						continue
					}
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
						log.Errorf("Curve - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Errorf("Curve - no pair found for address %s", swap.addr.Hex())
						continue
					}
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
						log.Errorf("Curve - no pairs found for address %s", swap.addr.Hex())
						continue
					}
					pair, ok := findPairByIdx(pairs, swap.soldID, swap.boughtID)
					if !ok {
						log.Errorf("Curve - no pair found for address %s", swap.addr.Hex())
						continue
					}
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
