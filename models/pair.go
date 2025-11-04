package models

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tkanos/gonfig"
)

// Separator for a pair ticker's assets, i.e. BTC-USDT.
const PAIR_TICKER_SEPARATOR = "-"

// ExchangePair is the container for a pair as used by exchanges.
// Across exchanges, these pairs cannot be uniquely mapped on asset pairs.
type ExchangePair struct {
	Symbol         string `json:"Symbol"`
	ForeignName    string `json:"ForeignName"`
	Exchange       string `json:"Exchange"`
	UnderlyingPair Pair   `json:"UnderlyingPair"`
	WatchDogDelay  int64  `json:"WatchDogDelay"`
}

// Pair is a pair of dia assets.
type Pair struct {
	QuoteToken Asset `json:"QuoteToken"`
	BaseToken  Asset `json:"BaseToken"`
}

type PairConfig struct {
	Pair          string `json:"pair"`
	WatchDogDelay int64  `json:"watchDogDelay"`
}

type ExchangeConfig struct {
	ExchangePairs []PairConfig `json:"exchangePairs"`
}

func (p *Pair) ExchangePairIdentifier(exchange string) string {
	return exchange + "-" + p.Identifier()
}

func (p *Pair) Identifier() string {
	return p.QuoteToken.Blockchain + "-" + p.QuoteToken.Address + "-" + p.BaseToken.Blockchain + "-" + p.BaseToken.Address
}

func getPath2Config(directory string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}
	dir := usr.HomeDir
	configPath := "/config/" + directory + "/"
	if dir == "/root" || dir == "/home" {
		return configPath, nil
	}
	return os.Getenv("GOPATH") + "/src/github.com/diadata-org/decentral-feeder" + configPath, nil
}

// According to pairs config file + symbol identifiers directory, construct []ExchangePair
// - exchangeList:        List of exchange names (e.g. ["GateIO", "Binance"])
// - return:              List of ExchangePairs with WatchDogDelay
// - return err:          Error if any
func ExchangePairsFromConfigFiles(exchangeList []string) (allExchangePairs []ExchangePair, err error) {

	if len(exchangeList) == 0 {
		return
	}
	if len(exchangeList) == 1 && len(strings.TrimSpace(exchangeList[0])) == 0 {
		return
	}

	for _, exchange := range exchangeList {
		// 1) Read exchange pairs config (e.g. GateIO.json -> map["ETH-USDT"]watchdogDelay)
		exPairMap, err := GetExchangePairMap(exchange)
		if err != nil {
			return nil, fmt.Errorf("GetExchangeConfig(%s): %w", exchange, err)
		}

		// 2) Read symbol identifiers mapping (directory + exchange.json)
		idMap, err := GetSymbolIdentificationMap(exchange)
		if err != nil {
			return nil, fmt.Errorf("GetSymbolIdentificationMap(%s): %w", exchange, err)
		}

		// 3) To ensure stable output, sort pairs by name (exPairMap's key is "QUOTE-BASE")
		pairs := make([]string, 0, len(exPairMap))
		for pair := range exPairMap {
			pairs = append(pairs, strings.TrimSpace(pair))
		}
		sort.Strings(pairs)

		// 4) Construct ExchangePair one by one
		exchangePairs := []ExchangePair{}
		for _, pair := range pairs {
			ep, err := ConstructExchangePair(exchange, pair, exPairMap[pair], idMap)
			if err != nil {
				return nil, fmt.Errorf("ConstructExchangePair(%s, %s, %v): %w", exchange, pair, exPairMap[pair], err)
			}
			exchangePairs = append(exchangePairs, ep)
		}
		allExchangePairs = append(allExchangePairs, exchangePairs...)
	}

	return allExchangePairs, nil
}

func ConstructExchangePair(exchange string, pair string, watchDogDelay int64, idMap map[string]Asset) (ExchangePair, error) {
	symbols := strings.Split(pair, PAIR_TICKER_SEPARATOR)
	if len(symbols) != 2 {
		return ExchangePair{}, fmt.Errorf("bad pair format: %q (separator=%q)", pair, PAIR_TICKER_SEPARATOR)
	}
	quote := strings.ToUpper(strings.TrimSpace(symbols[0]))
	base := strings.ToUpper(strings.TrimSpace(symbols[1]))

	quoteKey := ExchangeSymbolIdentifier(quote, exchange)
	baseKey := ExchangeSymbolIdentifier(base, exchange)

	qAsset, okQ := idMap[quoteKey]
	bAsset, okB := idMap[baseKey]

	if !okQ {
		return ExchangePair{}, fmt.Errorf("missing symbolId metadata for quote: %v", quoteKey)
	}
	if !okB {
		return ExchangePair{}, fmt.Errorf("missing symbolId metadata for base: %v", baseKey)
	}

	var ep ExchangePair
	ep.Exchange = exchange
	ep.ForeignName = pair
	ep.Symbol = quote                     // Keep consistent with existing function: Symbol use quote
	ep.UnderlyingPair.QuoteToken = qAsset // If missing, it will be zero value Asset
	ep.UnderlyingPair.BaseToken = bAsset
	ep.WatchDogDelay = watchDogDelay

	return ep, nil
}

func GetExchangePairMap(exchange string) (map[string]int64, error) {
	configPath, err := getPath2Config("exchangePairs")
	if err != nil {
		return nil, err
	}
	path := configPath + exchange + ".json"
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var cfg ExchangeConfig
	if err := json.NewDecoder(jsonFile).Decode(&cfg); err != nil {
		return nil, err
	}

	out := make(map[string]int64, len(cfg.ExchangePairs))
	for _, pair := range cfg.ExchangePairs {
		out[pair.Pair] = pair.WatchDogDelay
	}
	return out, nil
}

// MakeExchangepairMap returns a map in which exchangepairs are grouped by exchange string key.
func MakeExchangepairMap(exchangePairs []ExchangePair) map[string][]ExchangePair {
	exchangepairMap := make(map[string][]ExchangePair)
	for _, ep := range exchangePairs {
		exchangepairMap[ep.Exchange] = append(exchangepairMap[ep.Exchange], ep)
	}
	return exchangepairMap
}

// MakeExchangePairString returns a string of exchangepairs separated by comma.
// e.g. "GateIO:ETH-USDT,Binance:BTC-USDT,Binance:BTC-USDC"
func MakeExchangePairString(exchangePairs []ExchangePair) string {
	exchangePairString := ""
	for _, ep := range exchangePairs {
		exchangePairString += ep.Exchange + ":" + ep.ForeignName + ","
	}
	return exchangePairString
}

// MakeTickerPairMap returns a map that maps a pair ticker onto the underlying pair with full asset information.
func MakeTickerPairMap(exchangePairs []ExchangePair) map[string]Pair {
	tickerPairMap := make(map[string]Pair)
	for _, ep := range exchangePairs {
		symbols := strings.Split(ep.ForeignName, "-")
		if len(symbols) < 2 {
			continue
		}
		tickerPairMap[symbols[0]+symbols[1]] = ep.UnderlyingPair
	}
	return tickerPairMap
}

func ExchangeSymbolIdentifier(symbol string, exchange string) string {
	return symbol + "_" + exchange
}

func GetWhitelistedPoolsFromConfig(exchange string) (whitelistedPools []common.Address, err error) {
	path := utils.GetPath("whitelisted_pools/", exchange)
	type pool struct {
		Address string
	}
	type pools struct {
		Pools []pool
	}
	var p pools
	err = gonfig.GetConf(path, &p)
	if err != nil {
		return
	}
	for _, pool := range p.Pools {
		whitelistedPools = append(whitelistedPools, common.HexToAddress(pool.Address))
	}
	return
}
