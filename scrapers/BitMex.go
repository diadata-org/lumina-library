package scrapers

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

const (
	bitMexWSURL = "wss://ws.bitmex.com/realtime"

	// bitMexPingInterval is how often we send a ping to keep the connection alive.
	// BitMEX closes idle connections after ~60s, so 25s is a safe interval.
	bitMexPingInterval = 25 * time.Second
)

// bitMexSubscribeMsg is the outbound subscribe/unsubscribe request.
type bitMexSubscribeMsg struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

// bitMexWSResponse covers both subscription acks and trade pushes.
// BitMEX uses a table-diff protocol: the first push has action="partial",
// subsequent pushes have action="insert" (new trades).
type bitMexWSResponse struct {
	// trade push fields
	Table  string          `json:"table"`
	Action string          `json:"action"`
	Data   []bitMexWSTrade `json:"data"`

	// subscription ack fields
	Success   bool   `json:"success"`
	Subscribe string `json:"subscribe"`
	Error     string `json:"error"`
	Status    int    `json:"status"`
}

// bitMexWSTrade is a single trade row from the "trade" table push.
type bitMexWSTrade struct {
	Timestamp       time.Time `json:"timestamp"`
	Symbol          string    `json:"symbol"`          // BitMEX native symbol, e.g. "XBTUSD"
	Side            string    `json:"side"`            // "Buy" or "Sell"
	Size            float64   `json:"size"`            // contract count
	Price           float64   `json:"price"`
	HomeNotional    float64   `json:"homeNotional"`    // base asset value (XBT for inverse, quote for linear)
	ForeignNotional float64   `json:"foreignNotional"` // quote asset value
	TrdMatchID      string    `json:"trdMatchID"`      // unique trade ID
}

type bitMexHooks struct{}

func (bitMexHooks) ExchangeKey() string { return BITMEX_EXCHANGE }
func (bitMexHooks) WSURL() string       { return bitMexWSURL }

// OnOpen starts the periodic ping goroutine.
// BitMEX expects a plain text "ping" and replies "pong".
func (bitMexHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	go func() {
		ticker := time.NewTicker(bitMexPingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := bs.SafeWriteMessage(ws.TextMessage, []byte("ping")); err != nil {
					log.Warnf("BitMEX - ping error: %v", err)
				}
			}
		}
	}()
}

// Subscribe sends a subscribe or unsubscribe message for the trade channel
// of a single pair.
//
// ForeignName convention in lumina-library is "QUOTE-BASE" (e.g. "XBT-USD").
// BitMEX native symbols concatenate the two without a separator (e.g. "XBTUSD"),
// so we just strip the dash.
func (bitMexHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	op := "unsubscribe"
	if subscribe {
		op = "subscribe"
	}
	nativeSymbol := bitMexNativeSymbol(pair.ForeignName)
	msg := bitMexSubscribeMsg{
		Op:   op,
		Args: []string{"trade:" + nativeSymbol},
	}
	return bs.SafeWriteJSON(msg)
}

// OnMessage handles every raw WS message from BitMEX.
func (bitMexHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	// BitMEX sends plain "pong" in response to our "ping".
	if mt == ws.TextMessage && string(data) == "pong" {
		return
	}
	if mt != ws.TextMessage {
		return
	}

	var resp bitMexWSResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Warnf("BitMEX - failed to unmarshal message: %v", err)
		return
	}

	// Subscription ack — log errors, ignore success.
	if resp.Subscribe != "" {
		if !resp.Success {
			log.Errorf("BitMEX - subscription failed for %s: %s (status=%d)", resp.Subscribe, resp.Error, resp.Status)
		}
		return
	}

	// Trade push: both "partial" (snapshot on subscribe) and "insert" (live trades).
	if resp.Table == "trade" && len(resp.Data) > 0 {
		processBitMexTrades(bs, lock, resp.Data)
	}
}

func (bitMexHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false // use BaseCEXScraper default ReadMessage loop
}

// TickerKeyFromForeign converts "XBT-USD" -> "XBTUSD" for tickerPairMap lookup.
func (bitMexHooks) TickerKeyFromForeign(foreign string) string {
	return bitMexNativeSymbol(foreign)
}

// LastTradeTimeKeyFromForeign uses the ForeignName as-is ("XBT-USD").
func (bitMexHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// processBitMexTrades converts BitMEX trade rows to models.Trade and emits them.
func processBitMexTrades(bs *BaseCEXScraper, lock *sync.RWMutex, trades []bitMexWSTrade) {
	for _, d := range trades {
		// Look up the pair via the native symbol key (e.g. "XBTUSD").
		lock.RLock()
		pair, ok := bs.tickerPairMap[d.Symbol]
		lock.RUnlock()
		if !ok {
			log.Debugf("BitMEX - unknown symbol %s, skipping trade", d.Symbol)
			continue
		}

		// Volume sign convention: positive = buy, negative = sell.
		// Use homeNotional (base-asset denominated) as the volume, consistent
		// with how the old repo handled BitMEX inverse contracts.
		volume := d.HomeNotional
		if d.Side == "Sell" {
			volume = -volume
		}

		trade := models.Trade{
			QuoteToken:     pair.QuoteToken,
			BaseToken:      pair.BaseToken,
			Price:          d.Price,
			Volume:         volume,
			Time:           d.Timestamp,
			Exchange:       Exchanges[BITMEX_EXCHANGE],
			ForeignTradeID: d.TrdMatchID,
		}

		log.Tracef(
			"BitMEX - trade: %s/%s price=%.4f vol=%.8f id=%s",
			pair.QuoteToken.Symbol,
			pair.BaseToken.Symbol,
			trade.Price,
			trade.Volume,
			trade.ForeignTradeID,
		)

		// Update watchdog timestamp (keyed by "QUOTE-BASE" ForeignName).
		foreignName := pair.QuoteToken.Symbol + "-" + pair.BaseToken.Symbol
		bs.setLastTradeTime(lock, foreignName, trade.Time)

		bs.tradesChannel <- trade
	}
}

// bitMexNativeSymbol converts lumina ForeignName "XBT-USD" to BitMEX symbol "XBTUSD".
func bitMexNativeSymbol(foreignName string) string {
	return strings.ReplaceAll(foreignName, "-", "")
}

// NewBitMexScraper creates and starts a BitMEX scraper.
func NewBitMexScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, bitMexHooks{}, branchMarketConfig)
}