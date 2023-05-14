/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exits from the Infintiy Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		as := util.AgentStore()
		agentIDStr, err := as.Get("id")
		if err != nil {
			log.Fatal(err)
		}

		agentID, _ := new(big.Int).SetString(agentIDStr, 10)

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		conn := fevm.Connection()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		account, err := conn.PoolGetAccount(cmd.Context(), conn.InfinityPoolAddr, agentID)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		amountOwed, _, err := conn.AgentOwes(cmd, agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		payAmount := new(big.Int).Add(amountOwed, account.Principal)
		payAmount = addOnePercent(payAmount)

		tx, err := conn.AgentPay(cmd.Context(), agentAddr, poolID, payAmount, pk)
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

	// Add the newAmount to the original amount
	amount.Add(amount, newAmount)

	return amount
}

func init() {
	agentCmd.AddCommand(exitCmd)
	exitCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	exitCmd.Flags().String("from", "", "address to send the transaction from")
}
