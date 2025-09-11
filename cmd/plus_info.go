package cmd

import (
	"fmt"
	"math/big"
	"time"

	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints information about the GLIF Card",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		tokenID, err := getPlusTokenID()
		if err != nil {
			logFatal(err)
		}

		info, err := PoolsSDK.Query().SPPlusInfo(ctx, big.NewInt(tokenID), nil)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("GLIF Card Token ID: %d\n\n", tokenID)

		fmt.Printf("Tier: %s\n", tierName(info.Tier))
		fmt.Printf("Locked Amount: %.09f GLF\n", poolsutil.ToFIL(info.TierLockAmount))
		if info.WithdrawableExtraLockedFunds.Sign() == 1 {
			fmt.Printf("Withdrawable Extra Locked Funds: %.09f GLF\n", poolsutil.ToFIL(info.WithdrawableExtraLockedFunds))
		}

		penaltyWindow, penaltyFee, err := PoolsSDK.Query().SPPlusTierSwitchPenaltyInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		windowStart, windowEnd, days, hours := getTierSwitchWindow(info, penaltyWindow)

		if info.Tier > 0 {
			windowStartFormatted := windowStart.UTC().Format("January 2 2006 15:04")
			fmt.Printf("Tier activation timestamp: %v\n", windowStartFormatted)
			if windowEnd.After(time.Now()) {
				windowEndFormatted := windowEnd.UTC().Format("January 2 2006 15:04")
				fmt.Printf("Free downgrade after %v UTC (%d days, %d hours)\n", windowEndFormatted, days, hours)
				penaltyBasis, _ := penaltyFee.Float64()
				fmt.Printf("Early downgrade penalty fee: %.02f%%\n", penaltyBasis/100.00)
			} else {
				fmt.Println("Free downgrade available.")
			}
		}

		cashbackBasis, _ := info.PersonalCashBackPercent.Float64()
		fmt.Printf("\nPersonal Cashback Percentage: %.02f%%\n", cashbackBasis/100.00)
		fmt.Printf("Cashback earned: %.09f FIL\n", poolsutil.ToFIL(info.FilCashbackEarned))
		fmt.Printf("Vault balance: %.09f GLF\n", poolsutil.ToFIL(info.GLFVaultBalance))
		fmt.Printf("Base Conversion Rate: 1 FIL = %.09f GLF\n", poolsutil.ToFIL(info.BaseConversionRateFILtoGLF))
	},
}

func init() {
	plusCmd.AddCommand(plusInfoCmd)
}
