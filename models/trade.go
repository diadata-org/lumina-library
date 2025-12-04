package models

import (
	"errors"
	"time"
)

type Trade struct {
	QuoteToken     Asset
	BaseToken      Asset
	Price          float64
	Volume         float64
	Time           time.Time
	Exchange       Exchange
	PoolAddress    string
	ForeignTradeID string
	// Depending on the connection to the processing layer we might not need it here.
	EstimatedUSDPrice float64
}

// SwapTrade swaps base and quote token of a trade and inverts the price accordingly
func (t *Trade) SwapTrade() error {
	if t.Price == 0 {
		return errors.New("zero price. cannot swap trade")
	}
	t.BaseToken, t.QuoteToken = t.QuoteToken, t.BaseToken
	t.Volume = -t.Price * t.Volume
	t.Price = 1 / t.Price

	return nil
}

// Struct for decentralized scraper pools.
// TO DO: Revisit fields.
type TradesBlock struct {
	// Add field for Asset? So far, we only consider atomic tradesblocks.
	// Asset Asset
	Pair      Pair
	Trades    []Trade
	StartTime time.Time
	EndTime   time.Time
	// A tradesblock is atomic if trades all are from the same exchangepair.
	Atomic bool
	// Do we need this?
	ScraperID ScraperID
}

// ScraperID is the container identifying a scraper node.
type ScraperID struct {
	// ID could for instance be evm address.
	ID string
	// Human readable name of the entity that is running the scraper.
	Name             string
	RegistrationTime time.Time
}

// IsAtomic determines whether a tradesblock is atomic by looking at all trades.
func (tb TradesBlock) IsAtomic() bool {
	if len(tb.Trades) == 0 {
		return true
	}

	source := tb.Trades[0].Exchange.Name
	pair := Pair{QuoteToken: tb.Trades[0].QuoteToken, BaseToken: tb.Trades[0].BaseToken}

	for _, trade := range tb.Trades {
		if trade.Exchange.Name != source {
			return false
		}
		if (Pair{QuoteToken: trade.QuoteToken, BaseToken: trade.BaseToken}) != pair {
			return false
		}
	}
	return true
}

func (tb TradesBlock) GetSourceType() (SourceType, error) {
	if !tb.IsAtomic() {
		return SourceType(""), errors.New("block is not atomic")
	}
	if len(tb.Trades) == 0 {
		return SourceType(""), nil
	}
	return GetSourceType(tb.Trades[0].Exchange), nil
}

// GetLastTrade returns the latest trade from the slice @trades.
func (tb TradesBlock) GetLastTrade() (index int, lastTrade Trade) {
	for i, trade := range tb.Trades {
		if trade.Time.After(lastTrade.Time) {
			lastTrade = trade
			index = i
		}
	}
	return
}

func (tb *TradesBlock) RemoveTradeByIndex(index int) {
	tb.Trades = append(tb.Trades[:index], tb.Trades[index+1:]...)
}

// Transforms a @SimulatedTrade to a @Trade type so functions can be reused.
func SimulatedTradeToTrade(st SimulatedTrade) Trade {
	return Trade{
		BaseToken:         st.BaseToken,
		QuoteToken:        st.QuoteToken,
		Price:             st.Price,
		Volume:            st.Volume,
		Time:              st.Time,
		Exchange:          st.Exchange,
		PoolAddress:       st.PoolAddress,
		ForeignTradeID:    st.TXHash,
		EstimatedUSDPrice: st.EstimatedUSDPrice,
	}
}

func SimulatedTradesBlockToTradesBlock(stb SimulatedTradesBlock) (tb TradesBlock) {
	tb.Pair = stb.Pair
	for _, simulatedTrade := range stb.Trades {
		trade := SimulatedTradeToTrade(simulatedTrade)
		tb.Trades = append(tb.Trades, trade)
	}
	tb.StartTime = stb.StartTime
	tb.EndTime = stb.EndTime
	tb.Atomic = stb.Atomic
	tb.ScraperID = stb.ScraperID
	return
}
