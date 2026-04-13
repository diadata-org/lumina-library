package utils

import (
	"math"
	"math/big"
)

func ScaleFloat(f float64, decimals int) *big.Int {
	fBig := big.NewFloat(f)
	scaling := big.NewFloat(math.Pow10(decimals))
	priceScaled := new(big.Float).Mul(fBig, scaling)
	valueUSDInt := new(big.Int)
	priceScaled.Int(valueUSDInt)
	return valueUSDInt
}
