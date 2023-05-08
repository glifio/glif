/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

// borrowCmd represents the borrow command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw [amount]",
	Short: "Withdraw FIL from your Agent.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		receiver := agentAddr
		if cmd.Flag("recipient").Changed {
			receiver = common.HexToAddress(cmd.Flag("recipient").Value.String())
		}

		if len(args) != 1 {
			log.Fatal("Please provide an amount")
		}

		amount, err := parseFILAmount(args[1])
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		fmt.Printf("Withdrawing %s FIL from your Agent...", args[0])

		tx, err := fevm.Connection().AgentWithdraw(cmd.Context(), agentAddr, receiver, amount, ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		receipt := fevm.WaitReturnReceipt(tx.Hash())
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
	agentCmd.AddCommand(withdrawCmd)
	createCmd.Flags().String("recipient", "", "Receiver of the funds (agent's owner is chosen if not specified)")
}
