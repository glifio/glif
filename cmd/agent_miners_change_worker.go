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
	"github.com/glifio/glif/v2/events"
	"github.com/spf13/cobra"
)

// changeWorkerCmd represents the changeWorker command
var changeWorkerCmd = &cobra.Command{
	Use:   "change-worker <miner address> <worker address> [control addresses...]",
	Short: "Change the worker address of your miner",
	Long:  ``,
	Args:  cobra.RangeArgs(2, 5),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		minerAddr, err := ToMinerID(cmd.Context(), args[0])
		if err != nil {
			logFatal(err)
		}
		log.Println(minerAddr)

		workerAddr, err := ToMinerID(cmd.Context(), args[1])
		if err != nil {
			log.Print("Error parsing worker address")
			logFatal(err)
		}
		log.Println(workerAddr)

		var controlAddrs []address.Address
		for _, arg := range args[2:] {
			controlAddr, err := ToMinerID(cmd.Context(), arg)
			if err != nil {
				log.Print("Error parsing control address")
				logFatal(err)
			}
			controlAddrs = append(controlAddrs, controlAddr)
		}

		log.Printf("Changing worker address for miner %s to %s\n", minerAddr, workerAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		changeworkerevt := journal.RegisterEventType("miner", "changeworker")
		evt := &events.AgentMinerChangeWorker{
			AgentID:    agentAddr.String(),
			MinerID:    minerAddr.String(),
			NewWorker:  workerAddr.String(),
			NewControl: AddressesToStrings(controlAddrs),
		}
		defer journal.Close()
		defer journal.RecordEvent(changeworkerevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentChangeMinerWorker(cmd.Context(), auth, agentAddr, minerAddr, workerAddr, controlAddrs)
		if err != nil {
			evt.Error = err.Error()
			logFatalf("tx error: %s", err)
		}
		evt.Tx = tx.Hash().String()

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		s.Stop()

		fmt.Println("Successfully changed miner worker - you must confirm this change yourself using `glif agent miners confirm-worker`")
	},
}

func init() {
	minersCmd.AddCommand(changeWorkerCmd)
}
