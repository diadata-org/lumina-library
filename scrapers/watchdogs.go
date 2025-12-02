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
	for {
		select {
		case <-ticker.C:
			log.Debugf("%s - check liveliness of %s.", pair.Exchange, pair.ForeignName)

			// Make read lock for lastTradeTimeMap.
			lock.RLock()
			duration := time.Since(lastTradeTimeMap[pair.ForeignName])
			log.Debugf("%s - duration for %s: %v. Threshold: %v.", pair.Exchange, pair.ForeignName, duration, watchdogDelay)
			lock.RUnlock()
			if duration > time.Duration(watchdogDelay)*time.Second {
				log.Errorf("%s - watchdogTicker failover for %s.", pair.Exchange, pair.ForeignName)
				subscribeChannel <- pair
			}
		case <-ctx.Done():
			log.Debugf("%s - close watchdog for pair %s.", pair.Exchange, pair.ForeignName)
			return
		}
	}
}

// generic: start watchdog
func StartWatchdogForPair(
	ctx context.Context,
	lock *sync.RWMutex,
	pair models.ExchangePair,
	watchdogCancel map[string]context.CancelFunc,
	lastTradeTimeMap map[string]time.Time,
	subscribeCh chan models.ExchangePair,
) {
	lock.Lock()
	if cancel, exists := watchdogCancel[pair.ForeignName]; exists && cancel != nil {
		lock.Unlock()
		return
	}

	wdCtx, cancel := context.WithCancel(ctx)
	watchdogCancel[pair.ForeignName] = cancel
	lock.Unlock()

	ticker := time.NewTicker(time.Duration(pair.WatchDogDelay) * time.Second)
	go watchdog(wdCtx, pair, ticker, lastTradeTimeMap, pair.WatchDogDelay, subscribeCh, lock)
}

func StopWatchdogForPair(
	lock *sync.RWMutex,
	foreignName string,
	watchdogCancel map[string]context.CancelFunc,
) {
	lock.Lock()
	if cancel, ok := watchdogCancel[foreignName]; ok && cancel != nil {
		cancel()
		delete(watchdogCancel, foreignName)
	}
	lock.Unlock()
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
