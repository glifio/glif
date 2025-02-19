/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

var txSpeedUpCmd = &cobra.Command{
	Use:   "speed-up <tx hash>",
	Short: "Replaces a Eth transaction in the mempool with a higher premium",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		var cid cid.Cid

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		var ethHash common.Hash
		if strings.HasPrefix(args[0], "0x") {
			ethHashFil, err := ethtypes.ParseEthHash(args[0])
			if err != nil {
				logFatal(err)
			}
			msgCid, err := lapi.EthGetMessageCidByTransactionHash(ctx, &ethHashFil)
			if err != nil {
				logFatal(err)
			}
			if msgCid == nil {
				logFatal("Message not found.")
			}
			cid = *msgCid
			ethHash = common.HexToHash(ethHashFil.String())
		} else {
			cid.UnmarshalText([]byte(args[0]))
			ethHashPtr, err := lapi.EthGetTransactionHashByCid(ctx, cid)
			if err != nil {
				logFatal(err)
			}
			if ethHashPtr == nil {
				logFatal("No Eth hash found for CID.")
			}
			ethHash = common.HexToHash(ethHashPtr.String())
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

				oldTx, _, err := ethClient.TransactionByHash(ctx, ethHash)
				if err != nil {
					logFatal(err)
				}

				estimatedGas, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
					From:      from,
					To:        &from,
					GasTipCap: gasTipCap,
					GasFeeCap: gasFeeCap,
					Value:     oldTx.Value(),
					Data:      oldTx.Data(),
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
					Value:     oldTx.Value(),
					Data:      oldTx.Data(),
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
	txCmd.AddCommand(txSpeedUpCmd)
}
