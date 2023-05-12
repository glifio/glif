/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
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

		conn := fevm.Connection()

		bal, err := conn.PoolAvailableLiquidity(cmd.Context(), conn.InfinityPoolAddr)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		fmt.Printf("Total available liquidity in the Pool is %.02f FIL", bal)
	},
}

func init() {
	infinitypoolCmd.AddCommand(availLiquidityCmd)
}
