/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

var txCancelCmd = &cobra.Command{
	Use:   "cancel <tx hash or cid>",
	Short: "Replaces a transaction in the mempool with a dummy (zero value transfer to self)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		var cid cid.Cid

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		if strings.HasPrefix(args[0], "0x") {
			ethHash, err := ethtypes.ParseEthHash(args[0])
			if err != nil {
				logFatal(err)
			}
			msgCid, err := lapi.EthGetMessageCidByTransactionHash(ctx, &ethHash)
			if err != nil {
				logFatal(err)
			}
			if msgCid == nil {
				logFatal("Message not found.")
			}
			cid = *msgCid
		} else {
			cid.UnmarshalText([]byte(args[0]))
		}

		msgs, err := lapi.MpoolPending(ctx, types.EmptyTSK)
		if err != nil {
			logFatal(err)
		}

		for _, msg := range msgs {
			if msg.Cid() == cid {
				fromFilAddr := msg.Message.From

				fromEthAddr, err := lapi.FilecoinAddressToEthAddress(ctx, fromFilAddr)
				if err != nil {
					logFatal(err)
				}

				auth, senderAccount, err := commonGenericAccountSetup(cmd, fromEthAddr.String())
				if err != nil {
					logFatal(err)
				}

				from := senderAccount.Address

				ethClient, err := PoolsSDK.Extern().ConnectEthClient()
				if err != nil {
					logFatal(err)
				}

				mpoolCfg, err := lapi.MpoolGetConfig(ctx)
				if err != nil {
					logFatal(err)
				}

				gasTipCap := computeRBF(msg.Message.GasPremium, mpoolCfg.ReplaceByFeeRatio).Int

				ethHeader, err := ethClient.HeaderByNumber(ctx, nil)
				if err != nil {
					logFatal(err)
				}
				gasFeeCap := new(big.Int).Add(
					gasTipCap,
					new(big.Int).Mul(ethHeader.BaseFee, big.NewInt(2)),
				)

				estimatedGas, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
					From:      from,
					To:        &from,
					GasTipCap: gasTipCap,
					GasFeeCap: gasFeeCap,
					Value:     big.NewInt(0),
				})
				if err != nil {
					logFatal(err)
				}

				tx := ethcoretypes.NewTx(&ethcoretypes.DynamicFeeTx{
					ChainID:   big.NewInt(chainID),
					Nonce:     msg.Message.Nonce,
					GasTipCap: gasTipCap,
					GasFeeCap: gasFeeCap,
					Gas:       estimatedGas,
					To:        &from,
					Value:     big.NewInt(0),
				})

				signedTx, err := auth.Signer(from, tx)
				if err != nil {
					logFatal(err)
				}

				err = ethClient.SendTransaction(ctx, signedTx)
				if err != nil {
					logFatal(err)
				}

				fmt.Printf("Replacement transaction sent: %s\n", signedTx.Hash().Hex())

				return
			}
		}
		fmt.Println("No matching pending transactions found in mempool.")
	},
}

func init() {
	txCmd.AddCommand(txCancelCmd)
}
