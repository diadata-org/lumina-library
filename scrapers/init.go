package scrapers

import (
	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

const (
	BINANCE_EXCHANGE      = "Binance"
	COINBASE_EXCHANGE     = "CoinBase"
	BYBIT_EXCHANGE        = "ByBit"
	CRYPTODOTCOM_EXCHANGE = "Crypto.com"
	GATEIO_EXCHANGE       = "GateIO"
	KRAKEN_EXCHANGE       = "Kraken"
	KUCOIN_EXCHANGE       = "KuCoin"
	OKEX_EXCHANGE         = "OKEx"
	MEXC_EXCHANGE         = "MEXC"

	CURVE_EXCHANGE          = "Curve"
	UNISWAPV2_EXCHANGE      = "UniswapV2"
	UNISWAPV2_BASE_EXCHANGE = "UniswapV2_Base"
	UNISWAPV3_EXCHANGE      = "UniswapV3"
	UNISWAPV4_EXCHANGE      = "UniswapV4"
	UNISWAPV3_BASE_EXCHANGE = "UniswapV3_Base"
	PANCAKESWAPV3_EXCHANGE  = "PancakeswapV3"
	UNISWAP_SIMULATION      = "UniswapSimulation"
)

var (
	Exchanges = make(map[string]models.Exchange)
	log       *logrus.Logger
)

func init() {

	Exchanges[BYBIT_EXCHANGE] = models.Exchange{Name: BYBIT_EXCHANGE, Centralized: true}
	Exchanges[BINANCE_EXCHANGE] = models.Exchange{Name: BINANCE_EXCHANGE, Centralized: true}
	Exchanges[COINBASE_EXCHANGE] = models.Exchange{Name: COINBASE_EXCHANGE, Centralized: true}
	Exchanges[CRYPTODOTCOM_EXCHANGE] = models.Exchange{Name: CRYPTODOTCOM_EXCHANGE, Centralized: true}
	Exchanges[GATEIO_EXCHANGE] = models.Exchange{Name: GATEIO_EXCHANGE, Centralized: true}
	Exchanges[KRAKEN_EXCHANGE] = models.Exchange{Name: KRAKEN_EXCHANGE, Centralized: true}
	Exchanges[KUCOIN_EXCHANGE] = models.Exchange{Name: KUCOIN_EXCHANGE, Centralized: true}
	Exchanges[MEXC_EXCHANGE] = models.Exchange{Name: MEXC_EXCHANGE, Centralized: true}
	Exchanges[OKEX_EXCHANGE] = models.Exchange{Name: OKEX_EXCHANGE, Centralized: true}

	Exchanges[CURVE_EXCHANGE] = models.Exchange{Name: CURVE_EXCHANGE, Centralized: false, Blockchain: utils.ETHEREUM}
	Exchanges[PANCAKESWAPV3_EXCHANGE] = models.Exchange{Name: PANCAKESWAPV3_EXCHANGE, Centralized: false, Blockchain: utils.BINANCESMARTCHAIN}
	Exchanges[UNISWAP_SIMULATION] = models.Exchange{Name: UNISWAP_SIMULATION, Centralized: false, Simulation: true, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAPV2_EXCHANGE] = models.Exchange{Name: UNISWAPV2_EXCHANGE, Centralized: false, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAPV2_BASE_EXCHANGE] = models.Exchange{Name: UNISWAPV2_BASE_EXCHANGE, Centralized: false, Blockchain: utils.BASE}
	Exchanges[UNISWAPV3_EXCHANGE] = models.Exchange{Name: UNISWAPV3_EXCHANGE, Centralized: false, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAPV4_EXCHANGE] = models.Exchange{Name: UNISWAPV4_EXCHANGE, Centralized: false, Blockchain: utils.ETHEREUM}
	Exchanges[UNISWAPV3_BASE_EXCHANGE] = models.Exchange{Name: UNISWAPV3_BASE_EXCHANGE, Centralized: false, Blockchain: utils.BASE}

	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_SCRAPERS", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)

}
