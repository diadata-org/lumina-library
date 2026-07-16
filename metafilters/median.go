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
func MedianFilters(filters []models.FilterPoint) (fp models.FilterPoint) {

	filterValue := utils.Median(models.GetValuesFromFilterPoints(filters))

	fp.Value = filterValue
	fp.Type = medianFilterName
	fp.Time = models.GetLatestTimestampFromFilterPoints(filters)

	return
}

// Median returns the median value for all filter points that share the same quote asset.
// The input @filterPoints still consists of "atomic" filter points.
func Median(filterAssetMap map[models.AssetKey][]models.FilterPoint) (medianizedFilterPoints []models.FilterPoint) {

	for assetKey, filters := range filterAssetMap {
		filterValue := utils.Median(models.GetValuesFromFilterPoints(filters))
		var fp models.FilterPoint
		fp.Value = filterValue
		fp.Asset = assetKey.Key2Asset()
		fp.Type = medianFilterName
		fp.Time = models.GetLatestTimestampFromFilterPoints(filters)

		medianizedFilterPoints = append(medianizedFilterPoints, fp)
	}

	return
}
