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

var transferOwnershipCmd = &cobra.Command{
	Use:   "transfer-ownership <new-owner>",
	Short: "Proposes an ownership change to the Agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		newOwner, err := AddressOrAccountNameToEVM(ctx, args[0])
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
			Action:          "transfer-ownership",
			AgentID:         agentAddr.String(),
			NewAdminAddress: newOwner.Hex(),
		}
		defer journal.Close()
		defer journal.RecordEvent(exitevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentTransferOwnership(ctx, auth, agentAddr, newOwner)
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

		log.Printf("Successfully proposed ownership change to agent %s, new owner %s\n", agentAddr.String(), newOwner.Hex())
	},
}

func init() {
	adminCmd.AddCommand(transferOwnershipCmd)
}
