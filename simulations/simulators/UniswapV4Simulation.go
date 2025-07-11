package simulators

import (
	"context"
	"math"
	"math/big"
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
		"500":   big.NewInt(10),
		"3000":  big.NewInt(60),
		"10000": big.NewInt(200),
	}
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
}

func NewUniswapV4Simulator(exchangepairs []models.ExchangePair, tradesChannel chan models.SimulatedTrade) {
	s := &UniswapV4Simulator{
		exchangepairs: exchangepairs,
		priceMap:      make(map[models.Asset]models.AssetQuotation),
	}

	var err error
	s.restClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_URI_REST", restUrl))
	if err != nil {
		log.Fatal("Failed to connect to ETH node: ", err)
	} else {
		log.Info("Successfully connected to node")
	}
	defer s.restClient.Close()

	s.luminaClient, err = ethclient.Dial(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_LUMINA_URI_REST", restLuminaUrl))
	if err != nil {
		log.Error("init lumina client: ", err)
	} else {
		log.Info("Successfully connected to lumina node")
	}
	defer s.luminaClient.Close()

	quoterAddr := common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_QUOTER", "0x52f0e24d1c21c8a0cb1e5a5dd6198556bd9e1203"))
	s.quoter, err = v4quoter.NewV4Quoter(quoterAddr, s.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate V4Quoter: ", err)
	}

	s.poolManager, err = poolManager.NewPoolManager(common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_POOLMANAGER", "0x000000000004444c5dc75cB358380D2e3dE08A90")), s.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate PoolManager: ", err)
	}

	s.poolState, err = poolState.NewPoolState(common.HexToAddress(utils.Getenv(strings.ToUpper(UNISWAPV4_SIMULATION)+"_POOLSTATE", "0x7fFE42C4a5DEeA5b0feC41C94C136Cf115597227")), s.restClient)
	if err != nil {
		log.Fatal("Failed to instantiate PoolState: ", err)
	}

	s.slippageThreshold, err = strconv.ParseFloat(utils.Getenv("UNISWAPV4_SLIPPAGE_THRESHOLD", "0.1"), 64)
	if err != nil {
		log.Error("Failed to parse slippahe threshold: ", err)
	}

	err = s.getExchangePairs()
	if err != nil {
		log.Fatal("Failed to get exchange pairs: ", err)
	}

	var lock sync.RWMutex
	s.updatePriceMap(&lock)

	priceTicker := time.NewTicker(time.Duration(priceMapUpdateSeconds) * time.Second)
	go func() {
		for range priceTicker.C {
			s.updatePriceMap(&lock)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		s.simulateTrades(tradesChannel)
	}
}

func (s *UniswapV4Simulator) getExchangePairs() error {
	for i, ep := range s.exchangepairs {
		quote, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address), Exchanges[UNISWAPV4_SIMULATION].Blockchain, s.restClient)
		if err != nil {
			return err
		}
		base, err := models.GetAsset(common.HexToAddress(ep.UnderlyingPair.BaseToken.Address), Exchanges[UNISWAPV4_SIMULATION].Blockchain, s.restClient)
		if err != nil {
			return err
		}
		s.exchangepairs[i].UnderlyingPair.QuoteToken = quote
		s.exchangepairs[i].UnderlyingPair.BaseToken = base
		s.priceMap[quote] = models.AssetQuotation{}
		s.priceMap[base] = models.AssetQuotation{}
	}
	return nil
}

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

func (scraper *UniswapV4Simulator) getPriceFromAPI(asset models.Asset) float64 {
	log.Warnf("Could not determine price of %s on chain. Checking DIA API.", asset.Symbol)
	price, err := utils.GetPriceFromDiaAPI(asset.Address, asset.Blockchain)
	if err != nil {
		log.Errorf("Failed to get price of %s from DIA API: %v\n", asset.Symbol, err)
		log.Errorf("asset blockchain: %v\n", asset.Blockchain)
		log.Errorf("asset address: %v\n", asset.Address)
		price = 100
	} else {
		log.Infof("Fetched price of %s from DIA API: %.4f", asset.Symbol, price)
	}
	return price
}

func (s *UniswapV4Simulator) filterValidFeeTiers(ep models.ExchangePair) map[string]*big.Int {
	validFeeTiers := make(map[string]*big.Int)

	for feeStr, tickSpacing := range feeToTickSpacing {
		_, _, poolId, liquidity, params := s.getPoolState(ep, feeStr, tickSpacing)

		if poolId != (common.Hash{}) && liquidity != nil {
			_, err := s.quoter.QuoteExactInputSingle(&bind.CallOpts{Context: context.Background()}, params)
			if err != nil {
				log.Warnf("QuoteExactInputSingle failed for pool %s - %s - %s (%s): %v", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex(), err)
				continue
			}
		} else {
			log.Warnf("Pool %s - %s - %s (%s) is not valid", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr, poolId.Hex())
			continue
		}
		validFeeTiers[feeStr] = tickSpacing
	}

	return validFeeTiers
}

func (s *UniswapV4Simulator) getAmountIn(ep models.ExchangePair) (*big.Int, *big.Float) {
	basePrice := s.priceMap[ep.UnderlyingPair.BaseToken].Price
	amountInBase := amountInUSDConstant / basePrice // 100

	decimals := big.NewInt(int64(ep.UnderlyingPair.QuoteToken.Decimals)) // e.g. 18
	exponent := new(big.Int).Exp(big.NewInt(10), decimals, nil)          // e.g. 10^18
	exponentFloat := new(big.Float).SetInt(exponent)

	amountIn := new(big.Float).Mul(big.NewFloat(amountInBase), exponentFloat) // e.g. 10^20
	amountInInt := new(big.Int)
	amountIn.Int(amountInInt)

	amountInAfterDecimalAdjust := new(big.Float).Quo(amountIn, exponentFloat) // e.g. 10^2

	return amountInInt, amountInAfterDecimalAdjust
}

func (s *UniswapV4Simulator) getPoolState(ep models.ExchangePair, feeStr string, tickSpacing *big.Int) (*big.Int, *big.Float, common.Hash, *big.Int, v4quoter.IV4QuoterQuoteExactSingleParams) {
	fee := new(big.Int)
	fee.SetString(feeStr, 10)
	poolKey := v4quoter.PoolKey{
		Currency0:   common.HexToAddress(ep.UnderlyingPair.QuoteToken.Address),
		Currency1:   common.HexToAddress(ep.UnderlyingPair.BaseToken.Address),
		Fee:         fee,
		TickSpacing: tickSpacing,
		Hooks:       common.Address{},
	}

	amountInInt, amountInAfterDecimalAdjust := s.getAmountIn(ep)

	poolId, err := ComputePoolId(poolKey)
	if err != nil {
		log.Warnf("ComputePoolId failed for %s: %v", feeStr, err)
		return amountInInt, amountInAfterDecimalAdjust, common.Hash{}, nil, v4quoter.IV4QuoterQuoteExactSingleParams{}
	}

	liquidity, err := s.poolState.GetLiquidity(&bind.CallOpts{Context: context.Background()}, poolId)
	if err != nil || liquidity.Cmp(big.NewInt(0)) == 0 {
		log.Warnf("Liquidity check failed for pool %s - %s - %s", ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol, feeStr)
		return amountInInt, amountInAfterDecimalAdjust, poolId, nil, v4quoter.IV4QuoterQuoteExactSingleParams{}
	}

	params := v4quoter.IV4QuoterQuoteExactSingleParams{
		PoolKey:     poolKey,
		ZeroForOne:  true,
		ExactAmount: amountInInt, // small dummy amount
		HookData:    []byte{},
	}

	return amountInInt, amountInAfterDecimalAdjust, poolId, liquidity, params
}

func (s *UniswapV4Simulator) simulateTrades(tradesChannel chan models.SimulatedTrade) {
	var wg sync.WaitGroup

	for _, ep := range s.exchangepairs {
		validFeeTiers := s.filterValidFeeTiers(ep)
		for feeStr, tickSpacing := range validFeeTiers {
			wg.Add(1)
			go func(ep models.ExchangePair) {
				defer wg.Done()

				amountInInt, amountInAfterDecimalAdjust, poolId, liquidity, params := s.getPoolState(ep, feeStr, tickSpacing)

				log.Infof("amountIn after adjusting for decimals: %v\n", amountInAfterDecimalAdjust)

				log.Infof(
					"[UniswapV4 Simulate] PoolKey: {Token0: %s, Token1: %s, Fee: %d, TickSpacing: %s, Hooks: %s} | ZeroForOne: %v | ExactAmount: %s | HookData: %x",
					params.PoolKey.Currency0.Hex(),
					params.PoolKey.Currency1.Hex(),
					params.PoolKey.Fee,
					params.PoolKey.TickSpacing.String(),
					params.PoolKey.Hooks.Hex(),
					params.ZeroForOne,
					params.ExactAmount.String(),
					params.HookData,
				)

				amountOut, err := s.quoter.QuoteExactInputSingle(&bind.CallOpts{Context: context.Background()}, params)
				if err != nil {
					log.Warnf("QuoteExactInputSingle failed: %v", err)
					return
				}
				poolState, err := s.poolState.GetSlot0(&bind.CallOpts{Context: context.Background()}, poolId)
				if err != nil {
					log.Fatalf("Error getting sqrtPriceX96: %v", err)
					return
				}

				slippage := computeSlippage(poolState.SqrtPriceX96, amountInInt, amountOut.AmountOut, liquidity)
				log.Infof("Slippage: %v", slippage)

				if slippage > s.slippageThreshold {
					log.Warnf("Slippage is greater than threshold %v: %v", s.slippageThreshold, slippage)
					return
				}

				amountOutFloat := new(big.Float).SetInt(amountOut.AmountOut)
				power := ep.UnderlyingPair.BaseToken.Decimals
				divisor := new(big.Float).SetInt(
					new(big.Int).Exp(
						big.NewInt(10),
						big.NewInt(int64(power)),
						nil,
					),
				)
				amountOutFloat.Quo(amountOutFloat, divisor)
				amountOutAfterDecimalAdjustF64, _ := amountOutFloat.Float64()
				price, _ := new(big.Float).Quo(amountOutFloat, amountInAfterDecimalAdjust).Float64()

				trade := models.SimulatedTrade{
					Price:       price,
					Volume:      amountOutAfterDecimalAdjustF64,
					QuoteToken:  ep.UnderlyingPair.QuoteToken,
					BaseToken:   ep.UnderlyingPair.BaseToken,
					PoolAddress: poolId.Hex(),
					Time:        time.Now(),
					Exchange:    Exchanges[UNISWAPV4_SIMULATION],
				}
				tradesChannel <- trade
				log.Infof(
					"[base: %v/ quote: %v - fee: %v - tickSpacing: %v] amountIn: %s | amountOut: %s | price: %.6f | volume: %.6f",
					ep.UnderlyingPair.BaseToken.Symbol,
					ep.UnderlyingPair.QuoteToken.Symbol,
					feeStr,
					tickSpacing.String(),
					amountInAfterDecimalAdjust.String(),
					amountOut.AmountOut.String(),
					trade.Price,
					trade.Volume,
				)
			}(ep)
		}
		wg.Wait()
	}
}

func ComputePoolId(poolKey v4quoter.PoolKey) (common.Hash, error) {
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

func computeSlippage(sqrtPriceX96 *big.Int, amount0 *big.Int, amount1 *big.Int, liquidity *big.Int) (slippage float64) {
	log.Infof("sqrtPrice -- amount0 -- amount1 -- liquidity: %v -- %v -- %v -- %v", sqrtPriceX96, amount0, amount1, liquidity)

	price := new(big.Float).Quo(big.NewFloat(0).SetInt(sqrtPriceX96), new(big.Float).SetFloat64(math.Pow(2, 96)))

	// token0 -> token1 since ZeroForOne is true
	amount0Abs := big.NewInt(0).Abs(amount0)
	numerator := big.NewFloat(0).Mul(big.NewFloat(0).SetInt(amount0Abs), price)
	slippage, _ = new(big.Float).Quo(numerator, big.NewFloat(0).SetInt(liquidity)).Float64()

	return slippage
}
