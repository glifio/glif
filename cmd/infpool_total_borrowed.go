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

var inpoolTotalBorrowedCmd = &cobra.Command{
	Use:   "total-borrowed",
	Short: "Returns the amount of FIL currently borrowed by Agents from the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Querying the amount of FIL currently outstanding from the Infinity Pool")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		assets, err := PoolsSDK.Query().InfPoolTotalBorrowed(cmd.Context())
		if err != nil {
			logFatalf("Failed to get total borrowed %s", err)
		}

		fmt.Printf("Infinity Pool outstanding: %.04f FIL\n", assets)
	},
}

func init() {
	infinitypoolCmd.AddCommand(inpoolTotalBorrowedCmd)
}
