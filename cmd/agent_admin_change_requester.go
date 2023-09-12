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
	"github.com/glifio/cli/events"
	"github.com/spf13/cobra"
)

var changeRequesterCmd = &cobra.Command{
	Use:   "change-requester <new-requester-addr>",
	Short: "Changes the requester key on the Agent",
	Long:  "The `ADORequesterKey` is the key that is used to sign requests to the Agent Data Oracle. This command changes the key that signs requests for Signed Credentials from the Oracle.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		newRequester, err := AddressOrAccountNameToEVM(ctx, args[0])
		if err != nil {
			logFatal(err)
		}

		agentAddr, auth, _, _, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		exitevt := journal.RegisterEventType("agent", "admin")
		evt := &events.AgentAdmin{
			Action:          "change-requester",
			AgentID:         agentAddr.String(),
			NewAdminAddress: newRequester.Hex(),
		}
		defer journal.Close()
		defer journal.RecordEvent(exitevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentChangeRequester(ctx, auth, agentAddr, newRequester)
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

		log.Printf("Successfully changed requester key on the agent: %s, new requester: %s\n", agentAddr.String(), newRequester.Hex())
	},
}

func init() {
	adminCmd.AddCommand(changeRequesterCmd)
}
