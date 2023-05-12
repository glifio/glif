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

// borrowCmd represents the borrow command
var borrowCmd = &cobra.Command{
	Use:   "borrow <amount> [pool-name]",
	Short: "Borrow FIL from a Pool",
	Long:  "Borrow FIL from a Pool. If you do not pass a `pool-name` arg, the default pool is the Infinity Pool.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if amount.Cmp(util.WAD) == -1 {
			log.Fatal("Borrow amount must be greater than 1 FIL")
		}

		poolName := args[1]
		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Borrowing %s FIL from the %s into agent %s", amount, poolID, agentAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := fevm.Connection().AgentBorrow(cmd.Context(), agentAddr, poolID, amount, ownerKey)
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

		fmt.Printf("Successfully borrowed %s FIL", amount)
	},
}

func init() {
	agentCmd.AddCommand(borrowCmd)
}
