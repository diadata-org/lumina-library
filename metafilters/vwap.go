package metafilters

import (
	"sort"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_METAFILTERS", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)
}
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
			// This typically means the upstream filter type does not produce volume
			// (e.g. LastPrice). Operators should alarm on this log line.
			log.Warnf("VWAPMeta: zero volume for asset %s, falling back to equal-weight average", asset.Symbol)
			value = utils.Average(models.GetValuesFromFilterPoints(fps))
		} else {
			for _, fp := range fps {
				value += fp.Value * (fp.Volume / totalVolume)
			}
		}

		// Propagate SourceType from the first input filter point so that
		// downstream consumers (e.g. onchain/updater.go calling GetOracleKey)
		// can correctly apply source-specific key prefixes.
		// Note: SIMULATION_SOURCE data is handled by a separate pipeline and
		// is never expected to flow through VWAPMeta. This propagation is
		// defensive in case that assumption changes in the future.
		var sourceType models.SourceType
		if len(fps) > 0 {
			sourceType = fps[0].SourceType
		}

		result = append(result, models.FilterPointPair{
			Pair:       models.Pair{QuoteToken: asset},
			Value:      value,
			Name:       vwapMetaFilterName,
			Time:       models.GetLatestTimestampFromFilterPoints(fps),
			SourceType: sourceType,
		})
	}

	// Sort by QuoteToken address for deterministic output ordering.
	// VWAPMeta ranges over a map internally, so without this sort the result
	// slice order varies run-to-run, which affects oracle batch write order.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Pair.QuoteToken.Address < result[j].Pair.QuoteToken.Address
	})

	return
}