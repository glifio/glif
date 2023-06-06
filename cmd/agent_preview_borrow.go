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
	"github.com/filecoin-project/go-address"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/rpc"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var previewBorrowCmd = &cobra.Command{
	Use:   "borrow <amount> [flags]",
	Short: "Preview borrowing from the Infinity Pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		closer, err := PoolsSDK.Extern().ConnectAdoClient(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}
		defer closer()

		agentDataBefore, err := rpc.ADOClient.AgentData(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		attofil, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		agentDataAfter, err := rpc.ADOClient.PreviewAction(cmd.Context(), agentAddr, address.Undef, attofil, constants.MethodBorrow)
		if err != nil {
			log.Fatal(err)
		}

		rateAfter, err := PoolsSDK.Query().InfPoolRateFromGCRED(cmd.Context(), agentDataAfter.Gcred)
		if err != nil {
			log.Fatal(err)
		}

		wprFloat, _ := new(big.Float).Mul(rateAfter, big.NewFloat(constants.EpochsInWeek)).Float64()
		aprFloat, _ := new(big.Float).Mul(rateAfter, big.NewFloat(constants.EpochsInYear)).Float64()

		s.Stop()

		generateHeader("PREVIEW BORROW")
		fmt.Printf("Total borrowed before/after: %0.09f => %0.09f\n", util.ToFIL(agentDataBefore.Principal), util.ToFIL(agentDataAfter.Principal))
		fmt.Printf("GCRED before/after: %s => %s\n", agentDataBefore.Gcred, agentDataAfter.Gcred)
		fmt.Printf("The weekly/annual fee rate: %.03f%% / %.03f%%\n", wprFloat, aprFloat/52)
	},
}

func init() {
	previewCmd.AddCommand(previewBorrowCmd)
	previewBorrowCmd.Flags().String("agent-addr", "", "Agent address")
}
