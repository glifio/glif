/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var pushFundsCmd = &cobra.Command{
	Use:   "push-funds <miner address> <amount>",
	Short: "Push FIL from the Glif Agent to a specific Miner ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, senderKey, requesterKey, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := ToMinerID(cmd.Context(), args[0])
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[1])
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().AgentPushFunds(cmd.Context(), agentAddr, amount, minerAddr, senderKey, requesterKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully pushed funds down to miner %s", minerAddr)
	},
}

func init() {
	minersCmd.AddCommand(pushFundsCmd)
	pushFundsCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
