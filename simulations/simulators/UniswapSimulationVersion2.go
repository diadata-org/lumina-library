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

type PoolData struct {
	Caller *univ3pool.Univ3poolCaller
	Token0 common.Address
}

type PoolCallCache struct {
	mu    sync.RWMutex
	cache map[common.Address]*PoolData
}

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
	admissible_CountVersion2            = 10
	liquidity_Threshold                 = big.NewInt(30000000000000000)
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

	scraper.thresholdSlippage, err = strconv.ParseFloat(utils.Getenv("UNISWAPV3_THRESHOLD_SLIPPAGE", "0.1"), 64)
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

	data := &PoolData{
		Caller: caller,
		Token0: token0,
	}
	p.mu.Lock()
	p.cache[poolAddr] = data
	p.mu.Unlock()

	return data, nil
}

func (s *SimulationScraperVersion2) getAmountIn(ep models.ExchangePair, amountInBase float64) (*big.Int, *big.Float) {
	decimals := big.NewInt(int64(ep.UnderlyingPair.QuoteToken.Decimals)) // e.g. 18
	exponent := new(big.Int).Exp(big.NewInt(10), decimals, nil)          // e.g. 10^18
	exponentFloat := new(big.Float).SetInt(exponent)

	amountIn := new(big.Float).Mul(big.NewFloat(amountInBase), exponentFloat) // e.g. 10^20
	amountInInt := new(big.Int)
	amountIn.Int(amountInInt)

	amountInAfterDecimalAdjust := new(big.Float).Quo(amountIn, exponentFloat) // e.g. 10^2

	return amountInInt, amountInAfterDecimalAdjust
}

func (scraper *SimulationScraperVersion2) simulateTradesVersion2(tradesChannel chan models.SimulatedTrade) {
	var wg sync.WaitGroup

	for exchangePair, fees := range scraper.feesMap {
		time.Sleep(time.Duration(scraper.waitTime) * time.Millisecond)
		for _, poolFee := range fees {
			wg.Add(1)

			go func(ep models.ExchangePair, fee *big.Int, w *sync.WaitGroup) {
				defer w.Done()
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
				poolToken0 := poolData.Token0

				token0 := ep.UnderlyingPair.BaseToken
				token1 := ep.UnderlyingPair.QuoteToken
				power := ep.UnderlyingPair.QuoteToken.Decimals
				if ep.UnderlyingPair.QuoteToken.Address == poolToken0.Hex() {
					token0 = ep.UnderlyingPair.QuoteToken
					token1 = ep.UnderlyingPair.BaseToken
					power = ep.UnderlyingPair.BaseToken.Decimals
				}
				log.Infof("fee: %v, token0: %s, token1: %s, power: %d", fee, token0.Symbol, token1.Symbol, power)

				amountIn, amountInAfterDecimalAdjust := scraper.getAmountIn(ep, poolFee.amountIn)
				log.Infof("amountInAfterDecimalAdjust: %v", amountInAfterDecimalAdjust)

				amountOutInt, err := scraper.simulator.ExecuteVersion2(
					token0,
					token1,
					amountIn,
					fee,
				)
				if err != nil {
					log.Errorf("error getting price of %s - %s ", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
					return
				}

				// Retrieve post-trade price and liquidity
				slot0, err := caller.Slot0(&bind.CallOpts{})
				if err != nil || slot0.SqrtPriceX96.Cmp(big.NewInt(0)) == 0 {
					log.Errorf("Failed to get slot0: %v", err)
					return
				}

				liquidityBig, err := caller.Liquidity(&bind.CallOpts{})
				if err != nil || liquidityBig.Cmp(big.NewInt(0)) == 0 || liquidityBig.Cmp(liquidity_Threshold) < 0 {
					log.Errorf("Failed to get sufficient liquidity: %v", err)
					return
				}
				log.Infof("liquidityBig: %v", liquidityBig)

				amount0 := new(big.Int).Neg(amountIn)
				amount1 := amountOutInt
				if ep.UnderlyingPair.QuoteToken.Address == poolToken0.Hex() {
					amount0 = amountIn
					amount1 = new(big.Int).Neg(amountOutInt)
				}

				log.Infof("amount0: %v, amount1: %v\n", amount0, amount1)

				slippage := computeSlippageV3(
					slot0.SqrtPriceX96,
					amount0,
					amount1,
					liquidityBig,
				)

				log.Infof("slippage in pool %v: %v", address, slippage)

				amountOutFloat := new(big.Float).SetInt(amountOutInt)
				divisor := new(big.Float).SetInt(
					new(big.Int).Exp(
						big.NewInt(10),
						big.NewInt(int64(power)),
						nil,
					),
				)
				amountOutFloat.Quo(amountOutFloat, divisor)

				priceBig := new(big.Float).Quo(amountInAfterDecimalAdjust, amountOutFloat)

				if ep.UnderlyingPair.QuoteToken.Address == poolToken0.Hex() {
					priceBig = new(big.Float).Quo(amountOutFloat, amountInAfterDecimalAdjust)
				}

				price, _ := priceBig.Float64()
				volume, _ := amountOutFloat.Float64()

				if slippage > scraper.thresholdSlippage {
					log.Warnf("slippage above threshold: %v > %v", slippage, scraper.thresholdSlippage)
				} else {
					t := models.SimulatedTrade{
						Price:       price,
						Volume:      volume,
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

func computeSlippageV3(sqrtPriceX96 *big.Int, amount0 *big.Int, amount1 *big.Int, liquidity *big.Int) (slippage float64) {

	log.Infof("sqrtPrice -- amount0 -- amount1 -- liquidity: %s -- %s -- %s -- %s", sqrtPriceX96.String(), amount0.String(), amount1.String(), liquidity.String())

	// Convert sqrtPriceX96 to actual price using formula: sqrt(price) = (sqrtPriceX96 / 2^96)
	price := new(big.Float).Quo(
		new(big.Float).SetInt(sqrtPriceX96),
		new(big.Float).SetFloat64(math.Pow(2, 96)),
	)

	// Calculate slippage based on trade direction
	if amount0.Sign() < 0 { // Token0 -> Token1
		amount0Abs := new(big.Float).Abs(new(big.Float).SetInt(amount0))
		numerator := new(big.Float).Mul(amount0Abs, price)
		denominator := new(big.Float).SetInt(liquidity)
		slippage, _ := new(big.Float).Quo(numerator, denominator).Float64()

		return slippage
	} else if amount1.Sign() < 0 { // Token1 -> Token0
		amount1Abs := new(big.Float).Abs(new(big.Float).SetInt(amount1))
		numerator := amount1Abs
		denominator := new(big.Float).Mul(new(big.Float).SetInt(liquidity), price)
		slippage, _ := new(big.Float).Quo(numerator, denominator).Float64()
		return slippage
	}
	return 0
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
		quotation, err := asset.GetPrice(common.HexToAddress(DIA_Meta_Contract_AddressVersion2), DIA_Meta_Contract_PrecisionVersion2, scraper.luminaClient)
		if err != nil && quotation.Price == 1 {
			log.Errorf("Failed to GetOnchainPrice for %s -- %s: %v ; continue with default price of 1", asset.Symbol, asset.Address, err)
			continue
		} else {
			log.Infof("USD price for (base-)token %s: %v", asset.Symbol, quotation.Price)
		}
		lock.Lock()
		scraper.priceMap[asset] = quotation
		lock.Unlock()
	}
}

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

			// Check if pool is admissible, i.e. check balance(s)
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

			indexFeeMap[i] = fee
			poolMap[i] = poolAddress
		}

		for index, fee := range indexFeeMap {
			if !containsFeeVersion2(fee, scraper.feesMap[ep]) {
				scraper.feesMap[ep] = append(scraper.feesMap[ep], UniV3PoolFee{fee: fee, address: poolMap[index]})
			}
		}

		// Compute amountIn in such that it corresponds to @amountIn_USD amount in USD.
		baseTokenPrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price

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

func containsFeeVersion2(fee *big.Int, fees []UniV3PoolFee) bool {
	for _, f := range fees {
		if fee.Cmp(f.fee) == 0 {
			return true
		}
	}
	return false
}
