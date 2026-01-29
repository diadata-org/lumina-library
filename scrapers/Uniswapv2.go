package scrapers

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	uniswap "github.com/diadata-org/lumina-library/contracts/uniswap/pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type UniswapPair struct {
	Token0      models.Asset
	Token1      models.Asset
	Address     common.Address
	ForeignName string
	Order       int
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
	base    *BaseDEXScraper
	poolMap map[common.Address]UniswapPair
}

type uniswapV2Hooks struct {
	s            *UniswapV2Scraper
	exchangeName string
}

func (h *uniswapV2Hooks) ExchangeName() string {
	return h.exchangeName
}

func (h *uniswapV2Hooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {

	addr = common.HexToAddress(pool.Address)

	// if poolMap does not have this address, read token0/token1 metadata for the first time
	if _, ok := h.s.poolMap[addr]; !ok {
		caller, err := uniswap.NewUniswapV2PairCaller(addr, base.RESTClient())
		if err != nil {
			return common.Address{}, 0, err
		}

		t0, err := caller.Token0(&bind.CallOpts{Context: ctx})
		if err != nil {
			return common.Address{}, 0, err
		}
		t1, err := caller.Token1(&bind.CallOpts{Context: ctx})
		if err != nil {
			return common.Address{}, 0, err
		}

		a0, err := models.GetAsset(t0, base.exchange.Blockchain, base.RESTClient())
		if err != nil {
			return common.Address{}, 0, err
		}
		a1, err := models.GetAsset(t1, base.exchange.Blockchain, base.RESTClient())
		if err != nil {
			return common.Address{}, 0, err
		}

		lock.Lock()
		h.s.poolMap[addr] = UniswapPair{
			Address:     addr,
			Token0:      a0,
			Token1:      a1,
			ForeignName: a0.Symbol + "-" + a1.Symbol,
			Order:       pool.Order,
		}
		lock.Unlock()
	} else {
		// if poolMap already has this address, only update Order
		lock.Lock()
		up := h.s.poolMap[addr]
		up.Order = pool.Order
		h.s.poolMap[addr] = up
		lock.Unlock()
	}

	watchdogDelay = pool.WatchDogDelay
	if watchdogDelay <= 0 {
		watchdogDelay = 300
	}
	return addr, watchdogDelay, nil
}

func (h *uniswapV2Hooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	h.s.watchSwaps(ctx, base, addr, tradesChan, lock)
}

func (h *uniswapV2Hooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addr common.Address,
	lock *sync.RWMutex,
) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()
}

func (h *uniswapV2Hooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
	lock *sync.RWMutex,
) {
	lock.Lock()
	defer lock.Unlock()
	if up, ok := h.s.poolMap[addr]; ok {
		up.Order = newPool.Order
		h.s.poolMap[addr] = up
	}
}

func NewUniswapV2Scraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) *UniswapV2Scraper {
	log.Infof("Started %s scraper.", exchangeName)

	s := &UniswapV2Scraper{
		poolMap: make(map[common.Address]UniswapPair),
	}
	hooks := &uniswapV2Hooks{s: s, exchangeName: exchangeName}

	base := NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, branchMarketConfig, wg)
	s.base = base

	return s
}

// watchSwaps subscribes to a UniswapV2 pool and forwards trades to the trades channel.
func (s *UniswapV2Scraper) watchSwaps(
	ctx context.Context,
	base *BaseDEXScraper,
	address common.Address,
	tradesChannel chan models.Trade,
	lock *sync.RWMutex,
) {
	pair, ok := s.poolMap[address]
	if !ok {
		log.Warnf("UniswapV2 - watchSwaps called for unknown pool: %s", address.Hex())
		return
	}
	log.Infof("UniswapV2 - subscribe to %s with pair %s", address.Hex(), pair.Token0.Symbol+"-"+pair.Token1.Symbol)

	sink, sub, err := s.getSwapsChannel(base, address)
	if err != nil {
		log.Error("UniswapV2 - error fetching swaps channel: ", err)
		return
	}

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if ok {
					swap, err := s.normalizeUniswapSwap(*rawSwap)
					if err != nil {
						log.Error("UniswapV2 - error normalizing swap: ", err)
						continue
					}
					price, volume := getSwapDataV2(swap)
					if math.IsNaN(price) || math.IsInf(price, 0) ||
						math.IsNaN(volume) || math.IsInf(volume, 0) ||
						price == 0 || volume == 0 {
						log.Debugf("UniswapV2 - skip invalid price/volume, tx=%s", swap.ID)
						continue
					}

					t := makeTradeUniswapV2(
						pair,
						price,
						volume,
						time.Unix(swap.Timestamp, 0),
						rawSwap.Raw.Address,
						swap.ID,
						base.exchange.Name,
						base.exchange.Blockchain,
					)

					// update lastTradeTime
					base.UpdateLastTradeTime(rawSwap.Raw.Address, t.Time, lock)

					switch pair.Order {
					case 0:
						logTradeUniswapV2(t)
						tradesChannel <- t
					case 1:
						t.SwapTrade()
						logTradeUniswapV2(t)
						tradesChannel <- t
					case 2:
						logTradeUniswapV2(t)
						tradesChannel <- t
						t.SwapTrade()
						logTradeUniswapV2(t)
						tradesChannel <- t
					}
				}
			case err := <-sub.Err():
				log.Errorf("UniswapV2 - subscription error for pool %s: %v", address.Hex(), err)
				// trigger BaseDEX's resub handler
				base.SubscribeChannel() <- address
				return
			case <-ctx.Done():
				log.Infof("UniswapV2 - shutting down watchSwaps for %s", address.Hex())
				sub.Unsubscribe()
				return
			}
		}
	}()
}

func (s *UniswapV2Scraper) getSwapsChannel(
	base *BaseDEXScraper,
	pairAddress common.Address,
) (chan *uniswap.UniswapV2PairSwap, event.Subscription, error) {

	sink := make(chan *uniswap.UniswapV2PairSwap)

	pairFiltererContract, err := uniswap.NewUniswapV2PairFilterer(pairAddress, base.WSClient())
	if err != nil {
		log.Error("UniswapV2 - pair filterer: ", err)
		return nil, nil, err
	}

	sub, err := pairFiltererContract.WatchSwap(&bind.WatchOpts{}, sink, []common.Address{}, []common.Address{})
	if err != nil {
		log.Error("UniswapV2 - error in get swaps channel: ", err)
		return nil, nil, err
	}

	return sink, sub, nil
}

func (s *UniswapV2Scraper) normalizeUniswapSwap(
	swap uniswap.UniswapV2PairSwap,
) (normalizedSwap UniswapSwap, err error) {

	pair, ok := s.poolMap[swap.Raw.Address]
	if !ok {
		return UniswapSwap{}, nil
	}

	decimals0 := int(pair.Token0.Decimals)
	decimals1 := int(pair.Token1.Decimals)

	amount0In, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swap.Amount0In),
		new(big.Float).SetFloat64(math.Pow10(decimals0)),
	).Float64()
	amount0Out, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swap.Amount0Out),
		new(big.Float).SetFloat64(math.Pow10(decimals0)),
	).Float64()
	amount1In, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swap.Amount1In),
		new(big.Float).SetFloat64(math.Pow10(decimals1)),
	).Float64()
	amount1Out, _ := new(big.Float).Quo(
		new(big.Float).SetInt(swap.Amount1Out),
		new(big.Float).SetFloat64(math.Pow10(decimals1)),
	).Float64()

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

// getSwapDataV2 returns price, volume and sell/buy information of @swap
func getSwapDataV2(swap UniswapSwap) (price float64, volume float64) {
	if swap.Amount0In == float64(0) {
		volume = swap.Amount0Out
		price = swap.Amount1In / swap.Amount0Out
		return
	}
	volume = -swap.Amount0In
	price = swap.Amount1Out / swap.Amount0In
	return
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
		t.Time,
		t.QuoteToken.Symbol,
		t.QuoteToken.Symbol+"-"+t.BaseToken.Symbol,
		t.Price,
		t.Volume,
	)
}
