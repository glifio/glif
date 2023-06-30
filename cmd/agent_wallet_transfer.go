/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfers balances from owner or operator wallet to another FEVM address",
	Long:  "Transfers balances from owner or operator wallet to another FEVM address. This feature is in beta, and only certain types of addresses are supported at this time. Full support for GLIF CLI wallet is coming soon",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}
		defer ethClient.Close()

		ks := util.KeyStore()

		to, err := MustBeEVMAddr(cmd.Flag("to").Value.String())
		if err != nil {
			logFatal(errors.New("Unsupported `to` flag - only 0x addresses supported in this limited transfer cmd"))
		}

		value, err := parseFILAmount(cmd.Flag("value").Value.String())
		if err != nil {
			logFatal(err)
		}

		if value.Cmp(common.Big0) < 1 {
			logFatal(errors.New("Value must not be 0"))
		}

		// from must either be `owner` or `operator` in this limited transfer cmd
		keyToUse := util.KeyType(cmd.Flag("from").Value.String())
		if keyToUse != util.OwnerKey && keyToUse != util.OperatorKey {
			logFatal(errors.New("Unsupported `from` flag - must be `owner` or `operator` in this limited transfer cmd"))
		}

		pk, err := ks.GetPrivate(keyToUse)
		if err != nil {
			logFatal(err)
		}

		fromAddr, _, err := ks.GetAddrs(keyToUse)
		if err != nil {
			logFatal(err)
		}

		nonce, err := PoolsSDK.Query().ChainGetNonce(cmd.Context(), fromAddr)
		if err != nil {
			logFatal(err)
		}

		gasPrice, err := ethClient.SuggestGasPrice(cmd.Context())
		if err != nil {
			logFatal(err)
		}

		gasLimit := uint64(21000)

		tx := types.NewTransaction(nonce.Uint64(), to, value, gasLimit, gasPrice, []byte{})

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(PoolsSDK.Query().ChainID()), pk)
		if err != nil {
			logFatal(err)
		}

		err = ethClient.SendTransaction(cmd.Context(), signedTx)
		if err != nil {
			logFatal(err)
		}

		fmt.Println("Transfer sent: ", signedTx.Hash().Hex())
	},
}

func init() {
	walletCmd.AddCommand(transferCmd)
	transferCmd.Flags().String("from", "", "From address (`owner` or `operator`)")
	transferCmd.MarkFlagRequired("from")

	transferCmd.Flags().String("to", "", "To address (must be a 0x address)")
	transferCmd.MarkFlagRequired("to")

	transferCmd.Flags().String("value", "", "Value to send (in FIL)")
	transferCmd.MarkFlagRequired("value")
}
