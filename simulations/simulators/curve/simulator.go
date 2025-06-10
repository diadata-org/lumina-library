package uniswap

import (
	"math"
	"math/big"

	"github.com/sirupsen/logrus"

	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Simulator struct {
	Eth *ethclient.Client
	log *logrus.Logger
}

func New(client *ethclient.Client, log *logrus.Logger) *Simulator {
	c := Simulator{Eth: client, log: log}
	return &c
}

func (c *Simulator) Execute(pool *curvepool.Curvepool, i, j int, amountIn *big.Int) (*big.Int, float64, error) {
	// Fetch trading fee
	fee, _ := pool.Fee(&bind.CallOpts{})
	feePercent := parseCurveFee(fee, 8)

	// Run trade simulation
	amountOut, err := pool.GetDyUnderlying(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
	if err != nil {
		c.log.Printf("Trade simulation failed: %v", err)
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
