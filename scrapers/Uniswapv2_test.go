package scrapers

import (
	"math"
	"math/big"
	"testing"
	"time"

	uniswap "github.com/diadata-org/lumina-library/contracts/uniswap/pair"
	"github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// float comparison helper function
func almostEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) < eps
}

func TestNormalizeUniswapSwap(t *testing.T) {
	// prepare a virtual pool
	pairAddr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	pair := UniswapPair{
		Token0: models.Asset{
			Decimals: 18,
			Symbol:   "T0",
		},
		Token1: models.Asset{
			Decimals: 6,
			Symbol:   "T1",
		},
		Address: pairAddr,
	}

	scraper := &UniswapV2Scraper{
		poolMap: map[common.Address]UniswapPair{
			pairAddr: pair,
		},
	}

	// construct a Swap event:
	//  - amount0In = 1e18 -> 1.0
	//  - amount1Out = 2 * 1e6 -> 2.0
	amount0In := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)   // 1e18
	amount1Out := new(big.Int).Mul(big.NewInt(2), big.NewInt(1_000_000)) // 2e6
	txHash := common.HexToHash("0x1234")

	rawSwap := uniswap.UniswapV2PairSwap{
		Amount0In:  amount0In,
		Amount0Out: big.NewInt(0),
		Amount1In:  big.NewInt(0),
		Amount1Out: amount1Out,
		Raw: types.Log{
			Address: pairAddr,
			TxHash:  txHash,
		},
	}

	norm, err := scraper.normalizeUniswapSwap(rawSwap)
	if err != nil {
		t.Fatalf("normalizeUniswapSwap returned error: %v", err)
	}

	// check ID
	if norm.ID != txHash.Hex() {
		t.Fatalf("expected ID %s, got %s", txHash.Hex(), norm.ID)
	}

	// check pair binding
	if norm.Pair.Address != pairAddr {
		t.Fatalf("expected pair address %s, got %s", pairAddr.Hex(), norm.Pair.Address.Hex())
	}

	// check precision scaling
	if !almostEqual(norm.Amount0In, 1.0) {
		t.Fatalf("expected Amount0In = 1.0, got %f", norm.Amount0In)
	}
	if !almostEqual(norm.Amount1Out, 2.0) {
		t.Fatalf("expected Amount1Out = 2.0, got %f", norm.Amount1Out)
	}
}

func TestGetSwapDataV2_Amount0InZero(t *testing.T) {
	// Amount0In == 0
	swap := UniswapSwap{
		Amount0In:  0,
		Amount0Out: 1.0,
		Amount1In:  2.0,
		Amount1Out: 0,
	}

	price, volume := getSwapDataV2(swap)

	if !almostEqual(volume, 1.0) {
		t.Fatalf("expected volume = 1.0, got %f", volume)
	}
	if !almostEqual(price, 2.0) {
		t.Fatalf("expected price = 2.0, got %f", price)
	}
}

func TestGetSwapDataV2_Amount0InNonZero(t *testing.T) {
	// Amount0In != 0
	swap := UniswapSwap{
		Amount0In:  1.0,
		Amount0Out: 0,
		Amount1In:  0,
		Amount1Out: 3.0,
	}

	price, volume := getSwapDataV2(swap)

	if !almostEqual(volume, -1.0) {
		t.Fatalf("expected volume = -1.0, got %f", volume)
	}
	if !almostEqual(price, 3.0) {
		t.Fatalf("expected price = 3.0, got %f", price)
	}
}

func TestMakeTradeUniswapV2(t *testing.T) {
	pairAddr := common.HexToAddress("0x0000000000000000000000000000000000000002")

	token0 := models.Asset{
		Symbol:   "T0",
		Decimals: 18,
	}
	token1 := models.Asset{
		Symbol:   "T1",
		Decimals: 6,
	}

	pair := UniswapPair{
		Token0:  token0,
		Token1:  token1,
		Address: pairAddr,
		Order:   0,
	}

	price := 10.5
	volume := 123.45
	now := time.Now()
	exName := "UNISWAPV2"
	chain := "ethereum"
	txID := "0xdeadbeef"

	trade := makeTradeUniswapV2(
		pair,
		price,
		volume,
		now,
		pairAddr,
		txID,
		exName,
		chain,
	)

	// check basic fields
	if !almostEqual(trade.Price, price) {
		t.Fatalf("expected price = %f, got %f", price, trade.Price)
	}
	if !almostEqual(trade.Volume, volume) {
		t.Fatalf("expected volume = %f, got %f", volume, trade.Volume)
	}

	// BaseToken / QuoteToken
	if trade.BaseToken.Symbol != token1.Symbol {
		t.Fatalf("expected BaseToken = %s, got %s", token1.Symbol, trade.BaseToken.Symbol)
	}
	if trade.QuoteToken.Symbol != token0.Symbol {
		t.Fatalf("expected QuoteToken = %s, got %s", token0.Symbol, trade.QuoteToken.Symbol)
	}

	// PoolAddress
	if trade.PoolAddress != pairAddr.Hex() {
		t.Fatalf("expected PoolAddress = %s, got %s", pairAddr.Hex(), trade.PoolAddress)
	}

	// Exchange information
	if trade.Exchange.Name != exName {
		t.Fatalf("expected Exchange.Name = %s, got %s", exName, trade.Exchange.Name)
	}
	if trade.Exchange.Blockchain != chain {
		t.Fatalf("expected Exchange.Blockchain = %s, got %s", chain, trade.Exchange.Blockchain)
	}

	// ForeignTradeID
	if trade.ForeignTradeID != txID {
		t.Fatalf("expected ForeignTradeID = %s, got %s", txID, trade.ForeignTradeID)
	}
}
