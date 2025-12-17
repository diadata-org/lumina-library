package simulators

import (
	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

const (
	UNISWAP_SIMULATION      = "UniswapSimulation1"
	UNISWAP_SIMULATION_TEST = "UniswapSimulation"
	CURVE_SIMULATION        = "CurveSimulation"
	UNISWAPV4_SIMULATION    = "UniswapV4Simulation"
)

var (
	Exchanges = make(map[string]models.Exchange)
	log       *logrus.Logger
)

func init() {

	Exchanges[UNISWAP_SIMULATION] = models.Exchange{Name: UNISWAP_SIMULATION, Centralized: false, Simulation: true, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAP_SIMULATION_TEST] = models.Exchange{Name: UNISWAP_SIMULATION_TEST, Centralized: false, Simulation: true, Blockchain: utils.ETHEREUM}
	Exchanges[CURVE_SIMULATION] = models.Exchange{Name: CURVE_SIMULATION, Centralized: false, Simulation: true, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAPV4_SIMULATION] = models.Exchange{Name: UNISWAPV4_SIMULATION, Centralized: false, Simulation: true, Blockchain: utils.ETHEREUM}

	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_SCRAPERS", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)

}
