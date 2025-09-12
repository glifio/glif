package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusClaimCashBackCmd = &cobra.Command{
	Use:   "claim-cashback <receiver address>",
	Short: "Transfer earned FIL cashback to receiver address",
	Args:  cobra.ExactArgs(1),
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

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		receiver, err := AddressOrAccountNameToEVM(cmd.Context(), args[0])
		if err != nil {
			logFatal(err)
		}

		tx, err := PoolsSDK.Act().SPPlusClaimCashBack(ctx, auth, big.NewInt(tokenID), receiver)
		if err != nil {
			logFatalf("Failed to claim cashback %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to claim cashback %s", err)
		}

		s.Stop()

		fmt.Println("FIL cashback successfully claimed.")
	},
}

func init() {
	plusCashBackCmd.AddCommand(plusClaimCashBackCmd)
}
