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

var redeemFILCmd = &cobra.Command{
	Use:   "redeem <iFIL-amount> <receiver>",
	Short: "Redeem WFIL from the Infinity Pool by burning a specific number of iFIL tokens",
	Long:  "Redeem iFIL for WFIL from the Infinity Pool. The address of the SimpleRamp must be approved for the appropriate amount of iFIL in order for this call to go execute.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		from := cmd.Flag("from").Value.String()
		auth, senderAccount, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		receiver, err := AddressOrAccountNameToEVM(cmd.Context(), args[1])
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Burning %0.09f iFIL to receive wFIL\n", util.ToFIL(amount))

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().InfPoolRedeem(cmd.Context(), auth, amount, senderAccount.Address, receiver)
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

		fmt.Printf("Successfully redeemed WFIL for iFIL from the Infinity Pool\n")
	},
}

func init() {
	infinitypoolCmd.AddCommand(redeemFILCmd)
	redeemFILCmd.Flags().String("from", "", "account to send transaction")
}
