package scrapers

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	poolmanager "github.com/diadata-org/lumina-library/contracts/uniswapv4/poolManager"
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type UniV4Pair struct {
	Token0  models.Asset
	Token1  models.Asset
	PoolID  [32]byte
	PoolHex string

	AddrKey  common.Address // last 20 bytes of poolId
	Order    int
	Divisor0 *big.Float
	Divisor1 *big.Float
}

type UniswapV4Swap struct {
	ID        string
	Timestamp int64
	Pair      UniV4Pair
	Amount0   float64 // signed, scaled
	Amount1   float64 // signed, scaled
}

type UniswapV4Scraper struct {
	base       *BaseDEXScraper
	poolMap    map[common.Address]UniV4Pair          // key=addrKey
	subByAddr  map[common.Address]event.Subscription // key=addrKey
	managerAdr common.Address
	mu         sync.RWMutex
}

type uniswapV4Hooks struct {
	s            *UniswapV4Scraper
	exchangeName string
}

func NewUniswapV4Scraper(
	ctx context.Context,
	exchangeName string,
	blockchain string,
	pools []models.Pool,
	tradesChannel chan models.Trade,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) *BaseDEXScraper {
	s := &UniswapV4Scraper{
		poolMap:   make(map[common.Address]UniV4Pair),
		subByAddr: make(map[common.Address]event.Subscription),
	}

	hooks := &uniswapV4Hooks{
		s:            s,
		exchangeName: exchangeName,
	}

	// init poolManager address from env
	envKey := strings.ToUpper(exchangeName) + "_POOL_MANAGER_ADDRESS" // UNISWAPV4_POOL_MANAGER_ADDRESS
	pm := utils.Getenv(envKey, "")
	if pm == "" {
		log.Warnf("%s - %s not set; UniV4 will not subscribe until it's configured", exchangeName, envKey)
	}
	s.managerAdr = common.HexToAddress(pm)

	// BaseDEXScraper handles ws/rest, config watch, watchdog, resub/unsub
	s.base = NewBaseDEXScraper(ctx, exchangeName, blockchain, hooks, pools, tradesChannel, branchMarketConfig, wg)
	log.Infof("Started %s scraper.", exchangeName)
	return s.base
}

func (h *uniswapV4Hooks) ExchangeName() string { return h.exchangeName }

func (h *uniswapV4Hooks) EnsurePair(
	ctx context.Context,
	base *BaseDEXScraper,
	pool models.Pool,
	lock *sync.RWMutex,
) (addr common.Address, watchdogDelay int64, err error) {
	s := h.s

	if (s.managerAdr == common.Address{}) {
		return common.Address{}, 0, fmt.Errorf("%s_POOL_MANAGER_ADDRESS not set", strings.ToUpper(h.exchangeName))
	}

	pid, err := ParsePoolIDHex(pool.Address)
	if err != nil {
		return common.Address{}, 0, fmt.Errorf("invalid poolId in pool.Address=%q: %w", pool.Address, err)
	}

	// BaseDEXScraper uses common.Address as key => last 20 bytes of poolId
	addrKey := common.BytesToAddress(pid[12:])

	// get token addresses from config Assetvolumes[0/1].Asset.Address
	t0, t1, err := getTokenAddrsFromPoolConfig(pool)
	if err != nil {
		return common.Address{}, 0, fmt.Errorf("poolId=%s: %w", PoolIDToHex(pid), err)
	}

	// fetch asset metadata by ERC20 calls (decimals/symbol/name)
	if base == nil || base.RESTClient() == nil {
		return common.Address{}, 0, fmt.Errorf("REST client is nil (check %s_URI_REST)", strings.ToUpper(base.exchange.Name))
	}

	a0, err := models.GetAsset(t0, base.exchange.Blockchain, base.RESTClient())
	if err != nil {
		return common.Address{}, 0, fmt.Errorf("GetAsset token0=%s: %w", t0.Hex(), err)
	}
	a1, err := models.GetAsset(t1, base.exchange.Blockchain, base.RESTClient())
	if err != nil {
		return common.Address{}, 0, fmt.Errorf("GetAsset token1=%s: %w", t1.Hex(), err)
	}

	d0 := new(big.Float).SetFloat64(math.Pow10(int(a0.Decimals)))
	d1 := new(big.Float).SetFloat64(math.Pow10(int(a1.Decimals)))

	pair := UniV4Pair{
		Token0:   a0,
		Token1:   a1,
		PoolID:   pid,
		PoolHex:  PoolIDToHex(pid),
		AddrKey:  addrKey,
		Order:    pool.Order,
		Divisor0: d0,
		Divisor1: d1,
	}

	s.mu.Lock()
	s.poolMap[addrKey] = pair
	s.mu.Unlock()

	watchdogDelay = pool.WatchDogDelay
	if watchdogDelay <= 0 {
		watchdogDelay = 300
	}

	return addrKey, watchdogDelay, nil
}

func (h *uniswapV4Hooks) StartStream(
	ctx context.Context,
	base *BaseDEXScraper,
	addrKey common.Address,
	tradesChan chan models.Trade,
	lock *sync.RWMutex,
) {
	s := h.s

	// get pair from poolMap
	s.mu.RLock()
	pair, ok := s.poolMap[addrKey]
	s.mu.RUnlock()
	if !ok {
		log.Warnf("%s - StartStream: pool %s not found in poolMap", h.ExchangeName(), addrKey.Hex())
		return
	}

	log.Infof("%s - subscribe to %s with pair %s ",
		h.ExchangeName(), pair.PoolHex, pair.Token0.Symbol+"-"+pair.Token1.Symbol)

	sink, sub, err := s.GetSwapsChannel(base, pair)
	if err != nil {
		log.Errorf("%s - error fetching swaps channel for %s (poolId=%s): %v", h.ExchangeName(), addrKey.Hex(), pair.PoolHex, err)
		return
	}

	// store sub so OnPoolRemoved can unsubscribe
	s.mu.Lock()
	if old := s.subByAddr[addrKey]; old != nil {
		old.Unsubscribe()
	}
	s.subByAddr[addrKey] = sub
	s.mu.Unlock()

	go func() {
		for {
			select {
			case rawSwap, ok := <-sink:
				if !ok {
					log.Infof("%s - swaps channel closed for %s", h.ExchangeName(), addrKey.Hex())
					return
				}

				swap := s.normalizeUniV4Swap(*rawSwap, pair)
				price, volume := getSwapDataUniV4(swap)
				if math.IsNaN(price) || math.IsInf(price, 0) || price == 0 {
					log.Debugf("%s - skip zero/+inf price, tx=%s", h.ExchangeName(), swap.ID)
					continue
				}

				t := makeTradeUniV4(
					pair,
					price,
					volume,
					time.Unix(swap.Timestamp, 0),
					swap.ID,
					base.exchange.Name,
					base.exchange.Blockchain,
				)

				base.UpdateLastTradeTime(addrKey, t.Time, lock)

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
				log.Errorf("%s - subscription error for pool %s : %v ", h.ExchangeName(), pair.PoolHex, err)
				base.SubscribeChannel() <- addrKey
				return

			case <-ctx.Done():
				log.Infof("%s - shutting down stream for %s ", h.ExchangeName(), pair.PoolHex)
				sub.Unsubscribe()
				return
			}
		}
	}()
}

func (h *uniswapV4Hooks) OnPoolRemoved(
	base *BaseDEXScraper,
	addrKey common.Address,
	lock *sync.RWMutex,
) {
	s := h.s
	s.mu.Lock()
	if sub := s.subByAddr[addrKey]; sub != nil {
		sub.Unsubscribe()
		delete(s.subByAddr, addrKey)
	}
	delete(s.poolMap, addrKey)
	s.mu.Unlock()
}

func (h *uniswapV4Hooks) OnOrderChanged(
	base *BaseDEXScraper,
	addrKey common.Address,
	oldPool models.Pool,
	newPool models.Pool,
	branchMarketConfig string,
	lock *sync.RWMutex,
) {
	s := h.s
	s.mu.Lock()
	defer s.mu.Unlock()

	pair, ok := s.poolMap[addrKey]
	if !ok {
		return
	}
	pair.Order = newPool.Order
	s.poolMap[addrKey] = pair
}

// ============ UniV4 channel + normalize ============

func (s *UniswapV4Scraper) GetSwapsChannel(
	base *BaseDEXScraper,
	pair UniV4Pair,
) (chan *poolmanager.PoolManagerSwap, event.Subscription, error) {

	wsClient := base.WSClient()
	if wsClient == nil {
		return nil, nil, fmt.Errorf("WS client is nil for poolId %s", pair.PoolHex)
	}

	sink := make(chan *poolmanager.PoolManagerSwap)

	filterer, err := poolmanager.NewPoolManagerFilterer(s.managerAdr, wsClient)
	if err != nil {
		return nil, nil, err
	}

	sub, err := filterer.WatchSwap(
		&bind.WatchOpts{Context: context.Background()},
		sink,
		[][32]byte{pair.PoolID},
		[]common.Address{},
	)
	if err != nil {
		log.Errorf("%s - WatchSwap(poolId=%s) error: %v", base.exchange.Name, pair.PoolHex, err)
		return nil, nil, err
	}

	return sink, sub, nil
}

// normalizeUniV4Swap converts raw event amounts into scaled float64 amounts (signed).
func (s *UniswapV4Scraper) normalizeUniV4Swap(ev poolmanager.PoolManagerSwap, pair UniV4Pair) (normalizedSwap UniswapV4Swap) {
	amount0, _ := new(big.Float).Quo(new(big.Float).SetInt(ev.Amount0), pair.Divisor0).Float64()
	amount1, _ := new(big.Float).Quo(new(big.Float).SetInt(ev.Amount1), pair.Divisor1).Float64()

	id := ev.Raw.TxHash.Hex()

	id = fmt.Sprintf("%s:%d", id, ev.Raw.Index)

	normalizedSwap = UniswapV4Swap{
		ID:        id,
		Timestamp: time.Now().Unix(),
		Pair:      pair,
		Amount0:   amount0,
		Amount1:   amount1,
	}
	return
}

func getSwapDataUniV4(swap UniswapV4Swap) (price float64, volume float64) {
	volume = swap.Amount0
	if swap.Amount0 == 0 {
		log.Warnf("Invalid swap data: Amount0=0 for tx %s", swap.ID)
		return 0, 0
	}
	price = math.Abs(swap.Amount1 / swap.Amount0)
	return
}

func makeTradeUniV4(
	pair UniV4Pair,
	price float64,
	volume float64,
	timestamp time.Time,
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
		PoolAddress:    pair.PoolHex, // v4: poolId
		ForeignTradeID: foreignTradeID,
		Exchange:       models.Exchange{Name: exchangeName, Blockchain: blockchain},
	}
}

// ============ Config helpers ============

func getTokenAddrsFromPoolConfig(pool models.Pool) (common.Address, common.Address, error) {
	if len(pool.Assetvolumes) < 2 {
		return common.Address{}, common.Address{}, fmt.Errorf("pool has <2 assets in Assetvolumes")
	}

	a0 := strings.TrimSpace(pool.Assetvolumes[0].Asset.Address)
	a1 := strings.TrimSpace(pool.Assetvolumes[1].Asset.Address)
	if a0 == "" || a1 == "" {
		return common.Address{}, common.Address{}, fmt.Errorf("empty token address in Assetvolumes (a0=%q a1=%q)", a0, a1)
	}
	return common.HexToAddress(a0), common.HexToAddress(a1), nil
}

// ParsePoolIDHex parses "0x" + 64 hex chars into [32]byte poolId.
func ParsePoolIDHex(v string) ([32]byte, error) {
	var out [32]byte
	x := strings.TrimSpace(v)
	x = strings.TrimPrefix(x, "0x")
	if len(x) != 64 {
		return out, fmt.Errorf("poolId must be 32 bytes hex (64 chars), got len=%d", len(x))
	}
	b, err := hex.DecodeString(x)
	if err != nil {
		return out, fmt.Errorf("decode hex: %w", err)
	}
	copy(out[:], b)
	return out, nil
}

func PoolIDToHex(id [32]byte) string {
	return "0x" + hex.EncodeToString(id[:])
}
