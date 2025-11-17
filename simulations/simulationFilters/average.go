package simulationfilters

import (
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// Average returns the average price of all @trades.
// If @USDPrice=true it returns a USD price.
// basePrice is only evaluated for now, i.e. not for each trade's specific timestamp.
func Average(
	trades []models.SimulatedTrade,
	USDPrice bool,
	metacontractClient *ethclient.Client,
	metacontractAddress string,
	metacontractPrecision int,
) (avgPrice float64, timestamp time.Time, err error) {

	var prices []float64
	var basePriceMap = make(map[models.Asset]float64)

	// Fetch USD price of basetoken.
	if USDPrice {
		for _, t := range trades {

			if _, ok := basePriceMap[t.BaseToken]; !ok {
				basePrice, err := t.BaseToken.GetPrice(common.HexToAddress(metacontractAddress), metacontractPrecision, metacontractClient)
				if err != nil {
					log.Errorf("GetPrice for %s -- %s: %v ", t.BaseToken.Blockchain, t.BaseToken.Address, err)
					continue
				}
				prices = append(prices, basePrice.Price*t.Price)
				basePriceMap[t.BaseToken] = basePrice.Price
			} else {
				prices = append(prices, basePriceMap[t.BaseToken]*t.Price)
			}

			avgPrice = utils.Average(prices)

		}
	} else {

		for _, t := range trades {
			prices = append(prices, t.Price)
		}
		avgPrice = utils.Average(prices)

	}

	return
}
