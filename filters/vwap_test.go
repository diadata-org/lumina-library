package filters

import (
	"strings"
	"testing"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

// baseTime is a fixed reference point; all trade timestamps are derived from it.
var baseTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
var startTime = baseTime.Add(-1 * time.Hour)

// makeBlock constructs a TradesBlock with EndTime = baseTime and StartTime = baseTime - 1h.
func makeBlock(trades []models.Trade) models.TradesBlock {
	return models.TradesBlock{
		Pair: models.Pair{
			QuoteToken: models.Asset{Symbol: "ETH"},
			BaseToken:  models.Asset{Symbol: "USDT"},
		},
		Trades:    trades,
		StartTime: startTime,
		EndTime:   baseTime,
	}
}

// freshTrade returns a trade well within the tolerance window.
func freshTrade(price, volume float64, exchangeName string) models.Trade {
	return models.Trade{
		Price:    price,
		Volume:   volume,
		Time:     baseTime.Add(-30 * time.Second),
		Exchange: models.Exchange{Name: exchangeName},
	}
}

// staleTrade returns a trade older than toleranceSeconds before EndTime.
func staleTrade(price, volume float64, exchangeName string, toleranceSeconds int64) models.Trade {
	return models.Trade{
		Price:    price,
		Volume:   volume,
		Time:     startTime.Add(-time.Duration(toleranceSeconds+60) * time.Second),
		Exchange: models.Exchange{Name: exchangeName},
	}
}

func TestVWAPFilter(t *testing.T) {
	const tol = int64(300) // 5 minutes
	const epsilon = 1e-6
	// basePrice = 1.0 throughout: assumes USDT/USD ≈ 1, so VWAP output is
	// directly comparable to the raw trade prices.
	const basePrice = 1.0

	tests := []struct {
		name          string
		block         models.TradesBlock
		basePrice     float64
		tolerance     int64
		wantErr       bool
		wantErrSubstr string
		wantPrice     float64
		wantVolume    float64 // total absolute volume after stale exclusion and trim
	}{
		{
			name:          "no trades at all",
			block:         makeBlock(nil),
			basePrice:     basePrice,
			tolerance:     tol,
			wantErr:       true,
			wantErrSubstr: "no trades available",
		},
		{
			name: "all trades stale",
			block: makeBlock([]models.Trade{
				staleTrade(2000, 1.0, "Binance", tol),
				staleTrade(2100, 2.0, "Coinbase", tol),
			}),
			basePrice:     basePrice,
			tolerance:     tol,
			wantErr:       true,
			wantErrSubstr: "all trades are stale",
		},
		{
			name: "basePrice is zero",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Binance"),
				freshTrade(2050, 2.0, "Coinbase"),
				freshTrade(2100, 3.0, "Kraken"),
			}),
			basePrice:     0,
			tolerance:     tol,
			wantErr:       true,
			wantErrSubstr: "basePrice is zero",
		},
		{
			// TrimExtremesByVolume returns the slice unchanged when len <= 2,
			// so a single trade survives and VWAP == its price.
			name: "single fresh trade — no trim, returns trade price",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Binance"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  2000.0,
			wantVolume: 1.0,
		},
		{
			// Two trades also bypass trim (len <= MinTrimSize=2).
			// VWAP = (2000*1 + 2100*3) / (1+3) = 8300/4 = 2075
			name: "two fresh trades — no trim, weighted average",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Binance"),
				freshTrade(2100, 3.0, "Coinbase"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  2075.0,
			wantVolume: 4.0,
		},
		{
			// Three trades: below MinSizeForTrimming=5, no trim applied.
			// All three trades participate in VWAP.
			// VWAP = (2000*1 + 2050*2 + 2100*5) / (1+2+5) = 26600/8 = 2075
			name: "three trades — below trim threshold, all trades used",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Binance"),
				freshTrade(2050, 2.0, "Coinbase"),
				freshTrade(2100, 5.0, "Uniswap"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  2075.0,
			wantVolume: 8.0,
		},
		{
			// Four trades: below MinSizeForTrimming=5, no trim applied.
			// All four trades participate in VWAP.
			// VWAP = (2000*1 + 2020*2 + 2060*4 + 2100*10) / (1+2+4+10) = 35280/17 ≈ 2075.294
			name: "four trades — below trim threshold, all trades used",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Kraken"),
				freshTrade(2020, 2.0, "Binance"),
				freshTrade(2060, 4.0, "Coinbase"),
				freshTrade(2100, 10.0, "Uniswap"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  35280.0 / 17.0,
			wantVolume: 17.0,
		},
		{
			// Stale trades are excluded before VWAP; only the two fresh ones remain.
			// len=2 so no trim; VWAP = (2000*1 + 2100*3) / 4 = 2075
			name: "mixed stale and fresh trades — stale excluded before VWAP",
			block: makeBlock([]models.Trade{
				staleTrade(1000, 100.0, "OldExchange", tol),
				freshTrade(2000, 1.0, "Binance"),
				freshTrade(2100, 3.0, "Coinbase"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  2075.0,
			wantVolume: 4.0,
		},
		{
			// Negative volume (sell-side DEX trade): VWAP uses abs(volume).
			// Three trades, below MinSizeForTrimming=5, no trim applied.
			// abs volumes: 1, 2, 5 → total = 8
			// VWAP = (2000*1 + 2050*2 + 2100*5) / 8 = 2075
			name: "negative volume sell trade — abs value used in VWAP",
			block: makeBlock([]models.Trade{
				freshTrade(2000, -1.0, "Uniswap"),
				freshTrade(2050, 2.0, "Binance"),
				freshTrade(2100, 5.0, "Coinbase"),
			}),
			basePrice:  basePrice,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  2075.0,
			wantVolume: 8.0,
		},
		{
			// basePrice != 1: result should be basePrice * VWAP(all trades, no trim).
			// Three trades, below MinSizeForTrimming=5, no trim.
			// VWAP = (2000*1 + 2050*2 + 2100*5) / 8 = 2075
			// result = 0.999 * 2075 = 2072.925
			name: "basePrice not 1.0 — scales output correctly",
			block: makeBlock([]models.Trade{
				freshTrade(2000, 1.0, "Binance"),
				freshTrade(2050, 2.0, "Coinbase"),
				freshTrade(2100, 5.0, "Kraken"),
			}),
			basePrice:  0.999,
			tolerance:  tol,
			wantErr:    false,
			wantPrice:  0.999 * 2075.0,
			wantVolume: 8.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			price, vol, ts, err := VWAPFilter(tc.block, tc.basePrice, tc.tolerance)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil (price=%.6f)", tc.wantErrSubstr, price)
				}
				if tc.wantErrSubstr != "" && !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Fatalf("expected error containing %q, got: %v", tc.wantErrSubstr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := tc.wantPrice - price; diff > epsilon || diff < -epsilon {
				t.Errorf("price: want %.6f, got %.6f (diff=%.6f)", tc.wantPrice, price, diff)
			}
			if diff := tc.wantVolume - vol; diff > epsilon || diff < -epsilon {
				t.Errorf("volume: want %.6f, got %.6f (diff=%.6f)", tc.wantVolume, vol, diff)
			}
			if ts.IsZero() {
				t.Errorf("expected non-zero timestamp, got zero time")
			}
		})
	}
}
