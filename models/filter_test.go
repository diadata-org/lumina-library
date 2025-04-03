package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/utils"
)

// TO DO: Write more test cases.

var (
	ETH  = Asset{Address: "0x0000000000000000000000000000000000000000", Blockchain: utils.ETHEREUM}
	BTC  = Asset{Address: "0x0000000000000000000000000000000000000000", Blockchain: utils.BITCOIN}
	USDC = Asset{Address: "", Blockchain: utils.ETHEREUM}
	fpe1 = []FilterPointPair{
		{
			Pair:  Pair{QuoteToken: ETH, BaseToken: USDC},
			Value: 3388.34,
			Time:  time.Unix(1657961611, 0),
		},
		{
			Pair:  Pair{QuoteToken: ETH, BaseToken: USDC},
			Value: 3381.11,
			Time:  time.Unix(1689497611, 0),
		},
		{
			Pair:  Pair{QuoteToken: ETH, BaseToken: USDC},
			Value: 3179.78,
			Time:  time.Unix(1706846011, 0),
		},
	}

	fpe22 = FilterPointPair{Pair: Pair{QuoteToken: BTC, BaseToken: USDC}, Value: 63199.11, Time: time.Unix(0, 0)}
)

func testGroupFilterByAsset(t *testing.T) {
	map1 := make(map[Asset][]FilterPointPair)
	map2 := make(map[Asset][]FilterPointPair)

	map1[ETH] = fpe1

	fpe2 := fpe1
	fpe2 = append(fpe2, fpe22)
	map2[ETH] = fpe2
	map2[BTC] = []FilterPointPair{fpe22}

	cases := []struct {
		filterPoints []FilterPointPair
		fpMap        map[Asset][]FilterPointPair
	}{
		{
			filterPoints: fpe1,
			fpMap:        map1,
		},
		{
			filterPoints: fpe2,
			fpMap:        map2,
		},
	}

	for i, c := range cases {
		filterPointMap := GroupFiltersByAsset(c.filterPoints)
		if !reflect.DeepEqual(filterPointMap, c.fpMap) {
			t.Errorf("Map was incorrect, got: %v, expected: %v for set:%d", filterPointMap, c.fpMap, i)
		}

	}

}

func testGetValuesFromFilterPoints(t *testing.T) {

	fpe2 := fpe1
	fpe2 = append(fpe2, fpe22)

	cases := []struct {
		filterPoints []FilterPointPair
		filterValues []float64
	}{
		{
			fpe1,
			[]float64{3388.34, 3381.11, 3179.78},
		},
		{
			fpe2,
			[]float64{3388.34, 3381.11, 3179.78, 63199.11},
		},
	}

	for i, c := range cases {
		filterValues := GetValuesFromFilterPoints(c.filterPoints)
		if !reflect.DeepEqual(filterValues, c.filterValues) {
			t.Errorf("Values slice was incorrect, got: %v, expected: %v for set:%d", filterValues, c.filterValues, i)
		}
	}

}

func testGetLatestTimestampFromFilterPoints(t *testing.T) {
	cases := []struct {
		filterPoints []FilterPointPair
		timestamp    time.Time
	}{
		{
			fpe1,
			time.Unix(1706846011, 0),
		},
		{
			[]FilterPointPair{fpe22},
			time.Unix(0, 0),
		},
	}

	for i, c := range cases {
		latestTimestamp := GetLatestTimestampFromFilterPoints(c.filterPoints)
		if latestTimestamp.Unix() != c.timestamp.Unix() {
			t.Errorf("Timestamp was incorrect, got: %v, expected: %v for set:%d", latestTimestamp, c.timestamp, i)
		}
	}
}

func testRemoveOldFilters(t *testing.T) {
	cases := []struct {
		filterPoints        []FilterPointPair
		toleranceSeconds    int64
		timestamp           time.Time
		cleanedFilterPoints []FilterPointPair
		removedFilters      int
	}{
		{
			filterPoints:     fpe1,
			toleranceSeconds: int64(4),
			timestamp:        time.Unix(1689497606, 0),
			cleanedFilterPoints: []FilterPointPair{
				{
					Pair:  Pair{QuoteToken: ETH, BaseToken: USDC},
					Value: 3381.11,
					Time:  time.Unix(1689497611, 0),
				},
				{
					Pair:  Pair{QuoteToken: ETH, BaseToken: USDC},
					Value: 3179.78,
					Time:  time.Unix(1706846011, 0),
				},
			},
			removedFilters: 1,
		},
		{
			filterPoints:        fpe1,
			toleranceSeconds:    int64(4),
			timestamp:           time.Unix(1706846017, 0),
			cleanedFilterPoints: []FilterPointPair{},
			removedFilters:      3,
		},
	}

	for i, c := range cases {

		cleanedFilterPoints, removedFilters := RemoveOldFilters(c.filterPoints, c.toleranceSeconds, c.timestamp)
		if !reflect.DeepEqual(cleanedFilterPoints, c.cleanedFilterPoints) {
			t.Errorf("Cleaned filters was incorrect, got: %v, expected: %v for set:%d", cleanedFilterPoints, c.cleanedFilterPoints, i)
		}
		if removedFilters != c.removedFilters {
			t.Errorf("Number of removed filters was incorrect, got: %v, expected: %v for set:%d", removedFilters, c.removedFilters, i)
		}
	}
}
