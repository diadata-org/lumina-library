package filters

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

//	VWAP = Σ(price_i × |volume_i|) / Σ(|volume_i|)
//
// Returns the VWAP price (in USD), the timestamp of the most recent trade,
// and an error if no usable trades are available.
func VWAPFilter(tradesblock models.TradesBlock, basePrice float64, toleranceSeconds int64) (float64, float64, time.Time, error) {
	if len(tradesblock.Trades) == 0 {
		return 0, 0, time.Time{}, fmt.Errorf(
			"VWAPFilter: no trades available for %s-%s",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	if basePrice == 0 {
		return 0, 0, time.Time{}, fmt.Errorf(
			"VWAPFilter: basePrice is zero for %s-%s",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	tradeVolumes := make([]utils.TradeVolume, 0, len(tradesblock.Trades))
	var latestTime time.Time

	// Cutoff is relative to the block's EndTime rather than wall-clock time.
	// This is intentional: for merged blocks, EndTime reflects when the trigger
	// fired and defines the observation window consistently across all source blocks.
	cutoff := tradesblock.EndTime.Add(-time.Duration(toleranceSeconds) * time.Second)
	exchangeSet := make(map[string]struct{})
	for _, t := range tradesblock.Trades {
		if t.Time.Before(cutoff) {
			log.Debugf("VWAPFilter: skipping stale trade for %s-%s from %s (trade time: %v, cutoff:%v)",
				tradesblock.Pair.QuoteToken.Symbol,
				tradesblock.Pair.BaseToken.Symbol,
				t.Exchange.Name,
				t.Time,
				cutoff,
			)
			continue
		}

		tradeVolumes = append(tradeVolumes, utils.TradeVolume{
			Price:  t.Price,
			Volume: t.Volume,
		})
		if t.Time.After(latestTime) {
			latestTime = t.Time
		}
		exchangeSet[t.Exchange.Name] = struct{}{}
	}
	var exchanges []string
	for name := range exchangeSet {
		exchanges = append(exchanges, name)
	}
	sort.Strings(exchanges)
	exchangeList := strings.Join(exchanges, ", ")

	if len(tradeVolumes) == 0 {
		return 0, 0, time.Time{}, fmt.Errorf(
			"VWAPFilter: all trades are stale for %s-%s",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	// Sort by volume ascending.
	sorted := utils.SortByVolume(tradeVolumes)
	medianIdx := len(sorted) / 2
	log.Debugf(
		"VWAPFilter: %s-%s [%s] — median volume trade: price=%.6f volume=%.6f (%d total trades)",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		exchangeList,
		sorted[medianIdx].Price,
		sorted[medianIdx].Volume,
		len(sorted),
	)

	// Remove the single lowest-volume and single highest-volume trade.
	trimmed := utils.TrimExtremesByVolume(sorted)
	log.Debugf(
		"VWAPFilter: %s-%s [%s] — %d trades after trimming extremes",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		exchangeList,
		len(trimmed),
	)

	rawVwap, totalVolume := utils.VWAPWithVolume(trimmed)
	vwap := basePrice * rawVwap
	if vwap == 0 {
		return 0, 0, latestTime, fmt.Errorf(
			"VWAPFilter: VWAP is zero for %s-%s (all volumes may be zero)",
			tradesblock.Pair.QuoteToken.Symbol,
			tradesblock.Pair.BaseToken.Symbol,
		)
	}

	log.Infof(
		"VWAPFilter: %s-%s [%s] → %.6f USD (basePrice=%.6f, totalVolume=%.6f, %d/%d trades used)",
		tradesblock.Pair.QuoteToken.Symbol,
		tradesblock.Pair.BaseToken.Symbol,
		exchangeList,
		vwap,
		basePrice,
		totalVolume,
		len(trimmed),
		len(tradeVolumes),
	)

	return vwap, totalVolume, latestTime, nil
}
