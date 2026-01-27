package scrapers

import (
	"math"
	"math/big"
	"testing"
	"time"

	poolmanager "github.com/diadata-org/lumina-library/contracts/uniswapv4/poolManager"
	"github.com/diadata-org/lumina-library/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestParsePoolIDHex_ValidAndInvalid(t *testing.T) {
	valid := "0x" + repeatHex("11", 32) // 32 bytes => 64 hex chars
	id, err := ParsePoolIDHex(valid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if PoolIDToHex(id) != valid {
		t.Fatalf("roundtrip mismatch: got %s want %s", PoolIDToHex(id), valid)
	}

	// missing 0x is ok in your parser (it trims prefix)
	validNo0x := repeatHex("22", 32)
	id2, err := ParsePoolIDHex(validNo0x)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if PoolIDToHex(id2) != "0x"+validNo0x {
		t.Fatalf("roundtrip mismatch: got %s want %s", PoolIDToHex(id2), "0x"+validNo0x)
	}

	// wrong length
	_, err = ParsePoolIDHex("0x1234")
	if err == nil {
		t.Fatalf("expected error for short poolId")
	}

	// non-hex
	_, err = ParsePoolIDHex("0x" + repeat("zz", 32))
	if err == nil {
		t.Fatalf("expected error for non-hex input")
	}
}

func TestNormalizeUniV4Swap_ScalesAndKeepsSign(t *testing.T) {
	// USDC (6 decimals), WETH (18 decimals)
	usdc := models.Asset{Symbol: "USDC", Address: "0xA0b8...", Decimals: 6, Blockchain: "ethereum"}
	weth := models.Asset{Symbol: "WETH", Address: "0xC02a...", Decimals: 18, Blockchain: "ethereum"}

	poolHex := "0x" + repeatHex("ab", 32)
	poolID, err := ParsePoolIDHex(poolHex)
	if err != nil {
		t.Fatalf("ParsePoolIDHex: %v", err)
	}

	pair := UniV4Pair{
		Token0:   weth,
		Token1:   usdc,
		PoolID:   poolID,
		PoolHex:  poolHex,
		AddrKey:  common.BytesToAddress(poolID[12:]),
		Order:    0,
		Divisor0: new(big.Float).SetFloat64(math.Pow10(int(weth.Decimals))), // 1e18
		Divisor1: new(big.Float).SetFloat64(math.Pow10(int(usdc.Decimals))), // 1e6
	}

	// amount0 = -1 WETH, amount1 = +3000 USDC (signed int)
	amt0 := new(big.Int).Neg(new(big.Int).Mul(big.NewInt(1), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)))
	amt1 := new(big.Int).Mul(big.NewInt(3000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(6), nil))

	ev := poolmanager.PoolManagerSwap{
		Amount0: amt0,
		Amount1: amt1,
		Raw: types.Log{
			TxHash: common.HexToHash("0x" + repeatHex("12", 32)),
			Index:  7,
		},
	}

	s := &UniswapV4Scraper{}
	n := s.normalizeUniV4Swap(ev, pair)

	if n.ID == "" {
		t.Fatalf("expected non-empty ID")
	}
	if n.ID != ev.Raw.TxHash.Hex()+":7" {
		t.Fatalf("expected id to include log index: got %s", n.ID)
	}

	// scaled:
	// amount0 = -1
	// amount1 = 3000
	if !floatAlmostEq(n.Amount0, -1.0, 1e-12) {
		t.Fatalf("amount0 scale/sign wrong: got %v want %v", n.Amount0, -1.0)
	}
	if !floatAlmostEq(n.Amount1, 3000.0, 1e-9) {
		t.Fatalf("amount1 scale/sign wrong: got %v want %v", n.Amount1, 3000.0)
	}
}

func TestGetSwapDataUniV4_PriceVolume(t *testing.T) {
	swap := UniswapV4Swap{
		ID:      "x",
		Amount0: -1.0,
		Amount1: 3000.0,
	}

	price, vol := getSwapDataUniV4(swap)
	// your logic:
	// volume = Amount0
	// price = abs(Amount1 / Amount0)
	if !floatAlmostEq(vol, -1.0, 1e-12) {
		t.Fatalf("volume wrong: got %v want %v", vol, -1.0)
	}
	if !floatAlmostEq(price, 3000.0, 1e-9) {
		t.Fatalf("price wrong: got %v want %v", price, 3000.0)
	}

	// Amount0==0 => should return 0,0
	swap.Amount0 = 0
	price, vol = getSwapDataUniV4(swap)
	if price != 0 || vol != 0 {
		t.Fatalf("expected 0,0 for Amount0==0, got %v,%v", price, vol)
	}
}

func TestMakeTradeUniV4_MappingMatchesUniV3Style(t *testing.T) {
	// This test asserts your makeTradeUniV4 does the same mapping as UniV3:
	// BaseToken = token1, QuoteToken = token0
	token0 := models.Asset{Symbol: "WETH", Decimals: 18, Blockchain: "ethereum"}
	token1 := models.Asset{Symbol: "USDC", Decimals: 6, Blockchain: "ethereum"}

	poolHex := "0x" + repeatHex("cd", 32)
	poolID, _ := ParsePoolIDHex(poolHex)

	pair := UniV4Pair{
		Token0:  token0,
		Token1:  token1,
		PoolID:  poolID,
		PoolHex: poolHex,
	}

	now := time.Now()
	tr := makeTradeUniV4(pair, 3000.0, 1.0, now, "tx:1", "UniswapV4", "ethereum")

	if tr.QuoteToken.Symbol != "WETH" {
		t.Fatalf("QuoteToken should be token0(WETH), got %s", tr.QuoteToken.Symbol)
	}
	if tr.BaseToken.Symbol != "USDC" {
		t.Fatalf("BaseToken should be token1(USDC), got %s", tr.BaseToken.Symbol)
	}
	if tr.PoolAddress != poolHex {
		t.Fatalf("PoolAddress should be poolHex(poolId), got %s", tr.PoolAddress)
	}
	if tr.Exchange.Name != "UniswapV4" || tr.Exchange.Blockchain != "ethereum" {
		t.Fatalf("Exchange mismatch: %+v", tr.Exchange)
	}
	if tr.ForeignTradeID != "tx:1" {
		t.Fatalf("ForeignTradeID mismatch: %s", tr.ForeignTradeID)
	}
	if tr.Time.Unix() != now.Unix() {
		t.Fatalf("Time mismatch: got %v want %v", tr.Time, now)
	}
}

// ---------- helpers ----------

func floatAlmostEq(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func repeatHex(twoHex string, nBytes int) string {
	// twoHex like "ab" repeated nBytes times => 2*nBytes hex chars
	out := ""
	for i := 0; i < nBytes; i++ {
		out += twoHex
	}
	return out
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
