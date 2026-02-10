package uniswapv4

import (
	"context"
	"math/big"

	v4quoter "github.com/diadata-org/lumina-library/contracts/uniswapv4/V4Quoter"
	"github.com/sirupsen/logrus"

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

func (c *Simulator) Execute(caller *v4quoter.V4Quoter, params v4quoter.IV4QuoterQuoteExactSingleParams, blockNumber *big.Int) (struct {
	AmountOut   *big.Int
	GasEstimate *big.Int
}, error) {
	callOpts := &bind.CallOpts{Context: context.Background()}
	if blockNumber.Cmp(big.NewInt(0)) != 0 {
		callOpts.BlockNumber = blockNumber
	}
	amountOut, err := caller.QuoteExactInputSingle(callOpts, params)
	return amountOut, err
}
