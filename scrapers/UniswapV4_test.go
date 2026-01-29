package scrapers

import (
	"math"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/models"

	"github.com/ethereum/go-ethereum/common"
)

func TestParsePoolIDHex_OK(t *testing.T) {
	// 32 bytes hex (64 chars) with 0x prefix
	in := "0x" + repeatHex("ab", 32) // 64 chars => 32 bytes
	id, err := ParsePoolIDHex(in)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if PoolIDToHex(id) != "0x"+repeatHex("ab", 32) {
		t.Fatalf("hex mismatch: %s", PoolIDToHex(id))
	}
}

func TestParsePoolIDHex_BadLen(t *testing.T) {
	_, err := ParsePoolIDHex("0x1234")
	if err == nil {
		t.Fatalf("expected error for bad length")
	}
}

func TestGetTokenAddrsFromPoolConfig_OK(t *testing.T) {
	pool := models.Pool{
		Assetvolumes: []models.AssetVolume{
			{Asset: models.Asset{Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"}}, // WETH
			{Asset: models.Asset{Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"}}, // USDC
		},
	}
	a0, a1, err := getTokenAddrsFromPoolConfig(pool)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if a0 != common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") {
		t.Fatalf("token0 mismatch: %s", a0.Hex())
	}
	if a1 != common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") {
		t.Fatalf("token1 mismatch: %s", a1.Hex())
	}
}

func TestGetTokenAddrsFromPoolConfig_TooShort(t *testing.T) {
	pool := models.Pool{
		Assetvolumes: []models.AssetVolume{
			{Asset: models.Asset{Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"}},
		},
	}
	_, _, err := getTokenAddrsFromPoolConfig(pool)
	if err == nil {
		t.Fatalf("expected error for <2 assetvolumes")
	}
}

func TestGetSwapDataUniV4_PriceAndVolume(t *testing.T) {
	swap := UniswapV4Swap{
		ID:      "tx:1",
		Amount0: 1.0,
		Amount1: -3000.0,
	}

	price, vol := getSwapDataUniV4(swap)

	// volume = Amount0, price = abs(amount1/amount0)
	if !almostEqualUniV4(vol, 1.0, 1e-12) {
		t.Fatalf("volume mismatch: got=%v want=1", vol)
	}
	if !almostEqualUniV4(price, 3000.0, 1e-9) {
		t.Fatalf("price mismatch: got=%v want=3000", price)
	}
}

func TestMakeTradeUniV4_Order1_SwapTradeChangesDirection(t *testing.T) {
	// makeTrade uses BaseToken=Token1, QuoteToken=Token0,
	// and Order=1 will call t.SwapTrade() to flip.
	pair := UniV4Pair{
		Token0:  models.Asset{Symbol: "WETH", Decimals: 18, Blockchain: "ethereum"},
		Token1:  models.Asset{Symbol: "USDC", Decimals: 6, Blockchain: "ethereum"},
		PoolHex: "0x" + repeatHex("aa", 32),
		Order:   1,
	}

	t0 := makeTradeUniV4(pair, 3000, 1.0, time.Unix(100, 0), "tx:1", "UniswapV4", "ethereum")
	if t0.BaseToken.Symbol != "USDC" || t0.QuoteToken.Symbol != "WETH" {
		t.Fatalf("unexpected initial direction base=%s quote=%s", t0.BaseToken.Symbol, t0.QuoteToken.Symbol)
	}

	t0.SwapTrade()
	if t0.BaseToken.Symbol != "WETH" || t0.QuoteToken.Symbol != "USDC" {
		t.Fatalf("unexpected swapped direction base=%s quote=%s", t0.BaseToken.Symbol, t0.QuoteToken.Symbol)
	}
}

// ---------------- helpers ----------------

func almostEqualUniV4(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

// repeatHex repeats a 2-char hex byte string n times; e.g. repeatHex("ab",32) => 64 chars.
func repeatHex(twoHex string, n int) string {
	out := make([]byte, 0, 2*n)
	for i := 0; i < n; i++ {
		out = append(out, twoHex[0], twoHex[1])
	}
	return string(out)
}
