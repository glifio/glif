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

// pull represents the pull command
var pullFundsCmd = &cobra.Command{
	Use:   "pull-funds <amount> <miner address>",
	Short: "Pull FIL from a miner into your Glif Agent",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := fevm.Connection().ToMinerID(args[1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Pulling %s FIL from %s", amount.String(), minerAddr.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := fevm.Connection().AgentPullFunds(cmd.Context(), agentAddr, amount, minerAddr, pk)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		receipt, err := fevm.WaitReturnReceipt(tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		if receipt.Status == 0 {
			log.Fatal("Transaction failed")
		}

		s.Stop()

		fmt.Printf("Successfully pull funds up from miner %s", minerAddr)
	},
}

func init() {
	minersCmd.AddCommand(pullFundsCmd)
	pullFundsCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
