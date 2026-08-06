package scrapers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
)

const (
	coinexWSBaseURL         = "wss://socket.coinex.com/v2/spot"
	coinexPingPeriodDefault = 20
)

type coinexSubscribeParams struct {
	MarketList []string `json:"market_list"`
}

type coinexSubscribeMessage struct {
	Method string                `json:"method"`
	Params coinexSubscribeParams `json:"params"`
	ID     int64                 `json:"id"`
}

type coinexPingMessage struct {
	Method string   `json:"method"`
	Params struct{} `json:"params"`
	ID     int64    `json:"id"`
}

// coinexWSMessage covers both pushed updates (Method == "deals.update") and
// request acks / pong replies (Method == "", Code/Message set).
type coinexWSMessage struct {
	Method  string          `json:"method"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type coinexDealsData struct {
	Market   string       `json:"market"`
	DealList []coinexDeal `json:"deal_list"`
}

type coinexDeal struct {
	DealID    int64  `json:"deal_id"`
	CreatedAt int64  `json:"created_at"` // milliseconds
	Side      string `json:"side"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
}

// ---------------- hooks ----------------

type coinexHooks struct{}

func (coinexHooks) ExchangeKey() string { return COINEX_EXCHANGE }
func (coinexHooks) WSURL() string       { return coinexWSBaseURL }

func (coinexHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	pingPeriod, err := strconv.Atoi(utils.Getenv("COINEX_PING_PERIOD_SECONDS", strconv.Itoa(coinexPingPeriodDefault)))
	if err != nil || pingPeriod <= 0 {
		log.Errorf("COINEX - parse COINEX_PING_PERIOD_SECONDS: %v. Set to default %d.", err, coinexPingPeriodDefault)
		pingPeriod = coinexPingPeriodDefault
	}
	go func() {
		tick := time.NewTicker(time.Duration(pingPeriod) * time.Second)
		defer tick.Stop()
		var id int64
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				id++
				msg := coinexPingMessage{Method: "server.ping", ID: id}
				if err := bs.SafeWriteJSON(msg); err != nil {
					log.Errorf("COINEX - send server.ping: %v.", err)
				}
			}
		}
	}()
}

// Subscribe sends a deals.subscribe/deals.unsubscribe for the public spot
// market. CoinEx market names have no separator, e.g. BTCUSDT.
func (coinexHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	method := "deals.unsubscribe"
	if subscribe {
		method = "deals.subscribe"
	}
	msg := coinexSubscribeMessage{
		Method: method,
		Params: coinexSubscribeParams{
			MarketList: []string{strings.ReplaceAll(pair.ForeignName, "-", "")},
		},
		ID: time.Now().UnixNano(),
	}
	return bs.SafeWriteJSON(msg)
}

func (coinexHooks) OnMessage(bs *BaseCEXScraper, mt int, data []byte, lock *sync.RWMutex) {
	payload, err := coinexGunzip(data)
	if err != nil {
		log.Errorf("COINEX - decompress message: %v.", err)
		return
	}

	var msg coinexWSMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		log.Errorf("COINEX - unmarshal message: %v.", err)
		return
	}

	if msg.Method != "deals.update" {
		if msg.Code != 0 {
			log.Errorf("COINEX - server error: code=%d msg=%s.", msg.Code, msg.Message)
		}
		return
	}

	var deals coinexDealsData
	if err := json.Unmarshal(msg.Data, &deals); err != nil {
		log.Errorf("COINEX - unmarshal deals data: %v.", err)
		return
	}
	if len(deals.DealList) == 0 {
		return
	}

	lock.RLock()
	pair, ok := bs.tickerPairMap[deals.Market]
	lock.RUnlock()
	if !ok {
		return
	}

	for _, d := range deals.DealList {
		price, err := strconv.ParseFloat(d.Price, 64)
		if err != nil {
			log.Errorf("COINEX - parse price %q: %v.", d.Price, err)
			continue
		}
		volume, err := strconv.ParseFloat(d.Amount, 64)
		if err != nil {
			log.Errorf("COINEX - parse amount %q: %v.", d.Amount, err)
			continue
		}
		if d.Side == "sell" {
			volume = -volume
		}

		t := time.Unix(0, d.CreatedAt*int64(time.Millisecond))

		trade := models.Trade{
			Price:          price,
			Volume:         volume,
			Time:           t,
			Exchange:       Exchanges[COINEX_EXCHANGE],
			BaseToken:      pair.BaseToken,
			QuoteToken:     pair.QuoteToken,
			ForeignTradeID: strconv.FormatInt(d.DealID, 10),
		}

		bs.setLastTradeTime(lock, pair.QuoteToken.Symbol+"-"+pair.BaseToken.Symbol, t)
		bs.tradesChannel <- trade
	}
}

func (coinexHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// TickerKeyFromForeign maps "BTC-USDT" -> "BTCUSDT"
func (coinexHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

func (coinexHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign // "QUOTE-BASE"
}

func coinexGunzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func NewCoinExScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	return NewBaseCEXScraper(ctx, pairs, wg, coinexHooks{}, branchMarketConfig)
}
