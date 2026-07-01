package scrapers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

const (
	bitgetWSBaseURL = "wss://ws.bitget.com/v2/ws/public"
	// Bitget requires a string "ping" at least every 2 minutes, otherwise the
	// server closes the connection. Docs recommend a 30s timer.
	bitgetPingPeriodDefault = 30
)

// ---------------- wire types ----------------

type bitgetWSSubscribeArg struct {
	InstType string `json:"instType"`
	Channel  string `json:"channel"`
	InstID   string `json:"instId"`
}

type bitgetWSSubscribeMessage struct {
	OP   string                 `json:"op"`
	Args []bitgetWSSubscribeArg `json:"args"`
}

type bitgetWSResponse struct {
	Action string               `json:"action"` // "snapshot" | "update"
	Arg    bitgetWSSubscribeArg `json:"arg"`
	Data   []bitgetTradeData    `json:"data"`
	TS     int64                `json:"ts"`
	// present on subscribe/unsubscribe acks and errors
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

type bitgetTradeData struct {
	Timestamp string `json:"ts"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"` // "buy" | "sell"
	TradeID   string `json:"tradeId"`
}

// ---------------- hooks ----------------

type bitgetHooks struct{}

func (bitgetHooks) ExchangeKey() string { return BITGET_EXCHANGE }
func (bitgetHooks) WSURL() string       { return bitgetWSBaseURL }

// OnOpen starts the heartbeat. Bitget uses an application-level string "ping"
// (NOT a websocket control-frame ping) and replies with a string "pong".
func (bitgetHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	pingPeriod, err := strconv.Atoi(utils.Getenv("BITGET_PING_PERIOD_SECONDS", strconv.Itoa(bitgetPingPeriodDefault)))
	if err != nil || pingPeriod <= 0 {
		log.Errorf("BITGET - parse BITGET_PING_PERIOD_SECONDS: %v. Set to default %d.", err, bitgetPingPeriodDefault)
		pingPeriod = bitgetPingPeriodDefault
	}
	go func() {
		tick := time.NewTicker(time.Duration(pingPeriod) * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := bs.SafeWriteMessage(ws.TextMessage, []byte("ping")); err != nil {
					log.Errorf("BITGET - send ping: %v.", err)
				}
			}
		}
	}()
}

// Subscribe sends a subscribe/unsubscribe for the public SPOT "trade" channel.
// Bitget instId is the symbol without separator, e.g. BTCUSDT.
func (bitgetHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	op := "unsubscribe"
	if subscribe {
		op = "subscribe"
	}
	msg := bitgetWSSubscribeMessage{
		OP: op,
		Args: []bitgetWSSubscribeArg{
			{
				InstType: "SPOT",
				Channel:  "trade",
				InstID:   strings.ReplaceAll(pair.ForeignName, "-", ""),
			},
		},
	}
	return bs.SafeWriteJSON(msg)
}

// OnMessage parses a single raw ws message.
func (bitgetHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	// Heartbeat reply: server sends string "pong".
	if string(data) == "pong" {
		return
	}

	var resp bitgetWSResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}

	// Subscribe/unsubscribe acks and error frames carry "event".
	if resp.Event != "" {
		if resp.Event == "error" {
			log.Errorf("BITGET - server error: code=%s msg=%s.", resp.Code, resp.Msg)
		}
		return
	}

	// Only handle the trade channel.
	if resp.Arg.Channel != "trade" {
		return
	}

	// instId is e.g. BTCUSDT, which matches TickerKeyFromForeign(foreign).
	pairKey := resp.Arg.InstID
	lock.RLock()
	pair, ok := bs.tickerPairMap[pairKey]
	lock.RUnlock()
	if !ok {
		return
	}

	for _, d := range resp.Data {
		price, err := strconv.ParseFloat(d.Price, 64)
		if err != nil {
			log.Errorf("BITGET - parse price %q: %v.", d.Price, err)
			continue
		}
		volume, err := strconv.ParseFloat(d.Size, 64)
		if err != nil {
			log.Errorf("BITGET - parse size %q: %v.", d.Size, err)
			continue
		}
		if d.Side == "sell" {
			volume = -volume
		}

		var t time.Time
		if ms, err := strconv.ParseInt(d.Timestamp, 10, 64); err == nil {
			t = time.Unix(0, ms*int64(time.Millisecond))
		} else {
			t = time.Now()
		}

		trade := models.Trade{
			Price:          price,
			Volume:         volume,
			Time:           t,
			Exchange:       Exchanges[BITGET_EXCHANGE],
			BaseToken:      pair.BaseToken,
			QuoteToken:     pair.QuoteToken,
			ForeignTradeID: d.TradeID,
		}

		// lastTradeTime key convention: "QUOTE-BASE" (same as ByBit/MEXC).
		bs.setLastTradeTime(lock, pair.QuoteToken.Symbol+"-"+pair.BaseToken.Symbol, t)
		bs.tradesChannel <- trade
	}
}

// ReadLoop: use the default ReadMessage->OnMessage loop in BaseCEXScraper.
func (bitgetHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// TickerKeyFromForeign maps "BTC-USDT" -> "BTCUSDT" to match the instId in push data.
func (bitgetHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

func (bitgetHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign // "QUOTE-BASE"
}

func NewBitgetScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, bitgetHooks{}, branchMarketConfig)
}
