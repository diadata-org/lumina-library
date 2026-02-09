package scrapers

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	velodrome "github.com/diadata-org/lumina-library/contracts/velodrome"
	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestGetSwapDataAero_BuyToken0Branch(t *testing.T) {
	// amount0In == 0 => volume=amount0Out, price=amount1In/amount0Out
	swap := AerodromeV1Swap{
		Amount0In:  0,
		Amount0Out: 2,
		Amount1In:  10,
		Amount1Out: 0,
	}
	price, volume := getSwapDataAero(swap)

	if math.Abs(price-5.0) > 1e-12 {
		t.Fatalf("price mismatch: got %v want %v", price, 5.0)
	}
	if math.Abs(volume-2.0) > 1e-12 {
		t.Fatalf("volume mismatch: got %v want %v", volume, 2.0)
	}
}

func TestGetSwapDataAero_SellToken0Branch(t *testing.T) {
	// amount0In > 0 => volume=-amount0In, price=amount1Out/amount0In
	swap := AerodromeV1Swap{
		Amount0In:  3,
		Amount0Out: 0,
		Amount1In:  0,
		Amount1Out: 12,
	}
	price, volume := getSwapDataAero(swap)

	if math.Abs(price-4.0) > 1e-12 {
		t.Fatalf("price mismatch: got %v want %v", price, 4.0)
	}
	if math.Abs(volume-(-3.0)) > 1e-12 {
		t.Fatalf("volume mismatch: got %v want %v", volume, -3.0)
	}
}

func TestGetSwapDataAero_InvalidSwap(t *testing.T) {
	// amount0In==0 && amount0Out==0 => invalid
	swap := AerodromeV1Swap{
		Amount0In:  0,
		Amount0Out: 0,
		Amount1In:  123,
		Amount1Out: 456,
	}
	price, volume := getSwapDataAero(swap)
	if price != 0 || volume != 0 {
		t.Fatalf("expected (0,0), got (%v,%v)", price, volume)
	}
}

func TestMakeTradeAerodromeV1_MapsFieldsCorrectly(t *testing.T) {
	addr := common.HexToAddress("0x7f670f78b17dec44d5ef68a48740b6f8849cc2e6")
	pair := AerodromeV1Pair{
		Token0:  models.Asset{Symbol: "WETH", Decimals: 18},
		Token1:  models.Asset{Symbol: "AERO", Decimals: 18},
		Address: addr,
	}
	ts := time.Unix(1700000000, 0)

	tr := makeTradeAerodromeV1(pair, 123.45, -6.7, ts, addr, "0xtx", "AerodromeV1", "base")

	if tr.Price != 123.45 {
		t.Fatalf("Price mismatch: got %v", tr.Price)
	}
	if tr.Volume != -6.7 {
		t.Fatalf("Volume mismatch: got %v", tr.Volume)
	}
	if tr.Time != ts {
		t.Fatalf("Time mismatch: got %v want %v", tr.Time, ts)
	}
	if tr.PoolAddress != addr.Hex() {
		t.Fatalf("PoolAddress mismatch: got %s want %s", tr.PoolAddress, addr.Hex())
	}
	if tr.ForeignTradeID != "0xtx" {
		t.Fatalf("ForeignTradeID mismatch: got %s", tr.ForeignTradeID)
	}
	if tr.BaseToken.Symbol != "AERO" || tr.QuoteToken.Symbol != "WETH" {
		t.Fatalf("Token mapping mismatch: base=%s quote=%s", tr.BaseToken.Symbol, tr.QuoteToken.Symbol)
	}
	if tr.Exchange.Name != "AerodromeV1" || tr.Exchange.Blockchain != "base" {
		t.Fatalf("Exchange mismatch: %+v", tr.Exchange)
	}
}

func TestNormalizeSwap_PoolNotFound_ReturnsError(t *testing.T) {
	s := &AerodromeV1Scraper{
		poolMap: map[common.Address]AerodromeV1Pair{},
	}
	addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	raw := velodrome.IPoolSwap{
		Amount0In:  big.NewInt(1),
		Amount0Out: big.NewInt(0),
		Amount1In:  big.NewInt(0),
		Amount1Out: big.NewInt(1),
		Raw: types.Log{
			Address: addr,
			TxHash:  common.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
		},
	}

	_, err := s.normalizeSwap(context.Background(), nil, raw)
	if err == nil {
		t.Fatalf("expected error when pool not found")
	}
}

// helper: x * 10^dec
func mulPow10(x int64, dec int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(x), new(big.Int).Exp(big.NewInt(10), big.NewInt(dec), nil))
}
