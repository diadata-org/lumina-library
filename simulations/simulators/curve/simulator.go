package uniswap

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	"github.com/diadata-org/lumina-library/contracts/curve/curvelendingpool"
	"github.com/diadata-org/lumina-library/contracts/curve/curveplain"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func (c *Simulator) Execute(poolAddr common.Address, restClient *ethclient.Client, usedUnderlying bool, i, j int, amountIn *big.Int) (*big.Int, error) {
	type attempt struct {
		name string
		call func() (*big.Int, error)
	}

	attempts_underlying := []attempt{
		{
			name: "curvepool.GetDyUnderlying",
			call: func() (*big.Int, error) {
				p, err := curvepool.NewCurvepool(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDyUnderlying(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
		{
			name: "curvelendingpool.GetDyUnderlying",
			call: func() (*big.Int, error) {
				p, err := curvelendingpool.NewCurvelendingpool(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDyUnderlying(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
	}
	attempts := []attempt{
		{
			name: "curvepool.GetDy",
			call: func() (*big.Int, error) {
				p, err := curvepool.NewCurvepool(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
		{
			name: "curveplain.GetDy",
			call: func() (*big.Int, error) {
				p, err := curveplain.NewCurveplainCaller(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
		{
			name: "curvelendingpool.GetDy",
			call: func() (*big.Int, error) {
				p, err := curvelendingpool.NewCurvelendingpool(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
		{
			name: "curvefifactory.GetDy",
			call: func() (*big.Int, error) {
				p, err := curvefifactory.NewCurvefifactory(poolAddr, restClient)
				if err != nil {
					return nil, err
				}
				return p.GetDy(&bind.CallOpts{}, big.NewInt(int64(i)), big.NewInt(int64(j)), amountIn)
			},
		},
	}

	if usedUnderlying {
		attempts = attempts_underlying
	}
	// Run trade simulation - i - intoken (e.g. USDT)
	for _, a := range attempts {
		amountOut, err := a.call()
		if err == nil {
			c.log.Infof("Simulator.Execute succeeded with %s", a.name)
			return amountOut, nil
		}
		c.log.Infof("intoken index: %v, outtoken index: %v", big.NewInt(int64(i)), big.NewInt(int64(j)))
		c.log.Warnf("Simulator.Execute failed with %s: %v", a.name, err)
	}

	return nil, fmt.Errorf("Simulator.Execute: all ABI calls failed for pool %s", poolAddr.Hex())
}
