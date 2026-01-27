//go:build advanced
// +build advanced

/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/events"
	"github.com/spf13/cobra"
)

var transferOperatorCmd = &cobra.Command{
	Use:   "transfer-operator <new-operator>",
	Short: "Proposes an operator change to the Agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		newOperator, err := AddressOrAccountNameToEVM(ctx, args[0])
		if err != nil {
			logFatal(err)
		}

		agentAddr, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		exitevt := journal.RegisterEventType("agent", "admin")
		evt := &events.AgentAdmin{
			Action:          "transfer-operator",
			AgentID:         agentAddr.String(),
			NewAdminAddress: newOperator.Hex(),
		}
		defer journal.Close()
		defer journal.RecordEvent(exitevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentTransferOperator(ctx, auth, agentAddr, newOperator)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}
		evt.Tx = tx.Hash().String()

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		s.Stop()

		log.Printf("Successfully proposed operator change to agent %s, new operator %s\n", agentAddr.String(), newOperator.Hex())
	},
}

func init() {
	adminCmd.AddCommand(transferOperatorCmd)
}
