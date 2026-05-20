package metafilters

import (
	"math"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

const epsilon = 1e-6

// assetVWAPResult holds the expected output for a single asset after VWAPMeta.
type assetVWAPResult struct {
	value      float64
	name       string
	sourceType models.SourceType
}

// checkVWAPMetaResults verifies that the output of VWAPMeta matches the
// expected per-asset results, regardless of slice ordering.
func checkVWAPMetaResults(t *testing.T, got []models.FilterPointPair, want map[models.Asset]assetVWAPResult) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("result length: got %d, want %d", len(got), len(want))
	}

	// Index results by QuoteToken asset for order-independent comparison.
	gotMap := make(map[models.Asset]models.FilterPointPair, len(got))
	for _, fp := range got {
		gotMap[fp.Pair.QuoteToken] = fp
	}

	for asset, expected := range want {
		fp, ok := gotMap[asset]
		if !ok {
			t.Errorf("missing result for asset %v", asset)
			continue
		}
		if fp.Name != expected.name {
			t.Errorf("asset %v: Name = %q, want %q", asset, fp.Name, expected.name)
		}
		if fp.SourceType != expected.sourceType {
			t.Errorf("asset %v: SourceType = %q, want %q", asset, fp.SourceType, expected.sourceType)
		}
		if diff := math.Abs(fp.Value - expected.value); diff > epsilon {
			t.Errorf("asset %v: Value = %.8f, want %.8f (diff=%.8f)", asset, fp.Value, expected.value, diff)
		}
	}
}

func TestVWAPMeta(t *testing.T) {
	var (
		ETH  = models.Asset{Address: "0x0000000000000000000000000000000000000000", Blockchain: utils.ETHEREUM}
		BTC  = models.Asset{Address: "0x0000000000000000000000000000000000000000", Blockchain: utils.BITCOIN}
		USDC = models.Asset{Address: "", Blockchain: utils.ETHEREUM}
		now  = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	)

	makeFP := func(quote models.Asset, value, volume float64, t time.Time) models.FilterPointPair {
		return models.FilterPointPair{
			Pair:   models.Pair{QuoteToken: quote, BaseToken: USDC},
			Value:  value,
			Volume: volume,
			Name:   "vwap",
			Time:   t,
		}
	}

	tests := []struct {
		name         string
		filterPoints []models.FilterPointPair
		want         map[models.Asset]assetVWAPResult
	}{
		{
			// A single filter point passes through unchanged.
			name: "single source — value and name passed through",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 5.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2000.0, name: "vwap"},
			},
		},
		{
			// Two sources with equal volume → equal-weight average.
			// VWAP = (2000*1 + 2100*1) / (1+1) = 2050
			name: "two sources, equal volume — equal-weight average",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 1.0, now),
				makeFP(ETH, 2100.0, 1.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2050.0, name: "vwap"},
			},
		},
		{
			// Two sources with unequal volume → volume-weighted average.
			// VWAP = (2000*1 + 2100*3) / (1+3) = 8300/4 = 2075
			name: "two sources, unequal volume — volume-weighted average",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 1.0, now),
				makeFP(ETH, 2100.0, 3.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2075.0, name: "vwap"},
			},
		},
		{
			// Three sources with different volumes.
			// VWAP = (2000*1 + 2050*2 + 2100*4) / (1+2+4) = 14500/7 ≈ 2071.428...
			name: "three sources — correct weighted average",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 1.0, now),
				makeFP(ETH, 2050.0, 2.0, now),
				makeFP(ETH, 2100.0, 4.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 14500.0 / 7.0, name: "vwap"},
			},
		},
		{
			// Multiple assets: ETH and BTC filter points are aggregated independently.
			// ETH: (3000*2 + 3100*3) / 5 = 15300/5 = 3060
			// BTC: single point passthrough = 60000
			name: "multiple assets — aggregated independently",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 3000.0, 2.0, now),
				makeFP(ETH, 3100.0, 3.0, now),
				makeFP(BTC, 60000.0, 10.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 15300.0 / 5.0, name: "vwap"},
				BTC: {value: 60000.0, name: "vwap"},
			},
		},
		{
			// All filter points have volume=0 (e.g. produced by a non-VWAP filter type).
			// VWAPMeta falls back to equal-weight average.
			// avg(2000, 2100) = 2050
			name: "zero volume fallback — equal-weight average",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 0.0, now),
				makeFP(ETH, 2100.0, 0.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2050.0, name: "vwap"},
			},
		},
		{
			// Timestamp on output should be the latest among input filter points.
			// Verified indirectly: we check the output is not the zero time by
			// running the normal value check (time correctness is tested explicitly below).
			name: "latest timestamp is propagated",
			filterPoints: []models.FilterPointPair{
				makeFP(ETH, 2000.0, 1.0, now.Add(-10*time.Minute)),
				makeFP(ETH, 2100.0, 1.0, now),
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2050.0, name: "vwap"},
			},
		},
		{
			// SourceType must be propagated to the output so that downstream
			// consumers (e.g. GetOracleKey) apply the correct key prefix.
			// SIMULATION_SOURCE in particular triggers a "SIM:" prefix — losing
			// it would silently corrupt the simulation feed.
			name: "SourceType propagated from first input filter point",
			filterPoints: []models.FilterPointPair{
				{
					Pair:       models.Pair{QuoteToken: ETH, BaseToken: USDC},
					Value:      2000.0,
					Volume:     1.0,
					Time:       now,
					SourceType: models.SIMULATION_SOURCE,
				},
				{
					Pair:       models.Pair{QuoteToken: ETH, BaseToken: USDC},
					Value:      2100.0,
					Volume:     1.0,
					Time:       now,
					SourceType: models.SIMULATION_SOURCE,
				},
			},
			want: map[models.Asset]assetVWAPResult{
				ETH: {value: 2050.0, name: "vwap", sourceType: models.SIMULATION_SOURCE},
			},
		},
		{
			// Empty input returns an empty (nil) result slice without panicking.
			name:         "empty input — returns empty result",
			filterPoints: nil,
			want:         map[models.Asset]assetVWAPResult{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := VWAPMeta(tc.filterPoints)
			checkVWAPMetaResults(t, got, tc.want)
		})
	}
}

// TestVWAPMetaTimestamp verifies that VWAPMeta propagates the latest timestamp
// from each asset's input filter points to the corresponding output point.
func TestVWAPMetaTimestamp(t *testing.T) {
	USDC := models.Asset{Address: "", Blockchain: utils.ETHEREUM}
	ETH := models.Asset{Address: "0x0000000000000000000000000000000000000000", Blockchain: utils.ETHEREUM}

	earlier := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	later := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	filterPoints := []models.FilterPointPair{
		{Pair: models.Pair{QuoteToken: ETH, BaseToken: USDC}, Value: 2000.0, Volume: 1.0, Time: earlier},
		{Pair: models.Pair{QuoteToken: ETH, BaseToken: USDC}, Value: 2100.0, Volume: 1.0, Time: later},
	}

	got := VWAPMeta(filterPoints)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if !got[0].Time.Equal(later) {
		t.Errorf("Time = %v, want %v", got[0].Time, later)
	}
}

// TestVWAPMetaOrdering verifies that VWAPMeta returns results in a stable,
// deterministic order (sorted by QuoteToken.Address) regardless of map
// iteration order.
func TestVWAPMetaOrdering(t *testing.T) {
	USDC := models.Asset{Address: "", Blockchain: utils.ETHEREUM}
	// Three assets with addresses that sort in a known order.
	A := models.Asset{Address: "0x1111", Blockchain: utils.ETHEREUM}
	B := models.Asset{Address: "0x2222", Blockchain: utils.ETHEREUM}
	C := models.Asset{Address: "0x3333", Blockchain: utils.ETHEREUM}

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	filterPoints := []models.FilterPointPair{
		{Pair: models.Pair{QuoteToken: C, BaseToken: USDC}, Value: 300.0, Volume: 1.0, Time: now},
		{Pair: models.Pair{QuoteToken: A, BaseToken: USDC}, Value: 100.0, Volume: 1.0, Time: now},
		{Pair: models.Pair{QuoteToken: B, BaseToken: USDC}, Value: 200.0, Volume: 1.0, Time: now},
	}

	// Run multiple times to surface any non-determinism from map iteration.
	for i := 0; i < 20; i++ {
		got := VWAPMeta(filterPoints)
		if len(got) != 3 {
			t.Fatalf("run %d: expected 3 results, got %d", i, len(got))
		}
		if got[0].Pair.QuoteToken.Address != A.Address ||
			got[1].Pair.QuoteToken.Address != B.Address ||
			got[2].Pair.QuoteToken.Address != C.Address {
			t.Errorf("run %d: unexpected order: got [%s, %s, %s], want [%s, %s, %s]",
				i,
				got[0].Pair.QuoteToken.Address,
				got[1].Pair.QuoteToken.Address,
				got[2].Pair.QuoteToken.Address,
				A.Address, B.Address, C.Address,
			)
		}
	}
}