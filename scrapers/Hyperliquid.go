package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

const (
	hyperliquidWSURL   = "wss://api.hyperliquid.xyz/ws"
	hyperliquidInfoURL = "https://api.hyperliquid.xyz/info"
	// Hyperliquid closes connections that have not sent a message for 60s.
	hyperliquidPingPeriodDefault = 30
	// How long a spotMeta snapshot, which maps pair names onto subscription coins, is reused.
	// A pair that is not in the snapshot fails without refetching until the snapshot expires,
	// because the watchdog keeps retrying such pairs for as long as the process runs. Newly
	// listed pairs are picked up within this window.
	hyperliquidSpotMetaTTL = 10 * time.Minute
)

var hyperliquidHTTPClient = &http.Client{Timeout: 10 * time.Second}

// ---------------- wire types ----------------

type hyperliquidSubscription struct {
	Type string `json:"type"`
	Coin string `json:"coin"`
}

type hyperliquidWSRequest struct {
	Method       string                   `json:"method"`
	Subscription *hyperliquidSubscription `json:"subscription,omitempty"`
}

type hyperliquidWSMessage struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type hyperliquidTrade struct {
	Coin string `json:"coin"`
	Side string `json:"side"` // "B" | "A"
	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Time int64  `json:"time"`
	Tid  int64  `json:"tid"`
}

type hyperliquidSpotMeta struct {
	Tokens []struct {
		Name  string `json:"name"`
		Index int    `json:"index"`
	} `json:"tokens"`
	Universe []struct {
		Name   string `json:"name"`
		Tokens []int  `json:"tokens"`
	} `json:"universe"`
}

// ---------------- hooks ----------------

type hyperliquidHooks struct {
	infoURL       string
	mu            sync.RWMutex
	coins         map[string]string // "UBTC/USDC" -> "@142"
	fetchedAt     time.Time
	foreignByCoin map[string]string // "@142" -> "BTC-USDC"
}

func (h *hyperliquidHooks) ExchangeKey() string { return HYPERLIQUID_EXCHANGE }
func (h *hyperliquidHooks) WSURL() string       { return hyperliquidWSURL }

func (h *hyperliquidHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	pingPeriod, err := strconv.Atoi(utils.Getenv("HYPERLIQUID_PING_PERIOD_SECONDS", strconv.Itoa(hyperliquidPingPeriodDefault)))
	if err != nil || pingPeriod <= 0 {
		log.Errorf("HYPERLIQUID - parse HYPERLIQUID_PING_PERIOD_SECONDS: %v. Set to default %d.", err, hyperliquidPingPeriodDefault)
		pingPeriod = hyperliquidPingPeriodDefault
	}
	go func() {
		tick := time.NewTicker(time.Duration(pingPeriod) * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := bs.SafeWriteJSON(hyperliquidWSRequest{Method: "ping"}); err != nil {
					log.Errorf("HYPERLIQUID - send ping: %v.", err)
				}
			}
		}
	}()
}

func (h *hyperliquidHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	coin, err := h.coin(pair.ForeignName)
	if err != nil {
		return err
	}
	method := "unsubscribe"
	if subscribe {
		method = "subscribe"
	}
	return bs.SafeWriteJSON(hyperliquidWSRequest{
		Method:       method,
		Subscription: &hyperliquidSubscription{Type: "trades", Coin: coin},
	})
}

func (h *hyperliquidHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var msg hyperliquidWSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	switch msg.Channel {
	case "trades":
	case "error":
		log.Errorf("HYPERLIQUID - server error: %s.", msg.Data)
		return
	default:
		return
	}

	var trades []hyperliquidTrade
	if err := json.Unmarshal(msg.Data, &trades); err != nil {
		log.Errorf("HYPERLIQUID - unmarshal trades: %v.", err)
		return
	}

	for _, t := range trades {
		h.mu.RLock()
		foreignName, ok := h.foreignByCoin[t.Coin]
		h.mu.RUnlock()
		if !ok {
			continue
		}
		lock.RLock()
		pair, ok := bs.tickerPairMap[h.TickerKeyFromForeign(foreignName)]
		lock.RUnlock()
		if !ok {
			continue
		}

		price, err := strconv.ParseFloat(t.Px, 64)
		if err != nil {
			log.Errorf("HYPERLIQUID - parse price %q: %v.", t.Px, err)
			continue
		}
		volume, err := strconv.ParseFloat(t.Sz, 64)
		if err != nil {
			log.Errorf("HYPERLIQUID - parse size %q: %v.", t.Sz, err)
			continue
		}
		if math.IsNaN(price) || math.IsInf(price, 0) || price <= 0 ||
			math.IsNaN(volume) || math.IsInf(volume, 0) {
			log.Warnf("HYPERLIQUID - dropping %s trade with price %q and size %q.", foreignName, t.Px, t.Sz)
			continue
		}
		if t.Side == "A" {
			volume = -volume
		}

		trade := models.Trade{
			Price:          price,
			Volume:         volume,
			Time:           time.UnixMilli(t.Time),
			Exchange:       Exchanges[HYPERLIQUID_EXCHANGE],
			BaseToken:      pair.BaseToken,
			QuoteToken:     pair.QuoteToken,
			ForeignTradeID: strconv.FormatInt(t.Tid, 10),
		}

		bs.setLastTradeTime(lock, foreignName, trade.Time)
		bs.tradesChannel <- trade
	}
}

func (h *hyperliquidHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (h *hyperliquidHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

func (h *hyperliquidHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// coin maps a pair like "UBTC-USDC" onto its Hyperliquid spot coin, e.g. "@142".
func (h *hyperliquidHooks) coin(foreignName string) (string, error) {
	symbols := strings.Split(foreignName, "-")
	if len(symbols) != 2 {
		return "", fmt.Errorf("bad pair format: %q", foreignName)
	}
	name := symbols[0] + "/" + symbols[1]

	h.mu.RLock()
	coin, ok := h.coins[name]
	expired := time.Since(h.fetchedAt) > hyperliquidSpotMetaTTL
	h.mu.RUnlock()
	if !ok {
		if !expired {
			return "", fmt.Errorf("%s is not listed on Hyperliquid spot", name)
		}
		coins, err := fetchHyperliquidSpotCoins(h.infoURL)
		if err != nil {
			return "", fmt.Errorf("fetch spotMeta: %w", err)
		}
		h.mu.Lock()
		h.coins = coins
		h.fetchedAt = time.Now()
		h.mu.Unlock()
		if coin, ok = coins[name]; !ok {
			return "", fmt.Errorf("%s is not listed on Hyperliquid spot", name)
		}
	}

	h.mu.Lock()
	if previous, exists := h.foreignByCoin[coin]; exists && previous != foreignName {
		log.Errorf("HYPERLIQUID - coin %s is mapped to %s, overwriting with %s. One of them will receive no trades.", coin, previous, foreignName)
	}
	h.foreignByCoin[coin] = foreignName
	h.mu.Unlock()
	return coin, nil
}

func fetchHyperliquidSpotCoins(infoURL string) (map[string]string, error) {
	body, err := json.Marshal(map[string]string{"type": "spotMeta"})
	if err != nil {
		return nil, err
	}
	resp, err := hyperliquidHTTPClient.Post(infoURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var meta hyperliquidSpotMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, err
	}
	return hyperliquidSpotCoins(meta), nil
}

// hyperliquidSpotCoins maps "QUOTE/BASE" token names onto the coin used in subscriptions.
// Only PURR/USDC keeps its name, all other pairs are named "@<index>".
func hyperliquidSpotCoins(meta hyperliquidSpotMeta) map[string]string {
	names := make(map[int]string, len(meta.Tokens))
	for _, t := range meta.Tokens {
		names[t.Index] = t.Name
	}
	coins := make(map[string]string, len(meta.Universe))
	for _, p := range meta.Universe {
		if len(p.Tokens) != 2 {
			continue
		}
		quote, okQuote := names[p.Tokens[0]]
		base, okBase := names[p.Tokens[1]]
		if !okQuote || !okBase {
			continue
		}
		coins[quote+"/"+base] = p.Name
	}
	return coins
}

func NewHyperliquidScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	hooks := &hyperliquidHooks{
		infoURL:       hyperliquidInfoURL,
		coins:         make(map[string]string),
		foreignByCoin: make(map[string]string),
	}
	return NewBaseCEXScraper(ctx, pairs, wg, hooks, branchMarketConfig)
}
