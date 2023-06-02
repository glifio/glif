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
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var rmCmd = &cobra.Command{
	Use:   "remove <miner address> <new owner address>",
	Short: "Remove a miner from your agent",
	Long: `Removes a specific miner from your Agent by assigning its owner to "new owner address". 
	The new owner address must be a filecoin address, not a delegated address.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
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
		defer s.Stop()

		fmt.Printf("Removing miner %s from agent %s by changing its owner address to %s", minerAddr, agentAddr, newMinerOwnerAddr)

		tx, err := PoolsSDK.Act().AgentRemoveMiner(cmd.Context(), agentAddr, minerAddr, newMinerOwnerAddr, ownerKey, requesterKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully proposed an ownership change to miner %s, passing %s as the new owner", minerAddr, newMinerOwnerAddr)
	},
}

func init() {
	minersCmd.AddCommand(rmCmd)
}
