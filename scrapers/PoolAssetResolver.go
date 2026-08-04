package scrapers

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	uniswap "github.com/diadata-org/lumina-library/contracts/uniswap/pair"
	UniswapV3Pair "github.com/diadata-org/lumina-library/contracts/uniswapv3/uniswapV3Pair"
	poolState "github.com/diadata-org/lumina-library/contracts/uniswapv4/poolState"
	velodrome "github.com/diadata-org/lumina-library/contracts/velodrome"
	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PoolAssetResolver reads a pool's underlying assets and liquidity without
// running a full streaming scraper. Implementations are stateless: they only
// need an RPC client and the pool config. The pool's blockchain is looked up
// from the Exchanges registry (init.go) by pool.Exchange.Name, so callers don't
// pass it. This lets a consumer (e.g. the indexer's asset scraper) fill its
// asset table pool by pool without spinning up the DEX streaming machinery.
//
// Each DEX gets the two token addresses its own way (on-chain token0/token1,
// config, or Curve coin indices) and then resolves them via models.GetAsset.
type PoolAssetResolver interface {
	PoolAssets(ctx context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error)

	// PoolLiquidity returns how much of each underlying token the pool holds,
	// normalized by decimals. It is a coarse "is this pool basically empty"
	// signal (e.g. to deprecate a pool once liquidity drops below a threshold),
	// NOT exact active liquidity: it reads the pool's ERC20 balances, which for
	// Uniswap V3 include out-of-range positions. Uniswap V4 is not supported yet
	// (tokens are custodied by the singleton pool manager, not per pool).
	PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error)
}

// PoolAssetResolverFor returns the resolver for a config exchange name. The
// exchange -> resolver mapping mirrors the streaming dispatch in APIScraper.go:
// Pancakeswap V3 and Aerodrome Slipstream reuse the Uniswap V3 scraper, so they
// resolve with the V3 resolver here too. The blockchain is looked up per pool
// from the Exchanges registry, so chain variants (e.g. *_Base) share a resolver.
func PoolAssetResolverFor(exchangeName string) (PoolAssetResolver, bool) {
	switch exchangeName {
	case UNISWAPV2_EXCHANGE, UNISWAPV2_BASE_EXCHANGE:
		return uniswapV2AssetResolver{}, true
	case UNISWAPV3_EXCHANGE, UNISWAPV3_BASE_EXCHANGE, PANCAKESWAPV3_EXCHANGE, AERODROMESLIPSTREAM_EXCHANGE:
		return uniswapV3AssetResolver{}, true
	case UNISWAPV4_EXCHANGE, UNISWAPV4_BASE_EXCHANGE:
		return uniswapV4AssetResolver{}, true
	case AERODROMEV1_EXCHANGE:
		return aerodromeV1AssetResolver{}, true
	case CURVE_EXCHANGE:
		return curveAssetResolver{}, true
	default:
		return nil, false
	}
}

// blockchainForPool looks up the pool's chain from the Exchanges registry by
// exchange name (e.g. "UniswapV3_Base" -> "Base", "Curve" -> "Ethereum").
func blockchainForPool(pool models.Pool) (string, error) {
	ex, ok := Exchanges[pool.Exchange.Name]
	if !ok || ex.Blockchain == "" {
		return "", fmt.Errorf("no blockchain registered for exchange %q", pool.Exchange.Name)
	}
	return ex.Blockchain, nil
}

// pairFromAddrs is the shared last step across all DEXs: turn two token
// addresses into a fully-populated pair (symbol/decimals/blockchain via
// GetAsset). Quote is the first address, base the second.
func pairFromAddrs(client *ethclient.Client, blockchain string, quoteAddr, baseAddr common.Address) (models.Pair, error) {
	quote, err := models.GetAsset(quoteAddr, blockchain, client)
	if err != nil {
		return models.Pair{}, fmt.Errorf("resolve quote asset %s: %w", quoteAddr.Hex(), err)
	}
	base, err := models.GetAsset(baseAddr, blockchain, client)
	if err != nil {
		return models.Pair{}, fmt.Errorf("resolve base asset %s: %w", baseAddr.Hex(), err)
	}
	return models.Pair{QuoteToken: quote, BaseToken: base}, nil
}

// liquidityViaBalanceOf is the shared implementation of PoolLiquidity for DEXs
// that custody their tokens at the pool address (everything except Uniswap V4).
// It resolves the pair, then reads each token's ERC20 balance held by the pool.
func liquidityViaBalanceOf(ctx context.Context, client *ethclient.Client, pool models.Pool, resolver PoolAssetResolver) ([]models.AssetVolume, error) {
	pair, err := resolver.PoolAssets(ctx, client, pool)
	if err != nil {
		return nil, err
	}
	poolAddr := common.HexToAddress(pool.Address)
	out := make([]models.AssetVolume, 0, 2)
	for i, asset := range []models.Asset{pair.QuoteToken, pair.BaseToken} {
		bal, err := tokenBalanceOf(ctx, client, common.HexToAddress(asset.Address), poolAddr)
		if err != nil {
			return nil, fmt.Errorf("balanceOf %s in pool %s: %w", asset.Symbol, poolAddr.Hex(), err)
		}
		out = append(out, models.AssetVolume{
			Asset:  asset,
			Volume: normalizeByDecimals(bal, asset.Decimals),
			Index:  uint8(i),
		})
	}
	return out, nil
}

// tokenBalanceOf reads an ERC20 token's balance held by owner, reusing the
// generic token binding that models.GetAsset uses.
func tokenBalanceOf(ctx context.Context, client *ethclient.Client, token, owner common.Address) (*big.Int, error) {
	caller, err := models.NewTokenCaller(token, client)
	if err != nil {
		return nil, err
	}
	var res []interface{}
	if err := caller.Contract.Call(&bind.CallOpts{Context: ctx}, &res, "balanceOf", owner); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("balanceOf returned no value")
	}
	bal, ok := res[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("balanceOf returned unexpected type %T", res[0])
	}
	return bal, nil
}

// normalizeByDecimals turns a raw token amount into human units (amount / 10^decimals).
func normalizeByDecimals(amount *big.Int, decimals uint8) float64 {
	if amount == nil {
		return 0
	}
	denom := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	v, _ := new(big.Float).Quo(new(big.Float).SetInt(amount), denom).Float64()
	return v
}

// normalizeFloatByDecimals is normalizeByDecimals for an amount already in big.Float.
func normalizeFloatByDecimals(amount *big.Float, decimals uint8) float64 {
	if amount == nil {
		return 0
	}
	denom := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	v, _ := new(big.Float).Quo(amount, denom).Float64()
	return v
}

// virtualReserves expresses a concentrated-liquidity position (active liquidity
// L at price sqrtPriceX96) as the token0/token1 amounts a constant-product pool
// would hold at that price (in smallest units):
//
//	amount0 = L * 2^96 / sqrtPriceX96
//	amount1 = L * sqrtPriceX96 / 2^96
func virtualReserves(l, sqrtPriceX96 *big.Int) (amount0, amount1 *big.Float) {
	if l == nil || sqrtPriceX96 == nil || sqrtPriceX96.Sign() == 0 {
		return new(big.Float), new(big.Float)
	}
	q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96)) // 2^96
	lf := new(big.Float).SetInt(l)
	sp := new(big.Float).SetInt(sqrtPriceX96)
	amount0 = new(big.Float).Quo(new(big.Float).Mul(lf, q96), sp)
	amount1 = new(big.Float).Quo(new(big.Float).Mul(lf, sp), q96)
	return
}

// uniswapV2AssetResolver reads token0/token1 on-chain from a Uniswap V2 pair.
type uniswapV2AssetResolver struct{}

func (uniswapV2AssetResolver) PoolAssets(ctx context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error) {
	blockchain, err := blockchainForPool(pool)
	if err != nil {
		return models.Pair{}, err
	}
	caller, err := uniswap.NewUniswapV2PairCaller(common.HexToAddress(pool.Address), client)
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV2 pair caller: %w", err)
	}
	t0, err := caller.Token0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV2 token0: %w", err)
	}
	t1, err := caller.Token1(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV2 token1: %w", err)
	}
	return pairFromAddrs(client, blockchain, t0, t1)
}

func (r uniswapV2AssetResolver) PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error) {
	return liquidityViaBalanceOf(ctx, client, pool, r)
}

// uniswapV3AssetResolver reads token0/token1 on-chain from a Uniswap V3 pool.
type uniswapV3AssetResolver struct{}

func (uniswapV3AssetResolver) PoolAssets(ctx context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error) {
	blockchain, err := blockchainForPool(pool)
	if err != nil {
		return models.Pair{}, err
	}
	caller, err := UniswapV3Pair.NewUniswapV3PairCaller(common.HexToAddress(pool.Address), client)
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV3 pair caller: %w", err)
	}
	t0, err := caller.Token0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV3 token0: %w", err)
	}
	t1, err := caller.Token1(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV3 token1: %w", err)
	}
	return pairFromAddrs(client, blockchain, t0, t1)
}

func (r uniswapV3AssetResolver) PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error) {
	return liquidityViaBalanceOf(ctx, client, pool, r)
}

// uniswapV4AssetResolver takes the token addresses from the pool config, since a
// V4 pool is identified by a poolId rather than a callable pair contract.
type uniswapV4AssetResolver struct{}

func (uniswapV4AssetResolver) PoolAssets(_ context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error) {
	blockchain, err := blockchainForPool(pool)
	if err != nil {
		return models.Pair{}, err
	}
	t0, t1, err := getTokenAddrsFromPoolConfig(pool)
	if err != nil {
		return models.Pair{}, fmt.Errorf("uniswapV4 pool config: %w", err)
	}
	return pairFromAddrs(client, blockchain, t0, t1)
}

func (r uniswapV4AssetResolver) PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error) {
	// V4 tokens are custodied by the singleton pool manager, so a per-pool ERC20
	// balance is not meaningful. Instead we read the pool's active liquidity and
	// price from the StateView contract by poolId and express it as virtual
	// reserves — a coarse magnitude (active liquidity only), enough for the
	// "is this pool basically empty" threshold.
	stateViewEnv := strings.ToUpper(pool.Exchange.Name) + "_STATE_VIEW_ADDRESS"
	stateViewHex := utils.Getenv(stateViewEnv, "")
	if stateViewHex == "" {
		return nil, fmt.Errorf("%s not set", stateViewEnv)
	}

	pair, err := r.PoolAssets(ctx, client, pool)
	if err != nil {
		return nil, err
	}
	poolID, err := ParsePoolIDHex(pool.Address)
	if err != nil {
		return nil, fmt.Errorf("uniswapV4 poolId: %w", err)
	}

	stateView, err := poolState.NewPoolStateCaller(common.HexToAddress(stateViewHex), client)
	if err != nil {
		return nil, fmt.Errorf("uniswapV4 state view caller: %w", err)
	}
	opts := &bind.CallOpts{Context: ctx}
	liquidity, err := stateView.GetLiquidity(opts, poolID)
	if err != nil {
		return nil, fmt.Errorf("uniswapV4 getLiquidity: %w", err)
	}
	slot0, err := stateView.GetSlot0(opts, poolID)
	if err != nil {
		return nil, fmt.Errorf("uniswapV4 getSlot0: %w", err)
	}

	amount0, amount1 := virtualReserves(liquidity, slot0.SqrtPriceX96)

	// amount0/amount1 belong to the pool's canonical token0/token1 (lower address
	// = token0). Assign them to quote/base accordingly.
	q, b := pair.QuoteToken, pair.BaseToken
	qAmt, bAmt := amount0, amount1
	if bytes.Compare(common.HexToAddress(q.Address).Bytes(), common.HexToAddress(b.Address).Bytes()) > 0 {
		qAmt, bAmt = amount1, amount0
	}
	return []models.AssetVolume{
		{Asset: q, Volume: normalizeFloatByDecimals(qAmt, q.Decimals), Index: 0},
		{Asset: b, Volume: normalizeFloatByDecimals(bAmt, b.Decimals), Index: 1},
	}, nil
}

// aerodromeV1AssetResolver reads token0/token1 on-chain from an Aerodrome V1
// (velodrome-style) pool.
type aerodromeV1AssetResolver struct{}

func (aerodromeV1AssetResolver) PoolAssets(ctx context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error) {
	blockchain, err := blockchainForPool(pool)
	if err != nil {
		return models.Pair{}, err
	}
	caller, err := velodrome.NewIPoolCaller(common.HexToAddress(pool.Address), client)
	if err != nil {
		return models.Pair{}, fmt.Errorf("aerodromeV1 pool caller: %w", err)
	}
	t0, err := caller.Token0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("aerodromeV1 token0: %w", err)
	}
	t1, err := caller.Token1(&bind.CallOpts{Context: ctx})
	if err != nil {
		return models.Pair{}, fmt.Errorf("aerodromeV1 token1: %w", err)
	}
	return pairFromAddrs(client, blockchain, t0, t1)
}

func (r aerodromeV1AssetResolver) PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error) {
	return liquidityViaBalanceOf(ctx, client, pool, r)
}

// curveAssetResolver resolves the two coins of the pair encoded in pool.Order
// (out index = quote, in index = base), matching the streaming scraper.
type curveAssetResolver struct{}

func (curveAssetResolver) PoolAssets(ctx context.Context, client *ethclient.Client, pool models.Pool) (models.Pair, error) {
	blockchain, err := blockchainForPool(pool)
	if err != nil {
		return models.Pair{}, err
	}
	outIdx, inIdx, err := parseIndexCode(pool.Order)
	if err != nil {
		return models.Pair{}, fmt.Errorf("curve index code: %w", err)
	}
	poolAddr := common.HexToAddress(pool.Address)
	outAddr, err := coinAddressFromPool(ctx, client, poolAddr, outIdx)
	if err != nil {
		return models.Pair{}, err
	}
	inAddr, err := coinAddressFromPool(ctx, client, poolAddr, inIdx)
	if err != nil {
		return models.Pair{}, err
	}
	return pairFromAddrs(client, blockchain, outAddr, inAddr)
}

func (r curveAssetResolver) PoolLiquidity(ctx context.Context, client *ethclient.Client, pool models.Pool) ([]models.AssetVolume, error) {
	return liquidityViaBalanceOf(ctx, client, pool, r)
}
