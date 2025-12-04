package models

import (
	"testing"
	"time"
)

func testGetLastTrade(t *testing.T) {
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
