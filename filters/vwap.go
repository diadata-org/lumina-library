package filters

import (
	"fmt"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

//	VWAP = Σ(price_i × |volume_i|) / Σ(|volume_i|)
//
// Returns the VWAP price (in USD), the timestamp of the most recent trade,
// and an error if no usable trades are available.
func VWAPFilter(tradesblock models.TradesBlock, basePrice float64) (float64, time.Time, error) {
	if len(tradesblock.Trades) == 0 {
		return 0, time.Now(), fmt.Errorf(
			"VWAPFilter: no trades available for %s-%s",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	var tradeVolumes []utils.TradeVolume
	var latestTime time.Time

	for _, t := range tradesblock.Trades {
		tradeVolumes = append(tradeVolumes, utils.TradeVolume{
			Price:  t.Price,
			Volume: t.Volume,
		})
		if t.Time.After(latestTime) {
			latestTime = t.Time
		}
	}

	// Sort by volume ascending.
	sorted := utils.SortByVolume(tradeVolumes)
	medianIdx := len(sorted) / 2
	log.Debugf(
		"VWAPFilter: %s-%s on %s — median volume trade: price=%.6f volume=%.6f (%d total trades)",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		tradesblock.Trades[0].Exchange.Name,
		sorted[medianIdx].Price,
		sorted[medianIdx].Volume,
		len(sorted),
	)

	// Remove the single lowest-volume and single highest-volume trade.
	trimmed := utils.TrimExtremesByVolume(sorted)
	log.Debugf(
		"VWAPFilter: %s-%s on %s — %d trades after trimming extremes",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		tradesblock.Trades[0].Exchange.Name,
		len(trimmed),
	)

	vwap := basePrice * utils.VWAP(trimmed)
	if vwap == 0 {
		return 0, latestTime, fmt.Errorf(
			"VWAPFilter: VWAP is zero for %s-%s (all volumes may be zero)",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	log.Infof(
		"VWAPFilter: %s-%s on %s → %.6f USD (basePrice=%.6f, %d/%d trades used)",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		tradesblock.Trades[0].Exchange.Name,
		vwap,
		basePrice,
		len(trimmed),
		len(tradeVolumes),
	)

	return vwap, latestTime, nil
}
