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
		BaseToken: Asset{Symbol: "weth"}, // irrelevant for GetOracleKey()
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
			name:       "empty sourceType no prefix",
			sourceType: SourceType(""),
			want:       quoteKey,
		},
		{
			name:       "CEX_SOURCE prefixes with CEX_SOURCE:",
			sourceType: CEX_SOURCE,
			want:       string(CEX_SOURCE) + ":" + quoteKey,
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

func TestPair_GetOracleKey_UnknownSourceType_ReturnsEmpty(t *testing.T) {
	p := Pair{
		QuoteToken: Asset{
			Symbol:     "USDC",
			Blockchain: "ETH",
			Address:    "0xabc",
		},
	}

	got := p.GetOracleKey(SourceType("UNKNOWN_SOURCE"))
	if got != "" {
		t.Fatalf("GetOracleKey(UNKNOWN_SOURCE) = %q, want empty string", got)
	}
}

func TestPair_GetOracleKey_WhenQuoteKeyEmpty_ReturnsEmptyOrPrefixedEmpty(t *testing.T) {
	// QuoteToken missing chain/address -> QuoteToken.GetOracleKey() == ""
	p := Pair{
		QuoteToken: Asset{
			Symbol:     "USDC",
			Blockchain: "",
			Address:    "",
		},
	}

	// case SourceType("") returns empty
	if got := p.GetOracleKey(SourceType("")); got != "" {
		t.Fatalf("GetOracleKey(empty) = %q, want empty string", got)
	}

	// Prefixed cases will return "<SRC>:" because QuoteToken.GetOracleKey() is "".
	// This test documents current behavior (might be something you later change).
	if got := p.GetOracleKey(DEX_SOURCE); got != string(DEX_SOURCE)+":" {
		t.Fatalf("GetOracleKey(DEX_SOURCE) = %q, want %q", got, string(DEX_SOURCE)+":")
	}
}
