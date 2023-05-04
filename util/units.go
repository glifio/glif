package util

import (
	"math"
	"math/big"
)

func ToAtto(bal *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(bal.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}
