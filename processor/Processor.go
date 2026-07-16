package processor

import (
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/metafilters"
	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/scrapers"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Processor handles blocks from @tradesblockChannel.
// More precisley, it does so in a 2 step procedure:
// 1. Aggregate trades for each (atomic) block.
// 2. Aggregate filter values obtained in step 1.
func Processor(
	exchangePairs []models.ExchangePair,
	pools []models.Pool,
	customFeeds []models.CustomFeed,
	tradesblockChannel chan map[string]models.TradesBlock,
	filtersChannel chan []models.FilterPoint,
	triggerChannel chan time.Time,
	failoverChannel chan string,
	metacontractClient *ethclient.Client,
	metacontractAddress string,
	metacontractPrecision int,
	branchMarketConfig string,
	branchFeedConfig string,
	wg *sync.WaitGroup,
) {

	log.Info("Processor - Start......")

	// Collector starts collecting trades in the background and sends atomic tradesblocks to @tradesblockChannel.
	go scrapers.Collector(exchangePairs, pools, tradesblockChannel, triggerChannel, failoverChannel, branchMarketConfig, wg)

	// Periodically fetch customFeeds (managed via shared state).
	feedsState := newCustomFeedsState(customFeeds)
	stopReload := startCustomFeedsReloader(
		feedsState,
		time.Duration(watchFeedConfigSeconds)*time.Second,
		branchFeedConfig,
		exchangePairs,
	)
	defer stopReload()

	// As soon as the trigger channel receives input a processing step is initiated.
	for tradesblocks := range tradesblockChannel {

		currentCustomFeeds := feedsState.Snapshot()

		var atomicFilterPoints []models.FilterPointPair
		// Renew the price cache in each iteration. Could be refined by adjusting to the frequency of the source.
		// Always initialize with USD price 1.
		priceCacheMap := make(map[string]float64)
		priceCacheMap[usd.AssetIdentifier()] = 1

		// --------------------------------------------------------------------------------------------
		// 1. Compute an aggregated value for each pair on a given exchange using all collected trades.
		// --------------------------------------------------------------------------------------------
		for _, tb := range tradesblocks {

			// Get price of base asset from cache if possible.
			basePrice, err := models.GetPriceBaseAsset(tb, priceCacheMap, metacontractClient, metacontractAddress, metacontractPrecision)
			if err != nil {
				log.Errorf("Processor - GetPriceBaseAsset: %v", err)
				continue
			}

			// All blocks from Collector are atomic (single source). SourceType is always available.
			sourceType, err := tb.GetSourceType()
			if err != nil {
				log.Warnf("Processor - GetSourceType for pair %s-%s: %v, skipping.",
					tb.Pair.QuoteToken.Symbol,
					tb.Pair.BaseToken.Symbol,
					err,
				)
				continue
			}

			atomicFilterValue, atomicVolume, err := computeAtomicFilterValue(tb, filterTypeGlobal, basePrice, toleranceSeconds)
			if err != nil {
				log.Warn("Processor - computeAtomicFilterValue: ", err)
				continue
			}

			atomicFilterPoint := models.FilterPointPair{
				Pair:       tb.Pair,
				Value:      atomicFilterValue,
				Volume:     atomicVolume,
				Time:       tb.EndTime,
				SourceType: sourceType,
			}

			atomicFilterPoints = append(atomicFilterPoints, atomicFilterPoint)
		}

		var removedFilterPoints int
		atomicFilterPoints, removedFilterPoints = models.RemoveOldFilters(atomicFilterPoints, toleranceSeconds, time.Now())
		if removedFilterPoints > 0 {
			log.Warnf("Processor - Removed %v old filter points.", removedFilterPoints)
		}

		// --------------------------------------------------------------------------------------------
		// 2. Compute an aggregated value across exchanges for each asset obtained from the aggregated
		// filter values in Step 1.
		// --------------------------------------------------------------------------------------------

		// metafilter set by environment variable. For instance Median, Average, Minimum, etc.
		var filterPointsAggregated []models.FilterPoint
		filterAssetMap := models.GroupFiltersByAsset(atomicFilterPoints)

		switch metaFilterTypeGlobal {
		case models.METAFILTER_MEDIAN:
			filterPointsAggregated = metafilters.Median(filterAssetMap)
			for _, fpm := range filterPointsAggregated {
				log.Debugf("Processor - filter %s for %s: %v.", fpm.Type, fpm.Asset.Symbol, fpm.Value)
			}
		case models.METAFILTER_VWAP:
			filterPointsAggregated = metafilters.VWAPMeta(filterAssetMap)
			for _, fpm := range filterPointsAggregated {
				log.Debugf("Processor - meta VWAP for %s: %v.", fpm.Asset.Symbol, fpm.Value)
			}
		default:
			log.Warnf("Processor - no metafilter matched for metaFilterType=%q, skipping update", metaFilterTypeGlobal)
			continue
		}

		// --------------------------------------------------------------------------------------------
		// 3. Compute filter values for custom feeds in a similar 2-step procedure as above.
		// --------------------------------------------------------------------------------------------

		log.Debugf("number of customFeeds: %v", len(currentCustomFeeds))
		for _, customFeed := range currentCustomFeeds {

			atomicCustomFilterPoints := []models.FilterPoint{}
			// --------------------------------------------------------------------------------------------
			// a. As above, compute an aggregated value for each admissible pair on a given exchange.
			// --------------------------------------------------------------------------------------------
			for _, tb := range tradesblocks {
				if len(tb.Trades) == 0 {
					continue
				}
				if !customFeed.MatchingBlock(tb, false) {
					continue
				}

				basePrice, ok := priceCacheMap[tb.Trades[0].BaseToken.AssetIdentifier()]
				if !ok {
					log.Debugf("base price for %s -- %s:%s not available",
						tb.Trades[0].QuoteToken.Symbol,
						tb.Trades[0].QuoteToken.Blockchain,
						tb.Trades[0].QuoteToken.Address,
					)
					continue
				}

				atomicFilterValue, atomicVolume, err := computeAtomicFilterValue(tb, customFeed.Filter, basePrice, toleranceSeconds)
				if err != nil {
					log.Warn("Processor - computeAtomicFilterValue: ", err)
					continue
				}

				filterPoint := models.FilterPoint{
					Asset:  tb.Pair.QuoteToken,
					Value:  atomicFilterValue,
					Volume: atomicVolume,
					Time:   tb.EndTime,
				}

				atomicCustomFilterPoints = append(atomicCustomFilterPoints, filterPoint)
			}

			var removedFilterPoints int
			atomicCustomFilterPoints, removedFilterPoints = models.RemoveOldFiltersFp(atomicCustomFilterPoints, toleranceSeconds, time.Now())
			if removedFilterPoints > 0 {
				log.Warnf("Processor - Removed %v old filter points.", removedFilterPoints)
			}

			if len(atomicCustomFilterPoints) == 0 {
				continue
			}

			// --------------------------------------------------------------------------------------------
			// b. Compute an aggregated value across exchanges for each asset obtained from the aggregated
			// filter values in step a.
			// --------------------------------------------------------------------------------------------

			// metafilter set by custom feed selection.
			var filterPointAggregated models.FilterPoint
			switch customFeed.MetaFilter {
			case models.METAFILTER_MEDIAN:
				filterPointAggregated = metafilters.MedianFilters(atomicCustomFilterPoints)
				log.Infof("Processor - Median custom filter %s for %s: %v.", customFeed.Symbol, filterPointAggregated.Asset.Symbol, filterPointAggregated.Value)
			case models.METAFILTER_VWAP:
				filterPointAggregated = metafilters.VWAPFilters(atomicCustomFilterPoints)
				log.Infof("Processor - VWAP custom filter %s for %s: %v.", customFeed.Symbol, filterPointAggregated.Asset.Symbol, filterPointAggregated.Value)
			default:
				log.Warnf("Processor - no metafilter matched for metaFilterType=%q, skipping update", customFeed.MetaFilter)
				continue
			}
			filterPointAggregated.Name = customFeed.Symbol

			// This appends a final filter point for each custom feed.
			filterPointsAggregated = append(filterPointsAggregated, filterPointAggregated)
		}

		filtersChannel <- filterPointsAggregated
	}
}
