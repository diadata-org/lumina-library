package scrapers

import (
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
)

func TestGetSwapDataSlipstream(t *testing.T) {
	tests := []struct {
		name          string
		amount0       float64
		amount1       float64
		wantPrice     float64
		wantVolume    float64
		expectZeroOut bool
	}{
		{
			name:       "normal positive amounts",
			amount0:    2.0,
			amount1:    4000.0,
			wantPrice:  2000.0,
			wantVolume: 2.0,
		},
		{
			name:       "negative amount0 keeps sign for volume, price uses abs",
			amount0:    -2.0,
			amount1:    4000.0,
			wantPrice:  2000.0,
			wantVolume: -2.0,
		},
		{
			name:          "zero amount0 -> returns (0,0)",
			amount0:       0.0,
			amount1:       123.0,
			wantPrice:     0.0,
			wantVolume:    0.0,
			expectZeroOut: true,
		},
		{
			name:       "negative amount1 also ok",
			amount0:    2.0,
			amount1:    -4000.0,
			wantPrice:  2000.0,
			wantVolume: 2.0,
		},
	}

	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			swap := AerodromeSlipstreamSwap{
				ID:      "0xdeadbeef",
				Amount0: tt.amount0,
				Amount1: tt.amount1,
			}

			price, vol := getSwapDataSlipstream(swap)

			if price != tt.wantPrice {
				t.Fatalf("price mismatch: got=%v want=%v", price, tt.wantPrice)
			}
			if vol != tt.wantVolume {
				t.Fatalf("volume mismatch: got=%v want=%v", vol, tt.wantVolume)
			}

			if tt.expectZeroOut && (price != 0 || vol != 0) {
				t.Fatalf("expected zero output, got price=%v vol=%v", price, vol)
			}
		})
	}
}

func TestMakeTradeSlipstream(t *testing.T) {
	addr := common.HexToAddress("0xb2cc224c1c9fee385f8ad6a55b4d94e92359dc59")
	ts := time.Unix(1700000000, 0)

	// makeTradeSlipstream logic is:
	// BaseToken  = token1
	// QuoteToken = token0
	pair := AerodromeSlipstreamPair{
		Token0: models.Asset{Symbol: "USDC", Decimals: 6},
		Token1: models.Asset{Symbol: "WETH", Decimals: 18},
	}

	tTrade := makeTradeSlipstream(
		pair,
		2000.0,  // price
		-1.2345, // volume (signed allowed)
		ts,
		addr,
		"0xtxhash",
		"AERODROME_SLIPSTREAM",
		"base",
	)

	if tTrade.Price != 2000.0 {
		t.Fatalf("Price mismatch: got=%v want=%v", tTrade.Price, 2000.0)
	}
	if tTrade.Volume != -1.2345 {
		t.Fatalf("Volume mismatch: got=%v want=%v", tTrade.Volume, -1.2345)
	}

	// QuoteToken should be token0
	if tTrade.QuoteToken.Symbol != "USDC" {
		t.Fatalf("QuoteToken mismatch: got=%s want=%s", tTrade.QuoteToken.Symbol, "USDC")
	}
	// BaseToken should be token1
	if tTrade.BaseToken.Symbol != "WETH" {
		t.Fatalf("BaseToken mismatch: got=%s want=%s", tTrade.BaseToken.Symbol, "WETH")
	}

	if tTrade.Time.Unix() != ts.Unix() {
		t.Fatalf("Time mismatch: got=%v want=%v", tTrade.Time.Unix(), ts.Unix())
	}
	if tTrade.PoolAddress != addr.Hex() {
		t.Fatalf("PoolAddress mismatch: got=%s want=%s", tTrade.PoolAddress, addr.Hex())
	}
	if tTrade.ForeignTradeID != "0xtxhash" {
		t.Fatalf("ForeignTradeID mismatch: got=%s want=%s", tTrade.ForeignTradeID, "0xtxhash")
	}
	if tTrade.Exchange.Name != "AERODROME_SLIPSTREAM" {
		t.Fatalf("Exchange.Name mismatch: got=%s want=%s", tTrade.Exchange.Name, "AERODROME_SLIPSTREAM")
	}
	if tTrade.Exchange.Blockchain != "base" {
		t.Fatalf("Exchange.Blockchain mismatch: got=%s want=%s", tTrade.Exchange.Blockchain, "base")
	}
}
