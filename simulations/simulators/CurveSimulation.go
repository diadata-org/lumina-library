package simulators

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefi"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	"github.com/diadata-org/lumina-library/models"
	simulation "github.com/diadata-org/lumina-library/simulations/simulators/curve"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CurveSimulator struct {
	restClient    *ethclient.Client
	simulator     *simulation.Simulator
	exchangepairs []models.ExchangePair
}

var (
	restDialUrl     = "https://ethereum-mainnet.core.chainstack.com/1ec9bd0dbe3a2d876bfc7bc6ee25d0e6"
	registryAddress = "0x6A8cbed756804B16E05E741eDaBd5cB544AE21bf"
)

func NewCurveSimulator(exchangepairs []models.ExchangePair, tradesChannel chan models.SimulatedTrade) {
	var (
		err     error
		scraper CurveSimulator
	)
	scraper.restClient, err = ethclient.Dial(restDialUrl)
	if err != nil {
		log.Error("init rest client: ", err)
	} else {
		log.Info("Successfully connected to node")
	}
	defer scraper.restClient.Close()

	scraper.simulator = simulation.New(scraper.restClient, log)
	scraper.exchangepairs = exchangepairs

	for _, ep := range scraper.exchangepairs {
		pools := scraper.getPools(ep)
		scraper.simulateTrades(ep, pools)
	}

}

func (scraper *CurveSimulator) getPools(ep models.ExchangePair) []common.Address {

	registry, err := curvefi.NewCurvefi(common.HexToAddress(registryAddress), scraper.restClient)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	} else {
		log.Info("Successfully created contract instance")
	}
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

	log.Infof("Found %d USDT/USDC liquidity pools:\n", len(pools))
	for i, pool := range pools {
		log.Infof("%d. %s\n", i+1, pool.Hex())
	}
	return pools
}

func (scraper *CurveSimulator) simulateTrades(ep models.ExchangePair, pools []common.Address) {
	for _, poolAddr := range pools {
		pool, err := curvepool.NewCurvepool(poolAddr, scraper.restClient)
		if err != nil {
			log.Printf("Failed to create pool contract for %s: %v", poolAddr.Hex(), err)
			continue
		}

		// Get token indices and decimals
		quoteTokenIndex, baseTokenIndex, ok := scraper.validatePoolTokens(ep, pool)
		if !ok {
			continue
		}

		// Prepare trade input (e.g., $100)
		amountIn := big.NewInt(10000000000)

		// Run trade simulation
		amountOut, fee, err := scraper.executeTradeSimulation(pool, quoteTokenIndex, baseTokenIndex, amountIn)
		if err != nil {
			continue
		}

		// Calculate slippage
		slippage, err := calculateCurveSlippage(ep, pool, quoteTokenIndex, baseTokenIndex, amountIn)
		if err != nil {
			continue
		}

		fmt.Printf("\nPool: %s\n", poolAddr.Hex())
		fmt.Printf("%s Index: %v\n", ep.UnderlyingPair.QuoteToken.Symbol, quoteTokenIndex)
		fmt.Printf("%s Index: %v\n", ep.UnderlyingPair.BaseToken.Symbol, baseTokenIndex)
		fmt.Printf("Input: %.2f %s\n", float64(amountIn.Int64()), ep.UnderlyingPair.QuoteToken.Symbol)
		fmt.Printf("Output: %.2f %s\n", float64(amountOut.Int64()), ep.UnderlyingPair.BaseToken.Symbol)
		fmt.Printf("Fee: %.4f%%\n", fee)
		fmt.Printf("Slippage: %.4f%%\n", slippage)
	}
}

func (scraper *CurveSimulator) validatePoolTokens(ep models.ExchangePair, pool *curvepool.Curvepool) (i, j int, valid bool) {
	// Iterate over all possible token indices (Curve supports up to 8 tokens)
	var tokens []common.Address
	maxCoins := 8
	for idx := 0; idx < maxCoins; idx++ {
		addr, err := pool.Coins(&bind.CallOpts{}, big.NewInt(int64(idx)))
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
		log.Infof(
			"The pool is missing either %s (%t) or %s (%t)",
			ep.UnderlyingPair.QuoteToken.Symbol, quoteTokenIndex != -1,
			ep.UnderlyingPair.BaseToken.Symbol, baseTokenIndex != -1,
		)
		return 0, 0, false
	}

	return quoteTokenIndex, baseTokenIndex, true
}

func (scraper *CurveSimulator) executeTradeSimulation(pool *curvepool.Curvepool, i, j int, amountIn *big.Int) (*big.Int, float64, error) {
	// Fetch trading fee
	fee, _ := pool.Fee(&bind.CallOpts{})
	feePercent := parseCurveFee(fee, 8)

	// Run trade simulation
	amountOut, err := pool.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
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

func calculateCurveSlippage(ep models.ExchangePair, pool *curvepool.Curvepool, i, j int, amountIn *big.Int, // 输入代币数量（已考虑小数位）
) (float64, error) {
	// 1. Retrieve key parameters
	amp, _ := pool.A(&bind.CallOpts{})   // Amplification coefficient A
	balances := getAllPoolBalances(pool) // Get all token balances
	fee, _ := pool.Fee(&bind.CallOpts{}) // Fee

	// 2. Adjust precision (assuming tokens have 6 decimals)
	amountInAdjusted := adjustDecimals(amountIn, int(ep.UnderlyingPair.QuoteToken.Decimals), 0) // Convert to contract precision

	// 3. Calculate theoretical output (ignoring fees)
	theoreticalOutNoFee := calcTheoreticalOutput(amp, balances, i, j, amountInAdjusted)

	// 4. Calculate actual output (including fees)
	actualOutWithFee, _ := pool.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountInAdjusted)

	feeRate := new(big.Float).Quo(
		new(big.Float).SetInt(fee),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	feeMultiplier := new(big.Float).Sub(big.NewFloat(1), feeRate)

	actualOutNoFee := new(big.Int)
	actualOutWithFeeFloat := new(big.Float).SetInt(actualOutWithFee)
	actualOutNoFeeFloat := actualOutWithFeeFloat.Quo(actualOutWithFeeFloat, feeMultiplier)
	actualOutNoFeeFloat.Int(actualOutNoFee)

	// 5. Calculate slippage
	slippage := new(big.Float).Sub(
		new(big.Float).SetInt(theoreticalOutNoFee),
		actualOutWithFeeFloat,
	)
	slippage.Quo(slippage, new(big.Float).SetInt(theoreticalOutNoFee))
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

func adjustDecimals(amount *big.Int, fromDecimals, toDecimals int) *big.Int {
	diff := fromDecimals - toDecimals
	if diff == 0 {
		return new(big.Int).Set(amount)
	}

	multiplier := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(math.Abs(float64(diff)))), nil)
	if diff > 0 {
		// Reduce precision: amount / 10^diff
		return new(big.Int).Div(amount, multiplier)
	} else {
		// Increase precision: amount * 10^diff
		return new(big.Int).Mul(amount, multiplier)
	}
}

func calcTheoreticalOutput(amp *big.Int, balances []*big.Int, i, j int, dx *big.Int) *big.Int {
	if amp.Cmp(big.NewInt(0)) == 0 || len(balances) < 2 {
		return big.NewInt(0)
	}

	// 1. Get pool balances (mind the precision adjustment)
	x := new(big.Int).Set(balances[i]) // Input token balance (e.g., crvUSD)
	y := new(big.Int).Set(balances[j]) // Output token balance (e.g., USDT)

	// 2.  Handle different token precisions (draft)
	precisionAdjust := big.NewInt(1e12)
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
	dy := new(big.Int).Div(numerator, denominator)

	// 6. Convert the result back to output token precision
	dy = dy.Div(dy, precisionAdjust)

	return dy
}
