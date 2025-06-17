package simulators

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefi"
	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	"github.com/diadata-org/lumina-library/contracts/curve/curveplain"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	"github.com/diadata-org/lumina-library/models"
	simulation "github.com/diadata-org/lumina-library/simulations/simulators/curve"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CurveSimulator struct {
	restClient        *ethclient.Client
	luminaClient      *ethclient.Client
	simulator         *simulation.Simulator
	exchangepairs     []models.ExchangePair
	thresholdSlippage float64
	priceMap          map[models.Asset]models.AssetQuotation
}

var (
	restDialUrl              = ""
	registryAddress          = "0x90E00ACe148ca3b23Ac1bC8C240C2a7Dd9c2d7f5"
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
	scraper.getExchangePairs()

	var lock sync.RWMutex
	scraper.updatePriceMap(&lock)

	priceTicker := time.NewTicker(time.Duration(priceMapUpdateSeconds) * time.Second)
	go func() {
		for range priceTicker.C {
			scraper.updatePriceMap(&lock)
		}
	}()

	var pools map[common.Address]string
	for _, ep := range scraper.exchangepairs {
		pools = scraper.getPools(ep)
	}

	ticker := time.NewTicker(time.Duration(simulationUpdateSeconds) * time.Second)
	for range ticker.C {
		for _, ep := range scraper.exchangepairs {
			scraper.simulateTrades(ep, pools, tradesChannel)
		}
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
		quotation, err := asset.GetOnchainPrice(common.HexToAddress(DIAMetaContractAddress), DIAMetaContractPrecision, scraper.luminaClient)
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
	price, err := utils.GetPriceFromDiaAPI(asset.Address, asset.Blockchain)
	if err != nil {
		log.Errorf("Failed to get price of %s from DIA API: %v\n", asset.Symbol, err)
		log.Errorf("asset blockchain: %v\n", asset.Blockchain)
		log.Errorf("asset address: %v\n", asset.Address)
		price = 100
	}
	return price
}

func (scraper *CurveSimulator) getPools(ep models.ExchangePair) map[common.Address]string {

	registry, err := curvefi.NewCurvefi(common.HexToAddress(registryAddress), scraper.restClient)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	} else {
		log.Info("Successfully created contract instance")
	}

	poolCount, err := registry.PoolCount(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Total # of stable swap pools: %d\n", poolCount)

	// Retrieve all pools
	pools := make(map[common.Address]string)
	for i := 0; ; i++ {
		pool, err := registry.FindPoolForCoins0(
			&bind.CallOpts{Context: context.Background()},
			common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
			common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
			big.NewInt(int64(i)),
		)
		if err != nil {
			log.Errorf("Error querying index %d: %v", i, err)
			break
		}
		if pool == (common.Address{}) {
			break
		}
		poolType, err := detectAndInitPoolType(scraper.restClient, pool)
		if err != nil {
			log.Warnf("Skipping pool %s: %v", pool.Hex(), err)
			continue
		}

		pools[pool] = poolType
		log.Infof("Pool #%d: %s | Type: %s", i+1, pool.Hex(), poolType)
	}

	log.Infof("Found %d %v/%v liquidity pools:\n", len(pools), ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)

	return pools
}

func detectAndInitPoolType(client *ethclient.Client, poolAddr common.Address) (string, error) {
	underlying, err := curvepool.NewCurvepool(poolAddr, client)
	if err == nil {
		addr, err := underlying.UnderlyingCoins(&bind.CallOpts{}, big.NewInt(0))
		if err == nil && addr != (common.Address{}) {
			return "underlying", nil
		}
	}

	plain, err := curveplain.NewCurveplainCaller(poolAddr, client)
	if err == nil {
		addr, err := plain.Coins(&bind.CallOpts{}, big.NewInt(0))
		if err == nil && addr != (common.Address{}) {
			return "plain", nil
		}
	}

	factory, err := curvefifactory.NewCurvefifactory(poolAddr, client)
	if err == nil {
		addr, err := factory.Coins(&bind.CallOpts{}, big.NewInt(0))
		if err == nil && addr != (common.Address{}) {
			return "factory", nil
		}
	}

	return "", fmt.Errorf("unable to determine pool type for %s", poolAddr.Hex())
}

func (scraper *CurveSimulator) simulateTrades(ep models.ExchangePair, pools map[common.Address]string, tradesChannel chan models.SimulatedTrade) {
	for poolAddr, poolType := range pools {
		var (
			pool interface{}
			err  error
		)

		switch poolType {
		case "plain":
			pool, err = curveplain.NewCurveplainCaller(poolAddr, scraper.restClient)
		case "underlying":
			pool, err = curvepool.NewCurvepool(poolAddr, scraper.restClient)
		case "factory":
			pool, err = curvefifactory.NewCurvefifactory(poolAddr, scraper.restClient)
		default:
			log.Warnf("Unsupported pool type %s for pool %s", poolType, poolAddr.Hex())
			continue
		}
		if err != nil {
			log.Infof("Failed to create pool contract for %s: %v", poolAddr.Hex(), err)
			continue
		}

		var quoteTokenIndex, baseTokenIndex int
		var ok bool

		// Get token indices and decimals
		switch p := pool.(type) {
		case *curveplain.CurveplainCaller:
			quoteTokenIndex, baseTokenIndex, ok = scraper.validatePoolTokens(ep, p)
		case *curvepool.Curvepool:
			quoteTokenIndex, baseTokenIndex, ok = scraper.validatePoolTokens(ep, p)
		case *curvefifactory.Curvefifactory:
			quoteTokenIndex, baseTokenIndex, ok = scraper.validatePoolTokens(ep, p)
		default:
			log.Warnf("Unknown pool contract type for pool %s", poolAddr.Hex())
			continue
		}

		if !ok {
			continue
		}
		log.Infof("Token validated! token0_index: %v, token1_index: %v\n", quoteTokenIndex, baseTokenIndex)
		// Prepare trade input (e.g., $100)
		baseTokenPrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price
		amountInBase := int64(amountIn_USD_constant / baseTokenPrice) // 100

		decimals := big.NewInt(int64(ep.UnderlyingPair.BaseToken.Decimals)) // e.g. 18
		exponent := new(big.Int).Exp(big.NewInt(10), decimals, nil)         // e.g. 10^18
		amountIn := new(big.Int).Mul(big.NewInt(amountInBase), exponent)    // e.g. 10^20
		amountInAfterDecimalAdjust := new(big.Int).Quo(amountIn, exponent)  // e.g. 10^2
		amountInAfterDecimalAdjustF64, _ := amountInAfterDecimalAdjust.Float64()
		log.Infof("amountIn after adjusting for decimals: %v\n", amountInAfterDecimalAdjust)

		// Run trade simulation
		var amountOutFloat *big.Float
		amountOut, err := scraper.simulator.Execute(pool, baseTokenIndex, quoteTokenIndex, amountIn)
		if err != nil {
			continue
		} else {
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

			liquidity_ok := scraper.hasSufficientLiquidity(ep, pool, baseTokenIndex, quoteTokenIndex)

			if liquidity_ok {
				t := models.SimulatedTrade{
					Price:       amountInAfterDecimalAdjustF64 / amountOutAfterDecimalAdjustF64,
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
		}
	}
}

func (scraper *CurveSimulator) validatePoolTokens(ep models.ExchangePair, pool interface{}) (i, j int, valid bool) {
	// Iterate over all possible token indices (Curve supports up to 8 tokens)
	var tokens []common.Address
	maxCoins := 8
	for idx := 0; idx < maxCoins; idx++ {
		var (
			addr common.Address
			err  error
		)

		switch p := pool.(type) {
		case *curveplain.CurveplainCaller:
			addr, err = p.Coins(&bind.CallOpts{}, big.NewInt(int64(idx)))
		case *curvepool.Curvepool:
			addr, err = p.UnderlyingCoins(&bind.CallOpts{}, big.NewInt(int64(idx)))
		case *curvefifactory.Curvefifactory:
			addr, err = p.Coins(&bind.CallOpts{}, big.NewInt(int64(idx)))
		default:
			log.Warnf("Unknown pool type or contract instance for token index %d", idx)
			return 0, 0, false
		}

		if err != nil || addr == (common.Address{}) {
			continue
		}

		tokens = append(tokens, addr)
	}
	log.Infof("tokens: %v \n", tokens)
	// Check if quote token and base token exist
	quoteTokenIndex := -1
	baseTokenIndex := -1
	for idx, token := range tokens {
		switch token {
		case common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address):
			quoteTokenIndex = idx
		case common.HexToAddress(ep.UnderlyingPair.BaseToken.Address):
			baseTokenIndex = idx
		}
	}

	if quoteTokenIndex == -1 || baseTokenIndex == -1 {
		log.Printf("The pool is missing either token0 (%t) or token1 (%t)", quoteTokenIndex != -1, baseTokenIndex != -1)
		return 0, 0, false
	}

	return quoteTokenIndex, baseTokenIndex, true
}

func (scraper *CurveSimulator) hasSufficientLiquidity(ep models.ExchangePair, pool interface{}, baseIdx, quoteIdx int) bool {
	var baseBalance, quoteBalance *big.Int
	var err error

	switch p := pool.(type) {
	case *curveplain.CurveplainCaller:
		baseBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(baseIdx)))
		if err != nil {
			log.Warnf("Failed to get base balance: %v", err)
			return false
		}
		quoteBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(quoteIdx)))
		if err != nil {
			log.Warnf("Failed to get quote balance: %v", err)
			return false
		}
	case *curvepool.Curvepool:
		baseBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(baseIdx)))
		if err != nil {
			log.Warnf("Failed to get base underlying balance: %v", err)
			return false
		}
		quoteBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(quoteIdx)))
		if err != nil {
			log.Warnf("Failed to get quote underlying balance: %v", err)
			return false
		}
	case *curvefifactory.Curvefifactory:
		baseBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(baseIdx)))
		if err != nil {
			log.Warnf("Failed to get base balance: %v", err)
			return false
		}
		quoteBalance, err = p.Balances(&bind.CallOpts{}, big.NewInt(int64(quoteIdx)))
		if err != nil {
			log.Warnf("Failed to get quote balance: %v", err)
			return false
		}
	default:
		log.Warnf("Unknown pool type in liquidity check")
		return false
	}

	baseDecimals := float64(ep.UnderlyingPair.BaseToken.Decimals)
	quoteDecimals := float64(ep.UnderlyingPair.QuoteToken.Decimals)

	baseBalanceF := new(big.Float).Quo(new(big.Float).SetInt(baseBalance), big.NewFloat(math.Pow10(int(baseDecimals))))
	quoteBalanceF := new(big.Float).Quo(new(big.Float).SetInt(quoteBalance), big.NewFloat(math.Pow10(int(quoteDecimals))))

	if baseBalanceF.Cmp(liquidityThresholdNative) < 0 || quoteBalanceF.Cmp(liquidityThresholdNative) < 0 {
		log.Warnf("Native liquidity not sufficient: base=%s %s, quote=%s %s",
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
		log.Warnf("Base token %s has insufficient USD liquidity: %s", ep.UnderlyingPair.BaseToken.Symbol, baseUSD.Text('f', 2))
		return false
	}
	if quoteUSD.Cmp(liquidityThresholdUSD) < 0 {
		log.Warnf("Quote token %s has insufficient USD liquidity: %s", ep.UnderlyingPair.QuoteToken.Symbol, quoteUSD.Text('f', 2))
		return false
	}

	return true
}
