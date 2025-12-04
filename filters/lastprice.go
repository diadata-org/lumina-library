package filters

import (
	"fmt"
	"math"
	"strconv"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	utils "github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

var (
	log                  *logrus.Logger
	tradeVolumeThreshold float64
)

func init() {
	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_FILTERS", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)

	tradeVolumeThreshold, err = strconv.ParseFloat(utils.Getenv("TRADE_VOLUME_THRESHOLD_LAST_TRADE_FILTER", "0.0001"), 64)
	if err != nil {
		log.Errorf("parse TRADE_VOLUME_THRESHOLD_LAST_TRADE_FILTER %v. Set to default 0.0001", err)
		tradeVolumeThreshold = float64(0.0001)
	}
}

// LastPrice returns the price of the latest trade.
// If price should be returned in terms of native currency @basePrice should be set to 1.
func LastPrice(tradesblock models.TradesBlock, basePrice float64) (float64, time.Time, error) {

	if len(tradesblock.Trades) == 0 {
		return 0, time.Now(), fmt.Errorf("no trades available")
	}

	index, lastTrade := tradesblock.GetLastTrade()

	for len(tradesblock.Trades) > 0 {

		if math.Abs(basePrice*lastTrade.Price*lastTrade.Volume) > tradeVolumeThreshold {
			return basePrice * lastTrade.Price, lastTrade.Time, nil
		} else {
			log.Warnf("discard low volume trade %s -- %s on %s: %v", lastTrade.QuoteToken.Symbol, lastTrade.BaseToken.Symbol, lastTrade.Exchange.Name, basePrice*lastTrade.Price*lastTrade.Volume)
			err := tradesblock.RemoveTradeByIndex(index)
			if err != nil {
				return 0, lastTrade.Time, fmt.Errorf("no trade above volume threshold available for %s -- %s on %s", lastTrade.QuoteToken.Symbol, lastTrade.BaseToken.Symbol, lastTrade.Exchange.Name)
			}
			if len(tradesblock.Trades) == 0 {
				return 0, lastTrade.Time, fmt.Errorf("no trade above volume threshold available for %s -- %s on %s", lastTrade.QuoteToken.Symbol, lastTrade.BaseToken.Symbol, lastTrade.Exchange.Name)
			}
			index, lastTrade = tradesblock.GetLastTrade()
		}
	}
	return 0, lastTrade.Time, fmt.Errorf("no trade above volume threshold available for %s -- %s on %s", lastTrade.QuoteToken.Symbol, lastTrade.BaseToken.Symbol, lastTrade.Exchange.Name)
}
