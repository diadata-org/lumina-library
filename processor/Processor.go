package processor

import (
	"sync"
	"time"

	"github.com/diadata-org/lumina-library/filters"
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
	tradesblockChannel chan map[string]models.TradesBlock,
	filtersChannel chan []models.FilterPointPair,
	triggerChannel chan time.Time,
	failoverChannel chan string,
	metacontractClient *ethclient.Client,
	metacontractAddress string,
	metacontractPrecision int,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) {

	log.Info("Processor - Start......")

	// Collector starts collecting trades in the background and sends atomic tradesblocks to @tradesblockChannel.
	go scrapers.Collector(exchangePairs, pools, tradesblockChannel, triggerChannel, failoverChannel, branchMarketConfig, wg)

	// As soon as the trigger channel receives input a processing step is initiated.
	for tradesblocks := range tradesblockChannel {

		var filterPoints []models.FilterPointPair
		// Renew the price cache in each iteration. Could be refined by adjusting to the frequency of the source.
		priceCacheMap := make(map[string]float64)

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

			var atomicFilterValue float64
			var atomicVolume float64

			switch filterType {
			case string(FILTER_LAST_PRICE):
				atomicFilterValue, _, err = filters.LastPrice(tb, basePrice)
				if err != nil {
					log.Warn("last price filter: ", err)
					continue
				}

				log.Infof(
					"Processor - Atomic filter value for market %s with %v trades: %v.",
					tb.Trades[0].Exchange.Name+":"+tb.Trades[0].QuoteToken.Symbol+"-"+tb.Trades[0].BaseToken.Symbol,
					len(tb.Trades),
					atomicFilterValue,
				)
			case string(FILTER_VWAP):
				atomicFilterValue, atomicVolume, _, err = filters.VWAPFilter(tb, basePrice, toleranceSeconds)
				if err != nil {
					log.Warn("VWAP filter: ", err)
					continue
				}
				log.Infof(
					"Processor - VWAP filter value for pair %s-%s [%s] with %v trades: %v (volume: %v).",
					tb.Pair.QuoteToken.Symbol,
					tb.Pair.BaseToken.Symbol,
					tb.Trades[0].Exchange.Name,
					len(tb.Trades),
					atomicFilterValue,
					atomicVolume,
				)
			default:
				log.Warnf("Processor - unknown filterType %q for pair %s-%s, skipping.",
					filterType,
					tb.Pair.QuoteToken.Symbol,
					tb.Pair.BaseToken.Symbol,
				)
				continue
			}

			filterPoint := models.FilterPointPair{
				Pair:       tb.Pair,
				Value:      atomicFilterValue,
				Volume:     atomicVolume,
				Time:       tb.EndTime,
				SourceType: sourceType,
			}

			filterPoints = append(filterPoints, filterPoint)

		}

		var removedFilterPoints int
		filterPoints, removedFilterPoints = models.RemoveOldFilters(filterPoints, toleranceSeconds, time.Now())
		if removedFilterPoints > 0 {
			log.Warnf("Processor - Removed %v old filter points.", removedFilterPoints)
		}

		// --------------------------------------------------------------------------------------------
		// 2. Compute an aggregated value across exchanges for each asset obtained from the aggregated
		// filter values in Step 1.
		// --------------------------------------------------------------------------------------------

		// metafilter set by environment variable. For instance Median, Average, Minimum, etc.
		var filterPointsAggregated []models.FilterPointPair

		switch metaFilterType {
		case string(METAFILTER_MEDIAN):
			filterPointsAggregated = metafilters.Median(filterPoints)
			for _, fpm := range filterPointsAggregated {
				log.Infof("Processor - filter %s for %s: %v.", fpm.Name, fpm.Pair.QuoteToken.Symbol, fpm.Value)
			}
		case string(METAFILTER_VWAP):
			filterPointsAggregated = metafilters.VWAPMeta(filterPoints)
			for _, fpm := range filterPointsAggregated {
				log.Infof("Processor - meta VWAP for %s: %v.", fpm.Pair.QuoteToken.Symbol, fpm.Value)
			}
		default:
			log.Warnf("Processor - no metafilter matched for metaFilterType=%q, skipping update", metaFilterType)
		}

		filtersChannel <- filterPointsAggregated
	}

}