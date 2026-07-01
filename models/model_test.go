package models

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// mustChecksumAddress returns the EIP-55 checksummed form of an EVM address.
// Used in tests to avoid hardcoding the checksum output.
func mustChecksumAddress(address string) string {
	return common.HexToAddress(address).Hex()
}

func TestAsset_GetOracleKey(t *testing.T) {
	cases := []struct {
		name     string
		asset    Asset
		expected string
	}{
		{
			name: "normal asset returns SYMBOL/USD:BLOCKCHAIN/address",
			asset: Asset{
				Symbol:     "BTC",
				Blockchain: "Ethereum",
				Address:    "0xabcdef1234567890",
			},
			expected: "BTC/USD:Ethereum/0xabcdef1234567890",
		},
		{
			name: "EVM address is EIP-55 checksummed, non-EVM left as-is",
			asset: Asset{
				Symbol:     "btc",
				Blockchain: "ethereum",
				Address:    "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expected: "BTC/USD:Ethereum/" + mustChecksumAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
		},
		{
			name: "whitespace is trimmed",
			asset: Asset{
				Symbol:     "  BTC  ",
				Blockchain: "  Ethereum  ",
				Address:    "  0xabcdef  ",
			},
			expected: "BTC/USD:Ethereum/0xabcdef",
		},
		{
			name:     "empty blockchain returns empty string",
			asset:    Asset{Symbol: "BTC", Blockchain: "", Address: "0xabcdef"},
			expected: "",
		},
		{
			name:     "empty address returns empty string",
			asset:    Asset{Symbol: "BTC", Blockchain: "Ethereum", Address: ""},
			expected: "",
		},
		{
			name:     "both empty returns empty string",
			asset:    Asset{Symbol: "BTC", Blockchain: "", Address: ""},
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

func TestPair_GetOracleKey(t *testing.T) {
	btc := Asset{
		Symbol:     "BTC",
		Blockchain: "Ethereum",
		Address:    "0xabcdef",
	}
	usdt := Asset{
		Symbol:     "USDT",
		Blockchain: "Ethereum",
		Address:    "0x123456",
	}
	pair := Pair{QuoteToken: btc, BaseToken: usdt}

	// Expected base key from Asset.GetOracleKey()
	baseKey := "BTC/USD:Ethereum/0xabcdef"

	cases := []struct {
		name       string
		sourceType SourceType
		expected   string
	}{
		{
			name:       "CEX source (empty string) returns key as-is",
			sourceType: SourceType(""),
			expected:   baseKey,
		},
		{
			name:       "DEX source (empty string) returns key as-is",
			sourceType: SourceType(""),
			expected:   baseKey,
		},
		{
			name:       "SIMULATION source adds SIM: prefix",
			sourceType: SIMULATION_SOURCE,
			expected:   "SIM:" + baseKey,
		},
		{
			name:       "unknown source type returns key unchanged",
			sourceType: SourceType("UNKNOWN"),
			expected:   baseKey,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pair.GetOracleKey(c.sourceType)
			if got != c.expected {
				t.Errorf("GetOracleKey(%q) = %q, want %q", c.sourceType, got, c.expected)
			}
		})
	}
}

func TestPair_GetOracleKey_EmptyAsset(t *testing.T) {
	// If the quote token has no blockchain/address, key should be empty
	// regardless of sourceType.
	emptyAsset := Asset{Symbol: "BTC", Blockchain: "", Address: ""}
	usdt := Asset{Symbol: "USDT", Blockchain: "Ethereum", Address: "0x123456"}
	pair := Pair{QuoteToken: emptyAsset, BaseToken: usdt}

	for _, sourceType := range []SourceType{SourceType(""), SIMULATION_SOURCE} {
		got := pair.GetOracleKey(sourceType)
		if got != "" {
			t.Errorf("GetOracleKey(%q) = %q, want empty string for asset with no blockchain/address", sourceType, got)
		}
	}
}