package models

import (
	"math"
	"time"
)

const (
	FILTER_LAST_PRICE = FilterType("LastPrice")
	FILTER_VWAP       = FilterType("VWAP")
	METAFILTER_MEDIAN = MetafilterType("Median")
	METAFILTER_VWAP   = MetafilterType("VWAP")
)

type FilterType string

type MetafilterType string

// FilterPoint contains the resulting value of a filter applied to an asset.
type FilterPoint struct {
	Asset      Asset
	Value      float64
	Name       string
	Time       time.Time
	Source     Exchange
	SourceType SourceType
}

type FilterPointPair struct {
	Pool  Pool
	Pair  Pair
	Value float64
	// Volume holds the total absolute trade volume used to compute Value.
	// Populated by the VWAP filter; used by the VWAP metafilter to weight
	// per-source filter points in the cross-source aggregation.
	// Zero for filter types that do not track volume (e.g. LastPrice).
	Volume     float64
	Name       string
	Time       time.Time
	Source     Exchange
	SourceType SourceType
}

// GroupFilterByAsset returns @fpMap which maps an assetKey on all filter points contained in @filterPoints.
func GroupFiltersByAsset(filterPoints []FilterPointPair) (fpMap map[AssetKey][]FilterPointPair) {
	fpMap = make(map[AssetKey][]FilterPointPair)
	for _, fp := range filterPoints {
		fpMap[fp.Pair.QuoteToken.Asset2Key()] = append(fpMap[fp.Pair.QuoteToken.Asset2Key()], fp)
	}
	return
}

// GetValuesFromFilterPoints returns a slice containing just the values from @filterPoints.
func GetValuesFromFilterPoints(filterPoints []FilterPointPair) (filterValues []float64) {
	for _, fp := range filterPoints {
		filterValues = append(filterValues, fp.Value)
	}
	return
}

// GetLatestTimestampFromFilterPoints returns the latest timstamp among all @filterPoints.
func GetLatestTimestampFromFilterPoints(filterPoints []FilterPointPair) (timestamp time.Time) {
	for _, fp := range filterPoints {
		if fp.Time.After(timestamp) {
			timestamp = fp.Time
		}
	}
	return
}

// RemoveOldFilters removes all filter points from @filterPoints whith timestamp more than
// @toleranceSeconds before @timestamp.
func RemoveOldFilters(filterPoints []FilterPointPair, toleranceSeconds int64, timestamp time.Time) (cleanedFilterPoints []FilterPointPair, removedFilters int) {
	for _, fp := range filterPoints {
		if fp.Time.After(timestamp.Add(-time.Duration(toleranceSeconds) * time.Second)) {
			cleanedFilterPoints = append(cleanedFilterPoints, fp)
		} else {
			removedFilters++
		}
	}
	return
}

func RemoveLargeDeviationPrices(filterPoints []FilterPointPair) (newFilterPoints []FilterPointPair) {
	groups := make(map[string][]FilterPointPair)

	for _, fp := range filterPoints {
		symbol := fp.Pair.QuoteToken.Symbol
		groups[symbol] = append(groups[symbol], fp)
	}

	for _, group := range groups {
		count := len(group)

		if count%2 != 0 {
			newFilterPoints = append(newFilterPoints, group...)
			continue
		}
		valid := true
		for i := 0; i < len(group); i++ {
			for j := i + 1; j < len(group); j++ {
				a, b := group[i].Value, group[j].Value
				base := math.Min(a, b)
				if base == 0 {
					valid = false
					break // Break: Invalid base price
				}
				diff := math.Abs(a-b) / base
				if diff > 0.01 {
					valid = false
					log.Warnf("Price difference %.2f%% (%.4f vs %.4f) exceeds 1%% threshold for QuoteToken %s",
						diff*100,                        // Convert to percentage (0.01 -> 1.00%)
						a,                               // First value (float64)
						b,                               // Second value (float64)
						group[i].Pair.QuoteToken.Symbol, // Token symbol (string)
					)
					break
				}
			}
			if !valid {
				log.Warnf("Price validation failed for QuoteToken %s\n", group[i].Pair.QuoteToken.Symbol)
				break
			}
		}
		if valid {
			newFilterPoints = append(newFilterPoints, group...)
		}
	}
	log.Infof("Length of old filterPoints: %v, Length of new filterPoints: %v\n", len(filterPoints), len(newFilterPoints))
	return
}
