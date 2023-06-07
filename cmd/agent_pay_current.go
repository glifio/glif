/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payToCurrentCmd = &cobra.Command{
	Use:   "to-current [flags]",
	Short: "Make your account current",
	Long:  "Pays off all fees owed",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, senderKey, requesterKey, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			logFatal(err)
		}

		amountOwed, _, err := PoolsSDK.Query().AgentOwes(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			logFatal(err)
		}

		log.Printf("Paying %s FIL to the %s", util.ToFIL(amountOwed).String(), poolName)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().AgentPay(cmd.Context(), agentAddr, poolID, amountOwed, senderKey, requesterKey)
		if err != nil {
			logFatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully paid %s FIL", util.ToFIL(amountOwed).String())
	},
}

func init() {
	payCmd.AddCommand(payToCurrentCmd)
	payToCurrentCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payToCurrentCmd.Flags().String("from", "", "address to send the transaction from")
}
