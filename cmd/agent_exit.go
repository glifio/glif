/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exits from the Infintiy Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, senderKey, requesterKey, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			logFatal(err)
		}

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		account, err := PoolsSDK.Query().InfPoolGetAccount(cmd.Context(), agentAddr)
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		amountOwed, _, err := PoolsSDK.Query().AgentOwes(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		payAmount := new(big.Int).Add(amountOwed, account.Principal)
		payAmount = addOnePercent(payAmount)

		tx, err := PoolsSDK.Act().AgentPay(cmd.Context(), agentAddr, poolID, payAmount, senderKey, requesterKey)
		if err != nil {
			logFatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		log.Println("Successfully exited from the Infinity Pool")
	},
}

func addOnePercent(amount *big.Int) *big.Int {
	// Convert the amount to big.Float
	amountFloat := new(big.Float).SetInt(amount)

	// Add 1%
	onePercent := new(big.Float).SetFloat64(1.01)
	amountFloat.Mul(amountFloat, onePercent)

	// Convert back to big.Int
	newAmount := new(big.Int)
	amountFloat.Int(newAmount)

	return newAmount
}

func init() {
	agentCmd.AddCommand(exitCmd)
	exitCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	exitCmd.Flags().String("from", "", "address to send the transaction from")
}
