package util

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func ToAtto(bal *big.Float) *big.Int {
	wei := new(big.Float).Mul(bal, big.NewFloat(math.Pow10(18)))
	result := new(big.Int)
	wei.Int(result)
	return result
}

func ToFIL(atto *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	return f.Quo(f.SetInt(atto), big.NewFloat(params.Ether))
}

var WAD = big.NewInt(1000000000000000000)
