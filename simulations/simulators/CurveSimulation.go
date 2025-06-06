package simulators

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefi"
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
	simulator         *simulation.Simulator
	exchangepairs     []models.ExchangePair
	thresholdSlippage float64
}

var (
	restDialUrl     = ""
	registryAddress = "0x90E00ACe148ca3b23Ac1bC8C240C2a7Dd9c2d7f5"
)

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

	scraper.thresholdSlippage, err = strconv.ParseFloat(utils.Getenv("CURVE_THRESHOLD_SLIPPAGE", "3"), 64)
	if err != nil {
		log.Error("Parse THRESHOLD_SLIPPAGE: ", err)
		scraper.thresholdSlippage = 3
	}

	scraper.simulator = simulation.New(scraper.restClient, log)
	scraper.exchangepairs = exchangepairs
	scraper.getExchangePairs()

	for _, ep := range scraper.exchangepairs {
		pools := scraper.getPools(ep)
		scraper.simulateTrades(ep, pools, tradesChannel)
	}

}

func (scraper *CurveSimulator) getExchangePairs() (err error) {
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
	}
	return
}

func (scraper *CurveSimulator) getPools(ep models.ExchangePair) []common.Address {

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
	var pools []common.Address
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
		pools = append(pools, pool)
	}

	log.Infof("Found %d %v/%v liquidity pools:\n", len(pools), ep.UnderlyingPair.QuoteToken.Symbol, ep.UnderlyingPair.BaseToken.Symbol)
	for i, pool := range pools {
		log.Infof("%d. %s\n", i+1, pool.Hex())
	}
	return pools
}

func (scraper *CurveSimulator) simulateTrades(ep models.ExchangePair, pools []common.Address, tradesChannel chan models.SimulatedTrade) {
	for _, poolAddr := range pools {
		pool, err := curvepool.NewCurvepool(poolAddr, scraper.restClient)
		if err != nil {
			log.Infof("Failed to create pool contract for %s: %v", poolAddr.Hex(), err)
			continue
		}

		// Get token indices and decimals
		quoteTokenIndex, baseTokenIndex, ok := scraper.validatePoolTokens(ep, pool)
		if !ok {
			continue
		}
		log.Infof("Token validated! token0_index: %v, token1_index: %v\n", quoteTokenIndex, baseTokenIndex)
		// Prepare trade input (e.g., $100)
		amountIn := big.NewInt(100000000)
		power := ep.UnderlyingPair.QuoteToken.Decimals
		exponent := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(power)), nil)
		amountInAfterDecimalAdjust := new(big.Int).Quo(amountIn, exponent)
		amountInAfterDecimalAdjustF64, _ := amountInAfterDecimalAdjust.Float64()
		log.Infof("amountIn: %s\n", amountInAfterDecimalAdjust)

		// Run trade simulation
		var amountOutFloat *big.Float
		amountOut, fee, err := scraper.executeTradeSimulation(pool, quoteTokenIndex, baseTokenIndex, amountIn)
		if err != nil {
			continue
		} else {
			amountOutFloat = new(big.Float).SetInt(amountOut)
			power := ep.UnderlyingPair.BaseToken.Decimals
			divisor := new(big.Float).SetInt(
				new(big.Int).Exp(
					big.NewInt(10),
					big.NewInt(int64(power)),
					nil,
				),
			)
			amountOutFloat.Quo(amountOutFloat, divisor)
			fmt.Printf("amountOut: %v\n", amountOutFloat)
			amountOutAfterDecimalAdjustF64, _ := amountOutFloat.Float64()

			// Calculate slippage
			slippage, err := calculateCurveSlippage(ep, pool, quoteTokenIndex, baseTokenIndex, amountIn, amountOut)
			if err != nil {
				continue
			}

			log.Infof("slippage in pool %v: %v | Fee: %.4f%%\n", poolAddr.Hex(), slippage, fee)

			if slippage > scraper.thresholdSlippage {
				log.Warnf("slippage above threshold: %v > %v", slippage, scraper.thresholdSlippage)
			} else {
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

func (scraper *CurveSimulator) validatePoolTokens(ep models.ExchangePair, pool *curvepool.Curvepool) (i, j int, valid bool) {
	// Iterate over all possible token indices (Curve supports up to 8 tokens)
	var tokens []common.Address
	maxCoins := 8
	for idx := 0; idx < maxCoins; idx++ {
		addr, err := pool.UnderlyingCoins(&bind.CallOpts{}, big.NewInt(int64(idx)))
		if err != nil || addr == (common.Address{}) {
			break
		}
		tokens = append(tokens, addr)
	}

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

func (scraper *CurveSimulator) executeTradeSimulation(pool *curvepool.Curvepool, i, j int, amountIn *big.Int) (*big.Int, float64, error) {
	// Fetch trading fee
	fee, _ := pool.Fee(&bind.CallOpts{})
	feePercent := parseCurveFee(fee, 8)

	// Run trade simulation
	amountOut, err := pool.GetDyUnderlying(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
	if err != nil {
		log.Printf("Trade simulation failed: %v", err)
		return nil, 0, err
	}
	return amountOut, feePercent, nil
}

func parseCurveFee(fee *big.Int, decimals int) float64 {
	feeFloat := new(big.Float).SetInt(fee)
	divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
	result, _ := new(big.Float).Quo(feeFloat, divisor).Float64()
	return result
}

func calculateCurveSlippage(ep models.ExchangePair, pool *curvepool.Curvepool, i, j int, amountIn *big.Int, actualOutWithFee *big.Int,
) (float64, error) {
	// 1. Retrieve key parameters
	amp, _ := pool.A(&bind.CallOpts{})   // Amplification coefficient A
	balances := getAllPoolBalances(pool) // Get all token balances
	fee, _ := pool.Fee(&bind.CallOpts{}) // Fee

	// 2. Calculate theoretical output (ignoring fees)
	theoreticalOutNoFee := calcTheoreticalOutput(ep, amp, balances, i, j, amountIn)

	feeRate := new(big.Float).Quo(
		new(big.Float).SetInt(fee),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	feeMultiplier := new(big.Float).Sub(big.NewFloat(1), feeRate)

	actualOutNoFee := new(big.Int)
	actualOutWithFeeFloat := new(big.Float).SetInt(actualOutWithFee)
	actualOutNoFeeFloat := actualOutWithFeeFloat.Quo(actualOutWithFeeFloat, feeMultiplier)
	actualOutNoFeeFloat.Int(actualOutNoFee)

	// 3. Calculate slippage
	slippage := new(big.Float).Sub(theoreticalOutNoFee, actualOutWithFeeFloat)
	slippage.Quo(slippage, theoreticalOutNoFee)
	slippagePercent, _ := slippage.Float64()

	return math.Abs(slippagePercent * 100), nil
}

func getAllPoolBalances(pool *curvepool.Curvepool) []*big.Int {
	var balances []*big.Int
	for idx := 0; idx < 8; idx++ { // Curve supports up to 8 tokens
		bal, err := pool.Balances(&bind.CallOpts{}, big.NewInt(int64(idx)))
		if err != nil || bal.Cmp(big.NewInt(0)) == 0 {
			break
		}
		balances = append(balances, bal)
	}
	return balances
}

func calcTheoreticalOutput(ep models.ExchangePair, amp *big.Int, balances []*big.Int, i, j int, dx *big.Int) *big.Float {
	if amp.Cmp(big.NewInt(0)) == 0 || len(balances) < 2 {
		return big.NewFloat(0)
	}

	// 1. Get pool balances (mind the precision adjustment)
	x := new(big.Int).Set(balances[i]) // Input token balance (e.g., USDC)
	y := new(big.Int).Set(balances[j]) // Output token balance (e.g., DAI)

	// 2.  Handle different token precisions (draft)
	diff := int64(math.Abs(float64(ep.UnderlyingPair.QuoteToken.Decimals - ep.UnderlyingPair.BaseToken.Decimals)))
	precisionAdjust := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
	y = y.Mul(y, precisionAdjust)

	// 3. Calculate numerator：4*A*dx*y
	numerator := new(big.Int).Mul(amp, big.NewInt(4)) // 4*A
	numerator.Mul(numerator, dx)                      // 4*A*dx
	numerator.Mul(numerator, y)                       // 4*A*dx*y

	// 4. Calculate denominator：(4*A + 1)*(x + dx)
	fourA := new(big.Int).Mul(amp, big.NewInt(4))        // 4*A
	fourAPlus1 := new(big.Int).Add(fourA, big.NewInt(1)) // 4*A +1

	xPlusDx := new(big.Int).Add(x, dx)                   // x + dx
	denominator := new(big.Int).Mul(fourAPlus1, xPlusDx) // (4A+1)(x+dx)

	// 5. Compute the result
	numeratorF := new(big.Float).SetInt(numerator)
	denominatorF := new(big.Float).SetInt(denominator)
	dy := new(big.Float).Quo(numeratorF, denominatorF)

	// 6. Convert the result back to output token precision
	precisionAdjustF := new(big.Float).SetInt(precisionAdjust)
	dy = new(big.Float).Quo(dy, precisionAdjustF)

	return dy
}
