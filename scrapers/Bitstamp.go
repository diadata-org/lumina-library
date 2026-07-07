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

// ---- WebSocket message types (Bitstamp WS API v2) ----
// Docs: https://www.bitstamp.net/websocket/v2/

// Subscribe / unsubscribe request.
// e.g. {"event":"bts:subscribe","data":{"channel":"live_trades_btcusd"}}
type bitstampWSSubscribeMessage struct {
	Event string                 `json:"event"`
	Data  bitstampWSChannelParam `json:"data"`
}

type bitstampWSChannelParam struct {
	Channel string `json:"channel"`
}

type bitstampWSPingMessage struct {
	Event string `json:"event"`
}

type bitstampWSResponse struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type bitstampWSTradeData struct {
	ID             int64   `json:"id"`
	Amount         float64 `json:"amount"`
	AmountStr      string  `json:"amount_str"`
	Price          float64 `json:"price"`
	PriceStr       string  `json:"price_str"`
	Type           int     `json:"type"` // 0 = buy, 1 = sell
	Timestamp      string  `json:"timestamp"`
	Microtimestamp string  `json:"microtimestamp"`
	BuyOrderID     uint64  `json:"buy_order_id"`
	SellOrderID    uint64  `json:"sell_order_id"`
}

// Heartbeat payload (event == "bts:heartbeat").
type bitstampWSHeartbeatData struct {
	Status string `json:"status"`
}

const bitstampWSBaseString = "wss://ws.bitstamp.net"

var bitstampPingInterval = 20 * time.Second

var bitstampTradeTimeoutSeconds = 120

type bitstampHooks struct{}

func (bitstampHooks) ExchangeKey() string {
	return BITSTAMP_EXCHANGE
}

func (bitstampHooks) WSURL() string {
	return bitstampWSBaseString
}

func (bitstampHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	go bitstampPingLoop(ctx, bs)
}

func bitstampPingLoop(ctx context.Context, bs *BaseCEXScraper) {
	ticker := time.NewTicker(bitstampPingInterval)
	defer ticker.Stop()

	ping := bitstampWSPingMessage{Event: "bts:ping"}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := bs.SafeWriteJSON(ping); err != nil {
				log.Debugf("Bitstamp - ping write failed: %v", err)
			}
		}
	}
}

// Subscribe/unsubscribe to a pair's live_trades channel.
// ForeignName is "BASE-QUOTE" (e.g. BTC-USD); Bitstamp's url_symbol is the
// lowercase, hyphen-stripped form (e.g. btcusd).
func (bitstampHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	event := "bts:unsubscribe"
	if subscribe {
		event = "bts:subscribe"
	}

	urlSymbol := bitstampForeignToURLSymbol(pair.ForeignName)
	msg := bitstampWSSubscribeMessage{
		Event: event,
		Data: bitstampWSChannelParam{
			Channel: "live_trades_" + urlSymbol,
		},
	}

	return bs.SafeWriteJSON(msg)
}

func (bitstampHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	// Use the Base default ReadMessage -> OnMessage loop.
	return false
}

// Parse a single WS text message.
func (bitstampHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var resp bitstampWSResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}

	switch resp.Event {
	case "bts:subscription_succeeded":
		log.Debugf("Bitstamp - subscription succeeded: %s", resp.Channel)
		return

	case "bts:unsubscription_succeeded":
		log.Debugf("Bitstamp - unsubscription succeeded: %s", resp.Channel)
		return

	case "bts:request_reconnect":
		// Server asks us to reconnect; the connection watchdog / read-error path
		// will re-establish the connection and resubscribe.
		log.Warnf("Bitstamp - server requested reconnect on %s", resp.Channel)
		return

	case "bts:pong":
		return

	case "bts:heartbeat":
		var hb bitstampWSHeartbeatData
		if err := json.Unmarshal(resp.Data, &hb); err != nil {
			log.Warnf("Bitstamp - unmarshal heartbeat: %v", err)
			return
		}
		if hb.Status != "success" {
			log.Warnf("Bitstamp - heartbeat status: %s", hb.Status)
		}
		return

	case "trade":
		bitstampHandleTrade(bs, resp, lock)
		return

	default:
		log.Tracef("Bitstamp - unhandled event %q on %s", resp.Event, resp.Channel)
	}
}

func bitstampHandleTrade(bs *BaseCEXScraper, resp bitstampWSResponse, lock *sync.RWMutex) {
	foreignName := bitstampURLSymbolFromChannel(resp.Channel)
	if foreignName == "" {
		return
	}

	var td bitstampWSTradeData
	if err := json.Unmarshal(resp.Data, &td); err != nil {
		log.Errorf("Bitstamp - unmarshal trade: %v", err)
		return
	}

	trade, err := bitstampParseTrade(td)
	if err != nil {
		log.Errorf("Bitstamp - parseTrade: %v", err)
		return
	}

	tickerKey := strings.ToUpper(foreignName)

	lock.RLock()
	pair, ok := bs.tickerPairMap[tickerKey]
	lock.RUnlock()
	if !ok {
		log.Warnf("Bitstamp - no tickerPairMap entry for %q (from channel %q); dropping trade.", tickerKey, resp.Channel)
		return
	}
	trade.QuoteToken = pair.QuoteToken
	trade.BaseToken = pair.BaseToken

	foreignKey := pair.QuoteToken.Symbol + "-" + pair.BaseToken.Symbol

	bs.setLastTradeTime(lock, foreignKey, time.Now())

	if trade.Time.Before(time.Now().Add(-time.Duration(bitstampTradeTimeoutSeconds) * time.Second)) {
		log.Debugf("Bitstamp - dropping stale trade on %s (age > %ds).", foreignKey, bitstampTradeTimeoutSeconds)
		return
	}

	log.Tracef(
		"Bitstamp - got trade: %v -- %s -- %v -- %v -- %s.",
		trade.Time,
		foreignKey,
		trade.Price,
		trade.Volume,
		trade.ForeignTradeID,
	)

	bs.tradesChannel <- trade
}

// Sell trades (type == 1) get a negative volume
func bitstampParseTrade(td bitstampWSTradeData) (models.Trade, error) {
	price := td.Price
	if td.PriceStr != "" {
		p, err := strconv.ParseFloat(td.PriceStr, 64)
		if err != nil {
			return models.Trade{}, err
		}
		price = p
	}

	amount := td.Amount
	if td.AmountStr != "" {
		a, err := strconv.ParseFloat(td.AmountStr, 64)
		if err != nil {
			return models.Trade{}, err
		}
		amount = a
	}

	volume := amount
	if td.Type == 1 {
		volume *= -1
	}

	microts, err := strconv.ParseInt(td.Microtimestamp, 10, 64)
	if err != nil {
		return models.Trade{}, err
	}

	return models.Trade{
		Price:          price,
		Volume:         volume,
		Time:           time.Unix(0, microts*int64(time.Microsecond)),
		Exchange:       Exchanges[BITSTAMP_EXCHANGE],
		ForeignTradeID: strconv.FormatInt(td.ID, 10),
	}, nil
}

// e.g. "BTC-USD" -> "BTCUSD"
func (bitstampHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ToUpper(strings.ReplaceAll(foreign, "-", ""))
}

// LastTradeTimeKeyFromForeign returns the raw ForeignName ("BTC-USD").
func (bitstampHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// bitstampForeignToURLSymbol: "BTC-USD" -> "btcusd".
func bitstampForeignToURLSymbol(foreign string) string {
	return strings.ToLower(strings.ReplaceAll(foreign, "-", ""))
}

// bitstampURLSymbolFromChannel: "live_trades_btcusd" -> "btcusd".
func bitstampURLSymbolFromChannel(channel string) string {
	const prefix = "live_trades_"
	if strings.HasPrefix(channel, prefix) {
		return strings.TrimPrefix(channel, prefix)
	}
	return ""
}

func NewBitstampScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) Scraper {
	if v, err := strconv.Atoi(utils.Getenv("BITSTAMP_TRADE_TIMEOUT", "120")); err != nil {
		log.Errorf("Bitstamp - parse BITSTAMP_TRADE_TIMEOUT: %v.", err)
	} else {
		bitstampTradeTimeoutSeconds = v
	}
	return NewBaseCEXScraper(ctx, pairs, wg, bitstampHooks{}, branchMarketConfig)
}
