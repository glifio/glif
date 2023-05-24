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
var changeWorkerCmd = &cobra.Command{
	Use:   "change-worker",
	Short: "Change the worker address of your miner",
	Long:  ``,
	Args:  cobra.RangeArgs(2, 5),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, _, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}

		workerAddr, err := address.NewFromString(args[1])
		if err != nil {
			log.Fatal(err)
		}

		var controlAddrs []address.Address
		for _, arg := range args[2:] {
			controlAddr, err := address.NewFromString(arg)
			if err != nil {
				log.Fatal(err)
			}
			controlAddrs = append(controlAddrs, controlAddr)
		}

		log.Printf("Changing worker address for miner %s to %s", minerAddr, workerAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().AgentChangeMinerWorker(cmd.Context(), agentAddr, minerAddr, workerAddr, controlAddrs, ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Println("Successfully changed miner worker - you must confirm this change yourself using `glif agent miners confirm-worker-change`")
	},
}

func init() {
	minersCmd.AddCommand(changeWorkerCmd)
}
