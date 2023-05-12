package util_test

import (
	"log"
	"math/big"
	"testing"

	"github.com/glif-confidential/cli/util"
)

func TestToFIL(t *testing.T) {
	var tests = []struct {
		bal *big.Int
		fil *big.Float
	}{
		{mustNewBigInt("0"), big.NewFloat(0)},
		{mustNewBigInt("1000000000000000000"), big.NewFloat(1)},
		{mustNewBigInt("100000000000000000000"), big.NewFloat(100)},
	}

	for _, test := range tests {
		fil := util.ToFIL(test.bal)
		if fil.Cmp(test.fil) != 0 {
			t.Errorf("ToFIL(%s) expected %d, got %d", test.bal, test.fil, fil)
		}
	}
}

func mustNewBigInt(val string) *big.Int {
	res, success := big.NewInt(0).SetString(val, 10)
	if !success {
		log.Panic("big int failed")
	}
	return res
}
