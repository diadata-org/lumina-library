package models

import (
	"encoding/json"

	"github.com/diadata-org/lumina-library/utils"
)

var (
	CEX_SOURCE        = SourceType("CEX")
	DEX_SOURCE        = SourceType("DEX")
	SIMULATION_SOURCE = SourceType("SIM")
)

// @SourceType describes how the source(s) constituting the filter point.
// CEX for a filter point computed from trades originating from CEXes.
// DEX for a filter point computed from trades originating from DEXes.
// SIM for a filter point computed from trades originating from simulation of DEX pools.
type SourceType string

type Exchange struct {
	Name        string `json:"Name"`
	Centralized bool   `json:"Centralized"`
	Simulation  bool   `json:"Simulation"`
	Bridge      bool   `json:"Bridge"`
	Contract    string `json:"Contract"`
	Blockchain  string `json:"Blockchain"`
}

type IdentifiedAsset struct {
	Exchange   string
	Symbol     string
	Blockchain string
	Address    string
	Decimals   uint8
}

type IdentifiedAssets struct {
	Tokens []IdentifiedAsset
}

// GetSymbolIdentificationMap returns a map which maps an asset's symbol ticker from @exchange onto the underlying asset.
// It assumes symbol mappings can be found in the file exchange.json at @configPath.
// e.g. identificationMap["ETH_GateIO"] = Asset{Symbol: "ETH", Blockchain: "ETH", Address: "0x0000000000000000000000000000000000000000", Decimals: 18}
func GetSymbolIdentificationMap(exchange string, branch string) (map[string]Asset, error) {
	jsonFile, err := utils.GetConfig("symbolIdentification", exchange, branch)
	if err != nil {
		return nil, err
	}

	var ia IdentifiedAssets
	if err := json.Unmarshal(jsonFile, &ia); err != nil {
		return nil, err
	}

	identificationMap := make(map[string]Asset)

	for _, t := range ia.Tokens {
		identificationMap[ExchangeSymbolIdentifier(t.Symbol, t.Exchange)] = Asset{
			Symbol:     t.Symbol,
			Blockchain: t.Blockchain,
			Address:    t.Address,
			Decimals:   t.Decimals,
		}
	}

	return identificationMap, nil
}

// GetSourceType returns the type of @exchange such as DEX. For CEX an empty SourceType is returned.
func GetSourceType(exchange Exchange) SourceType {
	if exchange.Simulation {
		return SIMULATION_SOURCE
	}
	if !exchange.Simulation && !exchange.Centralized {
		return DEX_SOURCE
	}
	return SourceType("")
}

// GetOracleKey returns a key for an asset (or a pair) that can be used for calling the
// corresponding value in an oracle.
func GetOracleKey(sourceType SourceType, pair Pair) string {
	switch sourceType {
	case SourceType(""):
		return pair.QuoteToken.Symbol + "/USD"
	case SIMULATION_SOURCE:
		return string(SIMULATION_SOURCE) + ":" + pair.QuoteToken.Symbol + "/USD"
	case DEX_SOURCE:
		return string(DEX_SOURCE) + ":" + pair.QuoteToken.Symbol + "/USD"
	default:
		return ""
	}
}

// GetOracleKeySimulation returns a key for an asset (or a pool) that can be used for calling the
// corresponding value in an oracle. It is restricted to values originating from simulated DEX pools.
func GetOracleKeySimulation(pair Pair) string {

	return pair.QuoteToken.Symbol + "/USD"

}
