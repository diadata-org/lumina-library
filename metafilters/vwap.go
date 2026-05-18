package metafilters

import (
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

const vwapMetaFilterName = "vwap"

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
func VWAPMeta(filterPoints []models.FilterPointPair) (result []models.FilterPointPair) {
	filterAssetMap := models.GroupFiltersByAsset(filterPoints)

	for asset, fps := range filterAssetMap {
		var totalVolume float64
		for _, fp := range fps {
			totalVolume += fp.Volume
		}

		var value float64
		if totalVolume == 0 {
			// Fallback: equal-weight average when no volume information is available.
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
	return
}