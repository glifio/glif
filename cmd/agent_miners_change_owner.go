/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var changeOwnerCmd = &cobra.Command{
	Use:   "change-owner [miner address]",
	Short: "Proposes an ownership change to your miner to prepare it for pledging to the Agent.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, _, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) != 1 {
			log.Fatal("Please provide a miner address")
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}

		lapi, closer, err := fevm.Connection().ConnectLotusClient()
		if err != nil {
			log.Fatal(err)
		}
		defer closer()

		ethAddr, err := ethtypes.ParseEthAddress(agentAddr.String())
		if err != nil {
			log.Fatal(err)
		}

		delegated, err := ethAddr.ToFilecoinAddress()
		if err != nil {
			log.Fatal(err)
		}

		id, err := lapi.StateLookupID(cmd.Context(), delegated, types.EmptyTSK)
		if err != nil {
			log.Fatal(err)
		}

		mi, err := lapi.StateMinerInfo(cmd.Context(), minerAddr, types.EmptyTSK)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Miner Owner:", mi.Owner)

		sp, err := actors.SerializeParams(&id)
		if err != nil {
			log.Fatal(err)
		}

		smsg, err := lapi.MpoolPushMessage(cmd.Context(), &types.Message{
			From:   mi.Owner,
			To:     minerAddr,
			Method: builtin.MethodsMiner.ChangeOwnerAddress,
			Value:  big.Zero(),
			Params: sp,
		}, nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Message CID:", smsg.Cid())

		wait, err := lapi.StateWaitMsg(cmd.Context(), smsg.Cid(), build.MessageConfidence, 900, true)
		if err != nil {
			log.Fatal(err)
		}

		// check it executed successfully
		if wait.Receipt.ExitCode != 0 {
			log.Fatal(err)
		}

		fmt.Println("message succeeded!")
	},
}

func init() {
	minersCmd.AddCommand(changeOwnerCmd)
}
