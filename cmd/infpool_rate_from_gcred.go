/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/constants"
	"github.com/spf13/cobra"
)

var rateFromGcredCmd = &cobra.Command{
	Use:   "rate-from-gcred <gcred>",
	Short: "Get a preview of the Pool's fee rate from the current GCRED score.",
	Long:  "Get a preview of the Pool's fee rate from the current GCRED score. For example, `20%` means an Agent who borrows 100 will pay 20 in fees over 1 years worth of epochs.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Querying the rate from a GCRED of %s from the Infinity Pool...\n", args[0])

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		gcred, ok := new(big.Int).SetString(args[0], 10)
		if !ok {
			log.Fatalf("Failed to parse GCRED %s", args[0])
		}

		rate, err := PoolsSDK.Query().InfPoolRateFromGCRED(cmd.Context(), gcred)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		wprFloat, _ := new(big.Float).Mul(rate, big.NewFloat(constants.EpochsInWeek)).Float64()
		aprFloat, _ := new(big.Float).Mul(rate, big.NewFloat(constants.EpochsInYear)).Float64()

		s.Stop()

		fmt.Printf("%.03f%% annually, %.03f%% weekly", aprFloat*100, wprFloat*100)
	},
}

func init() {
	infinitypoolCmd.AddCommand(rateFromGcredCmd)
}
