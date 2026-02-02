package scrapers

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	UniswapV3Pair "github.com/diadata-org/lumina-library/contracts/uniswapv3/uniswapV3Pair"
	"github.com/diadata-org/lumina-library/models"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type AerodromeSlipstreamPair struct {
	Token0      models.Asset
	Token1      models.Asset
	ForeignName string
	Address     common.Address
	Order       int
	Divisor0    *big.Float
	Divisor1    *big.Float
}

type AerodromeSlipstreamSwap struct {
	ID        string
	Timestamp int64
	Pair      AerodromeSlipstreamPair
	Amount0   float64 // signed
	Amount1   float64 // signed
}

type AerodromeSlipstreamScraper struct {
	base    *BaseDEXScraper
	poolMap map[common.Address]AerodromeSlipstreamPair
}

// implements DEXHooks
type aerodromeSlipstreamHooks struct {
	s            *AerodromeSlipstreamScraper
	exchangeName string
}

func NewAerodromeSlipstreamScraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) *AerodromeSlipstreamScraper {
	s := &AerodromeSlipstreamScraper{
		poolMap: make(map[common.Address]AerodromeSlipstreamPair),
	}
	hooks := &aerodromeSlipstreamHooks{
		s:            s,
		exchangeName: exchangeName,
	}

	// BaseDEXScraper handles:
	//  - initialize exchange / WS client
	//  - subscribe/unsubscribe/watchdog logic
	//  - dynamically add/remove pools from config file
	s.base = NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, branchMarketConfig, wg)
	log.Infof("Started %s scraper.", exchangeName)
	return s
}

// 1) ExchangeName

func (h *aerodromeSlipstreamHooks) ExchangeName() string {
	return h.exchangeName
}

// 2) EnsurePair: create / update poolMap, return addr + watchdogDelay

func (h *aerodromeSlipstreamHooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {
	s := h.s
	addr = common.HexToAddress(pool.Address)

	lock.RLock()
	pair, ok := s.poolMap[addr]
	lock.RUnlock()

	if !ok {
		caller, err := UniswapV3Pair.NewUniswapV3PairCaller(addr, base.RESTClient())
		if err != nil {
			return addr, 0, err
		}

		token0Addr, err := caller.Token0(&bind.CallOpts{Context: ctx})
		if err != nil {
			return addr, 0, err
		}
		token1Addr, err := caller.Token1(&bind.CallOpts{Context: ctx})
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

		d0 := new(big.Float).SetFloat64(math.Pow10(int(t0.Decimals)))
		d1 := new(big.Float).SetFloat64(math.Pow10(int(t1.Decimals)))

		pair = AerodromeSlipstreamPair{
			Token0:   t0,
			Token1:   t1,
			Address:  addr,
			Order:    pool.Order,
			Divisor0: d0,
			Divisor1: d1,
		}

		lock.Lock()
		s.poolMap[addr] = pair
		lock.Unlock()
	} else {
		// update order if pool already exists
		pair.Order = pool.Order
		lock.Lock()
		s.poolMap[addr] = pair
		lock.Unlock()
	}

	watchdogDelay = pool.WatchDogDelay
	return addr, watchdogDelay, nil
}

// 3) StartStream: start the event subscription + processing logic for a single pool.

func (h *aerodromeSlipstreamHooks) StartStream(
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

	sink, sub, err := s.GetSwapsChannel(base, addr)
	if err != nil {
		log.Errorf("%s - error fetching swaps channel for %s: %v", h.ExchangeName(), addr.Hex(), err)
		return
	}

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if !ok {
					log.Infof("%s - swaps channel closed for %s", h.ExchangeName(), addr.Hex())
					return
				}

				swap := s.normalizeSlipstreamSwap(ctx, base, *rawSwap)

				price, volume := getSwapDataSlipstream(swap)
				if math.IsNaN(price) || math.IsInf(price, 0) {
					log.Debugf("%s - skip zero/+inf price, tx=%s", h.ExchangeName(), swap.ID)
					continue
				}

				t := makeTradeSlipstream(
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
				log.Errorf("%s - subscription error for pool %s: %v", h.ExchangeName(), addr.Hex(), err)
				// Let BaseDEXScraper resubscribe
				base.SubscribeChannel() <- addr
				return

			case <-ctx.Done():
				log.Infof("%s - shutting down stream for %s", h.ExchangeName(), addr.Hex())
				sub.Unsubscribe()
				return
			}
		}
	}()
}

// 4) OnPoolRemoved: called after stopPool to clean up poolMap.

func (h *aerodromeSlipstreamHooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addr common.Address,
	lock *sync.RWMutex,
) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()
}

// 5) OnOrderChanged: called when pool.Order changed in the config diff.

func (h *aerodromeSlipstreamHooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
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

// GetSwapsChannel subscribes to UniswapV3Pair Swap events (WS).
func (s *AerodromeSlipstreamScraper) GetSwapsChannel(
	base *BaseDEXScraper,
	poolAddress common.Address,
) (chan *UniswapV3Pair.UniswapV3PairSwap, event.Subscription, error) {

	wsClient := base.WSClient()
	if wsClient == nil {
		return nil, nil, fmt.Errorf("WS client is nil for pool %s", poolAddress.Hex())
	}

	sink := make(chan *UniswapV3Pair.UniswapV3PairSwap)

	filterer, err := UniswapV3Pair.NewUniswapV3PairFilterer(poolAddress, wsClient)
	if err != nil {
		return nil, nil, err
	}

	// Swap(sender indexed, recipient indexed, ...)
	sub, err := filterer.WatchSwap(
		&bind.WatchOpts{},
		sink,
		[]common.Address{}, // sender filter
		[]common.Address{}, // recipient filter
	)
	if err != nil {
		log.Errorf("%s - WatchSwap(%s) error: %v", base.exchange.Name, poolAddress.Hex(), err)
		return nil, nil, err
	}

	return sink, sub, nil
}

func (s *AerodromeSlipstreamScraper) normalizeSlipstreamSwap(ctx context.Context, base *BaseDEXScraper, ev UniswapV3Pair.UniswapV3PairSwap) (normalized AerodromeSlipstreamSwap) {
	pair := s.poolMap[ev.Raw.Address]

	// Scale amounts by token decimals, keep sign.
	amount0f, _ := new(big.Float).Quo(new(big.Float).SetInt(ev.Amount0), pair.Divisor0).Float64()
	amount1f, _ := new(big.Float).Quo(new(big.Float).SetInt(ev.Amount1), pair.Divisor1).Float64()

	normalized = AerodromeSlipstreamSwap{
		ID:        ev.Raw.TxHash.Hex(),
		Timestamp: time.Now().Unix(),
		Pair:      pair,
		Amount0:   amount0f,
		Amount1:   amount1f,
	}
	block, err := base.RESTClient().BlockByNumber(ctx, big.NewInt(int64(ev.Raw.BlockNumber)))
	if err == nil {
		normalized.Timestamp = int64(block.Time())
	}
	return
}

func getSwapDataSlipstream(swap AerodromeSlipstreamSwap) (price float64, volume float64) {
	volume = swap.Amount0
	if swap.Amount0 == 0 {
		log.Warnf("Invalid swap data: Amount0=0 for tx %s", swap.ID)
		return 0, 0
	}

	if swap.Amount1 == 0 {
		log.Warnf("Invalid swap data: Amount1=0 for tx %s", swap.ID)
		return 0, 0
	}

	price = math.Abs(swap.Amount1 / swap.Amount0)
	return
}

func makeTradeSlipstream(
	pair AerodromeSlipstreamPair,
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
