package models

import (
	"testing"
	"time"
)

func TestGetLastTrade(t *testing.T) {
	cases := []struct {
		tradesblock TradesBlock
		lastTrade   Trade
	}{
		{
			tradesblock: TradesBlock{
				Trades: []Trade{
					{
						Time: time.Unix(1721209858, 0),
					},
					{
						Time: time.Unix(1657961611, 0),
					},
					{
						Time: time.Unix(1689497611, 0),
					},
				},
			},
			lastTrade: Trade{Time: time.Unix(1721209858, 0)},
		},
		{
			tradesblock: TradesBlock{
				Trades: []Trade{
					{
						Time: time.Unix(0, 0),
					},
				},
			},
			lastTrade: Trade{Time: time.Unix(0, 0)},
		},
	}

	for i, c := range cases {
		_, lastTrade := c.tradesblock.GetLastTrade()
		if lastTrade != c.lastTrade {
			t.Errorf("Trade was incorrect, got: %v, expected: %v for set:%d", lastTrade, c.lastTrade, i)
		}
	}
}

func TestMergeTradesBlocksByPair(t *testing.T) {
	// Assets need Blockchain + Address for Pair.Identifier() to distinguish pairs.
	btc := Asset{Symbol: "BTC", Blockchain: "Ethereum", Address: "0x1"}
	usdt := Asset{Symbol: "USDT", Blockchain: "Ethereum", Address: "0x2"}
	eth := Asset{Symbol: "ETH", Blockchain: "Ethereum", Address: "0x3"}

	btcPair := Pair{QuoteToken: btc, BaseToken: usdt}
	ethPair := Pair{QuoteToken: eth, BaseToken: usdt}

	binance := Exchange{Name: "Binance"}
	kucoin := Exchange{Name: "KuCoin"}
	uniswap := Exchange{Name: "Uniswap"}

	t0 := time.Unix(1000, 0)
	t1 := time.Unix(1100, 0)
	t2 := time.Unix(1200, 0)
	t3 := time.Unix(1300, 0)

	input := map[string]TradesBlock{
		"Binance-BTC": {
			Pair:      btcPair,
			StartTime: t0,
			EndTime:   t1,
			Atomic:    true,
			Trades: []Trade{
				{QuoteToken: btc, BaseToken: usdt, Price: 50000, Volume: 1, Exchange: binance},
			},
		},
		"KuCoin-BTC": {
			Pair:      btcPair,
			StartTime: t1,
			EndTime:   t2,
			Atomic:    true,
			Trades: []Trade{
				{QuoteToken: btc, BaseToken: usdt, Price: 50100, Volume: 2, Exchange: kucoin},
			},
		},
		"Uniswap-ETH": {
			Pair:      ethPair,
			StartTime: t2,
			EndTime:   t3,
			Atomic:    true,
			Trades: []Trade{
				{QuoteToken: eth, BaseToken: usdt, Price: 3000, Volume: 5, Exchange: uniswap},
			},
		},
	}

	result := MergeTradesBlocksByPair(input)

	// Should produce 2 keys: one for BTC pair, one for ETH pair.
	if len(result) != 2 {
		t.Fatalf("expected 2 merged blocks, got %d", len(result))
	}

	btcKey := btcPair.Identifier()
	ethKey := ethPair.Identifier()

	// --- BTC block ---
	btcBlock, ok := result[btcKey]
	if !ok {
		t.Fatalf("expected BTC block in result, not found (key=%s)", btcKey)
	}
	if btcBlock.Atomic {
		t.Error("merged BTC block should have Atomic=false")
	}
	if len(btcBlock.Trades) != 2 {
		t.Errorf("expected 2 BTC trades, got %d", len(btcBlock.Trades))
	}
	if !btcBlock.StartTime.Equal(t0) {
		t.Errorf("BTC StartTime: expected %v, got %v", t0, btcBlock.StartTime)
	}
	if !btcBlock.EndTime.Equal(t2) {
		t.Errorf("BTC EndTime: expected %v, got %v", t2, btcBlock.EndTime)
	}

	// --- ETH block ---
	ethBlock, ok := result[ethKey]
	if !ok {
		t.Fatalf("expected ETH block in result, not found (key=%s)", ethKey)
	}
	if ethBlock.Atomic {
		t.Error("merged ETH block should have Atomic=false")
	}
	if len(ethBlock.Trades) != 1 {
		t.Errorf("expected 1 ETH trade, got %d", len(ethBlock.Trades))
	}

	// --- Original input not modified ---
	if !input["Binance-BTC"].Atomic {
		t.Error("MergeTradesBlocksByPair should not modify original input blocks")
	}
}

func TestMergeTradesBlocksByPair_EmptyInput(t *testing.T) {
	result := MergeTradesBlocksByPair(map[string]TradesBlock{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestMergeTradesBlocksByPair_SingleBlock(t *testing.T) {
	btc := Asset{Symbol: "BTC", Blockchain: "Ethereum", Address: "0x1"}
	usdt := Asset{Symbol: "USDT", Blockchain: "Ethereum", Address: "0x2"}
	pair := Pair{QuoteToken: btc, BaseToken: usdt}

	input := map[string]TradesBlock{
		"Binance-BTC": {
			Pair:   pair,
			Atomic: true,
			Trades: []Trade{{Price: 50000, Volume: 1}},
		},
	}

	result := MergeTradesBlocksByPair(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 block, got %d", len(result))
	}
	block := result[pair.Identifier()]
	if block.Atomic {
		t.Error("expected Atomic=false even for single-exchange merge")
	}
	if len(block.Trades) != 1 {
		t.Errorf("expected 1 trade, got %d", len(block.Trades))
	}
}
