package scrapers

import (
	"testing"
)

func TestGetSwapData(t *testing.T) {
	cases := []struct {
		name  string
		swap  UniswapSwap
		wantP float64
		wantV float64
	}{
		{
			name: "buy swap (Amount0In > 0)",
			swap: UniswapSwap{
				Amount0In:  10,
				Amount0Out: 0,
				Amount1In:  0,
				Amount1Out: 200,
			},
			wantP: 20,  // 200/10
			wantV: -10, // negative Amount0In
		},
		{
			name: "sell swap (Amount0Out > 0)",
			swap: UniswapSwap{
				Amount0In:  0,
				Amount0Out: 5,
				Amount1In:  100,
				Amount1Out: 0,
			},
			wantP: 20, // 100/5
			wantV: 5,  // Amount0Out
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			price, volume := getSwapData(tc.swap)
			if price != tc.wantP {
				t.Errorf("price: got %v, want %v", price, tc.wantP)
			}
			if volume != tc.wantV {
				t.Errorf("volume: got %v, want %v", volume, tc.wantV)
			}
		})
	}
}
