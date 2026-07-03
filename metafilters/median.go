package metafilters

import (
	models "github.com/diadata-org/lumina-library/models"
	utils "github.com/diadata-org/lumina-library/utils"
)

const (
	medianFilterName = "median"
)

// Median returns the median value for all filter points that share the same quote asset.
// The input @filterPoints still consists of "atomic" filter points.
func MedianFilters(assetKey models.AssetKey, filters []models.FilterPointPair) (fp models.FilterPointPair) {

	filterValue := utils.Median(models.GetValuesFromFilterPoints(filters))

	fp.Value = filterValue
	fp.Pair.QuoteToken = assetKey.Key2Asset()
	fp.Name = medianFilterName
	fp.Time = models.GetLatestTimestampFromFilterPoints(filters)

	return
}

// Median returns the median value for all filter points that share the same quote asset.
// The input @filterPoints still consists of "atomic" filter points.
func Median(filterAssetMap map[models.AssetKey][]models.FilterPointPair) (medianizedFilterPoints []models.FilterPointPair) {

	for assetKey, filters := range filterAssetMap {
		filterValue := utils.Median(models.GetValuesFromFilterPoints(filters))
		var fp models.FilterPointPair
		fp.Value = filterValue
		fp.Pair.QuoteToken = assetKey.Key2Asset()
		fp.Name = medianFilterName
		fp.Time = models.GetLatestTimestampFromFilterPoints(filters)

		medianizedFilterPoints = append(medianizedFilterPoints, fp)
	}

	return
}
