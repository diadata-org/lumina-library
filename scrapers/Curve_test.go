package scrapers

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/diadata-org/lumina-library/contracts/curve/curvefifactory"
	curvefitwocryptooptimized "github.com/diadata-org/lumina-library/contracts/curve/curvefitwocrypto"
	"github.com/diadata-org/lumina-library/contracts/curve/curvepool"
	models "github.com/diadata-org/lumina-library/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestParseIndexCode_Valid(t *testing.T) {
	tests := []struct {
		code    int
		wantOut int
		wantIn  int
	}{
		{1, 0, 1},  // "01"
		{10, 1, 0}, // "10"
		{12, 1, 2}, // "12"
		{21, 2, 1}, // "21"
	}

	for _, tt := range tests {
		outIdx, inIdx, err := parseIndexCode(tt.code)
		if err != nil {
			t.Fatalf("parseIndexCode(%d) unexpected error: %v", tt.code, err)
		}
		if outIdx != tt.wantOut || inIdx != tt.wantIn {
			t.Fatalf("parseIndexCode(%d) = (%d,%d), want (%d,%d)",
				tt.code, outIdx, inIdx, tt.wantOut, tt.wantIn)
		}
	}
}

func TestParseIndexCode_Invalid(t *testing.T) {
	// here we only test the cases that we "determine should fail":
	// - 99: outIdx == inIdx
	// - 123: not two digits
	cases := []int{99, 123}

	for _, code := range cases {
		_, _, err := parseIndexCode(code)
		if err == nil {
			t.Fatalf("parseIndexCode(%d) expected error, got nil", code)
		}
	}
}

func TestFindPairByIdx(t *testing.T) {
	pairs := []CurvePair{
		{InIndex: 0, OutIndex: 1},
		{InIndex: 1, OutIndex: 2},
		{InIndex: 2, OutIndex: 0},
	}

	// found
	p, ok := findPairByIdx(pairs, 1, 2)
	if !ok {
		t.Fatalf("expected to find pair for (soldID=1,boughtID=2)")
	}
	if p.InIndex != 1 || p.OutIndex != 2 {
		t.Fatalf("unexpected pair: InIndex=%d OutIndex=%d, want (1,2)", p.InIndex, p.OutIndex)
	}

	// not found
	_, ok = findPairByIdx(pairs, 2, 2)
	if ok {
		t.Fatalf("expected not to find pair for (soldID=2,boughtID=2)")
	}
}

func TestGetSwapDataCurve_Valid(t *testing.T) {
	swap := CurveSwap{
		Amount0: 100,
		Amount1: 25,
	}

	price, volume := getSwapDataCurve(swap)
	if volume != 25 {
		t.Fatalf("volume = %f, want 25", volume)
	}
	if price != 4 {
		t.Fatalf("price = %f, want 4 (100/25)", price)
	}
}

func TestGetSwapDataCurve_InvalidVolume(t *testing.T) {
	// volume == 0 -> price should be NaN to signal invalid
	swap := CurveSwap{
		Amount0: 100,
		Amount1: 0,
	}
	price, volume := getSwapDataCurve(swap)
	if volume != 0 {
		t.Fatalf("volume = %f, want 0", volume)
	}
	if !math.IsNaN(price) {
		t.Fatalf("price = %f, want NaN when volume==0", price)
	}
}

func TestMakeCurveTrade(t *testing.T) {
	pair := CurvePair{
		OutAsset: models.Asset{
			Symbol:   "USDC",
			Decimals: 6,
		},
		InAsset: models.Asset{
			Symbol:   "USDT",
			Decimals: 6,
		},
		Address: common.HexToAddress("0x0000000000000000000000000000000000000001"),
	}

	price := 1.001
	volume := 123.45
	now := time.Now()
	addr := pair.Address
	txID := "0x1234"
	exName := "Curve"
	chain := "Ethereum"

	trade := makeCurveTrade(pair, price, volume, now, addr, txID, exName, chain)

	if trade.Price != price {
		t.Fatalf("trade.Price = %f, want %f", trade.Price, price)
	}
	if trade.Volume != volume {
		t.Fatalf("trade.Volume = %f, want %f", trade.Volume, volume)
	}
	if trade.BaseToken.Symbol != "USDT" || trade.QuoteToken.Symbol != "USDC" {
		t.Fatalf("unexpected base/quote: base=%s quote=%s",
			trade.BaseToken.Symbol, trade.QuoteToken.Symbol)
	}
	if trade.PoolAddress != addr.Hex() {
		t.Fatalf("trade.PoolAddress = %s, want %s", trade.PoolAddress, addr.Hex())
	}
	if trade.ForeignTradeID != txID {
		t.Fatalf("trade.ForeignTradeID = %s, want %s", trade.ForeignTradeID, txID)
	}
	if trade.Exchange.Name != exName || trade.Exchange.Blockchain != chain {
		t.Fatalf("unexpected exchange: %+v", trade.Exchange)
	}
}

func TestIsNative(t *testing.T) {
	native := common.HexToAddress(NativeETHSentinel)
	if !isNative(native) {
		t.Fatalf("isNative(%s) = false, want true", native.Hex())
	}

	nonNative := common.HexToAddress("0x0000000000000000000000000000000000000001")
	if isNative(nonNative) {
		t.Fatalf("isNative(%s) = true, want false", nonNative.Hex())
	}
}

func TestNativeAsset(t *testing.T) {
	chain := "Ethereum"
	a := nativeAsset(chain)

	if a.Symbol != "ETH" {
		t.Fatalf("nativeAsset.Symbol = %s, want ETH", a.Symbol)
	}
	if a.Decimals != 18 {
		t.Fatalf("nativeAsset.Decimals = %d, want 18", a.Decimals)
	}
	if a.Blockchain != chain {
		t.Fatalf("nativeAsset.Blockchain = %s, want %s", a.Blockchain, chain)
	}
	if !common.IsHexAddress(a.Address) {
		t.Fatalf("nativeAsset.Address = %s is not a hex address", a.Address)
	}
	if common.HexToAddress(a.Address) != common.HexToAddress(NativeETHSentinel) {
		t.Fatalf("nativeAsset.Address = %s, want %s", a.Address, NativeETHSentinel)
	}
}

func TestExtractSwapData_CurvepoolTokenExchange(t *testing.T) {
	s := &CurveScraper{}

	addr := common.HexToAddress("0x0000000000000000000000000000000000000002")
	txHash := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000aa")

	ev := curvepool.CurvepoolTokenExchange{
		Raw: types.Log{
			Address: addr,
			TxHash:  txHash,
		},
		SoldId:       big.NewInt(1),
		BoughtId:     big.NewInt(2),
		TokensSold:   big.NewInt(1000),
		TokensBought: big.NewInt(2000),
	}

	data, err := s.extractSwapData(ev)
	if err != nil {
		t.Fatalf("extractSwapData curvepool: unexpected error: %v", err)
	}
	if data.addr != addr {
		t.Fatalf("addr = %s, want %s", data.addr.Hex(), addr.Hex())
	}
	if data.soldID != 1 || data.boughtID != 2 {
		t.Fatalf("soldID/boughtID = (%d,%d), want (1,2)", data.soldID, data.boughtID)
	}
	if data.sold.Cmp(big.NewInt(1000)) != 0 || data.bought.Cmp(big.NewInt(2000)) != 0 {
		t.Fatalf("sold/bought = (%s,%s), want (1000,2000)", data.sold, data.bought)
	}
	if data.swapID != txHash.Hex() {
		t.Fatalf("swapID = %s, want %s", data.swapID, txHash.Hex())
	}
}

func TestExtractSwapData_CurvefifactoryTokenExchange(t *testing.T) {
	s := &CurveScraper{}

	addr := common.HexToAddress("0x0000000000000000000000000000000000000003")
	txHash := common.HexToHash("0xbb000000000000000000000000000000000000000000000000000000000000bb")

	ev := curvefifactory.CurvefifactoryTokenExchange{
		Raw: types.Log{
			Address: addr,
			TxHash:  txHash,
		},
		SoldId:       big.NewInt(0),
		BoughtId:     big.NewInt(1),
		TokensSold:   big.NewInt(111),
		TokensBought: big.NewInt(222),
	}

	data, err := s.extractSwapData(ev)
	if err != nil {
		t.Fatalf("extractSwapData curvefifactory: unexpected error: %v", err)
	}
	if data.addr != addr {
		t.Fatalf("addr = %s, want %s", data.addr.Hex(), addr.Hex())
	}
	if data.soldID != 0 || data.boughtID != 1 {
		t.Fatalf("soldID/boughtID = (%d,%d), want (0,1)", data.soldID, data.boughtID)
	}
	if data.sold.Cmp(big.NewInt(111)) != 0 || data.bought.Cmp(big.NewInt(222)) != 0 {
		t.Fatalf("sold/bought = (%s,%s), want (111,222)", data.sold, data.bought)
	}
	if data.swapID != txHash.Hex() {
		t.Fatalf("swapID = %s, want %s", data.swapID, txHash.Hex())
	}
}

func TestExtractSwapData_TwoCryptoTokenExchange(t *testing.T) {
	s := &CurveScraper{}

	addr := common.HexToAddress("0x0000000000000000000000000000000000000004")
	txHash := common.HexToHash("0xcc000000000000000000000000000000000000000000000000000000000000cc")

	ev := curvefitwocryptooptimized.CurvefitwocryptooptimizedTokenExchange{
		Raw: types.Log{
			Address: addr,
			TxHash:  txHash,
		},
		SoldId:       big.NewInt(2),
		BoughtId:     big.NewInt(0),
		TokensSold:   big.NewInt(3333),
		TokensBought: big.NewInt(4444),
	}

	data, err := s.extractSwapData(ev)
	if err != nil {
		t.Fatalf("extractSwapData twocrypto: unexpected error: %v", err)
	}
	if data.addr != addr {
		t.Fatalf("addr = %s, want %s", data.addr.Hex(), addr.Hex())
	}
	if data.soldID != 2 || data.boughtID != 0 {
		t.Fatalf("soldID/boughtID = (%d,%d), want (2,0)", data.soldID, data.boughtID)
	}
	if data.sold.Cmp(big.NewInt(3333)) != 0 || data.bought.Cmp(big.NewInt(4444)) != 0 {
		t.Fatalf("sold/bought = (%s,%s), want (3333,4444)", data.sold, data.bought)
	}
	if data.swapID != txHash.Hex() {
		t.Fatalf("swapID = %s, want %s", data.swapID, txHash.Hex())
	}
}

func TestExtractSwapData_CurvepoolTokenExchangeUnderlying(t *testing.T) {
	s := &CurveScraper{}

	addr := common.HexToAddress("0x0000000000000000000000000000000000000005")
	txHash := common.HexToHash("0xdd000000000000000000000000000000000000000000000000000000000000dd")

	ev := curvepool.CurvepoolTokenExchangeUnderlying{
		Raw: types.Log{
			Address: addr,
			TxHash:  txHash,
		},
		SoldId:       big.NewInt(1),
		BoughtId:     big.NewInt(0),
		TokensSold:   big.NewInt(555),
		TokensBought: big.NewInt(666),
	}

	data, err := s.extractSwapData(ev)
	if err != nil {
		t.Fatalf("extractSwapData underlying: unexpected error: %v", err)
	}
	if data.addr != addr {
		t.Fatalf("addr = %s, want %s", data.addr.Hex(), addr.Hex())
	}
	if data.soldID != 1 || data.boughtID != 0 {
		t.Fatalf("soldID/boughtID = (%d,%d), want (1,0)", data.soldID, data.boughtID)
	}
	if data.sold.Cmp(big.NewInt(555)) != 0 || data.bought.Cmp(big.NewInt(666)) != 0 {
		t.Fatalf("sold/bought = (%s,%s), want (555,666)", data.sold, data.bought)
	}
	if data.swapID != txHash.Hex() {
		t.Fatalf("swapID = %s, want %s", data.swapID, txHash.Hex())
	}
}

func TestExtractSwapData_UnknownType(t *testing.T) {
	s := &CurveScraper{}

	_, err := s.extractSwapData(struct{}{})
	if err == nil {
		t.Fatalf("extractSwapData(struct{}) expected error, got nil")
	}
}
