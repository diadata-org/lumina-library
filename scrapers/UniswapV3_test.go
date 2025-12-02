package scrapers

import (
	"math"
	"math/big"
	"testing"
	"time"

	PancakeswapV3Pair "github.com/diadata-org/lumina-library/contracts/pancakeswapv3"
	UniswapV3Pair "github.com/diadata-org/lumina-library/contracts/uniswapv3/uniswapV3Pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// helper: construct 10^n
func pow10Int(n int64) *big.Int {
	base := big.NewInt(10)
	exp := big.NewInt(n)
	return new(big.Int).Exp(base, exp, nil)
}

// helper: float comparison
func almostEqualUniV3(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

// ---------------------------------------------------------
// 1) normalizeUniV3Swap - UniswapV3PairSwap
// ---------------------------------------------------------

func TestNormalizeUniV3Swap_UniswapV3(t *testing.T) {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000123")

	// Token0: 18 decimals, Token1: 6 decimals
	token0 := models.Asset{Symbol: "T0", Decimals: 18}
	token1 := models.Asset{Symbol: "T1", Decimals: 6}

	// Divisor0 = 10^18, Divisor1 = 10^6
	div0 := new(big.Float).SetInt(pow10Int(18))
	div1 := new(big.Float).SetInt(pow10Int(6))

	pair := UniV3Pair{
		Token0:   token0,
		Token1:   token1,
		Address:  addr,
		Order:    0,
		Divisor0: div0,
		Divisor1: div1,
	}

	s := &UniswapV3Scraper{
		poolMap: map[common.Address]UniV3Pair{
			addr: pair,
		},
	}

	// Amount0 = 2 * 10^18 -> 2.0
	// Amount1 = 5 * 10^6  -> 5.0
	amt0 := new(big.Int).Mul(big.NewInt(2), pow10Int(18))
	amt1 := new(big.Int).Mul(big.NewInt(5), pow10Int(6))

	swap := UniswapV3Pair.UniswapV3PairSwap{
		Amount0: amt0,
		Amount1: amt1,
		Raw: types.Log{
			Address: addr,
			TxHash:  common.HexToHash("0xabc"),
		},
	}

	normalized := s.normalizeUniV3Swap(swap)

	if normalized.ID != swap.Raw.TxHash.Hex() {
		t.Fatalf("expected ID %s, got %s", swap.Raw.TxHash.Hex(), normalized.ID)
	}
	if normalized.Pair.Address != addr {
		t.Fatalf("expected pair address %s, got %s", addr.Hex(), normalized.Pair.Address.Hex())
	}

	if !almostEqualUniV3(normalized.Amount0, 2.0, 1e-9) {
		t.Fatalf("expected Amount0 ~ 2.0, got %f", normalized.Amount0)
	}
	if !almostEqualUniV3(normalized.Amount1, 5.0, 1e-9) {
		t.Fatalf("expected Amount1 ~ 5.0, got %f", normalized.Amount1)
	}
}

// ---------------------------------------------------------
// 2) normalizeUniV3Swap - Pancakev3pairSwap
// ---------------------------------------------------------

func TestNormalizeUniV3Swap_PancakeswapV3(t *testing.T) {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000456")

	token0 := models.Asset{Symbol: "P0", Decimals: 18}
	token1 := models.Asset{Symbol: "P1", Decimals: 6}

	div0 := new(big.Float).SetInt(pow10Int(18))
	div1 := new(big.Float).SetInt(pow10Int(6))

	pair := UniV3Pair{
		Token0:   token0,
		Token1:   token1,
		Address:  addr,
		Order:    0,
		Divisor0: div0,
		Divisor1: div1,
	}

	s := &UniswapV3Scraper{
		poolMap: map[common.Address]UniV3Pair{
			addr: pair,
		},
	}

	// Amount0 = 3 * 10^18 -> 3.0
	// Amount1 = 7 * 10^6  -> 7.0
	amt0 := new(big.Int).Mul(big.NewInt(3), pow10Int(18))
	amt1 := new(big.Int).Mul(big.NewInt(7), pow10Int(6))

	swap := PancakeswapV3Pair.Pancakev3pairSwap{
		Amount0: amt0,
		Amount1: amt1,
		Raw: types.Log{
			Address: addr,
			TxHash:  common.HexToHash("0xdef"),
		},
	}

	normalized := s.normalizeUniV3Swap(swap)

	if normalized.ID != swap.Raw.TxHash.Hex() {
		t.Fatalf("expected ID %s, got %s", swap.Raw.TxHash.Hex(), normalized.ID)
	}
	if normalized.Pair.Address != addr {
		t.Fatalf("expected pair address %s, got %s", addr.Hex(), normalized.Pair.Address.Hex())
	}

	if !almostEqualUniV3(normalized.Amount0, 3.0, 1e-9) {
		t.Fatalf("expected Amount0 ~ 3.0, got %f", normalized.Amount0)
	}
	if !almostEqualUniV3(normalized.Amount1, 7.0, 1e-9) {
		t.Fatalf("expected Amount1 ~ 7.0, got %f", normalized.Amount1)
	}
}

// ---------------------------------------------------------
// 3) getSwapDataUniV3
// ---------------------------------------------------------

func TestGetSwapDataUniV3_NormalCase(t *testing.T) {
	swap := UniswapV3Swap{
		ID:      "tx1",
		Amount0: 2.0,
		Amount1: 10.0,
	}

	price, volume := getSwapDataUniV3(swap)

	if volume != 2.0 {
		t.Fatalf("expected volume 2.0, got %f", volume)
	}
	if price != 5.0 {
		t.Fatalf("expected price 5.0 (=10/2), got %f", price)
	}
}

func TestGetSwapDataUniV3_ZeroAmount0(t *testing.T) {
	swap := UniswapV3Swap{
		ID:      "tx_zero",
		Amount0: 0.0,
		Amount1: 10.0,
	}

	price, volume := getSwapDataUniV3(swap)

	if price != 0 || volume != 0 {
		t.Fatalf("expected (0,0) when Amount0=0, got (price=%f, volume=%f)", price, volume)
	}
}

// ---------------------------------------------------------
// 4) makeTrade
// ---------------------------------------------------------

func TestMakeTrade(t *testing.T) {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000abc")
	token0 := models.Asset{Symbol: "QUOTE", Decimals: 18}
	token1 := models.Asset{Symbol: "BASE", Decimals: 18}

	pair := UniV3Pair{
		Token0:  token0,
		Token1:  token1,
		Address: addr,
		Order:   0,
	}

	now := time.Now().Truncate(time.Second)
	price := 123.45
	volume := 6.78
	txID := "0xtxid"

	tr := makeTrade(
		pair,
		price,
		volume,
		now,
		addr,
		txID,
		"UniswapV3",
		"ethereum",
	)

	if tr.Price != price {
		t.Fatalf("expected price %f, got %f", price, tr.Price)
	}
	if tr.Volume != volume {
		t.Fatalf("expected volume %f, got %f", volume, tr.Volume)
	}
	if tr.BaseToken.Symbol != "BASE" || tr.QuoteToken.Symbol != "QUOTE" {
		t.Fatalf("unexpected tokens: base=%s, quote=%s", tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
	}
	if tr.PoolAddress != addr.Hex() {
		t.Fatalf("expected pool address %s, got %s", addr.Hex(), tr.PoolAddress)
	}
	if tr.ForeignTradeID != txID {
		t.Fatalf("expected ForeignTradeID %s, got %s", txID, tr.ForeignTradeID)
	}
	if tr.Exchange.Name != "UniswapV3" || tr.Exchange.Blockchain != "ethereum" {
		t.Fatalf("unexpected exchange: %+v", tr.Exchange)
	}
}
