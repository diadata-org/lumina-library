package simulators

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"

	v4quoter "github.com/diadata-org/lumina-library/contracts/uniswapv4/V4Quoter"
	poolManager "github.com/diadata-org/lumina-library/contracts/uniswapv4/poolManager"
	poolState "github.com/diadata-org/lumina-library/contracts/uniswapv4/poolState"
	"github.com/diadata-org/lumina-library/models"
	simulation "github.com/diadata-org/lumina-library/simulations/simulators/uniswapV4"
	"github.com/diadata-org/lumina-library/utils"
)

type UniswapV4Simulator struct {
	restClient        *ethclient.Client
	luminaClient      *ethclient.Client
	quoter            *v4quoter.V4Quoter
	poolManager       *poolManager.PoolManager
	poolState         *poolState.PoolState
	exchangepairs     []models.ExchangePair
	priceMap          map[models.Asset]models.AssetQuotation
	slippageThreshold float64
	simulator         *simulation.Simulator
	allPools          map[models.ExchangePair]map[common.Hash]PoolState
}

type PoolState struct {
	FeeStr      string
	TickSpacing *big.Int
	Params      v4quoter.IV4QuoterQuoteExactSingleParams
	Liquidity   *big.Int
}

var (
	restUrl       = ""
	restLuminaUrl = ""
	// Amount in USD that is used to simulate trades.
	amountInUSDConstant = float64(100)
	// TO DO: Put the following variables to environment variables.
	DIA_Meta_Contract_Address_V4   = "0x0087342f5f4c7AB23a37c045c3EF710749527c88"
	DIA_Meta_Contract_Precision_V4 = 8
	priceMap_Update_Seconds_V4     = 30 * 60
	simulation_Update_Seconds_V4   = 30
	feeToTickSpacing               = map[string]*big.Int{
		"10":    big.NewInt(1),
		"100":   big.NewInt(1),
		"500":   big.NewInt(10),
		"3000":  big.NewInt(60),
		"10000": big.NewInt(200),
	}
	liquidityThresholdV4    = big.NewInt(100000000)
	scraperMode             string
	blockRangeRecompute     int64
	historicalDataWritePath string
)

func init() {
	var err error
	// Import and cast environment variables.
	DIA_Meta_Contract_Address_V4 = utils.Getenv("DIA_META_CONTRACT_ADDRESS_V4", DIA_Meta_Contract_Address_V4)
	DIA_Meta_Contract_Precision_V4, err = strconv.Atoi(utils.Getenv("DIA_META_CONTRACT_PRECISION_V4", strconv.Itoa(DIA_Meta_Contract_Precision_V4)))
	if err != nil {
		log.Errorf("DIA_META_CONTRACT_PRECISION: %v", err)
	}
	priceMap_Update_Seconds_V4, err = strconv.Atoi(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_PRICE_MAP_UPDATE_SECONDS", strconv.Itoa(priceMap_Update_Seconds_V4)))
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAPV4_SIMULATION)+"_PRICE_MAP_UPDATE_SECONDS: %v", err)
	}
	simulation_Update_Seconds_V4, err = strconv.Atoi(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_SIMULATION_UPDATE_SECONDS", strconv.Itoa(simulation_Update_Seconds_V4)))
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAPV4_SIMULATION)+"_SIMULATION_UPDATE_SECONDS: %v", err)
	}
	scraperMode = utils.Getenv("SCRAPER_MODE", "")
	blockRangeRecompute, err = strconv.ParseInt(utils.Getenv("BLOCK_RANGE_RECOMPUTE", "30"), 10, 64)
	if err != nil {
		log.Error("parse BLOCK_RANGE_RECOMPUTE: ", err)
		blockRangeRecompute = 3000
	}
	historicalDataWritePath = utils.Getenv("HISTORICAL_DATA_WRITE_PATH", "/tmp/uniV4HistoricalTrades.json")

}

func NewUniswapV4Simulator(exchangepairs []models.ExchangePair, tradesChannel chan models.SimulatedTrade) {
	scraper := &UniswapV4Simulator{
		exchangepairs: exchangepairs,
		priceMap:      make(map[models.Asset]models.AssetQuotation),
	}

	var err error
	scraper.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_URI_REST", restUrl))
	if err != nil {
		log.Fatal("Failed to connect to ETH node: ", err)
	} else {
		log.Info("Successfully connected to node")
	}
	defer scraper.restClient.Close()

	scraper.luminaClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_LUMINA_URI_REST", restLuminaUrl))
	if err != nil {
		log.Error("init lumina client: ", err)
	} else {
		log.Info("Successfully connected to lumina node")
	}
	defer scraper.luminaClient.Close()

	liquidityThresholdV4Int, err := strconv.ParseInt(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_LIQUIDITY_THRESHOLD", "1000000000"), 10, 64)
	if err != nil {
		log.Errorf(strings.ToUpper(UNISWAPV4_SIMULATION)+"_LIQUIDITY_THRESHOLD: %v", err)
	} else {
		liquidityThresholdV4 = big.NewInt(liquidityThresholdV4Int)
	}

	scraper.simulator = simulation.New(scraper.restClient, log)

	quoterAddr := common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_QUOTER", "0x52f0e24d1c21c8a0cb1e5a5dd6198556bd9e1203"))
	scraper.quoter, err = v4quoter.NewV4Quoter(quoterAddr, scraper.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate V4Quoter: ", err)
	}

	scraper.poolManager, err = poolManager.NewPoolManager(common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_POOLMANAGER", "0x000000000004444c5dc75cB358380D2e3dE08A90")), scraper.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate PoolManager: ", err)
	}

	scraper.poolState, err = poolState.NewPoolState(common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_POOLSTATE", "0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227")), scraper.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate PoolState: ", err)
	}

	scraper.slippageThreshold, err = strconv.ParseFloat(utils.Getenv("UNISWAPV4_SLIPPAGE_THRESHOLD", "0.1"), 64)
	if err != nil {
		log.Error("Failed to parse slippahe threshold: ", err)
	}

	err = scraper.getExchangePairs()
	if err != nil {
		log.Fatal("Failed to get exchange pairs: ", err)
	}
	var lock sync.RWMutex

	if scraperMode != "history" {
		scraper.updatePriceMap(&lock)
		scraper.getSufficientLiquidityPools(&lock, big.NewInt(0))

		priceTicker := time.NewTicker(time.Duration(priceMapUpdateSeconds) * time.Second)
		go func() {
			for range priceTicker.C {
				scraper.updatePriceMap(&lock)
				scraper.getSufficientLiquidityPools(&lock, big.NewInt(0))
			}
		}()
	}

	if scraperMode == "history" {
		log.Infof("scraping in %s mode", scraperMode)

		startblock, err := strconv.ParseInt(utils.Getenv("UNIV4_SIMULATION_STARTBLOCK", "23224056"), 10, 64)
		if err != nil {
			log.Fatal("parse UNIV4_SIMULATION_STARTBLOCK: ", err)
		}
		endblock, err := strconv.ParseInt(utils.Getenv("UNIV4_SIMULATION_ENDBLOCK", "23224256"), 10, 64)
		if err != nil {
			log.Fatal("parse UNIV4_SIMULATION_ENDBLOCK: ", err)
		}
		if endblock == 0 {
			latestBlock, _ := scraper.restClient.BlockByNumber(context.Background(), big.NewInt(10))
			endblock = int64(latestBlock.NumberU64())
		}

		var allTrades []models.SimulatedTrade
		go func() {
			for t := range tradesChannel {
				allTrades = append(allTrades, t)
			}
		}()

		// TO DO: Can we remove the price map for history mode?
		scraper.updatePriceMap(&lock)
		scraper.getSufficientLiquidityPools(&lock, big.NewInt(startblock))

		for startblock < endblock {
			scraper.simulateTrades(tradesChannel, big.NewInt(startblock))
			startblock++

			// recompute liquidity every n blocks
			if startblock%blockRangeRecompute == 0 {
				log.Info("recompute prices and liquidity...")
				scraper.updatePriceMap(&lock)
				scraper.getSufficientLiquidityPools(&lock, big.NewInt(startblock))
			}
		}

		// Create (or overwrite) a file
		file, err := os.Create(historicalDataWritePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// Encode the struct slice into JSON and write it to the file
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ") // pretty print with indentation
		if err := encoder.Encode(allTrades); err != nil {
			panic(err)
		}
		log.Info("allTrades written.")
		return

	}

	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		scraper.simulateTrades(tradesChannel, big.NewInt(0))
	}
}

// getExchangePairs maps the underlying assets on @scraper.exchangepairs slice, including the
// specific handling of ETH asset. It also initializes the @scraper.priceMap.
func (scraper *UniswapV4Simulator) getExchangePairs() error {
	for i, ep := range scraper.exchangepairs {
		var quote, base models.Asset
		var err error
		var zeroAddressHex = common.Address{}.Hex()

		if ep.UnderlyingPair.QuoteToken.Address == zeroAddressHex {
			quote.Address = "0x0000000000000000000000000000000000000000"
			quote.Decimals = 18
			quote.Symbol = "ETH"
			quote.Blockchain = utils.ETHEREUM
		} else {
			quote, err = models.GetAsset(common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address), Exchanges[UNISWAPV4_SIMULATION].Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
		}

		if ep.UnderlyingPair.BaseToken.Address == zeroAddressHex {
			base.Address = "0x0000000000000000000000000000000000000000"
			base.Decimals = 18
			base.Symbol = "ETH"
			quote.Blockchain = utils.ETHEREUM
		} else {
			base, err = models.GetAsset(common.HexToAddress(ep.UnderlyingPair.BaseToken.Address), Exchanges[UNISWAPV4_SIMULATION].Blockchain, scraper.restClient)
			if err != nil {
				return err
			}
		}

		scraper.exchangepairs[i].UnderlyingPair.QuoteToken = quote
		scraper.exchangepairs[i].UnderlyingPair.BaseToken = base
		scraper.priceMap[quote] = models.AssetQuotation{}
		scraper.priceMap[base] = models.AssetQuotation{}
	}
	return nil
}

// updatePriceMap fetches the latest price of all observed assets from the diadata meta contract.
// If not existing in the meta contract, if fetches the price from diadata API.
func (scraper *UniswapV4Simulator) updatePriceMap(lock *sync.RWMutex) {
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

// getSufficientLiquidityPools fetches all pools given through the POOLS env var which
// fulfill given liquidity criteria.
func (scraper *UniswapV4Simulator) getSufficientLiquidityPools(lock *sync.RWMutex, blockNumber *big.Int) {

	scraper.allPools = make(map[models.ExchangePair]map[common.Hash]PoolState)
	for _, ep := range scraper.exchangepairs {
		pools := scraper.getPoolsBasedOnExchangePair(ep, blockNumber)
		if scraper.allPools[ep] == nil {
			scraper.allPools[ep] = make(map[common.Hash]PoolState)
		}
		for poolId, poolState := range pools {
			lock.Lock()
			scraper.allPools[ep][poolId] = poolState
			lock.Unlock()
		}
	}
}

// getPoolsBasedOnExchangePair returns all UniV4 pools for a given asset pair in a @models.Exchangepair.
func (scraper *UniswapV4Simulator) getPoolsBasedOnExchangePair(ep models.ExchangePair, blockNumber *big.Int) map[common.Hash]PoolState {
	pools := make(map[common.Hash]PoolState)

	for feeStr, tickSpacing := range feeToTickSpacing {
		_, _, poolId, liquidity, params := scraper.getPoolState(ep, feeStr, tickSpacing)

		// Check if the pool exists
		if poolId != (common.Hash{}) && liquidity != nil {
			// Check if the pool has sufficient liquidity
			if liquidity.Cmp(liquidityThresholdV4) >= 0 {

				// Check if potential simulation is valid
				callOpts := &bind.CallOpts{Context: context.Background()}
				if blockNumber.Cmp(big.NewInt(0)) != 0 {
					callOpts.BlockNumber = blockNumber
				}
				_, err := scraper.quoter.QuoteExactInputSingle(callOpts, params)
				if err != nil {
					log.Warnf("Invalid Pool - QuoteExactInputSingle failed for pool %s - %s - %s (%s): %v", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex(), err)
					continue
				} else {
					log.Infof("Pool %s - %s - %s (%s) has sufficient liquidity: %s", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex(), liquidity.String())
					pools[poolId] = PoolState{
						FeeStr:      feeStr,
						TickSpacing: tickSpacing,
						Params:      params,
						Liquidity:   liquidity,
					}
				}

			} else {
				log.Warnf("Pool %s - %s - %s (%s) does not have sufficient liquidity: %s", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex(), liquidity.String())
				// continue
			}
		} else {
			log.Warnf("Pool %s - %s - %s (%s) does not exist", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex())
			continue
		}
	}

	return pools
}

// getAmountIn returns the amountIn for a trade simulation corresponding to the base token amount
// equivalent to @amountInUSDConstant (usually 100$).
func (scraper *UniswapV4Simulator) getAmountIn(ep models.ExchangePair) (*big.Int, float64) {

	// log.Warnf("compute amount in for %s-%s", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
	basePrice := scraper.priceMap[ep.UnderlyingPair.BaseToken].Price
	amountInFloat := amountInUSDConstant / basePrice // 100
	// log.Warnf("basePrice -- amountInBase: %v -- %v", basePrice, amountInFloat)

	decimals := big.NewInt(int64(ep.UnderlyingPair.BaseToken.Decimals)) // e.g. 18
	exponent := new(big.Int).Exp(big.NewInt(10), decimals, nil)         // e.g. 10^18
	exponentFloat := new(big.Float).SetInt(exponent)

	amountIn := new(big.Float).Mul(big.NewFloat(amountInFloat), exponentFloat) // e.g. 10^20
	amountInBig := new(big.Int)
	amountIn.Int(amountInBig)

	// log.Warnf("getAmountIn -------------- amountInInt -- amountInAfterDecimalsAdjust: %v -- %s", amountInBig.String(), amountIn.String())

	return amountInBig, amountInFloat
}

// getPoolState returns all parameters needed for a simulation of the unique pool
// given by an asset pair, a fee value and a tick spacing.
func (scraper *UniswapV4Simulator) getPoolState(ep models.ExchangePair, feeStr string, tickSpacing *big.Int) (*big.Int, float64, common.Hash, *big.Int, v4quoter.IV4QuoterQuoteExactSingleParams) {
	fee := new(big.Int)
	fee.SetString(feeStr, 10)

	// Sort tokens in lexicographical order w.r.t. address. Also amend zeroForOne var below as needed.
	var poolKey v4quoter.PoolKey
	swapOrder := sortExchangepairLexicographically(ep)
	if swapOrder {
		poolKey = v4quoter.PoolKey{
			Currency0:   common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
			Currency1:   common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
			Fee:         fee,
			TickSpacing: tickSpacing,
			Hooks:       common.Address{},
		}
	} else {
		poolKey = v4quoter.PoolKey{
			Currency0:   common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
			Currency1:   common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
			Fee:         fee,
			TickSpacing: tickSpacing,
			Hooks:       common.Address{},
		}
	}

	amountInInt, amountInAfterDecimalAdjust := scraper.getAmountIn(ep)

	// In Uniswap V4, the poolId is computed from the poolKey. There is no pool Address for each pool.
	poolId, err := computePoolId(poolKey)
	if err != nil {
		log.Warnf("computePoolId failed for %s: %v", feeStr, err)
		return amountInInt, amountInAfterDecimalAdjust, common.Hash{}, nil, v4quoter.IV4QuoterQuoteExactSingleParams{}
	}

	liquidity, err := scraper.poolState.GetLiquidity(&bind.CallOpts{Context: context.Background()}, poolId)
	if err != nil || liquidity.Cmp(big.NewInt(0)) == 0 {
		log.Warnf("Liquidity check failed for pool %s - %s - %s", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr)
		return amountInInt, amountInAfterDecimalAdjust, poolId, nil, v4quoter.IV4QuoterQuoteExactSingleParams{}
	}

	params := v4quoter.IV4QuoterQuoteExactSingleParams{
		PoolKey:     poolKey,
		ZeroForOne:  swapOrder, // if assets are in correct order, we should swap asset1 for asset0.
		ExactAmount: amountInInt,
		HookData:    []byte{},
	}

	return amountInInt, amountInAfterDecimalAdjust, poolId, liquidity, params
}

// simulateTrades simulates trades in all filtered pools.
func (scraper *UniswapV4Simulator) simulateTrades(tradesChannel chan models.SimulatedTrade, blockNumber *big.Int) {
	var wg sync.WaitGroup
	// Sample allPools : {WBTC/USDC: {poolId: {params, liquidity, feeStr, tickSpacing}}}
	for ep, pools := range scraper.allPools {
		quoteToken := ep.UnderlyingPair.QuoteToken
		baseToken := ep.UnderlyingPair.BaseToken
		for poolId, poolStates := range pools {
			params := poolStates.Params
			// log.Infof("poolkey -- zeroforone --  exactamount -- hookdata: %v -- %v -- %v -- %v ", params.PoolKey, params.ZeroForOne, params.ExactAmount, params.HookData)
			liquidity := poolStates.Liquidity
			feeStr := poolStates.FeeStr
			wg.Add(1)
			go func(ep models.ExchangePair) {
				defer wg.Done()
				amountInInt, amountIn := scraper.getAmountIn(ep)

				amountOut, err := scraper.simulator.Execute(scraper.quoter, params, blockNumber)
				if err != nil {
					log.Warnf("QuoteExactInputSingle failed: %v", err)
					return
				}
				// log.Infof("amountOut for fee %s: %v", feeStr, amountOut.AmountOut.String())
				// log.Infof("params for fee %s: %v -- %v -- %s", feeStr, params.ZeroForOne, params.PoolKey, params.ExactAmount.String())

				callOptsPoolState := &bind.CallOpts{Context: context.Background()}
				if blockNumber.Cmp(big.NewInt(0)) != 0 {
					callOptsPoolState.BlockNumber = blockNumber
				}
				poolState, err := scraper.poolState.GetSlot0(callOptsPoolState, poolId)
				if err != nil || poolState.SqrtPriceX96.Cmp(big.NewInt(0)) == 0 {
					log.Errorf("Error getting sqrtPriceX96: %v", err)
					return
				}

				// since the trade is from token0 to token1, the slippage is computed as the amount of token1 received / amount of token0 sent
				slippage := computeSlippage(poolState.SqrtPriceX96, amountInInt, liquidity)
				log.Infof("Slippage for pool %v | %v - %v - fee %v: %v", poolId.Hex(), ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, slippage)

				if slippage > scraper.slippageThreshold {
					log.Warnf("Slippage for pool %v - %v - %v is greater than threshold %v: %v", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, scraper.slippageThreshold, slippage)
					return
				}

				amountOutFloat, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(amountOut.AmountOut), new(big.Float).SetFloat64(math.Pow10(int(quoteToken.Decimals)))).Float64()
				price := amountIn / amountOutFloat

				trade := models.SimulatedTrade{
					Price:       price,
					Volume:      amountOutFloat,
					QuoteToken:  quoteToken,
					BaseToken:   baseToken,
					PoolAddress: poolId.Hex(),
					Time:        time.Now(),
					Exchange:    Exchanges[UNISWAPV4_SIMULATION],
				}

				// Determine timestamp for historical blocks
				if blockNumber.Cmp(big.NewInt(0)) != 0 {
					block, err := scraper.restClient.BlockByNumber(context.Background(), blockNumber)
					if err != nil {
						log.Error("BlockByNumber: ", err)
					}

					trade.Time = time.Unix(int64(block.Time()), 0)
				}

				tradesChannel <- trade
				fee, _ := strconv.ParseFloat(feeStr, 64)
				log.Infof("Got trade in pool %v%%: %s-%s -- %v -- %v", fee/float64(10000), quoteToken.Symbol, baseToken.Symbol, trade.Price, trade.Volume)
			}(ep)
		}
		wg.Wait()
	}
}

// getPriceFromAPI returns the price of @asset in diadata API.
func (scraper *UniswapV4Simulator) getPriceFromAPI(asset models.Asset) float64 {
	log.Warnf("Could not determine price of %s on chain. Checking DIA API.", asset.Symbol)
	price, err := utils.GetPriceFromDiaAPI(asset.Address, asset.Blockchain)
	if err != nil {
		price = 1
		log.Errorf("Failed to get price of %s (blockchain -- address: %s -- %s) from DIA API: %v. Set to default %v", asset.Symbol, asset.Blockchain, asset.Address, err, price)
	} else {
		log.Infof("Fetched price of %s from DIA API: %.4f", asset.Symbol, price)
	}
	return price
}

func computePoolId(poolKey v4quoter.PoolKey) (common.Hash, error) {
	uint24Type, err := abi.NewType("uint24", "", nil)
	if err != nil {
		return common.Hash{}, err
	}
	int24Type, err := abi.NewType("int24", "", nil)
	if err != nil {
		return common.Hash{}, err
	}
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return common.Hash{}, err
	}

	arguments := abi.Arguments{
		{Type: addressType}, // currency0
		{Type: addressType}, // currency1
		{Type: uint24Type},
		{Type: int24Type},
		{Type: addressType}, // hooks
	}

	// ABI encode
	packed, err := arguments.Pack(
		poolKey.Currency0,
		poolKey.Currency1,
		poolKey.Fee,
		poolKey.TickSpacing,
		poolKey.Hooks,
	)
	if err != nil {
		return common.Hash{}, err
	}

	// keccak256
	hash := sha3.NewLegacyKeccak256()
	hash.Write(packed)
	var poolId common.Hash
	hash.Sum(poolId[:0])
	return poolId, nil
}

func computeSlippage(sqrtPriceX96 *big.Int, amountIn *big.Int, liquidity *big.Int) (slippage float64) {
	// log.Infof("sqrtPrice -- amount0 -- amount1 -- liquidity: %v -- %v -- %v -- %v", sqrtPriceX96, amount0, amount1, liquidity)

	price := new(big.Float).Quo(big.NewFloat(0).SetInt(sqrtPriceX96), new(big.Float).SetFloat64(math.Pow(2, 96)))

	// token0 -> token1 since ZeroForOne is true
	amount0Abs := big.NewInt(0).Abs(amountIn)
	numerator := big.NewFloat(0).Mul(big.NewFloat(0).SetInt(amount0Abs), price)
	slippage, _ = new(big.Float).Quo(numerator, big.NewFloat(0).SetInt(liquidity)).Float64()

	return slippage
}

// sortExchangepairLexicographically returns true, if the order of the assets has to be swapped in quoter contract.
// In UniV4 quoter contract, assets must be sorted in lexicographical order w.r.t. address when quoteExactInputSingle is queried.
// @swap can thereby also be used as @zeroForOne variable.
func sortExchangepairLexicographically(ep models.ExchangePair) (swap bool) {
	if ep.UnderlyingPair.QuoteToken.Address < ep.UnderlyingPair.BaseToken.Address {
		return
	}
	swap = true
	return
}
