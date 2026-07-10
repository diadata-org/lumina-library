package models

import (
	"encoding/json"
	"fmt"

	"github.com/diadata-org/lumina-library/utils"
)

type CustomFeed struct {
	Symbol     string         `json:"Symbol"`
	Asset      Asset          `json:"Asset,omitempty"`
	Filter     FilterType     `json:"Filter"`
	MetaFilter MetafilterType `json:"Metafilter"`
	Markets    []ExchangePair `json:"Markets"`
	Pools      []Pool         `json:"Pools"`
}

// @MatchingBlock returns true whenever asset, exchange and pair ticker are matching
// one of the markets in the custom feed.
func (cf *CustomFeed) MatchingBlock(tb TradesBlock, matchAsset bool) bool {

	if matchAsset {
		if tb.Pair.QuoteToken.MinimalAsset() != cf.Asset.MinimalAsset() {
			return false
		}
	}
	for _, market := range cf.Markets {
		if market.Exchange == tb.Exchange().Name && market.ForeignName == tb.ForeignName() {
			return true
		}
	}
	for _, pool := range cf.Pools {
		if pool.Blockchain.Name == tb.Pool.Blockchain.Name && pool.Address == tb.Pool.Address {
			return true
		}
	}

	return false
}

// @Admissible checks whether all markets within the custom feed are available.
func (cf *CustomFeed) Admissible(eps []ExchangePair) bool {
	for _, market := range cf.Markets {
		var marketExists bool
		for _, ep := range eps {
			if market.Exchange == ep.Exchange && market.ForeignName == ep.ForeignName {
				marketExists = true
				continue
			}
		}
		if !marketExists {
			return marketExists
		}
	}
	return true
}

func FeedsFromConfigFile(branchMarketConfig string) ([]CustomFeed, error) {

	jsonFile, err := utils.GetConfig("feeds", "feeds", branchMarketConfig)
	if err != nil {
		return nil, fmt.Errorf("GetConfig(feeds): %v", err)
	}

	var cfg []CustomFeed
	if err := json.Unmarshal(jsonFile, &cfg); err != nil {
		return nil, err
	}

	// TO DO: Assign assets to pair tokens using existing ExchangePairs and Pools?

	return cfg, nil

}
