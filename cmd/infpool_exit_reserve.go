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

var exitReserveCmd = &cobra.Command{
	Use:   "exit-reserve",
	Short: "Get the total FIL held aside for the exit reserve",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Querying the exit reserve from the Infinity Pool...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		reserve, err := PoolsSDK.Query().InfPoolExitReserve(cmd.Context(), nil)
		if err != nil {
			logFatalf("Failed to get exit reserve %s", err)
		}

		s.Stop()

		fmt.Printf("Total available liquidity in the Pool is %.08f FIL\n", reserve)
	},
}

func init() {
	infinitypoolCmd.AddCommand(exitReserveCmd)
}
