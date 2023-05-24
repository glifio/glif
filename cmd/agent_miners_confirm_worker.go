/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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

// changeWorkerCmd represents the changeWorker command
var confirmWorker = &cobra.Command{
	Use:   "confirm-worker <miner-addr>",
	Short: "Confirm the worker address change of your miner",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Confirming worker address change for miner %s", minerAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().AgentConfirmMinerWorkerChange(cmd.Context(), agentAddr, minerAddr, ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Println("Successfully confirmed worker change")
	},
}

func init() {
	minersCmd.AddCommand(confirmWorker)
}
