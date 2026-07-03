package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/diadata-org/lumina-library/utils"
)

type Feed struct {
	Asset      Asset
	Filter     FilterType
	MetaFilter FilterType
}

func FeedsFromConfigFile(branchMarketConfig string) (map[AssetKey]Feed, error) {
	type feed struct {
		Symbol     string `json:"Symbol"`
		Address    string `json:"Address"`
		Blockchain string `json:"Blockchain"`
		Filter     string `json:"Filter"`
		MetaFilter string `json:"MetaFilter"`
	}
	type fileSchema struct {
		Feeds []feed `json:"Feeds"`
	}

	jsonFile, err := utils.GetConfig("feeds", "feeds", branchMarketConfig)
	if err != nil {
		return nil, fmt.Errorf("GetConfig(feeds): %v", err)
	}

	var cfg fileSchema
	if err := json.Unmarshal(jsonFile, &cfg); err != nil {
		return nil, err
	}

	// construct the output
	feedMap := make(map[AssetKey]Feed)
	for i, feed := range cfg.Feeds {
		if strings.TrimSpace(feed.Symbol) == "" {
			return feedMap, fmt.Errorf("Feeds[%d] has empty symbol", i)
		}
		if strings.TrimSpace(feed.Address) == "" {
			return feedMap, fmt.Errorf("Feeds[%d] has empty address", i)
		}
		if strings.TrimSpace(feed.Blockchain) == "" {
			return feedMap, fmt.Errorf("Feeds[%d] has empty blockchain", i)
		}

		assetKey := AssetKey{Symbol: feed.Symbol, Address: feed.Address, Blockchain: feed.Blockchain}
		feedMap[assetKey] = Feed{Asset: assetKey.Key2Asset(), Filter: FilterType(feed.Filter), MetaFilter: FilterType(feed.MetaFilter)}

	}
	return feedMap, nil
}
