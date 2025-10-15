package scrapers

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

var (
	collectorOnce     sync.Once
	collectorUpdateCh chan []models.ExchangePair
)

func UpdateExchangePairs(newPairs []models.ExchangePair) {
	if collectorUpdateCh == nil {
		log.Warn("collectorUpdateCh not initialized, skipping update")
		return
	}
	select {
	case collectorUpdateCh <- newPairs:
		log.Infof("Collector - received hot-update request.")
	default:
		log.Warn("collectorUpdateCh is full, skipping update")
	}
}

func buildAllowedSet(m map[string][]models.ExchangePair) map[string]struct{} {
	allowed := make(map[string]struct{}, 64)
	for ex, list := range m {
		for _, ep := range list {
			id := ep.UnderlyingPair.ExchangePairIdentifier(ex)
			allowed[id] = struct{}{}
		}
	}
	return allowed
}

func pairEqual(a, b []models.ExchangePair) bool {
	if len(a) != len(b) {
		return false
	}
	key := func(ep models.ExchangePair) string {
		return ep.Exchange + ":" + ep.UnderlyingPair.QuoteToken.Symbol + "-" + ep.UnderlyingPair.BaseToken.Symbol
	}
	as := make([]string, 0, len(a))
	bs := make([]string, 0, len(b))
	for _, ep := range a {
		as = append(as, key(ep))
	}
	for _, ep := range b {
		bs = append(bs, key(ep))
	}
	sort.Strings(as)
	sort.Strings(bs)
	return reflect.DeepEqual(as, bs)
}

// Collector starts scrapers for all exchanges given by @exchangePairs.
func Collector(
	ctx context.Context,
	cancel context.CancelFunc,
	exchangePairs []models.ExchangePair,
	pools []models.Pool,
	tradesblockChannel chan map[string]models.TradesBlock,
	triggerChannel chan time.Time,
	failoverChannel chan string,
	wg *sync.WaitGroup,
) map[string]context.CancelFunc {
	collectorOnce.Do(func() {
		collectorUpdateCh = make(chan []models.ExchangePair, 1)
	})
	// exchangepairMap maps a centralized exchange onto the given pairs.
	exchangepairMap := models.MakeExchangepairMap(exchangePairs)
	log.Debugf("Collector - exchangepairMap: %v.", exchangepairMap)
	// poolMap maps a decentralized exchange onto the given pools.
	poolMap := models.MakePoolMap(pools)
	log.Debugf("Collector - poolMap: %v.", poolMap)

	cancelMap := make(map[string]context.CancelFunc)

	allowedSet := buildAllowedSet(exchangepairMap)

	// Start all needed scrapers.
	// @tradesChannelIn collects trades from the started scrapers.
	tradesChannelIn := make(chan models.Trade)
	for exchange := range exchangepairMap {
		// ctx, cancel := context.WithCancel(context.Background())
		cancelMap[exchange] = cancel
		wg.Add(1)
		go RunScraper(ctx, cancel, exchange, exchangepairMap[exchange], []models.Pool{}, tradesChannelIn, failoverChannel, wg)
	}
	for exchange := range poolMap {
		// ctx, cancel := context.WithCancel(context.Background())
		cancelMap[exchange] = cancel
		wg.Add(1)
		go RunScraper(ctx, cancel, exchange, []models.ExchangePair{}, poolMap[exchange], tradesChannelIn, failoverChannel, wg)
	}

	// tradesblockMap maps an exchangpair identifier onto a TradesBlock.
	// This also means that each value in the map consists of trades of only one exchangepair.
	// We call these blocks "atomic" tradesblocks.
	// TO DO: Make a dedicated type for atomic tradesblocks?
	tradesblockMap := make(map[string]models.TradesBlock)

	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, v := range cancelMap {
					v()
				}
				for k := range cancelMap {
					delete(cancelMap, k)
				}
				return
			case trade := <-tradesChannelIn:

				// Determine exchangepair and the corresponding identifier in order to assign the tradesBlockMap.
				exchangepair := models.Pair{QuoteToken: trade.QuoteToken, BaseToken: trade.BaseToken}
				exchangepairIdentifier := exchangepair.ExchangePairIdentifier(trade.Exchange.Name)

				if _, ok := allowedSet[exchangepairIdentifier]; !ok {
					log.Warnf("Collector - drop stale trade for %s", exchangepairIdentifier)
					continue
				}

				if _, ok := tradesblockMap[exchangepairIdentifier]; !ok {
					tradesblockMap[exchangepairIdentifier] = models.TradesBlock{
						Trades: []models.Trade{trade},
						Pair:   exchangepair,
					}
				} else {
					tradesblock := tradesblockMap[exchangepairIdentifier]
					tradesblock.Trades = append(tradesblock.Trades, trade)
					tradesblockMap[exchangepairIdentifier] = tradesblock
				}

			case timestamp := <-triggerChannel:

				log.Debugf("Collector - triggered at %v.", timestamp)
				for id := range tradesblockMap {
					tb := tradesblockMap[id]
					tb.Atomic = true
					tb.EndTime = timestamp
					tradesblockMap[id] = tb
				}

				tradesblockChannel <- tradesblockMap
				log.Infof("Collector - number of tradesblocks at %v: %v.", time.Now(), len(tradesblockMap))

				// Make a new tradesblockMap for the next trigger period.
				tradesblockMap = make(map[string]models.TradesBlock)

			case exchange := <-failoverChannel:
				log.Debugf("Collector - Restart scraper for %s.", exchange)
				if cancel, ok := cancelMap[exchange]; ok && cancel != nil {
					cancel()
					delete(cancelMap, exchange)
					time.Sleep(2 * time.Second)
				}
				// ctx, cancel := context.WithCancel(context.Background())
				cancelMap[exchange] = cancel
				wg.Add(1)
				go RunScraper(ctx, cancel, exchange, exchangepairMap[exchange], []models.Pool{}, tradesChannelIn, failoverChannel, wg)
			case newPairs := <-collectorUpdateCh:
				newMap := models.MakeExchangepairMap(newPairs)
				for ex := range exchangepairMap {
					if _, still := newMap[ex]; !still {
						log.Infof("Collector - stop scraper for removed exchange %s", ex)
						if cancel, ok := cancelMap[ex]; ok && cancel != nil {
							cancel()
							delete(cancelMap, ex)
							time.Sleep(2 * time.Second)
						}
					}
				}
				for ex, oldList := range exchangepairMap {
					if newList, ok := newMap[ex]; ok {
						if !pairEqual(oldList, newList) {
							log.Infof("Collector - reconfigure scraper for %s (pairs changed)", ex)
							if cancel, ok := cancelMap[ex]; ok && cancel != nil {
								cancel()
								delete(cancelMap, ex)
								time.Sleep(2 * time.Second)
							}
							// ctx, cancel := context.WithCancel(context.Background())
							cancelMap[ex] = cancel
							wg.Add(1)
							go RunScraper(ctx, cancel, ex, newList, []models.Pool{}, tradesChannelIn, failoverChannel, wg)
						}
					}
				}

				for ex, list := range newMap {
					if _, existed := exchangepairMap[ex]; !existed {
						log.Infof("Collector - start scraper for new exchange %s", ex)
						// ctx, cancel := context.WithCancel(context.Background())
						cancelMap[ex] = cancel
						wg.Add(1)
						go RunScraper(ctx, cancel, ex, list, []models.Pool{}, tradesChannelIn, failoverChannel, wg)
					}
				}
				exchangepairMap = newMap
				allowedSet = buildAllowedSet(exchangepairMap)
				log.Infof("Collector - exchangepairMap updated: %v", exchangepairMap)
			}
		}
	}()

	defer wg.Wait()
	return cancelMap
}
