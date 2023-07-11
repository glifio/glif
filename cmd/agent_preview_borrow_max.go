/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/rpc"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var previewBorrowMaxCmd = &cobra.Command{
	Use:   "borrow-max [flags]",
	Short: "Preview borrowing from the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		closer, err := PoolsSDK.Extern().ConnectAdoClient(cmd.Context())
		if err != nil {
			logFatal(err)
		}
		defer closer()

		agentData, err := rpc.ADOClient.AgentData(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		borrowNowMax, err := PoolsSDK.Query().InfPoolAgentMaxBorrow(cmd.Context(), agentAddr, agentData)
		if err != nil {
			logFatal(err)
		}

		account, err := PoolsSDK.Query().InfPoolGetAccount(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		chainHeight, err := PoolsSDK.Query().ChainHeight(cmd.Context())
		if err != nil {
			logFatal(err)
		}

		generateHeader("PREVIEW BORROW MAX")

		if agentData.Principal.Cmp(big.NewInt(0)) > 0 {
			fmt.Printf("Agent has existing principal of %0.09f FIL\n", util.ToFIL(agentData.Principal))
		}
		printWithBoldPreface("Now", fmt.Sprintf("Agent can immediately borrow %0.09f FIL", util.ToFIL(borrowNowMax)))
		printWithBoldPreface("Max", fmt.Sprintf("Agent can borrow up to %0.09f FIL", util.ToFIL(new(big.Int).Sub(agentData.AgentValue, agentData.Principal))))
		if account.EpochsPaid.Cmp(big.NewInt(0)) == 1 && new(big.Int).Add(account.EpochsPaid, big.NewInt(int64(constants.RepeatBorrowEpochTolerance))).Cmp(chainHeight) == -1 {
			printWithBoldPreface("\nWarning", "Agent must make a payment `to-current` before borrowing again.")
		}
	},
}

func printWithBoldPreface(preface, msg string) {
	fmt.Printf("\033[1m%s\033[0m: %s\n", preface, msg)
}

func init() {
	previewCmd.AddCommand(previewBorrowMaxCmd)
	previewBorrowMaxCmd.Flags().String("agent-addr", "", "Agent address")
}
