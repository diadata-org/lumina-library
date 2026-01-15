package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	ws "github.com/gorilla/websocket"
)

type kuCoinWSSubscribeMessage struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	PrivateChannel bool   `json:"privateChannel"`
	Response       bool   `json:"response"`
}

type kuCoinWSResponse struct {
	Type    string       `json:"type"`  // "pong", "message", etc.
	Topic   string       `json:"topic"` // "/market/match:BTC-USDT"
	Subject string       `json:"subject"`
	Data    kuCoinWSData `json:"data"`
}

type kuCoinWSData struct {
	Sequence string `json:"sequence"`
	Type     string `json:"type"`
	Symbol   string `json:"symbol"` // e.g. "BTC-USDT"
	Side     string `json:"side"`   // "buy" / "sell"
	Price    string `json:"price"`
	Size     string `json:"size"`
	TradeID  string `json:"tradeId"`
	Time     string `json:"time"` // ms timestamp as string
}

// Ping message
type kuCoinWSMessage struct {
	Id   string `json:"id"`
	Type string `json:"type"` // "ping"
}

type kuCoinPostResponse struct {
	Code string `json:"code"`
	Data struct {
		Token           string            `json:"token"`
		InstanceServers []instanceServers `json:"instanceServers"`
	} `json:"data"`
}

type instanceServers struct {
	PingInterval int64 `json:"pingInterval"`
}

var (
	kucoinWSBaseString    = "wss://ws-api-spot.kucoin.com/"
	kucoinTokenURL        = "https://api.kucoin.com/api/v1/bullet-public"
	kucoinPingIntervalFix = int64(10) // seconds
)

type kucoinHooks struct{}

func (kucoinHooks) ExchangeKey() string {
	return KUCOIN_EXCHANGE
}

// directly get token from WSURL, and concatenate with ws address
func (kucoinHooks) WSURL() string {
	token, _, err := getPublicKuCoinToken(kucoinTokenURL)
	if err != nil {
		log.Errorf("KuCoin - getPublicKuCoinToken: %v.", err)
		// return empty string, let upper layer Dial fail and trigger failover
		return ""
	}
	return kucoinWSBaseString + "?token=" + token
}

// after open connection, start ping routine
func (kucoinHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	go func() {
		tick := time.NewTicker(time.Duration(kucoinPingIntervalFix) * time.Second)
		defer tick.Stop()

		ping := kuCoinWSMessage{
			Type: "ping",
		}

		for {
			select {
			case <-tick.C:
				if err := bs.SafeWriteJSON(ping); err != nil {
					log.Errorf("KuCoin - Send ping: %s.", err.Error())
					return
				}
			case <-ctx.Done():
				log.Warn("KuCoin - Close ping.")
				return
			}
		}
	}()
}

// subscribe / unsubscribe
func (kucoinHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	subscribeType := "unsubscribe"
	if subscribe {
		subscribeType = "subscribe"
	}

	msg := &kuCoinWSSubscribeMessage{
		Type:           subscribeType,
		Topic:          "/market/match:" + pair.ForeignName,
		PrivateChannel: false,
		Response:       true,
	}

	return bs.SafeWriteJSON(msg)
}

func (kucoinHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// handle each ws text message
func (kucoinHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	if mt != ws.TextMessage {
		return
	}

	var message kuCoinWSResponse
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	// pong
	if message.Type == "pong" {
		log.Debug("KuCoin - Successful ping: received pong.")
		return
	}

	// actual trade message
	if message.Type != "message" {
		return
	}

	// Parse trade quantities.
	price, volume, timestamp, foreignTradeID, err := parseKuCoinTradeMessage(message)
	if err != nil {
		log.Errorf("KuCoin - parseTradeMessage: %v.", err)
		return
	}

	// Identify ticker symbols with underlying assets.
	pair := strings.Split(message.Data.Symbol, "-")
	if len(pair) < 2 {
		log.Warnf("KuCoin - Unexpected symbol format: %q", message.Data.Symbol)
		return
	}

	// tickerPairMap key: remove '-' (e.g. "BTC-USDT" -> "BTCUSDT")
	key := pair[0] + pair[1]

	lock.RLock()
	exchangepair, ok := bs.tickerPairMap[key]
	lock.RUnlock()
	if !ok {
		return
	}

	trade := models.Trade{
		QuoteToken:     exchangepair.QuoteToken,
		BaseToken:      exchangepair.BaseToken,
		Price:          price,
		Volume:         volume,
		Time:           timestamp,
		Exchange:       Exchanges[KUCOIN_EXCHANGE],
		ForeignTradeID: foreignTradeID,
	}

	// lastTradeTimeMap key: use "QUOTE-BASE" format (e.g. BTC-USDT)
	lastKey := pair[0] + "-" + pair[1]
	bs.setLastTradeTime(lock, lastKey, trade.Time)

	log.Tracef(
		"KuCoin - got trade: %s -- %v -- %v -- %s.",
		trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
		trade.Price,
		trade.Volume,
		trade.ForeignTradeID,
	)

	bs.tradesChannel <- trade
}

// tickerPairMap key: remove '-' (e.g. BTC-USDT -> BTCUSDT)
func (kucoinHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: use ForeignName (e.g. BTC-USDT)
func (kucoinHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

func parseKuCoinTradeMessage(message kuCoinWSResponse) (price float64, volume float64, timestamp time.Time, foreignTradeID string, err error) {
	price, err = strconv.ParseFloat(message.Data.Price, 64)
	if err != nil {
		return
	}
	volume, err = strconv.ParseFloat(message.Data.Size, 64)
	if err != nil {
		return
	}
	if message.Data.Side == "sell" {
		volume *= -1
	}
	timeMilliseconds, err := strconv.Atoi(message.Data.Time)
	if err != nil {
		return
	}
	timestamp = time.Unix(0, int64(timeMilliseconds))
	foreignTradeID = message.Data.TradeID
	return
}

// getPublicKuCoinToken returns a token for public market data along with the pingInterval in milliseconds.
func getPublicKuCoinToken(url string) (token string, pingInterval int64, err error) {
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)
	data, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return
	}
	defer data.Body.Close()
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return
	}

	var postResp kuCoinPostResponse
	err = json.Unmarshal(body, &postResp)
	if err != nil {
		return
	}
	if len(postResp.Data.InstanceServers) > 0 {
		pingInterval = postResp.Data.InstanceServers[0].PingInterval
	}
	token = postResp.Data.Token
	return
}

func NewKuCoinScraper(
	ctx context.Context,
	pairs []models.ExchangePair,
	branchMarketConfig string,
	wg *sync.WaitGroup,
) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, kucoinHooks{}, branchMarketConfig)
}
