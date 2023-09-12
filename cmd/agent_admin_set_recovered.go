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

var setRecoveredCmd = &cobra.Command{
	Use:   "set-recovered",
	Short: "Sets the Agent back into good standing",
	Long:  "If the Agent recovers from being in a faulty state, this command marks the Agent as healthy again.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		agentAddr, auth, _, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		exitevt := journal.RegisterEventType("agent", "admin")
		evt := &events.AgentAdmin{
			Action:  "set-recovered",
			AgentID: agentAddr.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(exitevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentSetRecovered(ctx, auth, agentAddr, requesterKey)
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

		log.Println("Successfully recovered agent: ", agentAddr.String())
	},
}

func init() {
	agentCmd.AddCommand(setRecoveredCmd)
}
