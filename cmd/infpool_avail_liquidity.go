/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var availLiquidityCmd = &cobra.Command{
	Use:   "avail-liquidity",
	Short: "Get the total FIL locked in the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Querying the available liquidity from the Infinity Pool...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		liquid, err := PoolsSDK.Query().InfPoolBorrowableLiquidity(cmd.Context())
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		liquidFIL, _ := liquid.Float64()

		s.Stop()

		fmt.Printf("Total available liquidity in the Pool is %.08f FIL", liquidFIL)
	},
}

func init() {
	infinitypoolCmd.AddCommand(availLiquidityCmd)
}
