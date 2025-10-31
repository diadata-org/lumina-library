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
	exchange           models.Exchange
	subscribeChannel   chan common.Address
	unsubscribeChannel chan common.Address
	lastTradeTimeMap   map[common.Address]time.Time
	poolMap            map[common.Address]UniswapPair
	wsClient           *ethclient.Client
	restClient         *ethclient.Client
	waitTime           int
	streamCancel       map[common.Address]context.CancelFunc
	watchdogCancel     map[string]context.CancelFunc
}

func NewUniswapV2Scraper(ctx context.Context, exchangeName string, blockchain string, pools []models.Pool, tradesChannel chan models.Trade, wg *sync.WaitGroup) {
	var err error
	var scraper UniswapV2Scraper
	log.Info("Started UniswapV2 scraper.")

	scraper.exchange.Name = exchangeName
	scraper.exchange.Blockchain = blockchain
	scraper.subscribeChannel = make(chan common.Address)
	scraper.unsubscribeChannel = make(chan common.Address)
	scraper.watchdogCancel = make(map[string]context.CancelFunc)
	scraper.streamCancel = make(map[common.Address]context.CancelFunc)
	scraper.lastTradeTimeMap = make(map[common.Address]time.Time)
	scraper.poolMap = make(map[common.Address]UniswapPair)
	scraper.waitTime, err = strconv.Atoi(utils.Getenv(strings.ToUpper(exchangeName)+"_WAIT_TIME", "500"))

	if err != nil {
		log.Errorf("Failed to parse waitTime for exchange %s: %v", exchangeName, err)
	}

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(exchangeName)+"_URI_REST", restDial))
	if err != nil {
		log.Error("UniswapV2 - init rest client: ", err)
	}
	scraper.wsClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(exchangeName)+"_URI_WS", wsDial))
	if err != nil {
		log.Error("UniswapV2 - init ws client: ", err)
	}

	scraper.makeUniPoolMap(pools)

	var lock sync.RWMutex
	// initial pools
	scraper.startInitialPools(ctx, pools, tradesChannel, &lock)

	// resub handler triggered by watchdog
	scraper.startResubHandler(ctx, tradesChannel, &lock)

	// unsubscribe handler
	scraper.startUnsubHandler(ctx, &lock)
	// watch config updates
	go scraper.watchConfig(ctx, exchangeName, tradesChannel, &lock)
}

func (scraper *UniswapV2Scraper) startInitialPools(ctx context.Context, pools []models.Pool, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	started := map[string]bool{}
	for _, p := range pools {
		lower := strings.ToLower(p.Address)
		if started[lower] {
			continue
		}
		if err := scraper.startPool(ctx, p, tradesChannel, lock); err != nil {
			log.Errorf("UniswapV2 - startPool %s: %v", p.Address, err)
		}
		started[lower] = true
	}
}

func (scraper *UniswapV2Scraper) startResubHandler(ctx context.Context, tradesChannel chan models.Trade, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-scraper.subscribeChannel:
				log.Infof("Resubscribing to pool: %s", addr.Hex())
				lock.Lock()
				scraper.lastTradeTimeMap[addr] = time.Now()
				if c, ok := scraper.streamCancel[addr]; ok && c != nil {
					c()
					delete(scraper.streamCancel, addr)
				}
				lock.Unlock()
				go scraper.ListenToPair(ctx, addr, tradesChannel, lock)
			case <-ctx.Done():
				log.Info("Stopping resubscription handler.")
				return
			}
		}
	}()
}

// stop pool when address is received
func (s *UniswapV2Scraper) startUnsubHandler(ctx context.Context, lock *sync.RWMutex) {
	go func() {
		for {
			select {
			case addr := <-s.unsubscribeChannel:
				log.Infof("UniswapV2 - Unsubscribing pool: %s", addr.Hex())
				s.stopPool(addr, lock)
			case <-ctx.Done():
				log.Info("UniswapV2 - unsubscribe handler stopped")
				return
			}
		}
	}()
}

// start pool: add to poolMap -> init lastTradeTime -> restart watchdog -> start listening
func (s *UniswapV2Scraper) startPool(ctx context.Context, p models.Pool, trades chan models.Trade, lock *sync.RWMutex) error {
	addr := common.HexToAddress(p.Address)

	// 1) add to poolMap (if exists, only update Order)
	if _, ok := s.poolMap[addr]; !ok {
		// first time: read token0/token1 metadata
		caller, err := uniswap.NewUniswapV2PairCaller(addr, s.restClient)
		if err != nil {
			return err
		}

		t0, err := caller.Token0(&bind.CallOpts{})
		if err != nil {
			return err
		}
		t1, err := caller.Token1(&bind.CallOpts{})
		if err != nil {
			return err
		}

		a0, err := models.GetAsset(t0, s.exchange.Blockchain, s.restClient)
		if err != nil {
			return err
		}
		a1, err := models.GetAsset(t1, s.exchange.Blockchain, s.restClient)
		if err != nil {
			return err
		}

		s.poolMap[addr] = UniswapPair{
			Address:     addr,
			Token0:      a0,
			Token1:      a1,
			ForeignName: a0.Symbol + "-" + a1.Symbol,
			Order:       p.Order,
		}
	} else {
		up := s.poolMap[addr]
		up.Order = p.Order
		s.poolMap[addr] = up
	}

	// 2) init lastTradeTime, avoid triggering watchdog immediately
	lock.Lock()
	if _, ok := s.lastTradeTimeMap[addr]; !ok {
		s.lastTradeTimeMap[addr] = time.Now()
	}
	lock.Unlock()

	// 3) restart watchdog (according to WatchDogDelay)
	delay := p.WatchDogDelay
	if delay <= 0 {
		delay = 300
	}
	s.restartWatchdogForAddr(ctx, addr, delay, lock)

	// 4) start listening (if already running, don't repeat)
	if _, ok := s.streamCancel[addr]; ok {
		return nil
	}
	pctx, cancel := context.WithCancel(ctx)
	s.streamCancel[addr] = cancel
	go s.ListenToPair(pctx, addr, trades, lock)
	return nil
}

func (s *UniswapV2Scraper) stopPool(addr common.Address, lock *sync.RWMutex) {
	if c, ok := s.streamCancel[addr]; ok && c != nil {
		c()
		delete(s.streamCancel, addr)
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

func (s *UniswapV2Scraper) restartWatchdogForAddr(ctx context.Context, addr common.Address, delay int64, lock *sync.RWMutex) {
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
					s.subscribeChannel <- addr // trigger resubscribe
				}
			case <-wdCtx.Done():
				return
			}
		}
	}()
}

func mapPoolsByAddrLower(pools []models.Pool) map[string]models.Pool {
	m := make(map[string]models.Pool, len(pools))
	for _, p := range pools {
		m[strings.ToLower(p.Address)] = p
	}
	return m
}

func (s *UniswapV2Scraper) watchConfig(ctx context.Context, exchangeName string, trades chan models.Trade, lock *sync.RWMutex) {
	tk := time.NewTicker(60 * time.Second)
	defer tk.Stop()

	// load initial pools
	cur, err := models.PoolsFromConfigFile(exchangeName)
	if err != nil {
		log.Errorf("UniswapV2 - load initial pools: %v", err)
		cur = nil
	}
	last := mapPoolsByAddrLower(cur)

	// start all initial addresses
	for _, p := range last {
		if err := s.startPool(ctx, p, trades, lock); err != nil {
			log.Errorf("UniswapV2 - startPool %s: %v", p.Address, err)
		}
	}

	for {
		select {
		case <-tk.C:
			nowList, err := models.PoolsFromConfigFile(exchangeName)
			if err != nil {
				log.Errorf("UniswapV2 - reload pools: %v", err)
				continue
			}
			now := mapPoolsByAddrLower(nowList)
			s.applyConfigDiff(ctx, last, now, trades, lock)
			last = now

		case <-ctx.Done():
			log.Info("UniswapV2 - watchConfig exit")
			return
		}
	}
}

func (s *UniswapV2Scraper) applyConfigDiff(
	ctx context.Context,
	last map[string]models.Pool,
	curr map[string]models.Pool,
	trades chan models.Trade,
	lock *sync.RWMutex,
) {
	// remove
	for k := range last {
		if _, ok := curr[k]; !ok {
			addr := common.HexToAddress(k)
			log.Infof("UniswapV2 - remove pool %s", addr.Hex())
			s.unsubscribeChannel <- addr
		}
	}

	// add + update
	for k, p := range curr {
		addr := common.HexToAddress(k)
		old, existed := last[k]
		if !existed {
			// add
			log.Infof("UniswapV2 - add pool %s (wd=%d, order=%d)", p.Address, p.WatchDogDelay, p.Order)
			if err := s.startPool(ctx, p, trades, lock); err != nil {
				log.Errorf("UniswapV2 - startPool %s: %v", p.Address, err)
			}
			continue
		}

		// update Order
		if up, ok := s.poolMap[addr]; ok && up.Order != p.Order {
			up.Order = p.Order
			s.poolMap[addr] = up
			log.Infof("UniswapV2 - update order %s: %d -> %d", addr.Hex(), old.Order, p.Order)
		}

		// update watchdog (restart directly, simple and reliable; or compare old/new and decide)
		newWD := p.WatchDogDelay
		oldWD := old.WatchDogDelay
		if newWD != oldWD {
			log.Infof("UniswapV2 - update watchdog %s: %d -> %d", addr.Hex(), oldWD, newWD)
			s.restartWatchdogForAddr(ctx, addr, newWD, lock)
		}
	}
}

// makeUniPoolMap returns a map with pool addresses as keys and the underlying UniswapPair as values.
func (scraper *UniswapV2Scraper) makeUniPoolMap(pools []models.Pool) error {
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
			Order:       p.Order,
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
