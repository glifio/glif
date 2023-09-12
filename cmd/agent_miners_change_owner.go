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
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glifio/cli/events"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var changeOwnerCmd = &cobra.Command{
	Use:   "change-owner <miner address>",
	Short: "Proposes an ownership change to your miner to prepare it for pledging to the Agent.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		ethAddr, err := ethtypes.ParseEthAddress(agentAddr.String())
		if err != nil {
			logFatal(err)
		}

		delegated, err := ethAddr.ToFilecoinAddress()
		if err != nil {
			logFatal(err)
		}

		id, err := lapi.StateLookupID(cmd.Context(), delegated, types.EmptyTSK)
		if err != nil {
			logFatal(err)
		}

		mi, err := lapi.StateMinerInfo(cmd.Context(), minerAddr, types.EmptyTSK)
		if err != nil {
			logFatal(err)
		}

		fmt.Println("Miner Owner:", mi.Owner)

		sp, err := actors.SerializeParams(&id)
		if err != nil {
			logFatal(err)
		}

		changeownerevt := journal.RegisterEventType("miner", "changeowner")
		evt := &events.AgentMinerChangeOwner{
			AgentID:  agentAddr.String(),
			MinerID:  minerAddr.String(),
			OldOwner: mi.Owner.String(),
			NewOwner: delegated.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(changeownerevt, func() interface{} { return evt })

		smsg, err := lapi.MpoolPushMessage(cmd.Context(), &types.Message{
			From:   mi.Owner,
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
	minersCmd.AddCommand(changeOwnerCmd)
}
