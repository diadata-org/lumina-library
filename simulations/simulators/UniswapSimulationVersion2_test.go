package simulators

import (
	"math"
	"math/big"
	"testing"

	"github.com/diadata-org/lumina-library/models"
)

// NOTE:
// getAmountIn() uses float64 -> big.Float scaling -> Int() truncation.
// float64 * 10^decimals can be slightly below the exact integer due to binary float precision,
// e.g. 12.345678*1e6 may become 12345677.99999999..., then Int() becomes 12345677.
// So for fractional inputs we allow an off-by-one tolerance on the scaled integer.

func TestGetAmountIn_TableDriven(t *testing.T) {
	type tc struct {
		name        string
		decimals    int
		amountIn    float64
		wantScaled  string // expected integer string after scaling by 10^decimals (ideal exact)
		allowOffBy1 bool
	}

	tests := []tc{
		{
			name:        "6 decimals: 12.345678 -> 12345678",
			decimals:    6,
			amountIn:    12.345678,
			wantScaled:  "12345678",
			allowOffBy1: true, // fractional + float64 precision risk
		},
		{
			name:        "6 decimals: 0.000001 -> 1",
			decimals:    6,
			amountIn:    0.000001,
			wantScaled:  "1",
			allowOffBy1: true,
		},
		{
			name:        "18 decimals: 1.0 -> 1e18",
			decimals:    18,
			amountIn:    1.0,
			wantScaled:  "1000000000000000000",
			allowOffBy1: false, // integer inputs should be stable
		},
		{
			name:        "8 decimals: 100.0 -> 100e8",
			decimals:    8,
			amountIn:    100.0,
			wantScaled:  "10000000000",
			allowOffBy1: false,
		},
		{
			name:        "0 decimals: 123.0 -> 123",
			decimals:    0,
			amountIn:    123.0,
			wantScaled:  "123",
			allowOffBy1: false,
		},
	}

	// Create a minimal scraper instance (adjust this if your struct needs extra fields).
	s := &SimulationScraperVersion2{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenIn := models.Asset{
				// Decimals in models.Asset is often uint8 in DIA codebases
				Decimals: uint8(tt.decimals),
			}

			gotInt, gotFloat := s.getAmountIn(tokenIn, tt.amountIn)

			wantInt, ok := new(big.Int).SetString(tt.wantScaled, 10)
			if !ok {
				t.Fatalf("invalid wantScaled: %s", tt.wantScaled)
			}

			// Check scaled integer with tolerance when needed
			if gotInt.Cmp(wantInt) != 0 {
				if tt.allowOffBy1 {
					diff := new(big.Int).Sub(gotInt, wantInt)
					diff.Abs(diff)
					if diff.Cmp(big.NewInt(1)) > 0 {
						t.Fatalf("amountInInt mismatch (tolerance ±1 exceeded): got=%s want=%s diff=%s",
							gotInt.String(), wantInt.String(), diff.String())
					}
				} else {
					t.Fatalf("amountInInt mismatch: got=%s want=%s", gotInt.String(), wantInt.String())
				}
			}

			// Check the returned decimal-adjusted float equals input (within a tiny tolerance)
			gotF64, _ := gotFloat.Float64()
			if math.Abs(gotF64-tt.amountIn) > 1e-12 {
				t.Fatalf("amountInAfterDecimalAdjust mismatch: got=%.18f want=%.18f", gotF64, tt.amountIn)
			}
		})
	}
}
