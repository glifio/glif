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

var withdrawCmd = &cobra.Command{
	Use:   "withdraw <amount> <receiver>",
	Short: "Withdraw FIL from your Agent.",
	Long:  "",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, auth, _, requesterKey, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		receiver, err := AddressOrAccountNameToEVM(cmd.Context(), args[1])
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		withdrawevt := journal.RegisterEventType("agent", "withdraw")
		evt := &events.AgentWithdraw{
			AgentID: agentAddr.String(),
			Amount:  amount.String(),
			To:      receiver.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(withdrawevt, func() interface{} { return evt })

		fmt.Printf("Withdrawing %s FIL from your Agent", args[0])

		tx, err := PoolsSDK.Act().AgentWithdraw(cmd.Context(), auth, agentAddr, receiver, amount, requesterKey)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}
		evt.Tx = tx.Hash().String()

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully withdrew %s FIL\n", args[0])
	},
}

func init() {
	agentCmd.AddCommand(withdrawCmd)
}
