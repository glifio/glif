/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/events"
	"github.com/spf13/cobra"
)

var pushFundsCmd = &cobra.Command{
	Use:   "push-funds <miner address> <amount>",
	Short: "Push FIL from the Glif Agent to a specific Miner ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		from := cmd.Flag("from").Value.String()
		agentAddr, auth, _, requesterKey, err := commonOwnerOrOperatorSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		minerAddr, err := ToMinerID(ctx, args[0])
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[1])
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		pushevt := journal.RegisterEventType("agent", "push")
		evt := &events.AgentMinerPush{
			AgentID: agentAddr.String(),
			MinerID: minerAddr.String(),
			Amount:  amount.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(pushevt, func() interface{} { return evt })

		tx, err := PoolsSDK.Act().AgentPushFunds(ctx, auth, agentAddr, amount, minerAddr, requesterKey)
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

		fmt.Printf("Successfully pushed funds down to miner %s\n", minerAddr)
	},
}

func init() {
	minersCmd.AddCommand(pushFundsCmd)
	pushFundsCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
