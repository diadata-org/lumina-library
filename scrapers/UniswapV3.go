package scrapers

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	PancakeswapV3Pair "github.com/diadata-org/lumina-library/contracts/pancakeswapv3"
	UniswapV3Pair "github.com/diadata-org/lumina-library/contracts/uniswapv3/uniswapV3Pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
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
	base    *BaseDEXScraper
	poolMap map[common.Address]UniV3Pair
}

// implements DEXHooks
type uniswapV3Hooks struct {
	s            *UniswapV3Scraper
	exchangeName string
}

func NewUniswapV3Scraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	wg *sync.WaitGroup,
) *UniswapV3Scraper {
	s := &UniswapV3Scraper{
		poolMap: make(map[common.Address]UniV3Pair),
	}
	hooks := &uniswapV3Hooks{
		s:            s,
		exchangeName: exchangeName,
	}

	// BaseDEXScraper handles:
	//  - initialize exchange / WS client
	//  - subscribe/unsubscribe/watchdog logic
	//  - dynamically add/remove pools from config file
	s.base = NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, wg)

	log.Infof("Started %s scraper.", exchangeName)
	return s
}

// 1) ExchangeName

func (h *uniswapV3Hooks) ExchangeName() string {
	return h.exchangeName
}

// 2) EnsurePair: create / update poolMap, return addr + watchdogDelay

func (h *uniswapV3Hooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {
	s := h.s
	addr = common.HexToAddress(pool.Address)

	lock.Lock()
	pair, ok := s.poolMap[addr]
	lock.Unlock()

	if !ok {
		caller, err := UniswapV3Pair.NewUniswapV3PairCaller(addr, base.RESTClient())
		if err != nil {
			return addr, 0, err
		}

		token0Addr, err := caller.Token0(&bind.CallOpts{})
		if err != nil {
			return addr, 0, err
		}
		token1Addr, err := caller.Token1(&bind.CallOpts{})
		if err != nil {
			return addr, 0, err
		}

		t0, err := models.GetAsset(token0Addr, base.exchange.Blockchain, base.RESTClient())
		if err != nil {
			return addr, 0, err
		}
		t1, err := models.GetAsset(token1Addr, base.exchange.Blockchain, base.RESTClient())
		if err != nil {
			return addr, 0, err
		}

		pair = UniV3Pair{
			Token0:  t0,
			Token1:  t1,
			Address: addr,
			Order:   pool.Order,
		}

		lock.Lock()
		s.poolMap[addr] = pair
		lock.Unlock()
	} else {
		pair.Order = pool.Order
		lock.Lock()
		s.poolMap[addr] = pair
		lock.Unlock()
	}

	watchdogDelay = pool.WatchDogDelay
	return addr, watchdogDelay, nil
}

// 3) StartStream: start the event subscription + processing logic for a single pool.
func (h *uniswapV3Hooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	s := h.s

	// get pair from poolMap
	lock.RLock()
	pair, ok := s.poolMap[addr]
	lock.RUnlock()
	if !ok {
		log.Warnf("%s - StartStream: pool %s not found in poolMap", h.ExchangeName(), addr.Hex())
		return
	}

	log.Infof("%s - subscribe to %s with pair %s",
		h.ExchangeName(), addr.Hex(), pair.Token0.Symbol+"-"+pair.Token1.Symbol)

	// -------- PancakeswapV3 --------
	if base.exchange.Name == PANCAKESWAPV3_EXCHANGE {
		sink, sub, err := s.GetPancakeSwapsChannel(base, addr)
		if err != nil {
			log.Errorf("PancakeswapV3 - error fetching swaps channel for %s: %v", addr.Hex(), err)
			return
		}

		go func() {
			for {
				select {
				case rawSwap, ok := <-sink:
					if !ok {
						log.Infof("PancakeswapV3 - swaps channel closed for %s", addr.Hex())
						return
					}

					swap := s.normalizeUniV3Swap(*rawSwap)
					price, volume := getSwapDataUniV3(swap)
					if math.IsNaN(price) || math.IsInf(price, 0) || price == 0 {
						log.Debugf("PancakeswapV3 - skip zero/+inf price, tx=%s", swap.ID)
						continue
					}

					t := makeTrade(
						pair,
						price,
						volume,
						time.Unix(swap.Timestamp, 0),
						addr,
						swap.ID,
						base.exchange.Name,
						base.exchange.Blockchain,
					)

					base.UpdateLastTradeTime(addr, t.Time, lock)

					switch pair.Order {
					case 0:
						tradesChan <- t
					case 1:
						t.SwapTrade()
						tradesChan <- t
					case 2:
						logTrade(t)
						tradesChan <- t
						t.SwapTrade()
						tradesChan <- t
					}
					logTrade(t)

				case err := <-sub.Err():
					log.Errorf("PancakeswapV3 - subscription error for pool %s: %v", addr.Hex(), err)
					base.SubscribeChannel() <- addr
					return

				case <-ctx.Done():
					log.Infof("PancakeswapV3 - shutting down stream for %s", addr.Hex())
					return
				}
			}
		}()

		return
	}

	// -------- UniswapV3 --------
	sink, sub, err := s.GetSwapsChannel(base, addr)
	if err != nil {
		log.Errorf("UniswapV3 - error fetching swaps channel for %s: %v", addr.Hex(), err)
		return
	}

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if !ok {
					log.Infof("UniswapV3 - swaps channel closed for %s", addr.Hex())
					return
				}

				swap := s.normalizeUniV3Swap(*rawSwap)
				price, volume := getSwapDataUniV3(swap)
				if math.IsNaN(price) || math.IsInf(price, 0) || price == 0 {
					log.Debugf("UniswapV3 - skip zero/+inf price, tx=%s", swap.ID)
					continue
				}

				t := makeTrade(
					pair,
					price,
					volume,
					time.Unix(swap.Timestamp, 0),
					addr,
					swap.ID,
					base.exchange.Name,
					base.exchange.Blockchain,
				)

				base.UpdateLastTradeTime(addr, t.Time, lock)

				switch pair.Order {
				case 0:
					tradesChan <- t
				case 1:
					t.SwapTrade()
					tradesChan <- t
				case 2:
					logTrade(t)
					tradesChan <- t
					t.SwapTrade()
					tradesChan <- t
				}
				logTrade(t)

			case err := <-sub.Err():
				log.Errorf("UniswapV3 - subscription error for pool %s: %v", addr.Hex(), err)
				base.SubscribeChannel() <- addr
				return

			case <-ctx.Done():
				log.Infof("UniswapV3 - shutting down stream for %s", addr.Hex())
				return
			}
		}
	}()
}

// 4) OnPoolRemoved: called after stopPool to clean up V3's poolMap.

func (h *uniswapV3Hooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addr common.Address,
	lock *sync.RWMutex,
) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()
}

// 5) OnOrderChanged: called when pool.Order changed in the config diff.

func (h *uniswapV3Hooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	lock *sync.RWMutex,
) {
	lock.Lock()
	defer lock.Unlock()

	pair, ok := h.s.poolMap[addr]
	if !ok {
		return
	}
	pair.Order = newPool.Order
	h.s.poolMap[addr] = pair
}

// Uniswap V3
func (s *UniswapV3Scraper) GetSwapsChannel(
	base *BaseDEXScraper,
	pairAddress common.Address,
) (chan *UniswapV3Pair.UniswapV3PairSwap, event.Subscription, error) {

	wsClient := base.WSClient()
	if wsClient == nil {
		return nil, nil, fmt.Errorf("WS client is nil for pair %s", pairAddress.Hex())
	}

	sink := make(chan *UniswapV3Pair.UniswapV3PairSwap)

	pairFiltererContract, err := UniswapV3Pair.NewUniswapV3PairFilterer(pairAddress, wsClient)
	if err != nil {
		return nil, nil, err
	}

	sub, err := pairFiltererContract.WatchSwap(
		&bind.WatchOpts{},
		sink,
		[]common.Address{},
		[]common.Address{},
	)
	if err != nil {
		log.Errorf("UniswapV3 - WatchSwap(%s) error: %v", pairAddress.Hex(), err)
		return nil, nil, err
	}

	return sink, sub, nil
}

// Pancakeswap V3
func (s *UniswapV3Scraper) GetPancakeSwapsChannel(
	base *BaseDEXScraper,
	pairAddress common.Address,
) (chan *PancakeswapV3Pair.Pancakev3pairSwap, event.Subscription, error) {

	wsClient := base.WSClient()
	if wsClient == nil {
		return nil, nil, fmt.Errorf("WS client is nil for pair %s", pairAddress.Hex())
	}

	sink := make(chan *PancakeswapV3Pair.Pancakev3pairSwap)

	pairFiltererContract, err := PancakeswapV3Pair.NewPancakev3pairFilterer(pairAddress, wsClient)
	if err != nil {
		return nil, nil, err
	}

	sub, err := pairFiltererContract.WatchSwap(
		&bind.WatchOpts{},
		sink,
		[]common.Address{},
		[]common.Address{},
	)
	if err != nil {
		log.Errorf("PancakeswapV3 - WatchSwap(%s) error: %v", pairAddress.Hex(), err)
		return nil, nil, err
	}

	return sink, sub, nil
}

// normalizeUniV3Swap takes a swap as returned by the swap contract's channel and converts it to a UniswapV3Swap type.
func (s *UniswapV3Scraper) normalizeUniV3Swap(swapI interface{}) (normalizedSwap UniswapV3Swap) {
	switch swap := swapI.(type) {
	case UniswapV3Pair.UniswapV3PairSwap:
		pair := s.poolMap[swap.Raw.Address]
		decimals0 := int(pair.Token0.Decimals)
		decimals1 := int(pair.Token1.Decimals)
		amount0, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount0), big.NewFloat(math.Pow10(decimals0))).Float64()
		amount1, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount1), big.NewFloat(math.Pow10(decimals1))).Float64()
		normalizedSwap = UniswapV3Swap{
			ID:        swap.Raw.TxHash.Hex(),
			Timestamp: time.Now().Unix(),
			Pair:      pair,
			Amount0:   amount0,
			Amount1:   amount1,
		}
	case PancakeswapV3Pair.Pancakev3pairSwap:
		pair := s.poolMap[swap.Raw.Address]
		decimals0 := int(pair.Token0.Decimals)
		decimals1 := int(pair.Token1.Decimals)
		amount0, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount0), big.NewFloat(math.Pow10(decimals0))).Float64()
		amount1, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount1), big.NewFloat(math.Pow10(decimals1))).Float64()
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
