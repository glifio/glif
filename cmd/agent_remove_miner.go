/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/filecoin-project/go-address"
	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var rmCmd = &cobra.Command{
	Use:   "rm-miner [miner address] [new owner address]",
	Short: "Remove a miner from your agent",
	Long:  "Removes a specific miner from your Agent by assigning its owner to `new owner address`",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) != 2 {
			log.Fatal("Please provide a miner and recipient address")
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}

		newMinerOwnerAddr, err := address.NewFromString(args[1])
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT: an ethereum address can not be an owner of a miner, this must be a filecoin address owner
		if newMinerOwnerAddr.Protocol() == address.Delegated {
			log.Fatal("New miner owner address must be a filecoin address, not a delegated address")
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		fmt.Printf("Removing miner %s from agent %s by changing its owner address to %s", minerAddr, agentAddr, newMinerOwnerAddr)

		tx, err := fevm.Connection().RemoveMiner(cmd.Context(), agentAddr, minerAddr, newMinerOwnerAddr, ownerKey)
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

		fmt.Printf("Successfully proposed an ownership change to miner %s, passing %s as the new owner", minerAddr, newMinerOwnerAddr)
	},
}

func init() {
	agentCmd.AddCommand(rmCmd)
}
