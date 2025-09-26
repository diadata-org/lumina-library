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
	exchange         models.Exchange
	poolMap          map[common.Address]UniV3Pair
	wsClient         *ethclient.Client
	restClient       *ethclient.Client
	subscribeChannel chan common.Address
	lastTradeTimeMap map[common.Address]time.Time
	waitTime         int
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
	go scraper.mainLoop(ctx, pools, tradesChannel, &lock)
}

// runs in a goroutine until scraper is closed.
func (scraper *UniswapV3Scraper) mainLoop(ctx context.Context, pools []models.Pool, tradesChannel chan models.Trade, lock *sync.RWMutex) {

	var wg sync.WaitGroup

	for _, pool := range pools {

		// Initialize lastTradeTimeMap.
		lock.Lock()
		scraper.lastTradeTimeMap[common.HexToAddress(pool.Address)] = time.Now()
		lock.Unlock()

		// Set up watchdog.
		envVar := strings.ToUpper(scraper.exchange.Name) + "_WATCHDOG_" + pool.Address
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "600"), 10, 64)
		if err != nil {
			log.Errorf("UniswapV3 - Parse coinbaseWatchdogDelay: %v.", err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdogPool(ctx, scraper.exchange.Name, common.HexToAddress(pool.Address), watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, lock)

		// spawn swaps watching for @pool.Address.
		time.Sleep(time.Duration(scraper.waitTime) * time.Millisecond)
		wg.Add(1)
		go func(ctx context.Context, address common.Address, w *sync.WaitGroup, lock *sync.RWMutex) {
			defer w.Done()
			scraper.watchSwaps(ctx, address, tradesChannel, lock)
		}(ctx, common.HexToAddress(pool.Address), &wg, lock)

	}

	// Routine for resubscription whenever no trades are registered on a pool
	// for longer than the corresponding watchdog time.
	go func() {
		for {
			select {
			case addr := <-scraper.subscribeChannel:

				log.Infof("Resubscribing to pool: %s", addr.Hex())

				// reset lastTradeTime to now to avoid immediate repeat.
				lock.Lock()
				scraper.lastTradeTimeMap[addr] = time.Now()
				lock.Unlock()

				go scraper.watchSwaps(ctx, addr, tradesChannel, lock)

			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				return
			}
		}
	}()

	wg.Wait()

}

// watchSwaps subscribes to a uniswap pool and forwards trades to the trades channel.
func (scraper *UniswapV3Scraper) watchSwaps(ctx context.Context, poolAddress common.Address, tradesChannel chan models.Trade, lock *sync.RWMutex) {

	// Relevant pair info is retrieved from @poolMap.
	pair := scraper.poolMap[poolAddress]

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
