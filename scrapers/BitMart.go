package scrapers

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	ws "github.com/gorilla/websocket"
)

// Bitmart public spot WebSocket (v2).
//
// There are two public endpoints:
//   - wss://ws-manager-compress.bitmart.com/api?protocol=1.1  (server sends raw-deflate
//     compressed binary frames; client sends "ping" / receives "pong")
//   - wss://ws-manager-spot-pub-sg.bitmart.com/api?protocol=1.1 (plain text JSON frames)
//
// We default to the compressed endpoint and transparently inflate binary frames in
// OnMessage. OnMessage also handles text frames, so switching to the plain endpoint
// requires only changing bitmartWSBaseURL.
const (
	bitmartWSBaseURL = "wss://ws-manager-compress.bitmart.com/api?protocol=1.1"

	// v2 spot public trade channel. Args are formatted as "spot/trade:BTC_USDT".
	bitmartTradeChannel = "spot/trade"

	bitmartPingMessage = "ping"
	bitmartPongMessage = "pong"
	bitmartSellSide    = "sell"

	// Bitmart expects a heartbeat roughly every 15s; closes idle connections after ~60s.
	bitmartPingInterval = 15 * time.Second

	// Ceiling on a single decompressed frame. Bitmart trade frames are far smaller than
	// this; the bound guards against a decompression bomb from a malicious/faulty upstream.
	bitmartMaxDecompressed = 8 << 20 // 8 MiB
)

// subscribe / unsubscribe request.
type bitmartWSRequest struct {
	Op   string   `json:"op"`   // "subscribe" / "unsubscribe"
	Args []string `json:"args"` // e.g. ["spot/trade:BTC_USDT"]
}

// bitmartTradeData is a single trade entry inside a spot/trade push.
type bitmartTradeData struct {
	Symbol       string `json:"symbol"` // e.g. "BTC_USDT"
	Price        string `json:"price"`
	Side         string `json:"side"` // "buy" / "sell"
	Size         string `json:"size"`
	TimestampSec int64  `json:"s_t"` // seconds
}

// trade push message.
type bitmartWSTradeResponse struct {
	Table string             `json:"table"` // "spot/trade"
	Data  []bitmartTradeData `json:"data"`
	// error envelope fields (present only on errors)
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
	Event        string `json:"event"`
}

type bitmartHooks struct{}

func (bitmartHooks) ExchangeKey() string { return BITMART_EXCHANGE }
func (bitmartHooks) WSURL() string       { return bitmartWSBaseURL }

// OnOpen starts the heartbeat. Bitmart's heartbeat is the literal text "ping"; the server
// answers with the literal text "pong" (ignored in OnMessage).
func (bitmartHooks) OnOpen(ctx context.Context, bs *BaseCEXScraper) {
	go func() {
		tick := time.NewTicker(bitmartPingInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := bs.SafeWriteMessage(ws.TextMessage, []byte(bitmartPingMessage)); err != nil {
					log.Errorf("%s - send ping: %v", BITMART_EXCHANGE, err)
					return
				}
			}
		}
	}()
}

// Subscribe sends a subscribe/unsubscribe for a single pair on the spot/trade channel.
// Bitmart's wire symbol uses an underscore (BTC_USDT). Our ExchangePair.ForeignName uses a
// hyphen (BTC-USDT), so we convert here.
func (bitmartHooks) Subscribe(bs *BaseCEXScraper, pair models.ExchangePair, subscribe bool, lock *sync.RWMutex) error {
	op := "unsubscribe"
	if subscribe {
		op = "subscribe"
	}
	wireSymbol := strings.ReplaceAll(pair.ForeignName, "-", "_")
	msg := bitmartWSRequest{
		Op:   op,
		Args: []string{bitmartTradeChannel + ":" + wireSymbol},
	}
	return bs.SafeWriteJSON(msg)
}

// ReadLoop: use the Base default read loop (ReadMessage -> OnMessage).
func (bitmartHooks) ReadLoop(ctx context.Context, bs *BaseCEXScraper, lock *sync.RWMutex) (handled bool) {
	return false
}

// OnMessage handles both compressed binary frames and plain text frames.
func (bitmartHooks) OnMessage(bs *BaseCEXScraper, messageType int, data []byte, lock *sync.RWMutex) {
	var payload []byte

	switch messageType {
	case ws.BinaryMessage:
		// Bitmart's compressed endpoint sends raw-deflate (no zlib header) frames.
		out, err := bitmartInflate(data)
		if err != nil {
			log.Errorf("%s - inflate message: %v", BITMART_EXCHANGE, err)
			return
		}
		payload = out
	case ws.TextMessage:
		// Heartbeat reply on the plain endpoint, or JSON.
		if string(data) == bitmartPongMessage {
			return
		}
		payload = data
	default:
		return
	}

	// The pong on the compressed endpoint arrives as text "pong" too; guard once more.
	if string(payload) == bitmartPongMessage {
		return
	}

	var resp bitmartWSTradeResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		// subscription acks and other control frames are not trade payloads; ignore quietly.
		log.Tracef("%s - non-trade message: %s", BITMART_EXCHANGE, string(payload))
		return
	}
	if resp.ErrorCode != "" {
		log.Errorf("%s - error code %s on %s event: %s", BITMART_EXCHANGE, resp.ErrorCode, resp.Event, resp.ErrorMessage)
		return
	}
	if resp.Table != bitmartTradeChannel {
		return
	}

	for _, d := range resp.Data {
		price, err := strconv.ParseFloat(d.Price, 64)
		if err != nil {
			continue
		}
		volume, err := strconv.ParseFloat(d.Size, 64)
		if err != nil {
			continue
		}
		if d.Side == bitmartSellSide {
			volume = -volume
		}

		// d.Symbol is "BTC_USDT"; tickerPairMap key is "BTCUSDT".
		tickerKey := strings.ReplaceAll(d.Symbol, "_", "")
		lock.RLock()
		pair, ok := bs.tickerPairMap[tickerKey]
		lock.RUnlock()
		if !ok {
			log.Tracef("%s - unknown ticker key: %s", BITMART_EXCHANGE, tickerKey)
			continue
		}

		timestamp := time.Now()
		if d.TimestampSec > 0 {
			timestamp = time.Unix(d.TimestampSec, 0)
		}

		trade := models.Trade{
			Price:      price,
			Volume:     volume,
			Time:       timestamp,
			Exchange:   Exchanges[BITMART_EXCHANGE],
			BaseToken:  pair.BaseToken,
			QuoteToken: pair.QuoteToken,
			// ForeignTradeID is NOT guaranteed unique: Bitmart's spot/trade payload exposes
			// only a second-granularity timestamp (s_t) and no trade id, so multiple trades
			// on the same symbol within one second share an ID. Do not dedup on this field
			// downstream. Prefer a finer field (ms timestamp / sequence) if the v2 payload
			// exposes one.
			ForeignTradeID: d.Symbol + "-" + strconv.FormatInt(d.TimestampSec, 10),
		}

		log.Tracef("%s - got trade: %s -- %v -- %v.",
			BITMART_EXCHANGE,
			trade.QuoteToken.Symbol+"-"+trade.BaseToken.Symbol,
			trade.Price, trade.Volume,
		)

		// lastTradeTimeMap key matches LastTradeTimeKeyFromForeign (hyphenated foreign name).
		bs.setLastTradeTime(lock, strings.ReplaceAll(d.Symbol, "_", "-"), timestamp)
		bs.tradesChannel <- trade
	}
}

// tickerPairMap key: hyphen removed, e.g. "BTC-USDT" -> "BTCUSDT".
func (bitmartHooks) TickerKeyFromForeign(foreign string) string {
	return strings.ReplaceAll(foreign, "-", "")
}

// lastTradeTimeMap key: hyphenated foreign name, e.g. "BTC-USDT".
func (bitmartHooks) LastTradeTimeKeyFromForeign(foreign string) string {
	return foreign
}

// bitmartInflate decompresses a raw-deflate (headerless) frame, bounded to
// bitmartMaxDecompressed to prevent a decompression bomb from exhausting memory.
func bitmartInflate(b []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(b))
	defer r.Close()
	// Read one byte past the ceiling so we can distinguish a legitimately large frame
	// from an over-limit one instead of silently truncating (LimitReader truncates).
	out, err := io.ReadAll(io.LimitReader(r, bitmartMaxDecompressed+1))
	if err != nil {
		return nil, err
	}
	if len(out) > bitmartMaxDecompressed {
		return nil, fmt.Errorf("bitmart: decompressed frame exceeds %d bytes", bitmartMaxDecompressed)
	}
	return out, nil
}

// ---------------- multi-connection sharding ----------------

// BitMart caps the number of topics per connection (error 90006: "Subscribed total topic
// quantity exceeds limit"), so a large pair list must be split across multiple connections.
//
// Each shard is a normal BaseCEXScraper using the same bitmartHooks, so all of Base's read
// loop / watchdog / resubscribe / reconnect logic is reused as-is.
type bitmartShardedScraper struct {
	tradesChannel chan models.Trade
	shards        []*BaseCEXScraper
	cancels       []context.CancelFunc
	closeOnce     sync.Once
}

func (s *bitmartShardedScraper) TradesChannel() chan models.Trade { return s.tradesChannel }

// Close cancels every shard. Safe to call multiple times.
func (s *bitmartShardedScraper) Close(cancel context.CancelFunc) error {
	var firstErr error
	s.closeOnce.Do(func() {
		if cancel != nil {
			cancel()
		}
		for i, shard := range s.shards {
			if err := shard.Close(s.cancels[i]); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	})
	return firstErr
}

// chunkPairs splits pairs into chunks of at most size. size<=0 means a single chunk.
func chunkPairs(pairs []models.ExchangePair, size int) [][]models.ExchangePair {
	if size <= 0 {
		size = len(pairs)
	}
	if size == 0 {
		return nil
	}
	var out [][]models.ExchangePair
	for i := 0; i < len(pairs); i += size {
		end := i + size
		if end > len(pairs) {
			end = len(pairs)
		}
		out = append(out, pairs[i:end])
	}
	return out
}

// NewBitMartScraper creates one BaseCEXScraper per shard of pairs and fans the shards' trade
// channels into one. Shard size is configurable via BITMART_SHARD_SIZE (default 70). With
// ~140 pairs and size 70 this produces 2 connections, each well under Bitmart's per-connection
// topic limit.
func NewBitMartScraper(ctx context.Context, pairs []models.ExchangePair, branchMarketConfig string, wg *sync.WaitGroup) Scraper {
	shardSize, err := strconv.Atoi(utils.Getenv("BITMART_SHARD_SIZE", "70"))
	if err != nil || shardSize <= 0 {
		log.Errorf("%s - parse BITMART_SHARD_SIZE: %v. Using default 70.", BITMART_EXCHANGE, err)
		shardSize = 70
	}

	chunks := chunkPairs(pairs, shardSize)
	merged := make(chan models.Trade)

	s := &bitmartShardedScraper{
		tradesChannel: merged,
		shards:        make([]*BaseCEXScraper, 0, len(chunks)),
		cancels:       make([]context.CancelFunc, 0, len(chunks)),
	}

	// No pairs -> no shard -> no NewBaseCEXScraper call -> nobody matches the caller's
	// wg.Add(1) in Collector.go. Balance it here and return an empty (but valid) scraper,
	// otherwise wg.Wait() hangs forever.
	if len(chunks) == 0 {
		log.Warnf("%s - no pairs to scrape; starting empty scraper.", BITMART_EXCHANGE)
		wg.Done()
		return s
	}

	// COUPLING: NewBaseCEXScraper runs `defer wg.Done()` exactly once per call, and the
	// caller (Collector.go) already did wg.Add(1) for this exchange. We add the extra shards
	// here so the total Add/Done balances. This depends on Base's wg lifecycle; if
	// NewBaseCEXScraper's wg handling ever changes, update this accounting too.
	if extra := len(chunks) - 1; extra > 0 {
		wg.Add(extra)
	}

	for i, chunk := range chunks {
		// Each shard gets its own cancelable context so Close can stop them independently,
		// while still cascading from the parent ctx.
		shardCtx, shardCancel := context.WithCancel(ctx)
		s.cancels = append(s.cancels, shardCancel)

		// Bitmart allows only one new connection per IP per second. NewBaseCEXScraper
		// connects synchronously, so creating shards serially with a >=1s gap respects that.
		if i > 0 {
			time.Sleep(2 * time.Second)
		}

		shard := NewBaseCEXScraper(shardCtx, chunk, wg, bitmartHooks{}, branchMarketConfig)
		s.shards = append(s.shards, shard)

		// Fan-in: drain this shard's trades into the merged channel.
		go func(src chan models.Trade) {
			for {
				select {
				case <-shardCtx.Done():
					return
				case t, ok := <-src:
					if !ok {
						return
					}
					select {
					case merged <- t:
					case <-shardCtx.Done():
						return
					}
				}
			}
		}(shard.TradesChannel())

		log.Infof("%s - shard %d/%d started with %d pairs.", BITMART_EXCHANGE, i+1, len(chunks), len(chunk))
	}

	return s
}