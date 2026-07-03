package metafilters

import (
	"sort"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

const vwapMetaFilterName = "vwap"

func VWAPFilters(assetKey models.AssetKey, fps []models.FilterPointPair) (fp models.FilterPointPair) {

	var totalVolume float64
	for _, fp := range fps {
		totalVolume += fp.Volume
	}
	asset := assetKey.Key2Asset()

	var value float64
	if totalVolume == 0 {
		// Fallback: equal-weight average when no volume information is available.
		// This typically means the upstream filter type does not produce volume
		// (e.g. LastPrice). Operators should alarm on this log line.
		log.Warnf("VWAPMeta: zero volume for asset %s, falling back to equal-weight average", asset.Symbol)
		value = utils.Average(models.GetValuesFromFilterPoints(fps))
	} else {
		for _, fp := range fps {
			value += fp.Value * (fp.Volume / totalVolume)
		}
	}

	fp = models.FilterPointPair{
		Pair:  models.Pair{QuoteToken: asset},
		Value: value,
		Name:  vwapMetaFilterName,
		Time:  models.GetLatestTimestampFromFilterPoints(fps),
	}

	return
}

// VWAPMeta aggregates per-source VWAP filter points into a single cross-source
// VWAP value for each quote asset.
//
// Each input filter point carries a Value (USD price) and a Volume (total
// absolute trade volume) produced by the per-source VWAPFilter. The cross-source
// price is then:
//
//	VWAP_cross = Σ(Value_i × Volume_i) / Σ(Volume_i)
//
// If all input filter points for a given asset have zero volume (e.g. because
// the filter type was not VWAP), the function falls back to a simple equal-weight
// average so the metafilter degrades gracefully rather than returning zero.
func VWAPMeta(filterAssetMap map[models.AssetKey][]models.FilterPointPair) (result []models.FilterPointPair) {

	for assetKey, fps := range filterAssetMap {
		var totalVolume float64
		for _, fp := range fps {
			totalVolume += fp.Volume
		}
		asset := assetKey.Key2Asset()

		var value float64
		if totalVolume == 0 {
			// Fallback: equal-weight average when no volume information is available.
			// This typically means the upstream filter type does not produce volume
			// (e.g. LastPrice). Operators should alarm on this log line.
			log.Warnf("VWAPMeta: zero volume for asset %s, falling back to equal-weight average", asset.Symbol)
			value = utils.Average(models.GetValuesFromFilterPoints(fps))
		} else {
			for _, fp := range fps {
				value += fp.Value * (fp.Volume / totalVolume)
			}
		}

		result = append(result, models.FilterPointPair{
			Pair:  models.Pair{QuoteToken: asset},
			Value: value,
			Name:  vwapMetaFilterName,
			Time:  models.GetLatestTimestampFromFilterPoints(fps),
		})
	}

	// Sort by (Blockchain, Address) for deterministic output ordering.
	// VWAPMeta ranges over a map internally, so without this sort the result
	// slice order varies run-to-run, which affects oracle batch write order.
	// Blockchain is the primary key to handle assets that share the same address
	// across chains (e.g. 0x0000... used as a native-asset sentinel).
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i].Pair.QuoteToken, result[j].Pair.QuoteToken
		if a.Blockchain != b.Blockchain {
			return a.Blockchain < b.Blockchain
		}
		return a.Address < b.Address
	})

	return
}
