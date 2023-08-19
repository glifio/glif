/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var withdrawFILCmd = &cobra.Command{
	Use:   "withdraw <wfil-amount> <receiver>",
	Short: "Withdraw WFIL from the Infinity Pool",
	Long:  "Withdraw WFIL from the Infinity Pool by burning the appropriate amount of iFIL tokens. The address of the SimpleRamp must be approved for the appropriate amount of iFIL in order for this call to go execute.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, pk, _, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		receiver, err := ParseAddressToEVM(cmd.Context(), args[1])
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Withdrawing %0.09f WFIL from the Infinity Pool\n", util.ToFIL(amount))

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().RampWithdraw(cmd.Context(), amount, receiver, pk)
		if err != nil {
			logFatal(err)
		}

		// transaction landed on chain or errored
		receipt, err := PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatal(err)
		}

		if receipt == nil {
			logFatal("Failed to get receipt")
		}

		if receipt.Status == 0 {
			logFatal("Transaction failed")
		}

		s.Stop()

		fmt.Printf("Successfully withdrew WFIL from the Infinity Pool\n")
	},
}

func init() {
	infinitypoolCmd.AddCommand(withdrawFILCmd)
	withdrawFILCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
