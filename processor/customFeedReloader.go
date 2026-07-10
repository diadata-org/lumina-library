package processor

import (
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
)

type customFeedsState struct {
	mu    sync.RWMutex
	feeds []models.CustomFeed
}

func newCustomFeedsState(initial []models.CustomFeed) *customFeedsState {
	cp := append([]models.CustomFeed(nil), initial...)
	return &customFeedsState{feeds: cp}
}

func (s *customFeedsState) Snapshot() []models.CustomFeed {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.CustomFeed(nil), s.feeds...)
}

func (s *customFeedsState) Replace(feeds []models.CustomFeed) {
	s.mu.Lock()
	s.feeds = append(s.feeds[:0], feeds...)
	s.mu.Unlock()
}

// startCustomFeedsReloader starts a goroutine and returns stop().
func startCustomFeedsReloader(
	state *customFeedsState,
	interval time.Duration,
	branchFeedConfig string,
	exchangePairs []models.ExchangePair,
) (stop func()) {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				newFeeds, err := models.FeedsFromConfigFile(branchFeedConfig)
				if err != nil {
					log.Warnf("Processor - FeedsFromConfigFile reload failed, keeping previous customFeeds: %v", err)
					continue
				}

				auxFeeds := make([]models.CustomFeed, 0, len(newFeeds))
				for _, cf := range newFeeds {
					if !cf.Admissible(exchangePairs) {
						log.Errorf("feed %v is not admissible.", cf.Symbol)
						continue
					}
					auxFeeds = append(auxFeeds, cf)
				}

				state.Replace(auxFeeds)
				log.Infof("Processor - customFeeds reloaded (%d entries).", len(auxFeeds))
			}
		}
	}()

	return func() { close(done) }
}
