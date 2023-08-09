package cmd

import (
	"math/big"
	"testing"
)

func TestPaymentDue(t *testing.T) {

}

func Test_paymentDue(t *testing.T) {
	type args struct {
		frequency       float64
		chainHeadHeight *big.Int
		epochsPaid      *big.Int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"due",
			args{
				0.1,
				big.NewInt(300),
				big.NewInt(5),
			},
			true,
		},
		{
			"not due",
			args{
				0.1,
				big.NewInt(30),
				big.NewInt(5),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paymentDue(tt.args.frequency, tt.args.chainHeadHeight, tt.args.epochsPaid); got != tt.want {
				t.Errorf("paymentDue() = %v, want %v", got, tt.want)
			}
		})
	}
}
