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

var depositFILCmd = &cobra.Command{
	Use:   "deposit-fil [amount]",
	Short: "Deposit FIL into the Infinity Pool",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, senderWallet, senderAccount, senderPassphrase, _, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			logFatal(err)
		}

		receiver := senderAccount.Address

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Depositing %s FIL into the Infinity Pool", amount.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().InfPoolDepositFIL(cmd.Context(), receiver, amount, senderWallet, senderAccount, senderPassphrase)
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

		fmt.Printf("Successfully deposited funds into the Infinity Pool\n")
	},
}

func init() {
	infinitypoolCmd.AddCommand(depositFILCmd)
	depositFILCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
