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
	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

var payPrincipalCmd = &cobra.Command{
	Use:   "principal <amount> [flags]",
	Short: "Pay down an amount of principal (will also pay fees if any are owed)",
	Long:  "<amount> is the amount of principal to pay down, in FIL. Any fees owed will be paid off as well in order to make the principal payment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		conn := fevm.Connection()

		amountOwed, _, err := conn.AgentOwes(cmd, agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		payAmt := new(big.Int).Add(amount, amountOwed)

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Paying %s FIL to the %s", amountOwed, poolName)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := conn.AgentPay(cmd.Context(), agentAddr, poolID, payAmt, pk)
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

		fmt.Printf("Successfully paid %s FIL", args[0])
	},
}

func init() {
	payCmd.AddCommand(payPrincipalCmd)
	payPrincipalCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payPrincipalCmd.Flags().String("from", "", "address to send the transaction from")
}
