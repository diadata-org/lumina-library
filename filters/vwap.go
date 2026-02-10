package filters

import (
	"math"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

// VWAP computes the Volume-Weighted Average Price for a slice of trades.
// - trades: the trades within one atomic block
// - basePrice: price of the base asset for denomination conversion
//   - If you need "native pricing" (e.g., "USDT per ETH" in USDT/ETH), pass 1
//   - If you need "USD pricing", pass the USD price of the base asset (GetPriceBaseAsset in processor)
//
// Returns: (vwapPriceDenominated, timestampOfResult)
// Timestamp: the last trade time in the block (consistent with LastPrice)
func VWAP(trades []models.Trade, basePrice float64) (float64, time.Time) {
	if len(trades) == 0 {
		return 0, time.Time{}
	}

	last := models.GetLastTrade(trades)

	var sumPV, sumV float64
	for _, tr := range trades {
		// use the absolute value of the volume to avoid the impact of the direction of the trade
		v := math.Abs(tr.Volume)
		if v == 0 || math.IsNaN(tr.Price) || math.IsInf(tr.Price, 0) {
			continue
		}
		sumPV += tr.Price * v
		sumV += v
	}

	// if no valid volume, return the last price to avoid division by zero
	if sumV == 0 {
		return basePrice * last.Price, last.Time
	}

	vwap := sumPV / sumV
	// if vwap is NaN or Inf, return the last price to avoid division by zero
	if math.IsNaN(vwap) || math.IsInf(vwap, 0) {
		return basePrice * last.Price, last.Time
	}

	return basePrice * vwap, last.Time
}
