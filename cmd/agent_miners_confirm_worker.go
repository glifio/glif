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
var confirmWorker = &cobra.Command{
	Use:   "confirm-worker <miner-addr>",
	Short: "Confirm the worker address change of your miner",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}

		log.Printf("Confirming worker address change for miner %s", minerAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		confirmworkerevt := journal.RegisterEventType("miner", "confirmworker")
		evt := &events.AgentMinerConfirmWorker{
			AgentID: agentAddr.String(),
			MinerID: minerAddr.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(confirmworkerevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentConfirmMinerWorkerChange(cmd.Context(), auth, agentAddr, minerAddr)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}
		evt.Tx = tx.Hash().String()

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		s.Stop()

		fmt.Println("Successfully confirmed worker change")
	},
}

func init() {
	minersCmd.AddCommand(confirmWorker)
}
