package scrapers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

type binanceWSSubscribeMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type binanceWSResponse struct {
	Timestamp      int64       `json:"T"`
	Price          string      `json:"p"`
	Volume         string      `json:"q"`
	ForeignTradeID int         `json:"t"`
	ForeignName    string      `json:"s"`
	Type           interface{} `json:"e"`
	Buy            bool        `json:"m"`
}

type binanceHooks struct {
	apiConnectRetries int
	proxyIndex        int
}

const (
	BINANCE_API_MAX_RETRIES = 5
	binanceApiWaitSeconds   = 5
)

var (
	binanceWSBaseString  = "wss://stream.binance.com:9443/ws"
	binanceWriteInterval = 500 * time.Millisecond
)

func (h *binanceHooks) ExchangeKey() string {
	return BINANCE_EXCHANGE
}

func (h *binanceHooks) WSURL() string {
	return binanceWSBaseString
}

func (h *binanceHooks) Dial(ctx context.Context, urlStr string) (wsConn, error) {
	var lastErr error

	for {
		// if over certain retry times, switch to alternative proxy
		if h.apiConnectRetries > BINANCE_API_MAX_RETRIES {
			log.Errorf("Binance - too many timeouts for proxy %v. Switch to alternative proxy.", h.proxyIndex)
			h.apiConnectRetries = 0
			h.proxyIndex = (h.proxyIndex + 1) % 2
		}

		username := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(h.proxyIndex)+"_USERNAME", "")
		password := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(h.proxyIndex)+"_PASSWORD", "")
		host := utils.Getenv("BINANCE_PROXY"+strconv.Itoa(h.proxyIndex)+"_HOST", "")

		var d ws.Dialer
		if host != "" {
			user := url.UserPassword(username, password)
			d = ws.Dialer{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "http",
					User:   user,
					Host:   host,
					Path:   "/",
				}),
			}
		}

		conn, _, err := d.Dial(urlStr, nil)
		if err != nil {
			h.apiConnectRetries++
			lastErr = err
			log.Errorf("Binance - Connect to API via proxy %d failed: %v", h.proxyIndex, err)

			// wait for a while and retry, or context is cancelled
			select {
			case <-ctx.Done():
				return nil, lastErr
			case <-time.After(time.Duration(binanceApiWaitSeconds) * time.Second):
				continue
			}
		}

		// success
		h.apiConnectRetries = 0
		return conn, nil
	}
}

// Binance does not need extra ping
func (h *binanceHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {}

// subscribe / unsubscribe
func (h *binanceHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	method := "UNSUBSCRIBE"
	if subscribe {
		method = "SUBSCRIBE"
	}

	pairTicker := strings.ToLower(strings.ReplaceAll(pair.ForeignName, "-", "")) // "BTC-USDT" -> "BTCUSDT"

	msg := &binanceWSSubscribeMessage{
		Method: method,
		Params: []string{pairTicker + "@trade"},
		ID:     1,
	}

	// simple rate limiting: avoid writing too fast
	time.Sleep(binanceWriteInterval)

	return bs.SafeWriteJSON(msg)
}

func (h *binanceHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

func (h *binanceHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var msg binanceWSResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	if msg.Type == nil {
		return
	}

	trade := binanceParseWSResponse(msg)

	// tickerPairMap key: use Binance symbol (e.g. "BTCUSDT")
	lock.RLock()
	pair, ok := bs.tickerPairMap[msg.ForeignName]
	lock.RUnlock()
	if !ok {
		log.Tracef("Binance - tickerPairMap not found for %s.", msg.ForeignName)
		return
	}

	trade.QuoteToken = pair.QuoteToken
	trade.BaseToken = pair.BaseToken

	log.Tracef("Binance - got trade %s -- %v -- %v -- %v.",
		trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
		trade.Price, trade.Volume, trade.ForeignTradeID,
	)

	// lastTradeTimeMap key: use symbol directly, same as TickerKeyFromForeign
	lastKey := h.LastTradeTimeKeyFromForeign(msg.ForeignName)
	bs.setLastTradeTime(lock, lastKey, trade.Time)

	bs.tradesChannel <- trade
}

func binanceParseWSResponse(message binanceWSResponse) (trade models.Trade) {
	var err error
	trade.Exchange = Exchanges[BINANCE_EXCHANGE]
	trade.Time = time.Unix(0, message.Timestamp*1000000)
	trade.Price, err = strconv.ParseFloat(message.Price, 64)
	if err != nil {
		log.Errorf("Binance - Parse price: %v.", err)
	}
	trade.Volume, err = strconv.ParseFloat(message.Volume, 64)
	if err != nil {
		log.Errorf("Binance - Parse volume: %v.", err)
	}
	if !message.Buy {
		trade.Volume *= -1
	}
	trade.ForeignTradeID = strconv.Itoa(int(message.ForeignTradeID))
	return
}

// tickerPairMap key: BTC-USDT -> BTCUSDT
func (h *binanceHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: BTC-USDT -> BTCUSDT
func (h *binanceHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

func NewBinanceScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	failoverChannel chan string,
	wg *sync.WaitGroup,
) Scraper {
	hooks := &binanceHooks{
		apiConnectRetries: 0,
		proxyIndex:        0,
	}
	return NewBaseCEXScraper(ctx, pairs, failoverChannel, wg, hooks)
}
