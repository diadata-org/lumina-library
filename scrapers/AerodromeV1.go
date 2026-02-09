package scrapers

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	velodrome "github.com/diadata-org/lumina-library/contracts/velodrome"
	"github.com/diadata-org/lumina-library/models"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type AerodromeV1Pair struct {
	Token0      models.Asset
	Token1      models.Asset
	Address     common.Address
	ForeignName string
	Order       int
}

type AerodromeV1Swap struct {
	ID         string
	Timestamp  int64
	Pair       AerodromeV1Pair
	Amount0In  float64
	Amount0Out float64
	Amount1In  float64
	Amount1Out float64
}

type AerodromeV1Scraper struct {
	base    *BaseDEXScraper
	poolMap map[common.Address]AerodromeV1Pair
}

type aerodromeV1Hooks struct {
	s            *AerodromeV1Scraper
	exchangeName string
}

func (h *aerodromeV1Hooks) ExchangeName() string { return h.exchangeName }

func (h *aerodromeV1Hooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {

	addr = common.HexToAddress(pool.Address)

	if _, ok := h.s.poolMap[addr]; !ok {
		caller, err := velodrome.NewIPoolCaller(addr, base.RESTClient())
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
		h.s.poolMap[addr] = AerodromeV1Pair{
			Address:     addr,
			Token0:      a0,
			Token1:      a1,
			ForeignName: a0.Symbol + "-" + a1.Symbol,
			Order:       pool.Order,
		}
		lock.Unlock()
	} else {
		lock.Lock()
		ap := h.s.poolMap[addr]
		ap.Order = pool.Order
		h.s.poolMap[addr] = ap
		lock.Unlock()
	}

	watchdogDelay = pool.WatchDogDelay
	if watchdogDelay <= 0 {
		watchdogDelay = 300
	}
	return addr, watchdogDelay, nil
}

func (h *aerodromeV1Hooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addr common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	h.s.watchSwaps(ctx, base, addr, tradesChan, lock)
}

func (h *aerodromeV1Hooks) OnPoolRemoved(base *BaseDEXScraper, addr common.Address, lock *sync.RWMutex) {
	lock.Lock()
	delete(h.s.poolMap, addr)
	lock.Unlock()
}

func (h *aerodromeV1Hooks) OnOrderChanged(
	base *BaseDEXScraper,
	addr common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
	lock *sync.RWMutex,
) {
	lock.Lock()
	defer lock.Unlock()
	if ap, ok := h.s.poolMap[addr]; ok {
		ap.Order = newPool.Order
		h.s.poolMap[addr] = ap
	}
}

func NewAerodromeV1Scraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) *AerodromeV1Scraper {
	log.Infof("Started %s scraper.", exchangeName)

	s := &AerodromeV1Scraper{
		poolMap: make(map[common.Address]AerodromeV1Pair),
	}
	hooks := &aerodromeV1Hooks{s: s, exchangeName: exchangeName}

	base := NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, branchMarketConfig, wg)
	s.base = base
	return s
}

func (s *AerodromeV1Scraper) watchSwaps(
	ctx context.Context,
	base *BaseDEXScraper,
	address common.Address,
	tradesChannel chan models.Trade,
	lock *sync.RWMutex,
) {
	pair, ok := s.poolMap[address]
	if !ok {
		log.Warnf("AerodromeV1 - watchSwaps called for unknown pool: %s", address.Hex())
		return
	}
	log.Infof("AerodromeV1 - subscribe to %s with pair %s", address.Hex(), pair.Token0.Symbol+"-"+pair.Token1.Symbol)

	sink, sub, err := s.getSwapsChannel(ctx, base, address)
	if err != nil {
		log.Error("AerodromeV1 - error fetching swaps channel: ", err)
		return
	}

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if ok {
					swap, err := s.normalizeSwap(ctx, base, *rawSwap)
					if err != nil {
						log.Error("AerodromeV1 - error normalizing swap: ", err)
						continue
					}
					price, volume := getSwapDataAero(swap)
					if math.IsNaN(price) || math.IsInf(price, 0) ||
						math.IsNaN(volume) || math.IsInf(volume, 0) ||
						price == 0 || volume == 0 {
						log.Debugf("AerodromeV1 - skip invalid price/volume, tx=%s", swap.ID)
						continue
					}

					t := makeTradeAerodromeV1(
						pair,
						price,
						volume,
						time.Unix(swap.Timestamp, 0),
						rawSwap.Raw.Address,
						swap.ID,
						base.exchange.Name,
						base.exchange.Blockchain,
					)

					base.UpdateLastTradeTime(rawSwap.Raw.Address, t.Time, lock)

					switch pair.Order {
					case 0:
						logTradeAerodromeV1(t)
						tradesChannel <- t
					case 1:
						t.SwapTrade()
						logTradeAerodromeV1(t)
						tradesChannel <- t
					case 2:
						logTradeAerodromeV1(t)
						tradesChannel <- t
						t.SwapTrade()
						logTradeAerodromeV1(t)
						tradesChannel <- t
					}
				}
			case err := <-sub.Err():
				log.Errorf("AerodromeV1 - subscription error for pool %s: %v", address.Hex(), err)
				base.SubscribeChannel() <- address
				return
			case <-ctx.Done():
				log.Infof("AerodromeV1 - shutting down watchSwaps for %s", address.Hex())
				sub.Unsubscribe()
				return
			}
		}
	}()
}

func (s *AerodromeV1Scraper) getSwapsChannel(
	ctx context.Context,
	base *BaseDEXScraper,
	pairAddress common.Address,
) (chan *velodrome.IPoolSwap, event.Subscription, error) {

	sink := make(chan *velodrome.IPoolSwap)

	filterer, err := velodrome.NewIPoolFilterer(pairAddress, base.WSClient())

	if err != nil {
		return nil, nil, err
	}

	sub, err := filterer.WatchSwap(&bind.WatchOpts{Context: ctx}, sink, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, nil, err
	}

	return sink, sub, nil
}

func (s *AerodromeV1Scraper) normalizeSwap(
	ctx context.Context,
	base *BaseDEXScraper,
	swap velodrome.IPoolSwap,
) (AerodromeV1Swap, error) {

	pair, ok := s.poolMap[swap.Raw.Address]
	if !ok {
		return AerodromeV1Swap{}, fmt.Errorf("pool not found in poolMap: %s", swap.Raw.Address.Hex())
	}

	decimals0 := int(pair.Token0.Decimals)
	decimals1 := int(pair.Token1.Decimals)

	amount0In, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount0In), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
	amount0Out, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount0Out), new(big.Float).SetFloat64(math.Pow10(decimals0))).Float64()
	amount1In, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount1In), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()
	amount1Out, _ := new(big.Float).Quo(new(big.Float).SetInt(swap.Amount1Out), new(big.Float).SetFloat64(math.Pow10(decimals1))).Float64()

	normalizedSwap := AerodromeV1Swap{
		ID:         swap.Raw.TxHash.Hex(),
		Timestamp:  time.Now().Unix(),
		Pair:       pair,
		Amount0In:  amount0In,
		Amount0Out: amount0Out,
		Amount1In:  amount1In,
		Amount1Out: amount1Out,
	}

	block, err := base.RESTClient().BlockByNumber(ctx, big.NewInt(int64(swap.Raw.BlockNumber)))
	if err == nil {
		normalizedSwap.Timestamp = int64(block.Time())
	}
	return normalizedSwap, nil
}

func getSwapDataAero(swap AerodromeV1Swap) (price float64, volume float64) {
	if swap.Amount0In == 0 && swap.Amount0Out == 0 {
		return 0, 0 // Invalid swap
	}
	if swap.Amount0In == 0 {
		volume = swap.Amount0Out
		price = swap.Amount1In / swap.Amount0Out
		return
	}
	volume = -swap.Amount0In
	price = swap.Amount1Out / swap.Amount0In
	return
}

func makeTradeAerodromeV1(
	pair AerodromeV1Pair,
	price float64,
	volume float64,
	timestamp time.Time,
	address common.Address,
	foreignTradeID string,
	exchangeName string,
	blockchain string,
) models.Trade {
	return models.Trade{
		Price:          price,
		Volume:         volume,
		BaseToken:      pair.Token1,
		QuoteToken:     pair.Token0,
		Time:           timestamp,
		PoolAddress:    address.Hex(),
		ForeignTradeID: foreignTradeID,
		Exchange:       models.Exchange{Name: exchangeName, Blockchain: blockchain},
	}
}

func logTradeAerodromeV1(t models.Trade) {
	log.Debugf(
		"AerodromeV1 - Got trade at time %v - symbol: %s, pair: %s, price: %v, volume:%v",
		t.Time, t.QuoteToken.Symbol, t.QuoteToken.Symbol+"-"+t.BaseToken.Symbol, t.Price, t.Volume,
	)
}
