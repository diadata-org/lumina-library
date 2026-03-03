package models

import "testing"

func TestAsset_GetOracleKey_NormalizesAndFormats(t *testing.T) {
	a := Asset{
		Symbol:     "usdc",         // should become "USDC"
		Blockchain: " base  ",      // should become "BASE"
		Address:    "0xAbCDeF123 ", // should become lower + trimmed
	}

	got := a.GetOracleKey()
	want := "USDC/USD:BASE/0xabcdef123"

	if got != want {
		t.Fatalf("GetOracleKey() = %q, want %q", got, want)
	}
}

func TestAsset_GetOracleKey_EmptyChainOrAddress_ReturnsEmpty(t *testing.T) {
	tests := []struct {
		name string
		a    Asset
	}{
		{
			name: "empty blockchain",
			a: Asset{
				Symbol:     "USDC",
				Blockchain: "",
				Address:    "0xabc",
			},
		},
		{
			name: "empty address",
			a: Asset{
				Symbol:     "USDC",
				Blockchain: "ETH",
				Address:    "",
			},
		},
		{
			name: "both empty",
			a: Asset{
				Symbol:     "USDC",
				Blockchain: "",
				Address:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.GetOracleKey(); got != "" {
				t.Fatalf("GetOracleKey() = %q, want empty string", got)
			}
		})
	}
}

func TestPair_GetOracleKey_PrefixBySourceType(t *testing.T) {
	p := Pair{
		QuoteToken: Asset{
			Symbol:     "usdc",
			Blockchain: "eth",
			Address:    "0xAbC",
		},
		BaseToken: Asset{Symbol: "weth"},
	}

	quoteKey := p.QuoteToken.GetOracleKey()
	if quoteKey == "" {
		t.Fatalf("precondition failed: quote token GetOracleKey() returned empty")
	}

	tests := []struct {
		name       string
		sourceType SourceType
		want       string
	}{
		{
			name:       "empty sourceType (CEX path) no prefix",
			sourceType: SourceType(""),
			want:       quoteKey,
		},
		{
			name:       "SIMULATION_SOURCE prefixes with SIMULATION_SOURCE:",
			sourceType: SIMULATION_SOURCE,
			want:       string(SIMULATION_SOURCE) + ":" + quoteKey,
		},
		{
			name:       "DEX_SOURCE prefixes with DEX_SOURCE:",
			sourceType: DEX_SOURCE,
			want:       string(DEX_SOURCE) + ":" + quoteKey,
		},
		{
			name:       "unknown sourceType returns empty",
			sourceType: SourceType("UNKNOWN_SOURCE"),
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.GetOracleKey(tt.sourceType)
			if got != tt.want {
				t.Fatalf("GetOracleKey(%v) = %q, want %q", tt.sourceType, got, tt.want)
			}
		})
	}
}

func TestPair_GetOracleKey_WhenQuoteKeyEmpty_ReturnsEmptyForAllSourceTypes(t *testing.T) {
	// QuoteToken missing chain/address -> QuoteToken.GetOracleKey() == ""
	p := Pair{
		QuoteToken: Asset{
			Symbol:     "USDC",
			Blockchain: "",
			Address:    "",
		},
		BaseToken: Asset{Symbol: "WETH"},
	}

	sourceTypes := []SourceType{
		SourceType(""), // CEX path
		SIMULATION_SOURCE,
		DEX_SOURCE,
		SourceType("UNKNOWN"), // default
	}

	for _, st := range sourceTypes {
		t.Run(string(st), func(t *testing.T) {
			if got := p.GetOracleKey(st); got != "" {
				t.Fatalf("GetOracleKey(%v) = %q, want empty string", st, got)
			}
		})
	}
}
