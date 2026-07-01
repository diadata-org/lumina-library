package filters

import (
	"math"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

// EMA computes an Exponential Moving Average over the trades INSIDE THIS BLOCK only.
// - trades: trades collected in one atomic block
// - basePrice: denomination multiplier (pass 1 for native quote, or pass USD price of base asset for USD result)
// - period: EMA period N (alpha = 2/(N+1))
// Returns (emaDenominated, timestampOfResult). Timestamp = last trade time in the block (consistent with LastPrice).
func EMA(trades []models.Trade, basePrice float64, period int) (float64, time.Time) {
	if len(trades) == 0 {
		return 0, time.Time{}
	}
	if period < 1 {
		period = 1
	}
	alpha := 2.0 / (float64(period) + 1.0)

	// Seed EMA with the first valid price in this block
	var (
		ema      float64
		seeded   bool
		lastTime time.Time
	)
	for _, tr := range trades {
		if p := tr.Price; !math.IsNaN(p) && !math.IsInf(p, 0) {
			ema = p
			seeded = true
			lastTime = tr.Time
			break
		}
	}
	if !seeded {
		// No valid price; fall back to last trade
		last := models.GetLastTrade(trades)
		return basePrice * last.Price, last.Time
	}

	// Update EMA over the rest trades in this block
	for i := 1; i < len(trades); i++ {
		p := trades[i].Price
		if math.IsNaN(p) || math.IsInf(p, 0) {
			continue
		}
		ema = alpha*p + (1.0-alpha)*ema
		lastTime = trades[i].Time
	}

	return basePrice * ema, lastTime
}

// EMAWithPrev computes EMA using an OPTIONAL previous EMA state (to continue across blocks).
// - prevEMA == nil means "no prior state"; it will seed with the first valid price in this block.
// Returns (emaDenominated, timestampOfResult). Timestamp = last trade time in the block.
func EMAWithPrev(trades []models.Trade, basePrice float64, period int, prevEMA *float64) (float64, time.Time) {
	if len(trades) == 0 {
		if prevEMA != nil {
			return basePrice * (*prevEMA), time.Time{}
		}
		return 0, time.Time{}
	}
	if period < 1 {
		period = 1
	}
	alpha := 2.0 / (float64(period) + 1.0)

	var (
		ema      float64
		inited   bool
		lastTime time.Time
	)
	if prevEMA != nil {
		ema = *prevEMA
		inited = true
	}

	for _, tr := range trades {
		p := tr.Price
		if math.IsNaN(p) || math.IsInf(p, 0) {
			continue
		}
		if !inited {
			ema = p // seed with first valid price
			inited = true
		} else {
			ema = alpha*p + (1.0-alpha)*ema
		}
		lastTime = tr.Time
	}

	return basePrice * ema, lastTime
}
