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
			log.Fatal(err)
		}

		newOwnerAddr, err := address.NewFromString(args[1])
		if err != nil {
			log.Fatal(err)
		}
		if newOwnerAddr.Protocol() != address.ID {
			log.Fatal("new owner address must be an ID address")
		}

		senderAddr, err := address.NewFromString(cmd.Flag("from").Value.String())
		if err != nil {
			log.Fatal(err)
		}

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			log.Fatal(err)
		}
		defer closer()

		sp, err := actors.SerializeParams(&newOwnerAddr)
		if err != nil {
			log.Fatal(err)
		}

		smsg, err := lapi.MpoolPushMessage(cmd.Context(), &types.Message{
			From:   senderAddr,
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
	minersCmd.AddCommand(reclaimMinerCmd)
	reclaimMinerCmd.Flags().String("from", "", "specify from address")
	reclaimMinerCmd.MarkFlagRequired("from")
}
