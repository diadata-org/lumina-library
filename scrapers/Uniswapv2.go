package scrapers

import (
	"context"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	uniswap "github.com/diadata-org/lumina-library/contracts/uniswap/pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

var (
	restDial = ""
	wsDial   = ""
)

type UniswapToken struct {
	Address  common.Address
	Symbol   string
	Decimals uint8
	Name     string
}

type UniswapPair struct {
	Token0      models.Asset
	Token1      models.Asset
	Address     common.Address
	ForeignName string
}

type UniswapSwap struct {
	ID         string
	Timestamp  int64
	Pair       UniswapPair
	Amount0In  float64
	Amount0Out float64
	Amount1In  float64
	Amount1Out float64
}

type UniswapV2Scraper struct {
	exchange         models.Exchange
	subscribeChannel chan common.Address
	lastTradeTimeMap map[common.Address]time.Time
	poolMap          map[common.Address]UniswapPair
	wsClient         *ethclient.Client
	restClient       *ethclient.Client
	waitTime         int
}

func NewUniswapV2Scraper(ctx context.Context, exchangeName string, blockchain string, pools []models.Pool, tradesChannel chan models.Trade, wg *sync.WaitGroup) {
	var err error
	var scraper UniswapV2Scraper
	log.Info("Started UniswapV2 scraper.")

	scraper.exchange.Name = exchangeName
	scraper.exchange.Blockchain = blockchain
	scraper.subscribeChannel = make(chan common.Address)
	scraper.lastTradeTimeMap = make(map[common.Address]time.Time)
	scraper.waitTime, err = strconv.Atoi(utils.Getenv(strings.ToUpper(exchangeName)+"_WAIT_TIME", "500"))
	if err != nil {
		log.Errorf("Failed to parse waitTime for exchange %s: %v", exchangeName, err)
	}

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV2_EXCHANGE)+"_URI_REST", restDial))
	if err != nil {
		log.Error("UniswapV2 - init rest client: ", err)
	}
	scraper.wsClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV2_EXCHANGE)+"_URI_WS", wsDial))
	if err != nil {
		log.Error("UniswapV2 - init ws client: ", err)
	}

	// Fetch all pool with given liquidity threshold from database.
	scraper.makeUniPoolMap(pools)

	var lock sync.RWMutex
	go scraper.mainLoop(ctx, pools, tradesChannel, &lock)
}

// runs in a goroutine until s is closed
func (scraper *UniswapV2Scraper) mainLoop(ctx context.Context, pools []models.Pool, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	var wg sync.WaitGroup
	for _, pool := range pools {
		lock.Lock()
		scraper.lastTradeTimeMap[common.HexToAddress(pool.Address)] = time.Now()
		lock.Unlock()

		envVar := strings.ToUpper(scraper.exchange.Name) + "_WATCHDOG_" + pool.Address
		watchdogDelay, err := strconv.ParseInt(utils.Getenv(envVar, "600"), 10, 64)
		if err != nil {
			log.Errorf("UniswapV2 - Parse %s: %v.", envVar, err)
		}
		watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
		go watchdogPool(ctx, scraper.exchange.Name, common.HexToAddress(pool.Address), watchdogTicker, scraper.lastTradeTimeMap, watchdogDelay, scraper.subscribeChannel, lock)

		time.Sleep(time.Duration(scraper.waitTime) * time.Millisecond)
		wg.Add(1)
		go func(ctx context.Context, address common.Address, w *sync.WaitGroup, lock *sync.RWMutex) {
			defer w.Done()
			scraper.ListenToPair(ctx, address, tradesChannel, lock)
		}(ctx, common.HexToAddress(pool.Address), &wg, lock)
	}

	go func() {
		for {
			select {
			case addr := <-scraper.subscribeChannel:
				log.Infof("Resubscribing to pool: %s", addr.Hex())
				lock.Lock()
				scraper.lastTradeTimeMap[addr] = time.Now()
				lock.Unlock()
				go scraper.ListenToPair(ctx, addr, tradesChannel, lock)
			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				return
			}
		}
	}()

	wg.Wait()
}

// makeUniPoolMap returns a map with pool addresses as keys and the underlying UniswapPair as values.
func (scraper *UniswapV2Scraper) makeUniPoolMap(pools []models.Pool) error {
	scraper.poolMap = make(map[common.Address]UniswapPair)
	var assetMap = make(map[common.Address]models.Asset)

	for _, p := range pools {
		univ2PairCaller, err := uniswap.NewUniswapV2PairCaller(common.HexToAddress(p.Address), scraper.restClient)
		if err != nil {
			return err
		}

		token0Address, err := univ2PairCaller.Token0(&bind.CallOpts{})
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

		token1Address, err := univ2PairCaller.Token1(&bind.CallOpts{})
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
		scraper.poolMap[common.HexToAddress(p.Address)] = UniswapPair{
			Address:     common.HexToAddress(p.Address),
			Token0:      assetMap[token0Address],
			Token1:      assetMap[token1Address],
			ForeignName: assetMap[token0Address].Symbol + "-" + assetMap[token1Address].Symbol,
		}
	}
	return nil
}

// ListenToPair subscribes to a uniswap pool.
// If @byAddress is true, it listens by pool address, otherwise by index.
func (scraper *UniswapV2Scraper) ListenToPair(ctx context.Context, address common.Address, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	var err error

	// Relevant pool info is retrieved from @poolMap.
	pair := scraper.poolMap[address]

	sink, sub, err := scraper.GetSwapsChannel(address)
	if err != nil {
		log.Error("UniswapV2 - error fetching swaps channel: ", err)
	}

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if ok {
					swap, err := scraper.normalizeUniswapSwap(*rawSwap)
					if err != nil {
						log.Error("UniswapV2 - error normalizing swap: ", err)
					}
					price, volume := getSwapData(swap)
					t := makeTradeUniswapV2(pair, price, volume, time.Unix(swap.Timestamp, 0), rawSwap.Raw.Address, swap.ID, scraper.exchange.Name, scraper.exchange.Blockchain)
					lock.Lock()
					scraper.lastTradeTimeMap[rawSwap.Raw.Address] = t.Time
					lock.Unlock()

					tradesChannel <- t
					logTradeUniswapV2(t)
				}
			case err := <-sub.Err():
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

func makeTradeUniswapV2(
	pair UniswapPair,
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

func logTradeUniswapV2(t models.Trade) {
	log.Debugf(
		"Got trade at time %v - symbol: %s, pair: %s, price: %v, volume:%v",
		t.Time, t.QuoteToken.Symbol, t.QuoteToken.Symbol+"-"+t.BaseToken.Symbol, t.Price, t.Volume)
}

// GetSwapsChannel returns a channel for swaps of the pair with address @pairAddress
func (scraper *UniswapV2Scraper) GetSwapsChannel(pairAddress common.Address) (chan *uniswap.UniswapV2PairSwap, event.Subscription, error) {

	sink := make(chan *uniswap.UniswapV2PairSwap)
	var pairFiltererContract *uniswap.UniswapV2PairFilterer
	pairFiltererContract, err := uniswap.NewUniswapV2PairFilterer(pairAddress, scraper.wsClient)
	if err != nil {
		log.Error("UniswapV2 - pair filterer: ", err)
	}

	sub, err := pairFiltererContract.WatchSwap(&bind.WatchOpts{}, sink, []common.Address{}, []common.Address{})
	if err != nil {
		log.Error("UniswapV2 - error in get swaps channel: ", err)
	}

	return sink, sub, nil
}

// normalizeUniswapSwap takes a swap as returned by the swap contract's channel and converts it to a UniswapSwap type
func (scraper *UniswapV2Scraper) normalizeUniswapSwap(swap uniswap.UniswapV2PairSwap) (normalizedSwap UniswapSwap, err error) {
	pair := scraper.poolMap[swap.Raw.Address]
	decimals0 := int(pair.Token0.Decimals)
	decimals1 := int(pair.Token1.Decimals)
	amount0In, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount0In), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
	amount0Out, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount0Out), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
	amount1In, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount1In), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()
	amount1Out, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(swap.Amount1Out), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()

	normalizedSwap = UniswapSwap{
		ID:         swap.Raw.TxHash.Hex(),
		Timestamp:  time.Now().Unix(),
		Pair:       pair,
		Amount0In:  amount0In,
		Amount0Out: amount0Out,
		Amount1In:  amount1In,
		Amount1Out: amount1Out,
	}
	return
}

// getSwapData returns price, volume and sell/buy information of @swap
func getSwapData(swap UniswapSwap) (price float64, volume float64) {
	if swap.Amount0In == float64(0) {
		volume = swap.Amount0Out
		price = swap.Amount1In / swap.Amount0Out
		return
	}
	volume = -swap.Amount0In
	price = swap.Amount1Out / swap.Amount0In
	return
}
