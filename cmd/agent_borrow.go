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

// borrowCmd represents the borrow command
var borrowCmd = &cobra.Command{
	Use:   "borrow <amount> [flags]",
	Short: "Borrow FIL from a Pool",
	Long:  "Borrow FIL from a Pool. If you do not pass a `pool-name` flag, the default pool is the Infinity Pool.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		if amount.Cmp(util.WAD) == -1 {
			logFatal("Borrow amount must be greater than 1 FIL")
		}

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Borrowing %s FIL from the %s into agent %s", amount, poolID, agentAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().AgentBorrow(cmd.Context(), agentAddr, poolID, amount, ownerKey, requesterKey)
		if err != nil {
			logFatal(err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully borrowed %s FIL", args[0])
	},
}

func init() {
	agentCmd.AddCommand(borrowCmd)
	borrowCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to borrow from")
}
