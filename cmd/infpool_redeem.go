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

var redeemFILCmd = &cobra.Command{
	Use:   "redeem <iFIL-amount> <receiver>",
	Short: "Redeem WFIL from the Infinity Pool by burning a specific number of iFIL tokens",
	Long:  "Redeem iFIL for WFIL from the Infinity Pool. The address of the SimpleRamp must be approved for the appropriate amount of iFIL in order for this call to go execute.",
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

		fmt.Printf("Withdrawing %s WFIL from the Infinity Pool", amount.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().RampRedeem(cmd.Context(), amount, receiver, pk)
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
	infinitypoolCmd.AddCommand(redeemFILCmd)
	redeemFILCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
