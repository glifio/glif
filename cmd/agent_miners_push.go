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

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push-funds [amount] [miner address]",
	Short: "Push FIL from the Glif Agent to a specific Miner ID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) != 2 {
			log.Fatal("Please provide an amount and a miner address")
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := fevm.Connection().ToMinerID(args[1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Pushing %s FIL to %s", amount.String(), minerAddr.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := fevm.Connection().AgentPushFunds(cmd.Context(), agentAddr, amount, minerAddr, pk)
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

		fmt.Printf("Successfully pushed funds down to miner %s", minerAddr)
	},
}

func init() {
	minersCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
