package simulators

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/diadata-org/lumina-library/contracts/uniswap/univ3factory"
	"github.com/diadata-org/lumina-library/contracts/uniswap/univ3pool"
	"github.com/diadata-org/lumina-library/models"
	simulation "github.com/diadata-org/lumina-library/simulations/simulators/uniswap"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SimulationScraperVersion2 struct {
	waitTime      int
	restClient    *ethclient.Client
	luminaClient  *ethclient.Client
	simulator     *simulation.Simulator
	exchangepairs []models.ExchangePair
	// priceMap maps an asset on the current quotation as published by DIA lumina meta contract.
	priceMap map[models.Asset]models.AssetQuotation
	// feesMap maps an exchangepair on the set of admissible fees. More precisely, the set of fees
	// such that a corresponding pool exists and fulfills liquidity requirements.
	feesMap           map[models.ExchangePair][]UniV3PoolFee
	thresholdSlippage float64
	poolCache         *PoolCallCache
}

type SwapEventsVersion2 struct {
	Amount0In  int64 `json:"Amount0In"`
	Amount0Out int64 `json:"Amount0Out"`
	Amount1In  int64 `json:"Amount1In"`
	Amount1Out int64 `json:"Amount1Out"`
}

type SimulationResponseVersion2 struct {
	Blocknumber string       `json:"blocknumber"`
	Events      []SwapEvents `json:"events"`
	Output      float64      `json:"output"`
	TokenIn     string       `json:"tokenInStr"`
	TokenOut    string       `json:"tokenOutStr"`
}

type UniV3PoolFeeVersion2 struct {
	fee      *big.Int
	address  common.Address
	amountIn float64
}

type PoolData struct {
	Caller *univ3pool.Univ3poolCaller
	Token0 common.Address
	Token1 common.Address
}

type PoolCallCache struct {
	mu    sync.RWMutex
	cache map[common.Address]*PoolData
}

const (
	// Maximal spacing in UniV3 pools with fees of 1%.
	max_spacingVersion2 = int32(200)
)

var (
	restDialVersion2       = ""
	restDialLuminaVersion2 = ""
	// Amount in USD that is used to simulate trades.
	amountIn_USDVersion2 = float64(100)
	// fees are ints with precision 6.
	allFeesVersion2 = []*big.Int{big.NewInt(100), big.NewInt(500), big.NewInt(3000), big.NewInt(10000)}

	// TO DO: Put the following variables to environment variables.
	DIA_Meta_Contract_AddressVersion2   = "0x0087342f5f4c7AB23a37c045c3EF710749527c88"
	DIA_Meta_Contract_PrecisionVersion2 = 8
	priceMap_Update_SecondsVersion2     = 30 * 60
	simulation_Update_SecondsVersion2   = 30
	liquidity_Threshold_USDVersion2     = float64(50000)
	liquidity_Threshold_NativeVersion2  = float64(2)
	threshold_Price_DeviationVersion2   = float64(0.05)
	// Minimal count of active ticks within tick range.
	admissible_CountVersion2 = 10
	// tick range around the current tick that is taken into account for tick check.
	// minWordPosition = int16(-3466)
	// maxWordPosition = int16(3465)
)

func init() {
	var err error
	// Import and cast environment variables.
	DIA_Meta_Contract_AddressVersion2 = utils.Getenv("DIA_META_CONTRACT_ADDRESS", DIA_Meta_Contract_AddressVersion2)
	DIA_Meta_Contract_PrecisionVersion2, err = strconv.Atoi(utils.Getenv("DIA_META_CONTRACT_PRECISION", strconv.Itoa(DIA_Meta_Contract_PrecisionVersion2)))
	if err != nil {
		log.Errorf("DIA_META_CONTRACT_PRECISION: %v", err)
	}
	priceMap_Update_SecondsVersion2, err = strconv.Atoi(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_PRICE_MAP_UPDATE_SECONDS", strconv.Itoa(priceMap_Update_SecondsVersion2)))
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_PRICE_MAP_UPDATE_SECONDS: %v", err)
	}
	simulation_Update_SecondsVersion2, err = strconv.Atoi(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_SIMULATION_UPDATE_SECONDS", strconv.Itoa(simulation_Update_SecondsVersion2)))
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_SIMULATION_UPDATE_SECONDS: %v", err)
	}
	liquidity_Threshold_USDVersion2, err = strconv.ParseFloat(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_LIQUIDITY_THRESHOLD_USD", strconv.Itoa(int(liquidity_Threshold_USDVersion2))), 64)
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_LIQUIDITY_THRESHOLD_USD: %v", err)
	}
	liquidity_Threshold_NativeVersion2, err = strconv.ParseFloat(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_LIQUIDITY_THRESHOLD_NATIVE", strconv.Itoa(int(liquidity_Threshold_NativeVersion2))), 64)
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_LIQUIDITY_THRESHOLD_NATIVE: %v", err)
	}
	threshold_Price_DeviationVersion2, err = strconv.ParseFloat(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_THRESHOLD_PRICE_DEVIATION", strconv.Itoa(int(threshold_Price_DeviationVersion2))), 64)
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_THRESHOLD_PRICE_DEVIATION: %v", err)
	}
	admissible_CountVersion2, err = strconv.Atoi(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_ADMISSIBLE_COUNT", strconv.Itoa(admissible_CountVersion2)))
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_ADMISSIBLE_COUNT: %v", err)
	}
}

func NewUniswapSimulatorVersion2(exchangepairs []models.ExchangePair, tradesChannel chan models.SimulatedTrade) {
	var (
		err     error
		scraper SimulationScraperVersion2
	)

	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_URI_REST", restDialVersion2))
	if err != nil {
		log.Error("init rest client: ", err)
	}
	scraper.luminaClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAP_SIMULATION_TEST)+"_LUMINA_URI_REST", restDialLuminaVersion2))
	if err != nil {
		log.Error("init lumina client: ", err)
	}

	scraper.simulator = simulation.New(scraper.restClient, log)
	scraper.exchangepairs = exchangepairs
	err = scraper.initAssetsAndMapsVersion2()
	if err != nil {
		log.Fatal("initAssetsAndMaps: ", err)
	}

	scraper.thresholdSlippage, err = strconv.ParseFloat(utils.Getenv("UNISWAPV3_THRESHOLD_SLIPPAGE", "5e-17"), 64)
	if err != nil {
		log.Error("Parse THRESHOLD_SLIPPAGE: ", err)
		scraper.thresholdSlippage = 5e-17
	}

	scraper.poolCache = NewPoolCallCache()

	var lock sync.RWMutex
	scraper.updatePriceMapVersion2(&lock)
	// map exchangepairs to list of fees corresponding to deployed pools such that all mapped pools are admissible.
	scraper.updateFeesMapVersion2(&lock)

	for key, poolFee := range scraper.feesMap {
		log.Infof("admissible pair %s -- %s : ", key.UnderlyingPair.QuoteToken.Symbol, key.UnderlyingPair.BaseToken.Symbol)
		for _, v := range poolFee {
			log.Infof("fee -- address: %s -- %s", v.fee.String(), v.address.Hex())
		}
		log.Info("-------------------------------------------------------")
	}

	priceTicker := time.NewTicker(time.Duration(priceMap_Update_SecondsVersion2) * time.Second)
	go func() {
		for range priceTicker.C {
			scraper.updatePriceMapVersion2(&lock)
			scraper.updateFeesMapVersion2(&lock)
		}
	}()

	for pair := range scraper.feesMap {
		log.Infof("Start Simulation scraper for pair: %s-%s", pair.UnderlyingPair.QuoteToken.Symbol, pair.UnderlyingPair.BaseToken.Symbol)
	}
	ticker := time.NewTicker(time.Duration(simulation_Update_SecondsVersion2) * time.Second)
	for range ticker.C {
		log.Info("Simulate trades.")
		scraper.simulateTradesVersion2(tradesChannel)
	}

}

func NewPoolCallCache() *PoolCallCache {
	return &PoolCallCache{
		cache: make(map[common.Address]*PoolData),
	}
}

func (p *PoolCallCache) Get(poolAddr common.Address, client *ethclient.Client) (*PoolData, error) {
	p.mu.RLock()
	if data, ok := p.cache[poolAddr]; ok {
		p.mu.RUnlock()
		return data, nil
	}
	p.mu.RUnlock()

	caller, err := univ3pool.NewUniv3poolCaller(poolAddr, client)
	if err != nil {
		return nil, err
	}
	token0, err := caller.Token0(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	token1, err := caller.Token1(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	data := &PoolData{
		Caller: caller,
		Token0: token0,
		Token1: token1,
	}
	p.mu.Lock()
	p.cache[poolAddr] = data
	p.mu.Unlock()

	return data, nil
}

func (scraper *SimulationScraperVersion2) simulateTradesVersion2(tradesChannel chan models.SimulatedTrade) {
	var wg sync.WaitGroup

	for exchangePair, fees := range scraper.feesMap {
		time.Sleep(time.Duration(scraper.waitTime) * time.Millisecond)
		for _, poolFee := range fees {
			wg.Add(1)

			go func(ep models.ExchangePair, fee *big.Int, w *sync.WaitGroup) {
				defer w.Done()
				amountIn := strconv.FormatFloat(poolFee.amountIn, 'f', -1, 64)
				log.Infof("amountIn: %s", amountIn)

				amountOutString, err := scraper.simulator.Execute(
					ep.UnderlyingPair.QuoteToken,
					ep.UnderlyingPair.BaseToken,
					amountIn,
					fee,
				)
				if err != nil {
					log.Errorf("error getting price of %s - %s ", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
					return
				}

				amountOut, _ := strconv.ParseFloat(amountOutString, 64)
				log.Infof("amountOut: %v", amountOut)

				var address common.Address
				poolFees := scraper.feesMap[ep]
				for _, pf := range poolFees {
					if pf.fee == fee {
						address = pf.address
					}
				}

				poolData, err := scraper.poolCache.Get(address, scraper.restClient)
				if err != nil {
					log.Errorf("Failed to load pool data: %v", err)
					return
				}
				caller := poolData.Caller
				PoolToken0 := poolData.Token0
				PoolToken1 := poolData.Token1

				// Retrieve post-trade price and liquidity
				slot0, err := caller.Slot0(&bind.CallOpts{})
				if err != nil {
					log.Errorf("Failed to get slot0: %v", err)
					return
				}
				liquidityBig, err := caller.Liquidity(&bind.CallOpts{})
				if err != nil {
					log.Errorf("Failed to get liquidity: %v", err)
					return
				}

				log.Infof("Token0 in Pool: %v , Token1 in Pool: %v\n", PoolToken0, PoolToken1)

				// Determine amount0/amount1 signs based on trade direction
				// Assuming simulated trade is Token0 -> Token1
				// When a user sells Token0, the pool receives Token0 → amount0 should be positive.
				// amount0In := big.NewInt(int64(poolFee.amountIn * math.Pow10(int(ep.UnderlyingPair.BaseToken.Decimals))))
				// amount0 := big.NewFloat(poolFee.amountIn) // Pool receives positive amount0

				amount0 := new(big.Float).Neg(big.NewFloat(poolFee.amountIn))
				amount1 := big.NewFloat(amountOut)
				if common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address) == PoolToken0 {
					amount0 = big.NewFloat(poolFee.amountIn)
					amount1 = new(big.Float).Neg(big.NewFloat(amountOut))
				}

				// amount1Out := big.NewInt(int64(amountOut * math.Pow10(int(ep.UnderlyingPair.QuoteToken.Decimals))))
				// amount1 := new(big.Int).Neg(amount1Out) //The pool pays out Token1 → amount1 should be negative.
				// amount1 := new(big.Float).Neg(big.NewFloat(amountOut))
				// amount1 := big.NewFloat(amountOut)

				log.Infof("amount0: %v, amount1: %v\n", amount0, amount1)

				// th_dy, err := CalculateDy(amount0, slot0.SqrtPriceX96, liquidityBig,
				// 	ep.UnderlyingPair.QuoteToken.Decimals, ep.UnderlyingPair.BaseToken.Decimals)
				// if err != nil {
				// 	log.Errorf("Failed to get theoretical dy: %v", err)
				// 	return
				// }
				// th_slippage, err := CalculateSlippage(th_dy, amount1)
				// if err != nil {
				// 	log.Errorf("Failed to get Slippage: %v", err)
				// 	return
				// }

				// log.Warnf("Theoretical dy: %v, Actual dy: %v, Slippage: %v", th_dy, amount1, th_slippage)

				slippage := computeSlippageVersion2(
					slot0.SqrtPriceX96,
					amount0,
					amount1,
					liquidityBig,
				)

				log.Infof("slippage in pool %v: %v", address, slippage)

				if slippage > scraper.thresholdSlippage {
					log.Warnf("slippage above threshold: %v > %v", slippage, scraper.thresholdSlippage)
				} else {
					t := models.SimulatedTrade{
						Price:       poolFee.amountIn / amountOut,
						Volume:      amountOut,
						QuoteToken:  ep.UnderlyingPair.QuoteToken,
						BaseToken:   ep.UnderlyingPair.BaseToken,
						PoolAddress: address.Hex(),
						Time:        time.Now(),
						Exchange:    Exchanges[UNISWAP_SIMULATION_TEST],
					}

					log.Infof("Got trade in pool %v%%: %s-%s -- %v -- %v", float64(fee.Int64())/float64(10000), t.QuoteToken.Symbol, t.BaseToken.Symbol, t.Price, t.Volume)
					tradesChannel <- t
				}

			}(exchangePair, poolFee.fee, &wg)
		}
	}
	wg.Wait()
}

func computeSlippageVersion2(sqrtPriceX96 *big.Int, amount0 *big.Float, amount1 *big.Float, liquidity *big.Int) float64 {
	log.Infof("sqrtPrice -- amount0 -- amount1 -- liquidity: %s -- %s -- %s -- %s", sqrtPriceX96.String(), amount0.String(), amount1.String(), liquidity.String())

	// Convert sqrtPriceX96 to actual price using formula: price = (sqrtPriceX96 / 2^96)^2
	price := new(big.Float).Quo(
		new(big.Float).SetInt(sqrtPriceX96),
		new(big.Float).SetFloat64(math.Pow(2, 96)),
	)

	// Calculate slippage based on trade direction
	if amount0.Sign() < 0 { // Token0 -> Token1
		amount0Abs := new(big.Float).Abs(amount0)
		numerator := new(big.Float).Mul(amount0Abs, price)
		denominator := new(big.Float).SetInt(liquidity)
		slippage, _ := new(big.Float).Quo(numerator, denominator).Float64()

		return slippage
	} else if amount1.Sign() < 0 { // Token1 -> Token0
		amount1Abs := new(big.Float).Abs(amount1)
		numerator := amount1Abs
		denominator := new(big.Float).Mul(new(big.Float).SetInt(liquidity), price)
		slippage, _ := new(big.Float).Quo(numerator, denominator).Float64()
		return slippage
	}
	return 0
}

func getSimulationSwapDataVersion2(events []SwapEvents, tokenInDecimal, tokenOutDecimal uint8) (float64, float64) {
	if len(events) == 0 {
		return 0, 0
	}

	decimalsout := int(tokenOutDecimal)
	decimalsin := int(tokenInDecimal)

	firstEvent := events[0]

	lastEvent := events[len(events)-1]

	var totalInput float64

	if firstEvent.Amount0In != int64(0) {
		totalInput, _ = new(big.Float).Quo(big.NewFloat(0).SetInt(big.NewInt(firstEvent.Amount0In)), new(big.Float).SetFloat64(math.Pow10(decimalsin))).Float64()
	} else {
		totalInput, _ = new(big.Float).Quo(big.NewFloat(0).SetInt(big.NewInt(firstEvent.Amount1In)), new(big.Float).SetFloat64(math.Pow10(decimalsin))).Float64()
	}

	var totalOutput float64
	if lastEvent.Amount1Out != int64(0) {
		totalOutput, _ = new(big.Float).Quo(big.NewFloat(0).SetInt(big.NewInt(lastEvent.Amount1Out)), new(big.Float).SetFloat64(math.Pow10(decimalsout))).Float64()
	} else {
		totalOutput, _ = new(big.Float).Quo(big.NewFloat(0).SetInt(big.NewInt(lastEvent.Amount0Out)), new(big.Float).SetFloat64(math.Pow10(decimalsout))).Float64()

	}

	if totalInput == 0 {
		return 0, 0
	}

	price := float64(totalInput) / float64(totalOutput)

	return price, 1000
}

// initAssets fetches complete asset data from on-chain for all assets in exchangepairs.
// It also initializes the keys of @scraper.priceMap and keys of @scraper.feesMap.
func (scraper *SimulationScraperVersion2) initAssetsAndMapsVersion2() error {

	// memoryMap prevents from iterating over assets twice.
	memoryMap := make(map[string]models.Asset)
	scraper.priceMap = make(map[models.Asset]models.AssetQuotation)
	scraper.feesMap = make(map[models.ExchangePair][]UniV3PoolFee)

	for i, ep := range scraper.exchangepairs {

		if _, ok := memoryMap[ep.UnderlyingPair.QuoteToken.Address]; !ok {

			quoteToken, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address), Exchanges[UNISWAP_SIMULATION_TEST].Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
			scraper.exchangepairs[i].UnderlyingPair.QuoteToken = quoteToken
			scraper.priceMap[quoteToken] = models.AssetQuotation{}
			memoryMap[ep.UnderlyingPair.QuoteToken.Address] = quoteToken

		} else {
			scraper.exchangepairs[i].UnderlyingPair.QuoteToken = memoryMap[ep.UnderlyingPair.QuoteToken.Address]
		}

		if _, ok := memoryMap[ep.UnderlyingPair.BaseToken.Address]; !ok {

			baseToken, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.BaseToken.Address), Exchanges[UNISWAP_SIMULATION_TEST].Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
			scraper.exchangepairs[i].UnderlyingPair.BaseToken = baseToken
			scraper.priceMap[baseToken] = models.AssetQuotation{}
			memoryMap[ep.UnderlyingPair.BaseToken.Address] = baseToken

		} else {
			scraper.exchangepairs[i].UnderlyingPair.BaseToken = memoryMap[ep.UnderlyingPair.BaseToken.Address]
		}
		// Initialize @scraper.feesMap.
		scraper.feesMap[scraper.exchangepairs[i]] = []UniV3PoolFee{}

	}
	return nil
}

// updatePriceMap fetches the current price of each of the involved assets from DIA lumina meta contract.
func (scraper *SimulationScraperVersion2) updatePriceMapVersion2(lock *sync.RWMutex) {
	for asset := range scraper.priceMap {
		quotation, err := asset.GetOnchainPrice(common.HexToAddress(DIA_Meta_Contract_AddressVersion2), DIA_Meta_Contract_PrecisionVersion2, scraper.luminaClient)
		if err != nil {
			log.Errorf("GetOnchainPrice for %s -- %s: %v", asset.Symbol, asset.Address, err)
			continue
		} else {
			log.Infof("USD price for (base-)token %s: %v", asset.Symbol, quotation.Price)
		}
		lock.Lock()
		scraper.priceMap[asset] = quotation
		lock.Unlock()
	}
}

// func CalculateSlippage(dyTheoretical, dyActual *big.Float) (float64, error) {
// 	if dyTheoretical.Sign() <= 0 {
// 		dyTheoretical = new(big.Float).Abs(dyTheoretical)
// 	}

// 	// Calculate the ratio: (actual - theoretical) / theoretical
// 	numerator := new(big.Float).Abs(new(big.Float).Sub(dyTheoretical, dyActual))
// 	ratio := new(big.Float).Quo(numerator, dyTheoretical)
// 	ratio64, _ := ratio.Float64()
// 	return ratio64, nil
// }

// Calculate theoretical dy (amount of output token1) given input token0 amount dx via Uniswap V3 formula
// func CalculateDy(
// 	dx *big.Float,
// 	sqrtPriceX96 *big.Int,
// 	liquidity *big.Int,
// 	token0Decimals, token1Decimals uint8,
// ) (*big.Float, error) {
// 	dxAdjusted := adjustDecimals(dx, int(token0Decimals), 18) // convert to 18-decimal precision

// 	log.Printf("dxAdjusted: %v", dxAdjusted)

// 	// 1. Price Conversion: √P = sqrtPriceX96 / 2^96
// 	sqrtPrice := new(big.Float).Quo(
// 		new(big.Float).SetInt(sqrtPriceX96),
// 		new(big.Float).SetFloat64(math.Pow(2, 96)),
// 	)

// 	L := new(big.Float).SetInt(liquidity)

// 	// 2. New Price: √P' = (L * √P) / (L + Δx * √P)
// 	numerator := L
// 	denominator := new(big.Float).Add(L, new(big.Float).Mul(dxAdjusted, sqrtPrice))
// 	sqrtPNew := new(big.Float).Quo(numerator, denominator)

// 	// 3. Output Amount: Δy = L * (√P - √P')
// 	deltaY := new(big.Float).Mul(L, new(big.Float).Abs(new(big.Float).Sub(sqrtPrice, sqrtPNew)))

// 	dyAdjusted := adjustDecimals(deltaY, 18, int(token1Decimals))

// 	return new(big.Float).Mul(dyAdjusted, big.NewFloat(1e-18)), nil
// }

// func adjustDecimals(amount *big.Float, fromDecimals, toDecimals int) *big.Float {
// 	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(fromDecimals-toDecimals)), nil)
// 	return new(big.Float).Quo(amount, new(big.Float).SetInt(factor))
// }

// updateFeesMap updates values in scraper.feesMap.
func (scraper *SimulationScraperVersion2) updateFeesMapVersion2(lock *sync.RWMutex) {

	whitelistedPools, err := models.GetWhitelistedPoolsFromConfig(UNISWAP_SIMULATION_TEST)
	if err != nil {
		log.Error("GetWhitelistedPoolsFromConfig: ", err)
	} else {
		log.Info("whitelisted pool addresses: ", whitelistedPools)
	}

	// Remark: In case initial load is too slow, this loop can be parallelized. Not sure if it works with ETH requests though.
	for _, ep := range scraper.exchangepairs {
		quoteToken := ep.UnderlyingPair.QuoteToken
		baseToken := ep.UnderlyingPair.BaseToken

		var indexFeeMap = make(map[int]*big.Int)
		// var prices0, prices1 []float64
		var poolMap = make(map[int]common.Address)

		for i, fee := range allFeesVersion2 {

			poolAddress, err := scraper.getPoolVersion2(ep, fee)
			if err != nil {
				log.Errorf("getPool with address %s and fee %v: %v", poolAddress.Hex(), float64(fee.Int64())/float64(10000), err)
				continue
			} else if (poolAddress == common.Address{}) {
				log.Warnf("pool for %s-%s with fees %v does not exist.", quoteToken.Symbol, baseToken.Symbol, float64(fee.Int64())/float64(10000))
				continue
			} else {
				log.Infof("Start checking admissibility for pool %s-%s with fee %v%% and address %s ", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, float64(fee.Int64())/float64(10000), poolAddress.Hex())
			}

			// Check if pool is admissible, i.e.
			// 1. check balance(s).
			// 2. check distribution of active ticks.
			// 3. check prices in current tick across pools/fees.

			balanceOk := scraper.checkBalancesVersion2(quoteToken, baseToken, poolAddress)
			if utils.ContainsAddress(whitelistedPools, poolAddress) {
				balanceOk = true
			}
			if !balanceOk {
				log.Warnf("Balances not ok for pool with fee %v%%", float64(fee.Int64())/float64(10000))
				// Remove from scraper.feesMap[ep] if existent.
				if containsAddress(poolAddress, scraper.feesMap[ep]) {
					log.Warn("low balance - remove pool from set of admissible pools: ", poolAddress)
					cleanedFees := removeFeeByAddress(poolAddress, scraper.feesMap[ep])
					scraper.feesMap[ep] = cleanedFees
				}
				continue
			}

			// ticksOk, currentTick := scraper.checkTicksVersion2(poolAddress, word_Range, considered_tick_range, admissible_Count)
			// log.Infof("currentTick:", currentTick)
			// if utils.containsAddressVersion2(whitelistedPools, poolAddress) {
			// 	ticksOk = true
			// }
			// if !ticksOk {
			// 	log.Warnf("ticks not ok for %s with fee %s", poolAddress.Hex(), fee.String())
			// 	// Remove from scraper.feesMap[ep] if existent.
			// 	if containsAddressVersion2(poolAddress, scraper.feesMap[ep]) {
			// 		log.Warn("poor tick distribution - remove pool from set of admissible pools: ", poolAddress)
			// 		cleanedFees := removeFeeByAddressVersion2(poolAddress, scraper.feesMap[ep])
			// 		scraper.feesMap[ep] = cleanedFees
			// 	}
			// 	continue
			// }

			indexFeeMap[i] = fee
			poolMap[i] = poolAddress

			// logging ------------
			// priceMin, _ := computeTickPrices(currentTick-considered_tick_range, int8(quoteToken.Decimals), int8(baseToken.Decimals))
			// priceMax, _ := computeTickPrices(currentTick+considered_tick_range, int8(quoteToken.Decimals), int8(baseToken.Decimals))
			// log.Infof("corresponding price range in tick range: %s -- %s", priceMin.String(), priceMax.String())
			// log.Infof("pool admitted %s for balances and ticks.", poolAddress.Hex())
			// // --------------------

			// p0, p1, err := scraper.getActivePricesVersion2(poolAddress, int8(quoteToken.Decimals), int8(baseToken.Decimals))
			// if err != nil {
			// 	log.Errorf("getActivePrices on %s: %v", poolAddress.Hex(), err)
			// }
			// prices0 = append(prices0, p0)
			// prices1 = append(prices1, p1)
			// log.Debugf("prices in current tick price0 -- price1: %v -- %v", prices0, prices1)
		}

		// Outlier detection on prices.
		// TO DO: Should we also check for prices1?
		// scraper.checkPricesVersion2(prices0, ep, indexFeeMap, poolMap)

		for index, fee := range indexFeeMap {
			if !containsFeeVersion2(fee, scraper.feesMap[ep]) {
				scraper.feesMap[ep] = append(scraper.feesMap[ep], UniV3PoolFee{fee: fee, address: poolMap[index]})
			}
		}

		// Compute amountIn in such that it corresponds to @amountIn_USD amount in USD.
		baseTokenPrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price

		if baseTokenPrice == 0 {
			log.Warnf("Could not determine price of base token on chain %s. Checking DIA API.", ep.UnderlyingPair.BaseToken.Symbol)
			baseTokenPrice, err = utils.GetPriceFromDiaAPI(ep.UnderlyingPair.BaseToken.Address, ep.UnderlyingPair.BaseToken.Blockchain)
			if err != nil {
				log.Errorf("Failed to get baseTokenPrice from DIA API: %v\n", err)
				log.Errorf("baseToken Blockchain: %v\n", ep.UnderlyingPair.BaseToken.Blockchain)
				log.Errorf("baseToken Address: %v\n", ep.UnderlyingPair.BaseToken.Address)
				log.Warnf("Could not determine price of base token from DIA API%s. Continue with native volume of 1.", ep.UnderlyingPair.BaseToken.Symbol)
				baseTokenPrice = 100
			}
		}
		for i := range scraper.feesMap[ep] {
			lock.Lock()
			// TODO: need to switch amountIn_USD to Big Float
			scraper.feesMap[ep][i].amountIn = amountIn_USDVersion2 / baseTokenPrice
			lock.Unlock()
		}

		// Remove from map in case no admissible pool/fee was found.
		if len(scraper.feesMap[ep]) == 0 {
			delete(scraper.feesMap, ep)
		}

	}

}

// --------------------------------------------------------------------------------------------------------------

func (scraper *SimulationScraperVersion2) checkBalancesVersion2(quoteToken models.Asset, baseToken models.Asset, poolAddress common.Address) bool {
	balance0, err := quoteToken.GetBalance(poolAddress, scraper.restClient)
	if err != nil {
		log.Errorf("GetBalance of %s: %v", quoteToken.Address, err)
	}
	balance1, err := baseToken.GetBalance(poolAddress, scraper.restClient)
	if err != nil {
		log.Errorf("GetBalance of %s: %v", baseToken.Address, err)
	}
	if balance0 < liquidity_Threshold_NativeVersion2 {
		log.Warnf("native liquidity for %s in %s-%s with address %s not sufficient: %v", quoteToken.Symbol, quoteToken.Symbol, baseToken.Symbol, poolAddress.Hex(), balance0)
		return false
	}
	if balance1 < liquidity_Threshold_NativeVersion2 {
		log.Warnf("native liquidity for %s in %s-%s  with address %s not sufficient: %v", baseToken.Symbol, quoteToken.Symbol, baseToken.Symbol, poolAddress.Hex(), balance1)
		return false
	}

	balance0USD := balance0 * scraper.priceMap[quoteToken].Price
	balance1USD := balance1 * scraper.priceMap[baseToken].Price
	if 0 < balance0USD && balance0USD < liquidity_Threshold_USDVersion2 {
		log.Warnf("USD liquidity for %s in %s-%s  with address %s not sufficient: %v", quoteToken.Symbol, quoteToken.Symbol, baseToken.Symbol, poolAddress.Hex(), balance0USD)
		return false
	}
	if 0 < balance1USD && balance1USD < liquidity_Threshold_USDVersion2 {
		log.Warnf("USD liquidity for %s in %s-%s  with address %s not sufficient: %v", baseToken.Symbol, quoteToken.Symbol, baseToken.Symbol, poolAddress.Hex(), balance1USD)
		return false
	}
	return true
}

// checkTicks counts active ticks within given @wordRange.
// @admissibleSteps is the number of steps within which range ticks are taken into account.
// @admissibleCount is the number of active ticks which have to lie within the given range.
func (scraper *SimulationScraperVersion2) checkTicksVersion2(poolAddress common.Address, wordRange int32, admissibleRange int32, admissibleCount int) (ok bool, currentTick int32) {
	ticks, err := scraper.getCurrentTicksVersion2(poolAddress, wordRange)
	if err != nil {
		log.Error("getCurrentTicks: ", err)
	} else {
		log.Debug("ticks: ", ticks)
	}

	var caller *univ3pool.Univ3poolCaller
	caller, err = univ3pool.NewUniv3poolCaller(poolAddress, scraper.restClient)
	if err != nil {
		return
	}
	slot0, err := caller.Slot0(&bind.CallOpts{})
	if err != nil {
		return
	}

	currentTick = int32(slot0.Tick.Int64())
	minTick := currentTick - admissibleRange
	maxTick := currentTick + admissibleRange
	ok = scraper.checkTickCountVersion2(ticks, minTick, maxTick, admissibleCount)
	return
}

// func (scraper *SimulationScraper) checkTicks(poolAddress common.Address, wordRange int32, admissibleRange int32, admissibleCount int) (ok bool, currentTick int32) {
// 	// Get the current valid ticks range
// 	ticks, err := scraper.getCurrentTicks(poolAddress, wordRange)
// 	if err != nil {
// 		log.Error("getCurrentTicks: ", err)
// 		return false, 0
// 	}
// 	log.Debug("active ticks: ", ticks)

// 	// Initialize the pool calling interface
// 	caller, err := univ3pool.NewUniv3poolCaller(poolAddress, scraper.restClient)
// 	if err != nil {
// 		log.Error("Pool caller initialization failed: ", err)
// 		return false, 0
// 	}

// 	// Get the current slot0 status (including current tick and price)
// 	slot0, err := caller.Slot0(&bind.CallOpts{})
// 	if err != nil {
// 		log.Error("Failed to obtain slot0 status: ", err)
// 		return false, 0
// 	}
// 	currentTick = int32(slot0.Tick.Int64())

// 	// Get pool parameters
// 	tickSpacing, err := caller.TickSpacing(&bind.CallOpts{}) // Dynamically obtain tick spacing
// 	if err != nil {
// 		log.Error("Failed to get tickSpacing: ", err)
// 		return false, 0
// 	}
// 	minSafeDistance := int32(tickSpacing.Int64()) * 2 // The safety distance is 2 times the tick spacing

// 	// Get the actual liquidity boundary of the pool
// 	poolMinTick := uniswap_utils.MinTick

// 	poolMaxTick := uniswap_utils.MaxTick

// 	// Boundary safety check (to prevent liquidity cliff)
// 	if (currentTick-int32(poolMinTick)) < minSafeDistance ||
// 		(int32(poolMaxTick)-currentTick) < minSafeDistance {
// 		log.Warnf("Liquidity boundary warning | Pool address: %s | Current tick: %d | Pool boundary: [%d,%d] | Safety distance: %d",
// 			poolAddress.Hex(), currentTick,
// 			poolMinTick, poolMaxTick,
// 			minSafeDistance)
// 		return false, currentTick
// 	}

// 	// Checking the number of active ticks (detecting liquidity fragmentation)
// 	activeTickMin := currentTick - admissibleRange
// 	activeTickMax := currentTick + admissibleRange
// 	ok = scraper.checkTickCount(ticks, activeTickMin, activeTickMax, admissibleCount)

// 	log.Infof("Pool detection results | Address: %s | Current tick: %d | Number of active ticks: %d/%d | Safety distance: %d",
// 		poolAddress.Hex(), currentTick,
// 		len(ticks), admissibleCount, minSafeDistance)
// 	return ok, currentTick
// }

func (scraper *SimulationScraperVersion2) checkPricesVersion2(prices []float64, ep models.ExchangePair, indexMap map[int]*big.Int, poolMap map[int]common.Address) {
	log.Infof("checking price outliers for %v pools %s-%s ", len(prices), ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
	// TO DO: Only append if not existent yet.
	if len(prices) == 1 {
		for index, fee := range indexMap {
			if !containsFeeVersion2(fee, scraper.feesMap[ep]) {
				scraper.feesMap[ep] = append(scraper.feesMap[ep], UniV3PoolFee{fee: fee, address: poolMap[index]})
			}
		}
	}
	if len(prices) == 2 {
		// If deviation too large store both and print warning.
		dist := utils.AvgDistances(prices)
		if dist[0] > threshold_Price_DeviationVersion2 {
			log.Warnf("prices in pools %s and %s differ by more than %v: %v -- %v ",
				"poolAddress",
				"poolAddress",
				threshold_Price_DeviationVersion2,
				prices[0],
				prices[1],
			)
		} else {
			for index, fee := range indexMap {
				if !containsFeeVersion2(fee, scraper.feesMap[ep]) {
					scraper.feesMap[ep] = append(scraper.feesMap[ep], UniV3PoolFee{fee: fee, address: poolMap[index]})
				}
			}
		}
	}
	if len(prices) > 2 {
		_, indices := utils.RemoveOutliers(prices, threshold_Price_DeviationVersion2)
		for _, ind := range indices {
			if !containsFeeVersion2(indexMap[ind], scraper.feesMap[ep]) {
				scraper.feesMap[ep] = append(scraper.feesMap[ep], UniV3PoolFee{fee: indexMap[ind], address: poolMap[ind]})
			}
		}
	}
}

// --------------------------------------------------------------------------------------------------------------

// checkTicks returns true, if @tickCount is greater or equal to @tickThreshold. tickCount is measured
// by counting, how many active ticks lie within @minTick and @maxTick.
// steps measured with respect to @tickSpacing.
func (scraper *SimulationScraperVersion2) checkTickCountVersion2(ticks []int32, minTick int32, maxTick int32, admissibleCount int) (ok bool) {
	var tickCount int
	for _, tick := range ticks {
		if minTick <= tick && tick <= maxTick {
			tickCount++
		}
	}

	log.Infof("%v active ticks out of total %v", tickCount, len(ticks))
	if tickCount > admissibleCount {
		ok = true
	}
	return
}

// getActivePrices returns prices of pool tokens in the currently active tick. Prices denomination is of course native.
func (scraper *SimulationScraperVersion2) getActivePricesVersion2(poolAddress common.Address, decimals0 int8, decimals1 int8) (price0 float64, price1 float64, err error) {
	var caller *univ3pool.Univ3poolCaller
	caller, err = univ3pool.NewUniv3poolCaller(poolAddress, scraper.restClient)
	if err != nil {
		return
	}
	slot0, err := caller.Slot0(&bind.CallOpts{})
	if err != nil {
		return
	}
	currentTick := slot0.Tick

	price0Big, price1Big := computeTickPricesVersion2(int32(currentTick.Int64()), decimals0, decimals1)
	price0, _ = price0Big.Float64()
	price1, _ = price1Big.Float64()
	return
}

// getCurrentTicks returns active ticks within a range of @wordRange around slot0.
func (scraper *SimulationScraperVersion2) getCurrentTicksVersion2(poolAddress common.Address, wordRange int32) (ticks []int32, err error) {

	var caller *univ3pool.Univ3poolCaller
	caller, err = univ3pool.NewUniv3poolCaller(poolAddress, scraper.restClient)
	if err != nil {
		return
	}
	slot0, err := caller.Slot0(&bind.CallOpts{})
	if err != nil {
		return
	}
	tickSpacing, err := caller.TickSpacing(&bind.CallOpts{})
	if err != nil {
		return
	}
	log.Debug("ticks spacing: ", tickSpacing.String())

	currentWordPosition := getWordPosition(int32(slot0.Tick.Int64()), int32(tickSpacing.Int64()))
	ticks, err = scraper.getActiveTicksInRangeVersion2(poolAddress, currentWordPosition-wordRange, currentWordPosition+wordRange)
	if err != nil {
		log.Fatal("getAllActiveTicks: ", err)
	}

	return
}

// getPool returns the unique pool containing both assets from @exchangepair with given @fee if it exists.
func (scraper *SimulationScraperVersion2) getPoolVersion2(exchangepair models.ExchangePair, fee *big.Int) (poolAddress common.Address, err error) {
	var caller *univ3factory.Univ3factoryCaller
	caller, err = univ3factory.NewUniv3factoryCaller(common.HexToAddress(helper.ContractV3Factory), scraper.restClient)
	if err != nil {
		return
	}
	log.Debugf("get pool for addresses %s -- %s", exchangepair.UnderlyingPair.QuoteToken.Address, exchangepair.UnderlyingPair.BaseToken.Address)
	poolAddress, err = caller.GetPool(
		&bind.CallOpts{},
		common.HexToAddress(exchangepair.UnderlyingPair.QuoteToken.Address),
		common.HexToAddress(exchangepair.UnderlyingPair.BaseToken.Address),
		fee,
	)
	return
}

// --------------------------------------------------------------------------------------------------------------

// getWordPosition calculates the correct word position in tickBitmap taking tickSpacing into account.
func getWordPositionVersion2(tickIndex int32, tickSpacing int32) int32 {
	// Normalize tickIndex by tickSpacing before computing word position
	scaledTickIndex := tickIndex / tickSpacing
	if scaledTickIndex < 0 {
		return (scaledTickIndex / 256) - 1 // Handle negative tick indices correctly
	}
	return scaledTickIndex / 256
}

// getActiveTicksInRange retrieves all active ticks using the tickBitmap structure
func (scraper *SimulationScraperVersion2) getActiveTicksInRangeVersion2(poolAddress common.Address, wordPositionLeft int32, wordPositionRight int32) (ticks []int32, err error) {

	caller, err := univ3pool.NewUniv3poolCaller(poolAddress, scraper.restClient)
	if err != nil {
		return
	}
	tickSpacing, err := caller.TickSpacing(&bind.CallOpts{})
	if err != nil {
		return
	}
	log.Debug("tickSpacing: ", tickSpacing)
	log.Debug("poolAddress for following ticks: ", poolAddress.Hex())

	for wordPosition := wordPositionLeft; wordPosition <= wordPositionRight; wordPosition++ {

		var tickBitmap *big.Int
		tickBitmap, err = caller.TickBitmap(&bind.CallOpts{}, int16(wordPosition))
		if err != nil {
			return
		}

		// Convert bitmap to active tick indices
		activeTicks := extractActiveTicksFromBitmap(wordPosition, int32(tickSpacing.Int64()), tickBitmap)
		log.Debugf("activeTicks at wordPosition %v: %v", wordPosition, activeTicks)
		ticks = append(ticks, activeTicks...)

	}
	return
}

// extractActiveTicksFromBitmap extracts active ticks from a Uniswap tickBitmap word
func extractActiveTicksFromBitmapVersion2(wordPosition int32, tickSpacing int32, tickBitmap *big.Int) []int32 {
	var activeTicks []int32

	for i := 0; i < 256; i++ {
		// Check if the bit at position `i` is set in the bitmap
		if tickBitmap.Bit(i) == 1 {
			// Compute the tick index correctly, considering tickSpacing
			scaledTickIndex := (wordPosition * 256) + int32(i)
			tickIndex := scaledTickIndex * tickSpacing // Scale back to actual tick index

			activeTicks = append(activeTicks, tickIndex)
		}
	}

	return activeTicks
}

func containsFeeVersion2(fee *big.Int, fees []UniV3PoolFee) bool {
	for _, f := range fees {
		if fee.Cmp(f.fee) == 0 {
			return true
		}
	}
	return false
}

func containsAddressVersion2(address common.Address, fees []UniV3PoolFee) bool {
	for _, f := range fees {
		if address == f.address {
			return true
		}
	}
	return false
}

func removeFeeByAddressVersion2(address common.Address, fees []UniV3PoolFee) []UniV3PoolFee {
	for i, f := range fees {
		if address == f.address {
			return append(fees[:i], fees[i+1:]...)
		}
	}
	return fees
}

// ------------------------------------------------------------------------------------------------------------------------

// computePrices calculates prices of token0 and token1 with respect to vice versa.
func computeTickPricesVersion2(tick int32, decimals0 int8, decimals1 int8) (*big.Float, *big.Float) {
	base := big.NewFloat(1.0001)
	numerator := powBigFloat(base, big.NewFloat(float64(tick)))
	denominator := powBigFloat(big.NewFloat(10), big.NewFloat(float64(decimals1-decimals0)))
	return big.NewFloat(0).Quo(numerator, denominator), big.NewFloat(0).Quo(denominator, numerator)
}

// powBigFloat calculates x^y using ln(x) * y and exp(result)
func powBigFloatVersion2(x, y *big.Float) *big.Float {
	// Ensure high precision
	precision := uint(100)
	x.SetPrec(precision)
	y.SetPrec(precision)

	// Convert x to float64 for math.Log
	xFloat64, _ := x.Float64()
	if xFloat64 <= 0 {
		panic("x must be greater than 0 for real-valued exponentiation")
	}

	// Compute ln(x) * y
	lnX := math.Log(xFloat64)  // Compute ln(x) as float64
	yFloat64, _ := y.Float64() // Convert y to float64
	result := big.NewFloat(0).SetPrec(precision)
	result.SetFloat64(lnX * yFloat64) // Multiply ln(x) * y

	// Compute exp(result)
	expFloat, _ := result.Float64()
	expResult := big.NewFloat(math.Exp(expFloat))

	return expResult
}
