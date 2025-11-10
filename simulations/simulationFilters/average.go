package simulationfilters

import (
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	log "github.com/sirupsen/logrus"
)

// Average returns the average price of all @trades.
// If @USDPrice=true it returns a USD price.
// basePrice is only evaluated for now, i.e. not for each trade's specific timestamp.
func Average(trades []models.SimulatedTrade, USDPrice bool) (avgPrice float64, timestamp time.Time, err error) {

	var prices []float64
	var basePriceMap = make(map[models.Asset]float64)

	// Fetch USD price of basetoken from DIA API.
	if USDPrice {
		for _, t := range trades {

			if _, ok := basePriceMap[t.BaseToken]; !ok {
				// TO DO: We can change this to GetOnchainPrice in order to fetch price from Lumina.
				basePrice, err := t.BaseToken.GetPriceFromDiaAPI()
				if err != nil {
					log.Errorf("GetPriceFromDiaAPI for %s -- %s: %v ", t.BaseToken.Blockchain, t.BaseToken.Address, err)
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
