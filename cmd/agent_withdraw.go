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
	"github.com/spf13/cobra"
)

// borrowCmd represents the borrow command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw <amount>",
	Short: "Withdraw FIL from your Agent.",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		var receiver common.Address
		if cmd.Flag("to") != nil && cmd.Flag("to").Changed {
			receiver = common.HexToAddress(cmd.Flag("to").Value.String())
		} else {
			// if no recipient is specified, use the agent's owner
			receiver, err = PoolsSDK.Query().AgentOwner(cmd.Context(), agentAddr)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("receiver: %s", receiver.String())

		if common.IsHexAddress(receiver.String()) {
			log.Fatal("Invalid withdraw address")
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		fmt.Printf("Withdrawing %s FIL from your Agent...", args[0])

		tx, err := PoolsSDK.Act().AgentWithdraw(cmd.Context(), agentAddr, receiver, amount, ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully withdrew %s FIL", args[0])
	},
}

func init() {
	agentCmd.AddCommand(withdrawCmd)
	withdrawCmd.Flags().String("to", "", "Receiver of the funds (agent's owner is chosen if not specified)")
}
