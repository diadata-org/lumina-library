package utils

import (
	"math"
	"testing"
)

// TO DO: Write more test cases.

func TestMedian(t *testing.T) {
	cases := []struct {
		samples []float64
		median  float64
	}{
		{
			[]float64{21.2415, 3.4421, 24.1490, 71.1216, 19.47, 11.3313, 24.77809, 13.3166, 22.9814},
			21.2415,
		},
		{
			[]float64{3.31, 2.33, 9.01, 3.24, 1.53, 1.14},
			2.785,
		},
		{
			[]float64{1},
			1.0,
		},
		{
			[]float64{},
			0.0,
		},
	}

	for i, c := range cases {
		median := Median(c.samples)
		if math.Abs(float64(median-c.median)) > 1e-4 {
			t.Errorf("Median was incorrect, got: %f, expected: %f for set:%d", median, c.median, i)
		}
	}

}

func TestVWAP(t *testing.T) {
	cases := []struct {
		name     string
		trades   []TradeVolume
		expected float64
	}{
		{
			name:     "empty slice returns 0",
			trades:   []TradeVolume{},
			expected: 0,
		},
		{
			name:     "single trade returns its price",
			trades:   []TradeVolume{{Price: 50000, Volume: 1}},
			expected: 50000,
		},
		{
			name: "equal volumes - VWAP equals simple average",
			trades: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 1},
				{Price: 300, Volume: 1},
			},
			expected: 200,
		},
		{
			name: "unequal volumes - higher volume trade has more weight",
			// price=100 vol=9, price=200 vol=1 → VWAP = (100*9 + 200*1) / 10 = 110
			trades: []TradeVolume{
				{Price: 100, Volume: 9},
				{Price: 200, Volume: 1},
			},
			expected: 110,
		},
		{
			name: "negative volume treated as absolute value",
			// DEX sell-side trades have negative volume; math.Abs should be applied.
			// price=100 |vol|=2, price=200 |vol|=2 → VWAP = 150
			trades: []TradeVolume{
				{Price: 100, Volume: -2},
				{Price: 200, Volume: 2},
			},
			expected: 150,
		},
		{
			name: "zero volume trades are ignored",
			trades: []TradeVolume{
				{Price: 99999, Volume: 0},
				{Price: 200, Volume: 2},
			},
			expected: 200,
		},
		{
			name: "all zero volume returns 0",
			trades: []TradeVolume{
				{Price: 100, Volume: 0},
				{Price: 200, Volume: 0},
			},
			expected: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := VWAP(c.trades)
			if math.Abs(got-c.expected) > 1e-9 {
				t.Errorf("VWAP = %.9f, want %.9f", got, c.expected)
			}
		})
	}
}

func TestSortByVolume(t *testing.T) {
	cases := []struct {
		name     string
		input    []TradeVolume
		wantVols []float64
	}{
		{
			name:     "empty slice",
			input:    []TradeVolume{},
			wantVols: []float64{},
		},
		{
			name:     "single trade unchanged",
			input:    []TradeVolume{{Price: 100, Volume: 5}},
			wantVols: []float64{5},
		},
		{
			name: "already sorted ascending",
			input: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 2},
				{Price: 300, Volume: 3},
			},
			wantVols: []float64{1, 2, 3},
		},
		{
			name: "reverse order gets sorted",
			input: []TradeVolume{
				{Price: 300, Volume: 3},
				{Price: 200, Volume: 2},
				{Price: 100, Volume: 1},
			},
			wantVols: []float64{1, 2, 3},
		},
		{
			name: "negative volumes sorted by absolute value",
			input: []TradeVolume{
				{Price: 100, Volume: -3},
				{Price: 200, Volume: 1},
				{Price: 300, Volume: -2},
			},
			wantVols: []float64{1, -2, -3},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			originalInput := make([]TradeVolume, len(c.input))
			copy(originalInput, c.input)

			got := SortByVolume(c.input)

			if len(got) != len(c.wantVols) {
				t.Fatalf("SortByVolume len = %d, want %d", len(got), len(c.wantVols))
			}
			for i, tv := range got {
				if tv.Volume != c.wantVols[i] {
					t.Errorf("SortByVolume[%d].Volume = %.2f, want %.2f", i, tv.Volume, c.wantVols[i])
				}
			}
			// Verify original slice was not modified.
			for i := range c.input {
				if c.input[i] != originalInput[i] {
					t.Errorf("SortByVolume modified original slice at index %d", i)
				}
			}
		})
	}
}

func TestTrimVolumeOutliers(t *testing.T) {
	cases := []struct {
		name     string
		input    []TradeVolume
		wantVols []float64
	}{
		{
			name:     "empty slice unchanged",
			input:    []TradeVolume{},
			wantVols: []float64{},
		},
		{
			name:     "single trade unchanged",
			input:    []TradeVolume{{Price: 100, Volume: 1}},
			wantVols: []float64{1},
		},
		{
			name: "two trades unchanged",
			input: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 2},
			},
			wantVols: []float64{1, 2},
		},
		{
			// With MinSizeForTrimming=5, three trades are below the threshold
			// and are returned unchanged (no trimming applied).
			name: "three trades - below threshold, unchanged",
			input: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 2},
				{Price: 300, Volume: 3},
			},
			wantVols: []float64{1, 2, 3},
		},
		{
			// Four trades are also below the threshold — unchanged.
			name: "four trades - below threshold, unchanged",
			input: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 2},
				{Price: 300, Volume: 3},
				{Price: 400, Volume: 4},
			},
			wantVols: []float64{1, 2, 3, 4},
		},
		{
			name: "five trades - removes lowest and highest volume",
			input: []TradeVolume{
				{Price: 100, Volume: 1},
				{Price: 200, Volume: 2},
				{Price: 300, Volume: 3},
				{Price: 400, Volume: 4},
				{Price: 500, Volume: 5},
			},
			wantVols: []float64{2, 3, 4},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := TrimVolumeOutliers(c.input)
			if len(got) != len(c.wantVols) {
				t.Fatalf("TrimVolumeOutliers len = %d, want %d", len(got), len(c.wantVols))
			}
			for i, tv := range got {
				if tv.Volume != c.wantVols[i] {
					t.Errorf("TrimVolumeOutliers[%d].Volume = %.2f, want %.2f", i, tv.Volume, c.wantVols[i])
				}
			}
		})
	}
}