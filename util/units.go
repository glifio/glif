package util

import (
	"math"
	"math/big"
)

func ToAtto(bal *big.Float) *big.Int {
	wei := new(big.Float).Mul(bal, big.NewFloat(math.Pow10(18)))
	result := new(big.Int)
	wei.Int(result)
	return result
}

func ToFIL(atto *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(atto.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

var WAD = big.NewInt(1000000000000000000)
