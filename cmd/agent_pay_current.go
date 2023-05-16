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
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payToCurrentCmd = &cobra.Command{
	Use:   "to-current [flags]",
	Short: "Make your account current",
	Long:  "Pays off all fees owed",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		amountOwed, _, err := PoolsSDK.Query().AgentOwes(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Paying %s FIL to the %s", amountOwed, poolName)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().AgentPay(cmd.Context(), agentAddr, poolID, amountOwed, pk)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		paidAmount := new(big.Int)
		paidAmount.SetString(args[0], 10)

		fmt.Printf("Successfully paid %s FIL", util.ToFIL(paidAmount))
	},
}

func init() {
	payCmd.AddCommand(payToCurrentCmd)
	payToCurrentCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payToCurrentCmd.Flags().String("from", "", "address to send the transaction from")
}
