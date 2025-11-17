package simulators

import (
	"context"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefactory"
	"github.com/diadata-org/lumina-library/contracts/curve/curvefi"
	"github.com/diadata-org/lumina-library/models"
	simulation "github.com/diadata-org/lumina-library/simulations/simulators/curve"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PoolMeta struct {
	GetUnderlyingCoins bool
	QuoteIdx           int
	BaseIdx            int
}

type CurveSimulator struct {
	restClient        *ethclient.Client
	luminaClient      *ethclient.Client
	simulator         *simulation.Simulator
	exchangepairs     []models.ExchangePair
	thresholdSlippage float64
	priceMap          map[models.Asset]models.AssetQuotation
}

var (
	restDialUrl       = ""
	registryAddresses = map[string]string{
		"base":       "0x90E00ACe148ca3b23Ac1bC8C240C2a7Dd9c2d7f5",
		"cryptoswap": "0x8F942C20D02bEfc377D41445793068908E2250D0",
		"meta":       "0xB9fC157394Af804a3578134A6585C0dc9cc990d4",
		"factory":    "0xF18056Bbd320E96A48e3Fbf8bC061322531aac99",
		"factory2":   "0x4F8846Ae9380B90d2E71D5e3D042dff3E7ebb40d",
	}
	DIAMetaContractAddress   = "0x0087342f5f4c7AB23a37c045c3EF710749527c88"
	DIAMetaContractPrecision = 8
	amountIn_USD_constant    = float64(100)
	simulationUpdateSeconds  = 30
	priceMapUpdateSeconds    = 30 * 60
	thresholdSlippage        = 3
	liquidityThresholdUSD    = big.NewFloat(50000)
	liquidityThresholdNative = big.NewFloat(2)
)

func init() {
	var err error

	simulationUpdateSeconds, err = strconv.Atoi(utils.Getenv(strings.ToUpper(CURVE_SIMULATION)+"_SIMULATION_UPDATE_SECONDS", strconv.Itoa(simulationUpdateSeconds)))
	if err != nil {
		log.Errorf(strings.ToUpper(CURVE_SIMULATION)+"_SIMULATION_UPDATE_SECONDS: %v", err)
	}

	priceMap_Update_SecondsVersion2, err = strconv.Atoi(utils.Getenv(strings.ToUpper(CURVE_SIMULATION)+"_PRICE_MAP_UPDATE_SECONDS", strconv.Itoa(priceMapUpdateSeconds)))
	if err != nil {
		log.Errorf(strings.ToUpper(CURVE_SIMULATION)+"_PRICE_MAP_UPDATE_SECONDS: %v", err)
	}

}

func NewCurveSimulator(exchangepairs []models.ExchangePair, tradesChannel chan models.SimulatedTrade) {
	var (
		err     error
		scraper CurveSimulator
	)

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(CURVE_SIMULATION)+"_URI_REST", restDialUrl))
	if err != nil {
		log.Error("init rest client: ", err)
	} else {
		log.Info("Successfully connected to node")
	}
	defer scraper.restClient.Close()

	scraper.luminaClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(CURVE_SIMULATION)+"_LUMINA_URI_REST", restDialLuminaVersion2))
	if err != nil {
		log.Error("init lumina client: ", err)
	} else {
		log.Info("Successfully connected to lumina node")
	}
	defer scraper.luminaClient.Close()

	scraper.thresholdSlippage, err = strconv.ParseFloat(utils.Getenv("CURVE_THRESHOLD_SLIPPAGE", strconv.Itoa(thresholdSlippage)), 64)
	if err != nil {
		log.Error("Parse THRESHOLD_SLIPPAGE: ", err)
	}

	scraper.simulator = simulation.New(scraper.restClient, log)
	scraper.exchangepairs = exchangepairs
	err = scraper.getExchangePairs()
	if err != nil {
		log.Fatal("Failed to get exchange pairs: ", err)
	}

	var lock sync.RWMutex
	scraper.updatePriceMap(&lock)

	priceTicker := time.NewTicker(time.Duration(priceMapUpdateSeconds) * time.Second)
	go func() {
		for range priceTicker.C {
			scraper.updatePriceMap(&lock)
		}
	}()

	allPools := make(map[models.ExchangePair]map[common.Address]PoolMeta)
	for registryName, registryAddress := range registryAddresses {
		registry := scraper.getRegistry(registryName, registryAddress)
		for _, ep := range scraper.exchangepairs {
			pools := scraper.getPools(registry, ep)
			if allPools[ep] == nil {
				allPools[ep] = make(map[common.Address]PoolMeta)
			}
			for k, v := range pools {
				allPools[ep][k] = v
			}
		}
	}

	ticker := time.NewTicker(time.Duration(simulationUpdateSeconds) * time.Second)
	for range ticker.C {
		scraper.simulateTrades(allPools, tradesChannel)
	}
}

func (scraper *CurveSimulator) getExchangePairs() (err error) {
	scraper.priceMap = make(map[models.Asset]models.AssetQuotation)
	for i, ep := range scraper.exchangepairs {
		quoteToken, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address), Exchanges[CURVE_SIMULATION].Blockchain, scraper.restClient)
		if err != nil {
			return err
		}
		scraper.exchangepairs[i].UnderlyingPair.QuoteToken = quoteToken
		baseToken, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.BaseToken.Address), Exchanges[CURVE_SIMULATION].Blockchain, scraper.restClient)
		if err != nil {
			return err
		}
		scraper.exchangepairs[i].UnderlyingPair.BaseToken = baseToken
		scraper.priceMap[quoteToken] = models.AssetQuotation{}
		scraper.priceMap[baseToken] = models.AssetQuotation{}
	}
	return
}

func (scraper *CurveSimulator) updatePriceMap(lock *sync.RWMutex) {
	for asset := range scraper.priceMap {
		quotation, err := asset.GetPrice(common.HexToAddress(DIAMetaContractAddress), DIAMetaContractPrecision, scraper.luminaClient)
		if err != nil {
			log.Errorf("GetOnchainPrice for %s -- %s: %v", asset.Symbol, asset.Address, err)
			quotation.Price = scraper.getPriceFromAPI(asset)
		} else {
			log.Infof("USD price for (base-)token %s: %v", asset.Symbol, quotation.Price)
		}
		if quotation.Price == 0 {
			quotation.Price = scraper.getPriceFromAPI(asset)
		}
		lock.Lock()
		scraper.priceMap[asset] = quotation
		lock.Unlock()
	}
}

func (scraper *CurveSimulator) getPriceFromAPI(asset models.Asset) float64 {
	log.Warnf("Could not determine price of %s on chain. Checking DIA API.", asset.Symbol)
	assetQuotation, err := asset.GetPriceFromDiaAPI()
	if err != nil {
		log.Errorf("Failed to get price of %s from DIA API: %v", asset.Symbol, err)
		log.Errorf("asset blockchain -- address: %v -- %v", asset.Blockchain, asset.Address)
		return 100
	}
	return assetQuotation.Price
}

func (scraper *CurveSimulator) getRegistry(registryName string, registryAddress string) interface{} {
	var (
		registry interface{}
		err      error
	)

	switch registryAddress {
	case "0x90E00ACe148ca3b23Ac1bC8C240C2a7Dd9c2d7f5", //base
		"0x8F942C20D02bEfc377D41445793068908E2250D0": // cryptoswap
		registry, err = curvefi.NewCurvefi(common.HexToAddress(registryAddress), scraper.restClient)

	case "0xB9fC157394Af804a3578134A6585C0dc9cc990d4", //meta
		"0xF18056Bbd320E96A48e3Fbf8bC061322531aac99", // factory
		"0x4F8846Ae9380B90d2E71D5e3D042dff3E7ebb40d": //factory2
		registry, err = curvefactory.NewCurvefactory(common.HexToAddress(registryAddress), scraper.restClient)

	default:
		log.Errorf("Unsupported registry address: %s", registryAddress)
	}

	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	} else {
		log.Info("Successfully created contract instance")
	}

	poolCount, err := scraper.getPoolCount(registry)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Total # of %v pools: %d\n", registryName, poolCount)

	return registry
}

func (scraper *CurveSimulator) getPools(registry interface{}, ep models.ExchangePair) map[common.Address]PoolMeta {
	// Retrieve all pools
	pools := make(map[common.Address]PoolMeta)
	for i := 0; ; i++ {
		pool, err := scraper.getPool(registry, ep, i)
		if err != nil {
			log.Errorf("Error querying index %d: %v", i, err)
			break
		}
		if pool == (common.Address{}) {
			break
		}

		balances, balanceErr := scraper.getBalances(registry, pool)
		if balanceErr != nil {
			log.Warnf("Skipping pool %s: GetBalances failed: %v", pool.Hex(), balanceErr)
			continue
		}

		tokens, hasWrappedCoins, tokenErr := scraper.getTokens(registry, pool, ep)
		if tokenErr != nil {
			log.Warnf("Skipping pool %s: tokens do not match or call failed", pool.Hex())
			continue
		}

		quoteIdx, baseIdx, indexOK := findTokenIndices(tokens, ep)
		if !indexOK {
			log.Warnf("Skipping pool %s: unable to determine token indices", pool.Hex())
			continue
		}

		if scraper.hasSufficientLiquidity(ep, pool, balances[baseIdx], balances[quoteIdx]) {
			pools[pool] = PoolMeta{GetUnderlyingCoins: hasWrappedCoins, QuoteIdx: quoteIdx, BaseIdx: baseIdx}
			functionUnsed := "registry.GetCoins"
			if hasWrappedCoins {
				functionUnsed = "registry.GetUnderlyingCoins"
			}
			log.Infof("Pool #%d: %s | Pool has wrapped coins: %v | Function Used: %s", i+1, pool.Hex(), hasWrappedCoins, functionUnsed)
			log.Infof("Token validated! token0_index: %v, token1_index: %v\n", quoteIdx, baseIdx)
		}
	}

	log.Infof("Found %d %v/%v liquidity pools:\n", len(pools), ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
	return pools
}

func (scraper *CurveSimulator) getPoolCount(registry interface{}) (*big.Int, error) {
	var poolCount *big.Int
	var err error
	switch r := registry.(type) {
	case *curvefi.Curvefi:
		poolCount, err = r.PoolCount(nil)
	case *curvefactory.Curvefactory:
		poolCount, err = r.PoolCount(nil)
	default:
		log.Fatal("getPoolCount - Unknown registry type")
	}
	return poolCount, err
}

func (scraper *CurveSimulator) getPool(registry interface{}, ep models.ExchangePair, i int) (common.Address, error) {
	var pool common.Address
	var err error
	switch r := registry.(type) {
	case *curvefi.Curvefi:
		pool, err = r.FindPoolForCoins0(
			&bind.CallOpts{Context: context.Background()},
			common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
			common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
			big.NewInt(int64(i)),
		)
	case *curvefactory.Curvefactory:
		pool, err = r.FindPoolForCoins0(
			&bind.CallOpts{Context: context.Background()},
			common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
			common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
			big.NewInt(int64(i)),
		)
	default:
		log.Fatal("getPool - Unknown registry type")
	}
	return pool, err
}

func (scraper *CurveSimulator) getBalances(registry interface{}, pool common.Address) ([]*big.Int, error) {
	var result []*big.Int
	var err error

	switch r := registry.(type) {
	case *curvefi.Curvefi:
		balances, balanceErr := r.GetBalances(&bind.CallOpts{Context: context.Background()}, pool)
		err = balanceErr
		result = balances[:]
	case *curvefactory.Curvefactory:
		balances, balanceErr := r.GetBalances(&bind.CallOpts{Context: context.Background()}, pool)
		err = balanceErr
		result = balances[:]
	default:
		log.Fatal("getBalances - Unknown registry type")
	}
	return result, err
}

func (scraper *CurveSimulator) getTokens(registry interface{}, pool common.Address, ep models.ExchangePair) ([]common.Address, bool, error) {
	var result []common.Address
	var err error
	var hasWrappedCoins bool = false

	switch r := registry.(type) {
	case *curvefi.Curvefi:
		tokens, tokenErr := r.GetCoins(&bind.CallOpts{Context: context.Background()}, pool)

		if tokenErr != nil || !matchTokens(tokens[:], ep) {
			tokens, tokenErr = r.GetUnderlyingCoins(&bind.CallOpts{Context: context.Background()}, pool)
			if tokenErr != nil || !matchTokens(tokens[:], ep) {
				return result, hasWrappedCoins, tokenErr
			}
			hasWrappedCoins = true
		}
		result = tokens[:]
	case *curvefactory.Curvefactory:
		tokens, tokenErr := r.GetCoins(&bind.CallOpts{Context: context.Background()}, pool)
		if tokenErr != nil || !matchTokens(tokens[:], ep) {
			return result, hasWrappedCoins, tokenErr
		}
		result = tokens[:]
	default:
		log.Fatal("getTokens - Unknown registry type")
	}
	return result, hasWrappedCoins, err
}

func matchTokens(tokens []common.Address, ep models.ExchangePair) bool {
	var baseMatch, quoteMatch bool
	for _, token := range tokens {
		if token == common.HexToAddress(ep.UnderlyingPair.BaseToken.Address) {
			baseMatch = true
		}
		if token == common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address) {
			quoteMatch = true
		}
	}
	return baseMatch && quoteMatch
}

func findTokenIndices(tokens []common.Address, ep models.ExchangePair) (quoteIndex, baseIndex int, ok bool) {
	quoteIndex, baseIndex = -1, -1
	for i, token := range tokens {
		if token == common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address) {
			quoteIndex = i
		}
		if token == common.HexToAddress(ep.UnderlyingPair.BaseToken.Address) {
			baseIndex = i
		}
	}
	return quoteIndex, baseIndex, quoteIndex != -1 && baseIndex != -1
}

func (scraper *CurveSimulator) simulateTrades(allPools map[models.ExchangePair]map[common.Address]PoolMeta, tradesChannel chan models.SimulatedTrade) {
	var wg sync.WaitGroup

	for ep, pools := range allPools {
		for poolAddr, meta := range pools {
			wg.Add(1)
			go func(ep models.ExchangePair, w *sync.WaitGroup) {
				defer w.Done()
				log.Infof("============== Pool Addr: %v=============\n", poolAddr.Hex())

				// Prepare trade input (e.g., $100)
				baseTokenPrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price
				amountInBase := amountIn_USD_constant / baseTokenPrice // 100

				decimals := big.NewInt(int64(ep.UnderlyingPair.BaseToken.Decimals)) // e.g. 18
				exponent := new(big.Int).Exp(big.NewInt(10), decimals, nil)         // e.g. 10^18
				exponentFloat := new(big.Float).SetInt(exponent)

				amountIn := new(big.Float).Mul(big.NewFloat(amountInBase), exponentFloat) // e.g. 10^20
				amountInInt := new(big.Int)
				amountIn.Int(amountInInt)
				amountInAfterDecimalAdjust := new(big.Float).Quo(amountIn, exponentFloat) // e.g. 10^2
				log.Infof("amountIn after adjusting for decimals: %v\n", amountInAfterDecimalAdjust)

				// Run trade simulation
				var amountOutFloat *big.Float
				amountOut, err := scraper.simulator.Execute(poolAddr, scraper.restClient, meta.GetUnderlyingCoins, meta.BaseIdx, meta.QuoteIdx, amountInInt)
				if err == nil {
					amountOutFloat = new(big.Float).SetInt(amountOut)
					power := ep.UnderlyingPair.QuoteToken.Decimals
					divisor := new(big.Float).SetInt(
						new(big.Int).Exp(
							big.NewInt(10),
							big.NewInt(int64(power)),
							nil,
						),
					)
					amountOutFloat.Quo(amountOutFloat, divisor)
					log.Infof("amountOut: %v\n", amountOutFloat)
					amountOutAfterDecimalAdjustF64, _ := amountOutFloat.Float64()
					price, _ := new(big.Float).Quo(amountInAfterDecimalAdjust, amountOutFloat).Float64()

					t := models.SimulatedTrade{
						Price:       price,
						Volume:      amountOutAfterDecimalAdjustF64,
						QuoteToken:  ep.UnderlyingPair.QuoteToken,
						BaseToken:   ep.UnderlyingPair.BaseToken,
						PoolAddress: poolAddr.Hex(),
						Time:        time.Now(),
						Exchange:    Exchanges[CURVE_SIMULATION],
					}

					log.Infof("Got trade in pool %v%%: %s-%s -- %v -- %v", poolAddr.Hex(), t.QuoteToken.Symbol, t.BaseToken.Symbol, t.Price, t.Volume)
					tradesChannel <- t
				}
			}(ep, &wg)
		}
	}
}

func (scraper *CurveSimulator) hasSufficientLiquidity(ep models.ExchangePair, pool common.Address, baseBalance *big.Int, quoteBalance *big.Int) bool {
	baseDecimals := float64(ep.UnderlyingPair.BaseToken.Decimals)
	quoteDecimals := float64(ep.UnderlyingPair.QuoteToken.Decimals)

	baseBalanceF := new(big.Float).Quo(new(big.Float).SetInt(baseBalance), big.NewFloat(math.Pow10(int(baseDecimals))))
	quoteBalanceF := new(big.Float).Quo(new(big.Float).SetInt(quoteBalance), big.NewFloat(math.Pow10(int(quoteDecimals))))

	if baseBalanceF.Cmp(liquidityThresholdNative) < 0 || quoteBalanceF.Cmp(liquidityThresholdNative) < 0 {
		log.Warnf("Native liquidity not sufficient for pool %v: base=%s %s, quote=%s %s",
			pool.Hex(),
			baseBalanceF.Text('f', 4), ep.UnderlyingPair.BaseToken.Symbol,
			quoteBalanceF.Text('f', 4), ep.UnderlyingPair.QuoteToken.Symbol)
		return false
	}

	// USD threshold check
	baseTokenPrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price
	quoteTokenPrice := scraper.priceMap[ep.UnderlyingPair.QuoteToken].Price

	baseUSD := new(big.Float).Mul(baseBalanceF, big.NewFloat(baseTokenPrice))
	quoteUSD := new(big.Float).Mul(quoteBalanceF, big.NewFloat(quoteTokenPrice))

	if baseUSD.Cmp(liquidityThresholdUSD) < 0 {
		log.Warnf("Base token %s has insufficient USD liquidity for pool %v: %s", ep.UnderlyingPair.BaseToken.Symbol, pool.Hex(), baseUSD.Text('f', 2))
		return false
	}
	if quoteUSD.Cmp(liquidityThresholdUSD) < 0 {
		log.Warnf("Quote token %s has insufficient USD liquidity for pool %v: %s", ep.UnderlyingPair.QuoteToken.Symbol, pool.Hex(), quoteUSD.Text('f', 2))
		return false
	}

	return true
}
