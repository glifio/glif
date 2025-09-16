package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusWithdrawExtraLockedFundsCmd = &cobra.Command{
	Use:   "withdraw-extra-locked-funds",
	Short: "Withdraw extra locked GLF when price of tier decreases",
	Args:  cobra.NoArgs,
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

		tx, err := PoolsSDK.Act().SPPlusWithdrawExtraLockedFunds(ctx, auth, big.NewInt(tokenID))
		if err != nil {
			logFatalf("Failed to withdraw extra locked funds %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to withdraw extra locked funds %s", err)
		}

		s.Stop()

		fmt.Println("Extra locked funds withdrawn.")
	},
}

func init() {
	plusAdvancedCmd.AddCommand(plusWithdrawExtraLockedFundsCmd)
}
