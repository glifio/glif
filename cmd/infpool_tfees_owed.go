/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var tfeesOwedCmd = &cobra.Command{
	Use:   "treasury-fees-owed",
	Short: "Gets the WFIL held in the Infinity Pool that is owed to the Protocol Treasury",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Querying the fees collected but not paid to the treasury...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		fees, err := PoolsSDK.Query().InfPoolFeesAccrued(cmd.Context())
		if err != nil {
			logFatalf("Failed to get fees accrued %s", err)
		}
		s.Stop()

		feesOwed := util.ToFIL(fees)

		log.Printf("Fees owed: %0.09f", feesOwed)
	},
}

func init() {
	infinitypoolCmd.AddCommand(tfeesOwedCmd)
}
