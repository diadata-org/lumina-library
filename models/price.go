package models

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AssetQuotation is the most recent price point information on an asset.
type AssetQuotation struct {
	Asset  Asset     `json:"Asset"`
	Price  float64   `json:"Price"`
	Source string    `json:"Source"`
	Time   time.Time `json:"Time"`
}

// GetPriceBaseAsset returns the price of the base asset in an atomic tradesblock.
// It updates @priceCacheMap if necessary.
func GetPriceBaseAsset(
	tb TradesBlock,
	priceCacheMap map[string]float64,
	client *ethclient.Client,
	metacontractAddress string,
	metacontractPrecision int,
) (basePrice float64, err error) {
	var ok bool

	if len(tb.Trades) > 0 {
		basetoken := tb.Trades[0].BaseToken

		if basetoken.Blockchain == "Fiat" && basetoken.Address == "840" {
			basePrice = 1
			return
		}

		basePrice, ok = priceCacheMap[basetoken.AssetIdentifier()]
		if !ok {
			// Get price from metacontract with fallback diadata API.
			var assetQuotation AssetQuotation
			assetQuotation, err = basetoken.GetPrice(common.HexToAddress(metacontractAddress), metacontractPrecision, client)
			if err != nil {
				return
			}
			basePrice = assetQuotation.Price
			priceCacheMap[basetoken.AssetIdentifier()] = basePrice
		}

		if ok {
			log.Debugf("took price from cache: %s", basetoken.AssetIdentifier())
		}

		return

	}

	err = errors.New("atomic tradesblock is empty")
	return
}
