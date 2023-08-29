/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var exitReserveCmd = &cobra.Command{
	Use:   "exit-reserve",
	Short: "Get the total FIL held aside for the exit reserve",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Querying the exit reserve from the Infinity Pool...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		reserveBal, reserveMax, err := PoolsSDK.Query().InfPoolExitReserve(cmd.Context(), nil)
		if err != nil {
			logFatalf("Failed to get exit reserve %s", err)
		}

		s.Stop()

		reserveBalFIL := util.ToFIL(reserveBal)

		if reserveBal.Cmp(reserveMax) == 0 {
			fmt.Printf("The exit reserve in the Pool is full at %.08f FIL\n", reserveBalFIL)
		} else {
			fmt.Printf("The exit reserve in the Pool is not full - it has %.08f FIL\n", reserveBalFIL)
			// get the deficit
			deficit := util.ToFIL(new(big.Int).Sub(reserveMax, reserveBal))
			// div out the wads from the big ints
			reserveBal.Div(reserveBal, big.NewInt(1e18))
			reserveMax.Div(reserveMax, big.NewInt(1e18))
			reservePerc := computePercentage(reserveBal, reserveMax)

			fmt.Printf("The exit reserve is %.2f%% full, it needs %.08f FIL to be full\n", reservePerc, deficit)
		}
	},
}

// computePercentage computes (numerator/denominator) * 100
func computePercentage(numerator, denominator *big.Int) float64 {
	// Convert big.Int to big.Rat
	rNumerator := new(big.Rat).SetInt(numerator)
	rDenominator := new(big.Rat).SetInt(denominator)

	// Compute the ratio
	ratio := new(big.Rat).Quo(rNumerator, rDenominator)

	percentageFloat := float64(ratio.Num().Int64()) / float64(ratio.Denom().Int64()) * 100

	return percentageFloat
}

func init() {
	infinitypoolCmd.AddCommand(exitReserveCmd)
}
