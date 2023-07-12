/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/cli/util"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var forwardFIL = &cobra.Command{
	Use:   "forward-fil",
	Short: "Transfers balances from owner or operator wallet to another address through the FilForwarder smart contract",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}
		defer ethClient.Close()

		ks := util.KeyStore()

		toStr := cmd.Flag("to").Value.String()

		to, err := ParseAddressToNative(cmd.Context(), toStr)
		if err != nil {
			logFatal(err)
		}

		value, err := parseFILAmount(cmd.Flag("value").Value.String())
		if err != nil {
			logFatal(err)
		}

		if value.Cmp(common.Big0) < 1 {
			logFatal(errors.New("Value must not be 0"))
		}

		if toStr == to.String() {
			fmt.Printf("Forwarding %0.09f FIL to %s\n", denoms.ToFIL(value), to.String())
		} else {
			fmt.Printf("Forwarding %0.09f FIL to %s (converted to %s)\n", denoms.ToFIL(value), toStr, to.String())
		}
		fmt.Println("(Note that on block explorers, the transaction's `to` address will be the FilForwarder smart contract address, which will forward the funds to the receiver address)")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		// from must either be `owner` or `operator` in this limited transfer cmd
		keyToUse := util.KeyType(cmd.Flag("from").Value.String())
		if keyToUse != util.OwnerKey && keyToUse != util.OperatorKey {
			logFatal(errors.New("Unsupported `from` flag - must be `owner` or `operator`"))
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

		var filForwardAddr common.Address
		switch PoolsSDK.Query().ChainID().Int64() {
		case constants.MainnetChainID:
			filForwardAddr = deploy.FilForwarder
		case constants.CalibnetChainID:
			filForwardAddr = deploy.TFilForwarder
		default:
			logFatal(errors.New("unsupported chain id for forward-fil command"))
		}

		// get the FilForwarder contract address
		filf, err := abigen.NewFilForwarderTransactor(filForwardAddr, ethClient)
		if err != nil {
			logFatal(err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(pk, PoolsSDK.Query().ChainID())
		if err != nil {
			logFatal(err)
		}

		auth.Nonce = nonce
		auth.Value = value

		tx, err := filf.Forward(auth, to.Bytes())
		if err != nil {
			logFatal(err)
		}
		s.Stop()

		fmt.Printf("Forward FIL transaction sent: %s\n", tx.Hash().Hex())
		fmt.Println("Waiting for transaction to confirm...")

		s.Start()

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		fmt.Println("Success!")
	},
}

func init() {
	walletCmd.AddCommand(forwardFIL)
	forwardFIL.Flags().String("from", "", "From address (`owner` or `operator`)")
	forwardFIL.MarkFlagRequired("from")

	forwardFIL.Flags().String("to", "", "To address")
	forwardFIL.MarkFlagRequired("to")

	forwardFIL.Flags().String("value", "", "Value to send (in FIL)")
	forwardFIL.MarkFlagRequired("value")
}
