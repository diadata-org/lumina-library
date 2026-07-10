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
			expected: "btc/USD",
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

	cases := []struct {
		name        string
		filterPoint FilterPoint
		expected    string
	}{
		{
			name:        "regular BTC feed",
			filterPoint: FilterPoint{Name: "", Asset: Asset{Symbol: "BTC"}},
			expected:    "BTC/USD",
		},
		{
			name:        "custom feed with given asset",
			filterPoint: FilterPoint{Name: "customFeed", Asset: Asset{Symbol: "BTC"}},
			expected:    "customFeed",
		},
		{
			name:        "custom feed without asset",
			filterPoint: FilterPoint{Name: "customFeed"},
			expected:    "customFeed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GetOracleKey(c.filterPoint)
			if got != c.expected {
				t.Errorf("GetOracleKey(%v) = %q, want %q", c.filterPoint, got, c.expected)
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
