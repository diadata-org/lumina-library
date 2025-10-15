package scrapers

import (
	"context"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
)

// watchdog checks for liveliness of a pair subscription.
// More precisely, if there is no trades for a period longer than @watchdogDelayMap[pair.ForeignName],
// the @runChannel receives the corresponding pair. The calling function can decide what to do, for
// instance resubscribe to the pair.
func watchdog(
	ctx context.Context,
	pair models.ExchangePair,
	ticker *time.Ticker,
	lastTradeTimeMap map[string]time.Time,
	watchdogDelay int64,
	subscribeChannel chan models.ExchangePair,
	lock *sync.RWMutex,
) {
	log.Infof("%s - start watching %s with watchdog %v.", pair.Exchange, pair.ForeignName, watchdogDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debugf("%s - check liveliness of %s.", pair.Exchange, pair.ForeignName)

			// Make read lock for lastTradeTimeMap.
			lock.RLock()
			lastTradeTime, ok := lastTradeTimeMap[pair.ForeignName]
			lock.RUnlock()
			if !ok || lastTradeTime.IsZero() {
				log.Errorf("%s - lastTradeTimeMap[%s] not found.", pair.Exchange, pair.ForeignName)
				continue
			}
			duration := time.Since(lastTradeTime)
			log.Debugf("%s - duration for %s: %v. Threshold: %v.", pair.Exchange, pair.ForeignName, duration, watchdogDelay)
			if duration > time.Duration(watchdogDelay)*time.Second {
				log.Errorf("%s - watchdogTicker failover for %s.", pair.Exchange, pair.ForeignName)
				select {
				case <-ctx.Done():
					log.Debugf("%s - watchdog exit (context cancelled) for pair %s.", pair.Exchange, pair.ForeignName)
					return
				case subscribeChannel <- pair:
				}
				// subscribeChannel <- pair
			}
		case <-ctx.Done():
			log.Debugf("%s - close watchdog for pair %s.", pair.Exchange, pair.ForeignName)
			return
		}
	}
}

// watchdog checks for liveliness of a pair subscription.
// More precisely, if there is no trades for a period longer than @watchdogDelayMap[pair.ForeignName],
// the @runChannel receives the corresponding pair. The calling function can decide what to do, for
// instance resubscribe to the pair.
func watchdogPool(
	ctx context.Context,
	exchange string,
	pool common.Address,
	ticker *time.Ticker,
	lastTradeTimeMap map[common.Address]time.Time,
	watchdogDelay int64,
	subscribeChannel chan common.Address,
	lock *sync.RWMutex,
) {
	log.Infof("%s - start watching %s with watchdog %v.", exchange, pool.Hex(), watchdogDelay)
	for {
		select {
		case <-ticker.C:
			log.Debugf("%s - check liveliness of %s.", exchange, pool.Hex())

			// Make read lock for lastTradeTimeMap.
			lock.RLock()
			duration := time.Since(lastTradeTimeMap[pool])
			log.Debugf("%s - duration for %s: %v. Threshold: %v.", exchange, pool.Hex(), duration, watchdogDelay)
			lock.RUnlock()
			if duration > time.Duration(watchdogDelay)*time.Second {
				log.Warnf("%s - watchdogTicker failover for %s.", exchange, pool.Hex())
				subscribeChannel <- pool
			}
		case <-ctx.Done():
			log.Debugf("%s - close watchdog for pair %s.", exchange, pool.Hex())
			return
		}
	}
}
