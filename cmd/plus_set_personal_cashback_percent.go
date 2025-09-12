package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusSetPersonalCashBackPercentCmd = &cobra.Command{
	Use:   "set-personal-cashback-percent <percent>",
	Short: "Sets the cashback percentage",
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

		cashBackPercentFloat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			logFatal(err)
		}
		cashBackPercent := int64(cashBackPercentFloat * 100.00)

		tx, err := PoolsSDK.Act().SPPlusSetPersonalCashBackPercent(ctx, auth, big.NewInt(tokenID), big.NewInt(cashBackPercent))
		if err != nil {
			logFatalf("Failed to set personal cashback percent %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to set personal cashback percent %s", err)
		}

		s.Stop()

		fmt.Println("Personal cashback percent set.")
	},
}

func init() {
	plusAdvancedCmd.AddCommand(plusSetPersonalCashBackPercentCmd)
}
