package models

import (
	"testing"
)

func TestAsset_GetOracleKey(t *testing.T) {
	cases := []struct {
		name     string
		asset    Asset
		expected string
	}{
		{
			name:     "normal asset returns SYMBOL/USD",
			asset:    Asset{Symbol: "BTC", Blockchain: "Ethereum", Address: "0xabcdef"},
			expected: "BTC/USD",
		},
		{
			name:     "symbol is uppercased",
			asset:    Asset{Symbol: "btc", Blockchain: "ethereum", Address: "0xabcdef"},
			expected: "BTC/USD",
		},
		{
			name:     "whitespace is trimmed",
			asset:    Asset{Symbol: "  BTC  ", Blockchain: "  Ethereum  ", Address: "  0xabcdef  "},
			expected: "BTC/USD",
		},
		{
			name:     "empty symbol returns empty string",
			asset:    Asset{Symbol: "", Blockchain: "Ethereum", Address: "0xabcdef"},
			expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.asset.GetOracleKey()
			if got != c.expected {
				t.Errorf("GetOracleKey() = %q, want %q", got, c.expected)
			}
		})
	}
}

func TestGetOracleKey(t *testing.T) {
	btc := Asset{Symbol: "BTC", Blockchain: "Ethereum", Address: "0xabcdef"}
	usdt := Asset{Symbol: "USDT", Blockchain: "Ethereum", Address: "0x123456"}
	pair := Pair{QuoteToken: btc, BaseToken: usdt}

	cases := []struct {
		name       string
		sourceType SourceType
		expected   string
	}{
		{
			name:       "empty source type (CEX/DEX) returns SYMBOL/USD",
			sourceType: SourceType(""),
			expected:   "BTC/USD",
		},
		{
			name:       "SIMULATION source adds SIM: prefix",
			sourceType: SIMULATION_SOURCE,
			expected:   "SIM:BTC/USD",
		},
		{
			name:       "DEX source adds DEX: prefix",
			sourceType: DEX_SOURCE,
			expected:   "DEX:BTC/USD",
		},
		{
			name:       "unknown source type returns empty string",
			sourceType: SourceType("UNKNOWN"),
			expected:   "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GetOracleKey(c.sourceType, pair)
			if got != c.expected {
				t.Errorf("GetOracleKey(%q, pair) = %q, want %q", c.sourceType, got, c.expected)
			}
		})
	}
}

func TestGetSourceType(t *testing.T) {
	cases := []struct {
		name     string
		exchange Exchange
		expected SourceType
	}{
		{
			name:     "simulation exchange returns SIMULATION_SOURCE",
			exchange: Exchange{Simulation: true, Centralized: false},
			expected: SIMULATION_SOURCE,
		},
		{
			name:     "DEX (non-simulation, non-centralized) returns DEX_SOURCE",
			exchange: Exchange{Simulation: false, Centralized: false},
			expected: DEX_SOURCE,
		},
		{
			name:     "CEX returns empty SourceType",
			exchange: Exchange{Simulation: false, Centralized: true},
			expected: SourceType(""),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GetSourceType(c.exchange)
			if got != c.expected {
				t.Errorf("GetSourceType() = %q, want %q", got, c.expected)
			}
		})
	}
}