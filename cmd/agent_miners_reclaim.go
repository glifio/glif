/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/glifio/glif/events"
	"github.com/spf13/cobra"
)

var reclaimMinerCmd = &cobra.Command{
	Use:   "reclaim <miner address> <new-owner-address> --from [from]",
	Short: "Proposes an ownership change to your miner to complete the removal process of a miner from your Agent.",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}

		newOwnerAddr, err := address.NewFromString(args[1])
		if err != nil {
			logFatal(err)
		}
		if newOwnerAddr.Protocol() != address.ID {
			logFatal("new owner address must be an ID address")
		}

		senderAddr, err := address.NewFromString(cmd.Flag("from").Value.String())
		if err != nil {
			logFatal(err)
		}

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		sp, err := actors.SerializeParams(&newOwnerAddr)
		if err != nil {
			logFatal(err)
		}

		reclaimevt := journal.RegisterEventType("agent", "reclaim")
		evt := &events.AgentMinerReclaim{
			MinerID:  minerAddr.String(),
			NewOwner: newOwnerAddr.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(reclaimevt, func() interface{} { return evt })

		smsg, err := lapi.MpoolPushMessage(cmd.Context(), &types.Message{
			From:   senderAddr,
			To:     minerAddr,
			Method: builtin.MethodsMiner.ChangeOwnerAddress,
			Value:  big.Zero(),
			Params: sp,
		}, nil)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		fmt.Println("Message CID:", smsg.Cid())

		wait, err := lapi.StateWaitMsg(cmd.Context(), smsg.Cid(), build.MessageConfidence, 900, true)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		// check it executed successfully
		if wait.Receipt.ExitCode != 0 {
			evt.Error = err.Error()
			logFatal(err)
		}

		fmt.Println("message succeeded!")
	},
}

func init() {
	minersCmd.AddCommand(reclaimMinerCmd)
	reclaimMinerCmd.Flags().String("from", "", "specify from address")
	reclaimMinerCmd.MarkFlagRequired("from")
}
