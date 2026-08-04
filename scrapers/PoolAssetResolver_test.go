package scrapers

import (
	"math"
	"math/big"
	"testing"

	"github.com/diadata-org/lumina-library/models"
)

func TestVirtualReserves(t *testing.T) {
	tests := []struct {
		name         string
		l            *big.Int
		sqrtPriceX96 *big.Int
		wantZero     bool
	}{
		{"nil liquidity", nil, big.NewInt(1 << 62), true},
		{"nil sqrtPrice", big.NewInt(1000), nil, true},
		{"zero sqrtPrice", big.NewInt(1000), big.NewInt(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a0, a1 := virtualReserves(tt.l, tt.sqrtPriceX96)
			if tt.wantZero {
				if a0.Sign() != 0 || a1.Sign() != 0 {
					t.Errorf("expected (0,0), got (%v,%v)", a0, a1)
				}
			}
		})
	}

	// Known values: L = 2^96 (so amount0 = 2^192/sqrtP, amount1 = sqrtP), and
	// sqrtPriceX96 = 2^96 (price = 1). At price 1, amount0 should equal amount1.
	q96 := new(big.Int).Lsh(big.NewInt(1), 96)
	l := new(big.Int).Set(q96)
	a0, a1 := virtualReserves(l, q96)
	f0, _ := a0.Float64()
	f1, _ := a1.Float64()
	if math.Abs(f0-f1) > 1e-6 {
		t.Errorf("at price 1, expected amount0 == amount1, got %v vs %v", f0, f1)
	}
	// amount0 = L * 2^96 / sqrtPriceX96 = q96 * q96 / q96 = q96.
	wantF, _ := new(big.Float).SetInt(q96).Float64()
	if math.Abs(f0-wantF) > wantF*1e-9 {
		t.Errorf("amount0 = %v, want ~%v", f0, wantF)
	}
}

func TestNormalizeByDecimals(t *testing.T) {
	// A trillion 18-decimal tokens (10^30 raw units): large enough to probe
	// float64 precision loss, so this uses relative rather than absolute
	// tolerance below.
	hugeSupply, _ := new(big.Int).SetString("1000000000000000000000000000000", 10)

	tests := []struct {
		name     string
		amount   *big.Int
		decimals uint8
		want     float64
	}{
		{"nil amount", nil, 18, 0},
		{"zero decimals", big.NewInt(1234), 0, 1234},
		{"6 decimals (USDC-like)", big.NewInt(1_000_000), 6, 1},
		{"18 decimals (ETH-like)", new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), 18, 1},
		{"large balance (10^30 raw, 18 decimals)", hugeSupply, 18, 1e12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeByDecimals(tt.amount, tt.decimals)
			tol := 1e-9
			if tt.want != 0 {
				tol = math.Abs(tt.want) * 1e-9
			}
			if math.Abs(got-tt.want) > tol {
				t.Errorf("normalizeByDecimals(%v, %d) = %v, want %v", tt.amount, tt.decimals, got, tt.want)
			}
		})
	}
}

func TestNormalizeFloatByDecimals(t *testing.T) {
	hugeSupply, _, _ := big.ParseFloat("1000000000000000000000000000000", 10, 200, big.ToNearestEven)

	tests := []struct {
		name     string
		amount   *big.Float
		decimals uint8
		want     float64
	}{
		{"nil amount", nil, 18, 0},
		{"18 decimals", new(big.Float).SetFloat64(math.Pow10(18)), 18, 1},
		{"large balance (10^30 raw, 18 decimals)", hugeSupply, 18, 1e12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeFloatByDecimals(tt.amount, tt.decimals)
			tol := 1e-9
			if tt.want != 0 {
				tol = math.Abs(tt.want) * 1e-9
			}
			if math.Abs(got-tt.want) > tol {
				t.Errorf("normalizeFloatByDecimals(...) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoolAssetResolverFor(t *testing.T) {
	tests := []struct {
		exchange string
		wantOK   bool
		wantType PoolAssetResolver
	}{
		{UNISWAPV2_EXCHANGE, true, uniswapV2AssetResolver{}},
		{UNISWAPV2_BASE_EXCHANGE, true, uniswapV2AssetResolver{}},
		{UNISWAPV3_EXCHANGE, true, uniswapV3AssetResolver{}},
		{UNISWAPV3_BASE_EXCHANGE, true, uniswapV3AssetResolver{}},
		{PANCAKESWAPV3_EXCHANGE, true, uniswapV3AssetResolver{}},
		{AERODROMESLIPSTREAM_EXCHANGE, true, uniswapV3AssetResolver{}},
		{UNISWAPV4_EXCHANGE, true, uniswapV4AssetResolver{}},
		{UNISWAPV4_BASE_EXCHANGE, true, uniswapV4AssetResolver{}},
		{AERODROMEV1_EXCHANGE, true, aerodromeV1AssetResolver{}},
		{CURVE_EXCHANGE, true, curveAssetResolver{}},
		{"NotARealExchange", false, nil},
		{"", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.exchange, func(t *testing.T) {
			got, ok := PoolAssetResolverFor(tt.exchange)
			if ok != tt.wantOK {
				t.Fatalf("PoolAssetResolverFor(%q) ok = %v, want %v", tt.exchange, ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != nil {
					t.Errorf("expected nil resolver for unknown exchange, got %T", got)
				}
				return
			}
			if got != tt.wantType {
				t.Errorf("PoolAssetResolverFor(%q) = %T, want %T", tt.exchange, got, tt.wantType)
			}
		})
	}
}

func TestBlockchainForPool(t *testing.T) {
	tests := []struct {
		name      string
		exchange  string
		wantChain string
		wantErr   bool
	}{
		{"unregistered exchange", "NotARealExchange", "", true},
		{"known DEX", UNISWAPV2_EXCHANGE, Exchanges[UNISWAPV2_EXCHANGE].Blockchain, false},
		{"known DEX on Base", UNISWAPV3_BASE_EXCHANGE, Exchanges[UNISWAPV3_BASE_EXCHANGE].Blockchain, false},
		// CEX exchanges (e.g. Binance) are registered but have no Blockchain set.
		{"CEX exchange has no blockchain", BINANCE_EXCHANGE, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := blockchainForPool(models.Pool{Exchange: models.Exchange{Name: tt.exchange}})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("blockchainForPool(%q) expected error, got chain %q", tt.exchange, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("blockchainForPool(%q) unexpected error: %v", tt.exchange, err)
			}
			if got != tt.wantChain {
				t.Errorf("blockchainForPool(%q) = %q, want %q", tt.exchange, got, tt.wantChain)
			}
		})
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid address", "0x5612599CF48032d7428399d5Fcb99eDcc75c06A7", false},
		{"native ETH sentinel", NativeETHSentinel, false},
		{"empty string", "", true},
		{"truncated address", "0x5612599CF48032d7428399d5Fcb99eDcc75c06", true},
		{"not hex", "not-an-address", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseAddress("test", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
