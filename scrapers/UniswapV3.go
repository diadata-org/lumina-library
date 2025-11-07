package scrapers

import (
	"context"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	PancakeswapV3Pair "github.com/diadata-org/lumina-library/contracts/pancakeswapv3"
	UniswapV3Pair "github.com/diadata-org/lumina-library/contracts/uniswapv3/uniswapV3Pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

var (
	restDialUniV3 = ""
	wsDialUniV3   = ""
)

type UniV3Pair struct {
	Token0      models.Asset
	Token1      models.Asset
	ForeignName string
	Address     common.Address
	Order       int
}

type UniswapV3Swap struct {
	ID        string
	Timestamp int64
	Pair      UniV3Pair
	Amount0   float64
	Amount1   float64
}

type UniswapV3Scraper struct {
	exchange           models.Exchange
	poolMap            map[common.Address]UniV3Pair
	wsClient           *ethclient.Client
	restClient         *ethclient.Client
	subscribeChannel   chan common.Address
	unsubscribeChannel chan common.Address
	swapStreamCancel   map[common.Address]context.CancelFunc
	watchdogCancel     map[string]context.CancelFunc
	lastTradeTimeMap   map[common.Address]time.Time
	waitTime           int
}

func NewUniswapV3Scraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	wg *sync.WaitGroup,
) {

	var err error
	var scraper UniswapV3Scraper
	log.Infof("Started %s scraper.", exchangeName)

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

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(exchangeName)+"_URI_REST", restDialUniV3))
	if err != nil {
		log.Error("UniswapV3 - init rest client: ", err)
	}
	scraper.wsClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(exchangeName)+"_URI_WS", wsDialUniV3))
	if err != nil {
		log.Error("UniswapV3 - init ws client: ", err)
	}

	scraper.makePoolMap(pools)

	var lock sync.RWMutex
	scraper.startResubHandler(ctx, tradesChannel, &lock)
	scraper.startUnsubHandler(ctx, &lock)
	go scraper.watchConfig(ctx, exchangeName, tradesChannel, &lock)
}

func mapPoolsByAddress(pools []models.Pool) map[string]models.Pool {
	poolMap := make(map[string]models.Pool)
	for _, pool := range pools {
		poolMap[strings.ToLower(pool.Address)] = pool
	}
	return poolMap
}

func (s *UniswapV3Scraper) startResubHandler(ctx context.Context, trades chan models.Trade, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-s.subscribeChannel:
				log.Infof("Resubscribing to pool: %s", addr.Hex())
				lock.Lock()
				s.lastTradeTimeMap[addr] = time.Now()
				lock.Unlock()
				if cancel, ok := s.swapStreamCancel[addr]; ok && cancel != nil {
					cancel()
					delete(s.swapStreamCancel, addr)
				}
				go s.watchSwaps(ctx, addr, trades, lock)
			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				return
			}
		}
	}()
}

func (scraper *UniswapV3Scraper) startUnsubHandler(ctx context.Context, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-scraper.unsubscribeChannel:
				log.Infof("Unsubscribing from pool: %s", addr.Hex())
				scraper.stopPool(addr, lock)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (scraper *UniswapV3Scraper) watchConfig(ctx context.Context, exchangeName string, trades chan models.Trade, lock *sync.RWMutex) {
	// Check for config changes every 60 minutes.
	envKey := strings.ToUpper(exchangeName) + "_WATCH_CONFIG_INTERVAL"
	interval, err := strconv.Atoi(utils.Getenv(envKey, "3600"))
	if err != nil {
		log.Errorf("UniswapV3 - Failed to parse %s: %v.", envKey, err)
		return
	}
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	initial, err := models.PoolsFromConfigFile(exchangeName)
	if err != nil {
		log.Errorf("UniswapV3 - GetPoolsFromConfigFile: %v.", err)
		initial = nil
		return
	}
	last := mapPoolsByAddress(initial)

	for _, pool := range last {
		_ = scraper.startPool(ctx, pool, trades, lock)
	}

	for {
		select {
		case <-t.C:
			currList, err := models.PoolsFromConfigFile(exchangeName)
			if err != nil {
				log.Errorf("UniswapV3 - GetPoolsFromConfigFile: %v.", err)
				continue
			}
			curr := mapPoolsByAddress(currList)
			scraper.applyConfigDiff(ctx, last, curr, trades, lock)
			last = curr
		case <-ctx.Done():
			log.Info("UniswapV3 - Close watchConfig routine of scraper.")
			return
		}
	}
}

func (scraper *UniswapV3Scraper) applyConfigDiff(ctx context.Context, last map[string]models.Pool, curr map[string]models.Pool, trades chan models.Trade, lock *sync.RWMutex) {
	for k := range last {
		if _, ok := curr[k]; !ok {
			addr := common.HexToAddress(k)
			log.Infof("UniswapV3 - Removed pool: %s.", addr.Hex())
			scraper.unsubscribeChannel <- addr
		}
	}

	for k, p := range curr {
		addr := common.HexToAddress(k)

		if old, ok := last[k]; !ok {
			log.Infof("UniswapV3 - add pool %s (wd=%d, order=%d)", p.Address, p.WatchDogDelay, p.Order)
			if err := scraper.startPool(ctx, p, trades, lock); err != nil {
				log.Errorf("UniswapV3 - startPool %s: %v", p.Address, err)
			}
			continue
		} else {
			newDelay := p.WatchDogDelay
			oldDelay := old.WatchDogDelay
			if newDelay != oldDelay {
				log.Infof("UniswapV3 - Changed pool %s (wd=%d -> %d)", p.Address, oldDelay, newDelay)
				scraper.restartWatchdogForPool(ctx, addr, newDelay, lock)
			}

			if pair, ok := scraper.poolMap[addr]; ok && pair.Order != p.Order {
				log.Infof("UniswapV3 - Update order %s: %d -> %d", p.Address, pair.Order, p.Order)
				pair.Order = p.Order
				lock.Lock()
				scraper.poolMap[addr] = pair
				lock.Unlock()
			}
		}
	}
}

func (scraper *UniswapV3Scraper) startPool(ctx context.Context, pool models.Pool, trades chan models.Trade, lock *sync.RWMutex) error {
	addr := common.HexToAddress(pool.Address)
	if _, ok := scraper.poolMap[addr]; !ok {
		caller, err := UniswapV3Pair.NewUniswapV3PairCaller(addr, scraper.restClient)
		if err != nil {
			return err
		}
		token0, err := caller.Token0(&bind.CallOpts{})
		if err != nil {
			return err
		}
		token1, err := caller.Token1(&bind.CallOpts{})
		if err != nil {
			return err
		}
		t0, err := models.GetAsset(token0, scraper.exchange.Blockchain, scraper.restClient)
		if err != nil {
			return err
		}
		t1, err := models.GetAsset(token1, scraper.exchange.Blockchain, scraper.restClient)
		if err != nil {
			return err
		}
		lock.Lock()
		scraper.poolMap[addr] = UniV3Pair{
			Token0:  t0,
			Token1:  t1,
			Address: addr,
			Order:   pool.Order,
		}
		lock.Unlock()
	} else {
		p := scraper.poolMap[addr]
		p.Order = pool.Order
		lock.Lock()
		scraper.poolMap[addr] = p
		lock.Unlock()
	}
	lock.Lock()
	if _, ok := scraper.lastTradeTimeMap[addr]; !ok {
		scraper.lastTradeTimeMap[addr] = time.Now()
	}
	lock.Unlock()

	watchdogDelay := pool.WatchDogDelay
	scraper.restartWatchdogForPool(ctx, addr, watchdogDelay, lock)
	if _, ok := scraper.swapStreamCancel[addr]; ok {
		return nil
	}

	pctx, cancel := context.WithCancel(ctx)
	scraper.swapStreamCancel[addr] = cancel
	go scraper.watchSwaps(pctx, addr, trades, lock)
	return nil
}

func (scraper *UniswapV3Scraper) stopPool(addr common.Address, lock *sync.RWMutex) {
	if cancel, ok := scraper.swapStreamCancel[addr]; ok && cancel != nil {
		cancel()
		delete(scraper.swapStreamCancel, addr)
	}
	key := addr.Hex()
	if cancel, ok := scraper.watchdogCancel[key]; ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, key)
	}
	lock.Lock()
	delete(scraper.lastTradeTimeMap, addr)
	lock.Unlock()
}

func (scraper *UniswapV3Scraper) restartWatchdogForPool(ctx context.Context, addr common.Address, watchdogDelay int64, lock *sync.RWMutex) {
	key := addr.Hex()
	if cancel, ok := scraper.watchdogCancel[key]; ok && cancel != nil {
		cancel()
		delete(scraper.watchdogCancel, key)
	}
	wdCtx, cancel := context.WithCancel(ctx)
	scraper.watchdogCancel[key] = cancel
	watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
	go func() {
		defer watchdogTicker.Stop()
		for {
			select {
			case <-wdCtx.Done():
				return
			case <-watchdogTicker.C:
				lock.RLock()
				lastTradeTime := scraper.lastTradeTimeMap[addr]
				lock.RUnlock()
				if time.Since(lastTradeTime) > time.Duration(watchdogDelay)*time.Second {
					log.Warnf("UniswapV3 - watchdog failover for %s.", addr.Hex())
					scraper.subscribeChannel <- addr
				}
			}
		}
	}()
}

// watchSwaps subscribes to a uniswap pool and forwards trades to the trades channel.
func (scraper *UniswapV3Scraper) watchSwaps(ctx context.Context, poolAddress common.Address, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	// Relevant pair info is retrieved from @poolMap.
	pair := scraper.poolMap[poolAddress]
	log.Infof("UniswapV3 - subscribe to %s with pair %s", poolAddress.Hex(), pair.Token0.Symbol+"-"+pair.Token1.Symbol)

	if scraper.exchange.Name != PANCAKESWAPV3_EXCHANGE {

		sink, subscription, err := scraper.GetSwapsChannel(poolAddress)
		if err != nil {
			log.Error("UniswapV3 - error fetching swaps channel: ", err)
		}

		go func() {
			for {
				select {
				case rawSwap, ok := <-sink:
					if ok {
						swap := scraper.normalizeUniV3Swap(*rawSwap)
						if err != nil {
							log.Error("UniswapV3 - error normalizing swap: ", err)
						}
						price, volume := getSwapDataUniV3(swap)
						t := makeTrade(pair, price, volume, time.Unix(swap.Timestamp, 0), poolAddress, swap.ID, scraper.exchange.Name, scraper.exchange.Blockchain)

						// Update lastTradeTimeMap
						lock.Lock()
						scraper.lastTradeTimeMap[poolAddress] = t.Time
						lock.Unlock()

						if pair.Order == 0 {
							tradesChannel <- t
						} else if pair.Order == 1 {
							t.SwapTrade()
							tradesChannel <- t
						} else if pair.Order == 2 {
							logTrade(t)
							tradesChannel <- t
							t.SwapTrade()
							tradesChannel <- t
						}
						logTrade(t)
					}

				case err := <-subscription.Err():
					log.Errorf("Subscription error for pool %s: %v", poolAddress.Hex(), err)
					scraper.subscribeChannel <- poolAddress
					return

				case <-ctx.Done():
					log.Infof("Shutting down watchSwaps for %s", poolAddress.Hex())
					return
				}

			}
		}()
	} else {

		sink, subscription, err := scraper.GetPancakeSwapsChannel(poolAddress)
		if err != nil {
			log.Error("PancakeswapV3 - error fetching swaps channel: ", err)
		}

		go func() {
			for {
				select {
				case rawSwap, ok := <-sink:
					if ok {
						swap := scraper.normalizeUniV3Swap(*rawSwap)
						if err != nil {
							log.Error("PancakeswapV3 - error normalizing swap: ", err)
						}
						price, volume := getSwapDataUniV3(swap)
						t := makeTrade(pair, price, volume, time.Unix(swap.Timestamp, 0), poolAddress, swap.ID, scraper.exchange.Name, scraper.exchange.Blockchain)

						// Update lastTradeTimeMap
						lock.Lock()
						scraper.lastTradeTimeMap[poolAddress] = t.Time
						lock.Unlock()

						if pair.Order == 0 {
							tradesChannel <- t
						} else if pair.Order == 1 {
							t.SwapTrade()
							tradesChannel <- t
						} else if pair.Order == 2 {
							logTrade(t)
							tradesChannel <- t
							t.SwapTrade()
							tradesChannel <- t
						}
						logTrade(t)

					}

				case err := <-subscription.Err():
					log.Errorf("Subscription error for pool %s: %v", poolAddress.Hex(), err)
					scraper.subscribeChannel <- poolAddress
					return

				case <-ctx.Done():
					log.Infof("Shutting down watchSwaps for %s", poolAddress.Hex())
					return
				}
			}
		}()
	}
}

// GetSwapsChannel returns a channel for swaps of the pair with address @pairAddress.
func (scraper *UniswapV3Scraper) GetSwapsChannel(pairAddress common.Address) (chan *UniswapV3Pair.UniswapV3PairSwap, event.Subscription, error) {
	sink := make(chan *UniswapV3Pair.UniswapV3PairSwap)
	var pairFiltererContract *UniswapV3Pair.UniswapV3PairFilterer

	pairFiltererContract, err := UniswapV3Pair.NewUniswapV3PairFilterer(pairAddress, scraper.wsClient)
	if err != nil {
		log.Fatal(err)
	}

	sub, err := pairFiltererContract.WatchSwap(&bind.WatchOpts{}, sink, []common.Address{}, []common.Address{})
	if err != nil {
		log.Error("error in get swaps channel: ", err)
	}

	return sink, sub, nil

}

// GetPancakeSwapsChannel returns a channel for swaps on Pancakeswap DEX of the pair with address @pairAddress.
func (scraper *UniswapV3Scraper) GetPancakeSwapsChannel(pairAddress common.Address) (chan *PancakeswapV3Pair.Pancakev3pairSwap, event.Subscription, error) {
	sink := make(chan *PancakeswapV3Pair.Pancakev3pairSwap)
	var pairFiltererContract *PancakeswapV3Pair.Pancakev3pairFilterer

	pairFiltererContract, err := PancakeswapV3Pair.NewPancakev3pairFilterer(pairAddress, scraper.wsClient)
	if err != nil {
		log.Fatal(err)
	}

	sub, err := pairFiltererContract.WatchSwap(&bind.WatchOpts{}, sink, []common.Address{}, []common.Address{})
	if err != nil {
		log.Error("error in get swaps channel: ", err)
	}
	return sink, sub, nil
}

// makePoolMap fills a map that maps pool addresses on the full pool asset's information.
func (scraper *UniswapV3Scraper) makePoolMap(pools []models.Pool) error {

	scraper.poolMap = make(map[common.Address]UniV3Pair)
	var assetMap = make(map[common.Address]models.Asset)

	for _, pool := range pools {

		univ3PairCaller, err := UniswapV3Pair.NewUniswapV3PairCaller(common.HexToAddress(pool.Address), scraper.restClient)
		if err != nil {
			return err
		}

		token0Address, err := univ3PairCaller.Token0(&bind.CallOpts{})
		if err != nil {
			return err
		}
		if _, ok := assetMap[token0Address]; !ok {
			token0, err := models.GetAsset(token0Address, scraper.exchange.Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
			assetMap[token0Address] = token0
		}

		token1Address, err := univ3PairCaller.Token1(&bind.CallOpts{})
		if err != nil {
			return err
		}
		if _, ok := assetMap[token1Address]; !ok {
			token1, err := models.GetAsset(token1Address, scraper.exchange.Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
			assetMap[token1Address] = token1
		}

		scraper.poolMap[common.HexToAddress(pool.Address)] = UniV3Pair{
			Token0:  assetMap[token0Address],
			Token1:  assetMap[token1Address],
			Address: common.HexToAddress(pool.Address),
			Order:   pool.Order,
		}

	}

	return nil

}

// normalizeUniswapSwap takes a swap as returned by the swap contract's channel and converts it to a UniswapSwap type.
func (scraper *UniswapV3Scraper) normalizeUniV3Swap(swapI interface{}) (normalizedSwap UniswapV3Swap) {
	switch swap := swapI.(type) {
	case UniswapV3Pair.UniswapV3PairSwap:
		pair := scraper.poolMap[swap.Raw.Address]
		decimals0 := int(pair.Token0.Decimals)
		decimals1 := int(pair.Token1.Decimals)
		amount0, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount0), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
		amount1, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount1), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()

		normalizedSwap = UniswapV3Swap{
			ID:        swap.Raw.TxHash.Hex(),
			Timestamp: time.Now().Unix(),
			Pair:      pair,
			Amount0:   amount0,
			Amount1:   amount1,
		}
	case PancakeswapV3Pair.Pancakev3pairSwap:
		pair := scraper.poolMap[swap.Raw.Address]
		decimals0 := int(pair.Token0.Decimals)
		decimals1 := int(pair.Token1.Decimals)
		amount0, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount0), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
		amount1, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount1), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()

		normalizedSwap = UniswapV3Swap{
			ID:        swap.Raw.TxHash.Hex(),
			Timestamp: time.Now().Unix(),
			Pair:      pair,
			Amount0:   amount0,
			Amount1:   amount1,
		}
	}

	return
}

func getSwapDataUniV3(swap UniswapV3Swap) (price float64, volume float64) {
	volume = swap.Amount0
	price = math.Abs(swap.Amount1 / swap.Amount0)
	return
}

func makeTrade(
	pair UniV3Pair,
	price float64,
	volume float64,
	timestamp time.Time,
	address common.Address,
	foreignTradeID string,
	exchangeName string,
	blockchain string,
) models.Trade {
	token0 := pair.Token0
	token1 := pair.Token1
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

func logTrade(t models.Trade) {
	log.Debugf(
		"Got trade at time %v - symbol: %s, pair: %s, price: %v, volume:%v",
		t.Time, t.QuoteToken.Symbol, t.QuoteToken.Symbol+"-"+t.BaseToken.Symbol, t.Price, t.Volume)
}
