package models

import (
	"strings"

	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tkanos/gonfig"
)

// ExchangePair is the container for a pair as used by exchanges.
// Across exchanges, these pairs cannot be uniquely mapped on asset pairs.
type ExchangePair struct {
	Symbol         string `json:"Symbol"`
	ForeignName    string `json:"ForeignName"`
	Exchange       string `json:"Exchange"`
	UnderlyingPair Pair   `json:"UnderlyingPair"`
	WatchDogDelay  int    `json:"WatchDogDelay"`
}

// Pair is a pair of dia assets.
type Pair struct {
	QuoteToken Asset `json:"QuoteToken"`
	BaseToken  Asset `json:"BaseToken"`
}

func (p *Pair) ExchangePairIdentifier(exchange string) string {
	return exchange + "-" + p.Identifier()
}

func (p *Pair) Identifier() string {
	return p.QuoteToken.Blockchain + "-" + p.QuoteToken.Address + "-" + p.BaseToken.Blockchain + "-" + p.BaseToken.Address
}

// ExchangePairsFromEnv parses the string @exchangePairsEnv consisting of pairs on exchanges
// and returns full asset information for the corresponding exchangepairs.
// It assumes mappings can be found in the files exchange.json at @configPath where
// @exchange is the corresponding exchange name.
func ExchangePairsFromEnv(
	exchangePairsEnv string,
	envSeparator string,
	exchangePairSeparator string,
	pairTickerSeparator string,
	configPath string,
) (exchangePairs []ExchangePair) {

	// epMap maps an exchange on a slice of the underlying pair symbol tickers.
	epMap := make(map[string][]string)
	list := strings.Split(exchangePairsEnv, envSeparator)
	if len(list) == 0 {
		return
	}
	if len(list) == 1 && len(strings.TrimSpace(list[0])) == 0 {
		return
	}

	for _, ep := range list {
		exchange := strings.TrimSpace(strings.Split(ep, exchangePairSeparator)[0])
		pairSymbol := strings.TrimSpace(strings.Split(ep, exchangePairSeparator)[1])
		epMap[exchange] = append(epMap[exchange], pairSymbol)
	}

	// Assign assets to pair symbols.
	for exchange := range epMap {
		symbolIdentificationMap, err := GetSymbolIdentificationMap(exchange, configPath)
		if err != nil {
			log.Fatal("GetSymbolIdentificationMap: ", err)
		}
		for _, pairSymbol := range epMap[exchange] {
			symbols := strings.Split(pairSymbol, pairTickerSeparator)
			var ep ExchangePair
			ep.Exchange = exchange
			ep.ForeignName = pairSymbol
			ep.Symbol = symbols[0]
			ep.UnderlyingPair.QuoteToken = symbolIdentificationMap[ExchangeSymbolIdentifier(symbols[0], exchange)]
			ep.UnderlyingPair.BaseToken = symbolIdentificationMap[ExchangeSymbolIdentifier(symbols[1], exchange)]
			exchangePairs = append(exchangePairs, ep)
		}
	}
	return
}

func ExchangePairsFromPairs(
	epMap map[string][]string,
	pairTickerSeparator string,
	configPath string,
	watchdog map[string]int,
) (exchangePairs []ExchangePair) {
	for exchange, symbolsList := range epMap {
		symbolIdentificationMap, err := GetSymbolIdentificationMap(exchange, configPath)
		if err != nil {
			log.Fatal("GetSymbolIdentificationMap: ", err)
		}
		for _, pairSymbol := range symbolsList {
			parts := strings.Split(pairSymbol, pairTickerSeparator)
			if len(parts) < 2 {
				log.Warnf("Invalid pair format: %s", pairSymbol)
				continue
			}
			quote := strings.TrimSpace(parts[0])
			base := strings.TrimSpace(parts[1])

			var ep ExchangePair
			ep.Exchange = exchange
			ep.ForeignName = pairSymbol
			ep.Symbol = quote

			qID := ExchangeSymbolIdentifier(quote, exchange)
			bID := ExchangeSymbolIdentifier(base, exchange)

			qAsset, okQ := symbolIdentificationMap[qID]

			if !okQ {
				qAsset = Asset{Symbol: quote}
				log.Warnf("[%s] missing quote asset mapping for %s", exchange, quote)
			}

			bAsset, okB := symbolIdentificationMap[bID]
			if !okB {
				bAsset = Asset{Symbol: base}
				log.Warnf("[%s] missing base asset mapping for %s", exchange, base)
			}

			ep.UnderlyingPair.QuoteToken = qAsset
			ep.UnderlyingPair.BaseToken = bAsset
			if watchdog != nil {
				if wd, ok := watchdog[exchange+":"+pairSymbol]; ok {
					ep.WatchDogDelay = wd
				}
			}
			exchangePairs = append(exchangePairs, ep)
		}
	}
	return exchangePairs
}

// MakeExchangepairMap returns a map in which exchangepairs are grouped by exchange string key.
func MakeExchangepairMap(exchangePairs []ExchangePair) map[string][]ExchangePair {
	exchangepairMap := make(map[string][]ExchangePair)
	for _, ep := range exchangePairs {
		exchangepairMap[ep.Exchange] = append(exchangepairMap[ep.Exchange], ep)
	}
	return exchangepairMap
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
