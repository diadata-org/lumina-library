package simulators

import (
	"context"
	"sync"

	"github.com/diadata-org/lumina-library/models"
)

type Simulator interface {
	TradesChannel() chan models.SimulatedTrade
	Close(cancel context.CancelFunc) error
	// Subscribe(pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error
}

func RunSimulator(
	ctx context.Context,
	exchange string,
	exchangePairs []models.ExchangePair,
	tradesChannel chan models.SimulatedTrade,
	wg *sync.WaitGroup,
) {
	switch exchange {
	// case UNISWAP_SIMULATION:
	// 	NewUniswapSimulator(exchangePairs, tradesChannel)
	case UNISWAP_SIMULATION_TEST:
		NewUniswapSimulatorVersion2(exchangePairs, tradesChannel)
	}
}
