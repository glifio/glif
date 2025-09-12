package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusWithdrawGlfVaultCmd = &cobra.Command{
	Use:   "withdraw <amount> <receiver address>",
	Short: "Transfer GLF tokens from vault to receiver address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		tokenID, err := getPlusTokenID()
		if err != nil {
			logFatal(err)
		}

		_, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		receiver, err := AddressOrAccountNameToEVM(cmd.Context(), args[1])
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().SPPlusWithdrawGlfVault(ctx, auth, big.NewInt(tokenID), amount, receiver)
		if err != nil {
			logFatalf("Failed to withdraw from GLF vault %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to withdraw from GLF vault %s", err)
		}

		s.Stop()

		fmt.Println("GLF successfully withdrawn from vault.")
	},
}

func init() {
	plusCashBackCmd.AddCommand(plusWithdrawGlfVaultCmd)
}
