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
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

var acceptOperatorCmd = &cobra.Command{
	Use:   "accept-operator",
	Short: "Approves an operator change on the Agent",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		as := util.AccountsStore()

		opEvm, opFevm, err := as.GetAddrs(string(util.OperatorKey))
		if err != nil {
			if err == util.ErrKeyNotFound {
				logFatal("agent operator not found in wallet")
			}
			logFatal(err)
		}

		agentAddr, auth, _, _, err := commonOwnerOrOperatorSetup(ctx, from)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		exitevt := journal.RegisterEventType("agent", "admin")
		evt := &events.AgentAdmin{
			Action:  "accept-operator",
			AgentID: agentAddr.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(exitevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentAcceptOperator(ctx, auth, agentAddr)
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

		log.Printf("Successfully accepted operator change on agent %s\n", agentAddr.String())
	},
}

func init() {
	adminCmd.AddCommand(acceptOperatorCmd)
}
