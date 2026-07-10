package processor

import (
	"fmt"

	"github.com/diadata-org/lumina-library/filters"
	models "github.com/diadata-org/lumina-library/models"
)

func computeAtomicFilterValue(tb models.TradesBlock, filterType models.FilterType, basePrice float64, toleranceSeconds int64) (atomicFilterValue float64, atomicVolume float64, err error) {
	switch filterType {
	case models.FILTER_VWAP:
		atomicFilterValue, atomicVolume, _, err = filters.VWAPFilter(tb, basePrice, toleranceSeconds)
		if err != nil {
			err = fmt.Errorf("VWAP filter: %w", err)
			return
		}
		log.Debugf(
			"Processor - VWAP filter value for pair %s-%s [%s] with %v trades: %v (volume: %v).",
			tb.Pair.QuoteToken.Symbol,
			tb.Pair.BaseToken.Symbol,
			tb.Trades[0].Exchange.Name,
			len(tb.Trades),
			atomicFilterValue,
			atomicVolume,
		)
	case models.FILTER_LAST_PRICE:
		atomicFilterValue, _, err = filters.LastPrice(tb, basePrice)
		if err != nil {
			err = fmt.Errorf("LastPrice filter: %w", err)
			return
		}
		log.Debugf(
			"Processor - Atomic filter value for market %s with %v trades: %v.",
			tb.Trades[0].Exchange.Name+":"+tb.Trades[0].QuoteToken.Symbol+"-"+tb.Trades[0].BaseToken.Symbol+"--"+tb.Trades[0].QuoteToken.Name,
			len(tb.Trades),
			atomicFilterValue,
		)
	default:
		err = fmt.Errorf("Processor - no filter matched for filterType=%q, skipping update", filterType)
		return
	}
	return
}
