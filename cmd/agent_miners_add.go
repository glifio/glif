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
	"github.com/glifio/cli/events"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <miner address>",
	Short: "Add a miner id to the agent",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}

		log.Printf("Adding miner %s to agent %s", minerAddr, agentAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		addminerevt := journal.RegisterEventType("agent", "addminer")
		evt := &events.AgentAddMiner{
			AgentID: agentAddr.String(),
			MinerID: minerAddr.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(addminerevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentAddMiner(cmd.Context(), agentAddr, minerAddr, ownerKey, requesterKey)
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

		fmt.Printf("Successfully added miner %s to agent %s", minerAddr, agentAddr)
	},
}

func init() {
	minersCmd.AddCommand(addCmd)
}
