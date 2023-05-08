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
	"github.com/filecoin-project/go-address"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add-miner [miner address]",
	Short: "Add a miner id to the agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AgentStore()
		ks := util.KeyStore()
		// Check if an agent already exists
		agentAddrStr, err := as.Get("address")
		if err != nil {
			log.Fatal(err)
		}

		if agentAddrStr == "" {
			log.Fatalf("Did you forget to create your agent? ")
		}

		agentAddr := common.HexToAddress(agentAddrStr)

		pk, err := ks.GetPrivate(util.OwnerKey)
		if err != nil {
			log.Fatal(err)
		}

		if pk == nil {
			log.Fatal("Owner key not found. Please check your `keys.toml` file. Only an Agent's owner can add a miner to it")
		}

		if len(args) != 1 {
			log.Fatal("Please provide a miner address")
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		fmt.Printf("Adding miner %s to agent %s", minerAddr, agentAddr)

		tx, err := fevm.Connection().AddMiner(cmd.Context(), agentAddr, minerAddr)
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

		fmt.Printf("Successfully added miner %s to agent %s", minerAddr, agentAddr)
	},
}

func init() {
	agentCmd.AddCommand(addCmd)
}
