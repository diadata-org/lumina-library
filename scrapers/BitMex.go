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

// bitMexWSResponse covers subscription acks, error frames, and trade pushes.
// BitMEX uses a table-diff protocol: the first push has action="partial"
// (snapshot), subsequent pushes have action="insert" (new trades).
type bitMexWSResponse struct {
	// trade push fields
	Table  string          `json:"table"`
	Action string          `json:"action"`
	Data   []bitMexWSTrade `json:"data"`

	// subscription ack fields
	Success   bool   `json:"success"`
	Subscribe string `json:"subscribe"`

	// error frame fields (sent on a bad subscribe; has no "subscribe" field)
	Error   string          `json:"error"`
	Status  int             `json:"status"`
	Request json.RawMessage `json:"request"`
}

// bitMexWSTrade is a single trade row from the spot "trade" table push.
// For spot symbols (e.g. "ETH_USDT"), size is the base-asset amount traded.
type bitMexWSTrade struct {
	Timestamp  time.Time `json:"timestamp"`
	Symbol     string    `json:"symbol"` // BitMEX native spot symbol, e.g. "ETH_USDT"
	Side       string    `json:"side"`   // "Buy" or "Sell"
	Size       float64   `json:"size"`   // base-asset amount (spot)
	Price      float64   `json:"price"`
	TrdMatchID string    `json:"trdMatchID"` // unique trade ID
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

// Subscribe sends a subscribe or unsubscribe message for the spot trade channel
// of a single pair.
//
// ForeignName convention in lumina-library is "QUOTE-BASE" (e.g. "ETH-USDT").
// BitMEX spot symbols use an underscore separator (e.g. "ETH_USDT"), so we
// convert the dash to an underscore.
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
	if mt != ws.TextMessage {
		return
	}
	// BitMEX sends plain "pong" in response to our "ping".
	if string(data) == "pong" {
		return
	}

	var resp bitMexWSResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Warnf("BitMEX - failed to unmarshal message: %v", err)
		return
	}

	// Subscription ack: {"success":true,"subscribe":"trade:ETH_USDT", ...}
	if resp.Subscribe != "" {
		if !resp.Success {
			log.Errorf("BitMEX - subscription failed for %s.", resp.Subscribe)
		}
		return
	}

	// Error frame: a bad subscribe returns {"status":...,"error":"...","request":{...}}
	// with no "subscribe" field, so handle it separately.
	if resp.Error != "" {
		log.Errorf("BitMEX - server error: %s (status=%d, request=%s).", resp.Error, resp.Status, string(resp.Request))
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

// TickerKeyFromForeign converts "ETH-USDT" -> "ETH_USDT" for tickerPairMap lookup.
func (bitMexHooks) TickerKeyFromForeign(foreign string) string {
	return bitMexNativeSymbol(foreign)
}

// LastTradeTimeKeyFromForeign uses the ForeignName as-is ("ETH-USDT").
func (bitMexHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// processBitMexTrades converts BitMEX spot trade rows to models.Trade and emits them.
func processBitMexTrades(bs *BaseCEXScraper, lock *sync.RWMutex, trades []bitMexWSTrade) {
	for _, d := range trades {
		// Look up the pair via the native spot symbol key (e.g. "ETH_USDT").
		lock.RLock()
		pair, ok := bs.tickerPairMap[d.Symbol]
		lock.RUnlock()
		if !ok {
			log.Debugf("BitMEX - unknown symbol %s, skipping trade.", d.Symbol)
			continue
		}

		// Spot: size is the base-asset amount traded.
		// Volume sign convention: positive = buy, negative = sell.
		volume := d.Size
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

// converts ForeignName "ETH-USDT" to BitMEX spot symbol "ETHUSDT".
func bitMexNativeSymbol(foreignName string) string {
	return strings.ReplaceAll(foreignName, "-", "")
}

// NewBitMexScraper creates and starts a BitMEX scraper.
func NewBitMexScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, bitMexHooks{}, branchMarketConfig)
}