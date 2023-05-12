/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var depositFILCmd = &cobra.Command{
	Use:   "deposit-fil [amount]",
	Short: "Deposit FIL into the Infinity Pool",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		receiver, err := util.DeriveAddressFromPk(pk)
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Depositing %s FIL into the Infinity Pool", amount.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		conn := fevm.Connection()

		tx, err := conn.PoolDepositFIL(cmd.Context(), conn.InfinityPoolAddr, receiver, amount, pk)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		receipt, err := fevm.WaitReturnReceipt(tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		if receipt == nil {
			log.Fatal("Failed to get receipt")
		}

		if receipt.Status == 0 {
			log.Fatal("Transaction failed")
		}

		s.Stop()

		fmt.Printf("Successfully deposited funds into the Infinity Pool")
	},
}

func init() {
	infinitypoolCmd.AddCommand(depositFILCmd)
	depositFILCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
