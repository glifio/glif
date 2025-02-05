/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/filecoin-project/lotus/lib/tablewriter"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var listPendingCmd = &cobra.Command{
	Use:   "list-pending <account or address>",
	Short: "Lists pending transactions in the mempool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		var addr common.Address
		var err error
		if strings.HasPrefix(args[0], "0x") {
			addr = common.HexToAddress(args[0])
		} else {
			as := util.AccountsStore()
			addr, _, err = as.GetAddrs(args[0])
			if err != nil {
				var e *util.ErrKeyNotFound
				if errors.As(err, &e) {
					logFatal("Account not found in wallet")
				}
				logFatal(err)
			}
		}

		filecoinAddr, err := ethtypes.ParseEthAddress(addr.String())
		if err != nil {
			logFatal(err)
		}

		delegatedAddr, err := filecoinAddr.ToFilecoinAddress()
		if err != nil {
			logFatal(err)
		}

		// fmt.Printf("From: %v (%v)\n", addr, delegatedAddr)

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		msgs, err := lapi.MpoolPending(ctx, types.EmptyTSK)
		if err != nil {
			logFatal(err)
		}

		tw := tablewriter.New(
			tablewriter.Col("Nonce"),
			tablewriter.Col("Transaction"),
			tablewriter.Col("Gas Premium"),
			tablewriter.Col("Gas Fee Cap"),
		)

		count := 0
		for _, msg := range msgs {
			if msg.Message.From == delegatedAddr {
				var txStr string
				cid := msg.Cid()
				ethHash, _ := lapi.EthGetTransactionHashByCid(ctx, cid)
				if ethHash != nil {
					txStr = ethHash.String()
				} else {
					txStr = cid.String()
				}
				tw.Write(map[string]interface{}{
					"Nonce":       msg.Message.Nonce,
					"Transaction": txStr,
					"Gas Premium": msg.Message.GasPremium,
					"Gas Fee Cap": msg.Message.GasFeeCap,
				})
				count++
			}
		}
		if count > 0 {
			tw.Flush(os.Stdout)
		} else {
			fmt.Println("No pending transactions found in mempool.")
		}
	},
}

func init() {
	txCmd.AddCommand(listPendingCmd)
}
