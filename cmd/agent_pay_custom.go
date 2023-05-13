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
	"github.com/spf13/cobra"
)

var payCustomCmd = &cobra.Command{
	Use:   "custom <amount> [flags]",
	Short: "Pay down a custom amount of FIL",
	Args:  cobra.ExactArgs(1),
	Long:  "",
	// Args:  cobra.RangeArgs(1, 2),
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

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Paying %s FIL to the %s", args[0], poolName)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := conn.AgentPay(cmd.Context(), agentAddr, poolID, amount, pk)
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
	payCmd.AddCommand(payCustomCmd)
	payCustomCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payCustomCmd.Flags().String("from", "", "address to send the transaction from")
}
