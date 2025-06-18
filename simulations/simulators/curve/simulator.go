package uniswap

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	"github.com/diadata-org/lumina-library/contracts/curve/curveplain"
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

func (c *Simulator) Execute(pool interface{}, i, j int, amountIn *big.Int) (*big.Int, error) {
	var (
		amountOut *big.Int
		err       error
	)

	// Run trade simulation - i - intoken (e.g. USDT)
	switch p := pool.(type) {
	case *curveplain.CurveplainCaller:
		amountOut, err = p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
	case *curvepool.Curvepool:
		amountOut, err = p.GetDyUnderlying(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
	case *curvefifactory.Curvefifactory:
		amountOut, err = p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
	default:
		return nil, fmt.Errorf("unsupported pool contract type")
	}

	if err != nil {
		c.log.Infof("intoken index: %v, outtoken index: %v", big.NewInt(int64(i)), big.NewInt(int64(j)))
		c.log.Printf("Trade simulation failed: %v", err)
		return nil, err
	}
	return amountOut, nil
}
