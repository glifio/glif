/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/cli/events"
	"github.com/spf13/cobra"
)

// borrowCmd represents the borrow command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw <amount>",
	Short: "Withdraw FIL from your Agent.",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		receiver, err := PoolsSDK.Query().AgentOwner(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		if !common.IsHexAddress(receiver.String()) {
			logFatal("Invalid withdraw address")
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

		tx, err := PoolsSDK.Act().AgentWithdraw(cmd.Context(), agentAddr, receiver, amount, ownerKey, requesterKey)
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

		fmt.Printf("Successfully withdrew %s FIL", args[0])
	},
}

func init() {
	agentCmd.AddCommand(withdrawCmd)
}
