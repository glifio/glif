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

		var tokenID int64
		var err error

		tokenIDFlag, err := cmd.Flags().GetInt64("token-id")
		if err != nil {
			logFatal(err)
		}

		if tokenIDFlag > 0 {
			tokenID = tokenIDFlag
		} else {
			tokenID, err = getPlusTokenID()
			if err != nil {
				logFatal(err)
			}
		}

		info, err := PoolsSDK.Query().SPPlusInfo(ctx, big.NewInt(tokenID), nil)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("GLIF Card Token ID: %d\n", tokenID)

		// TIER INFORMATION
		fmt.Printf("\n── Tier Information ──\n")
		fmt.Printf("Tier: %s\n", tierName(info.Tier))
		fmt.Printf("Locked Amount: %.09f GLF\n", poolsutil.ToFIL(info.TierLockAmount))
		if info.WithdrawableExtraLockedFunds.Sign() == 1 {
			fmt.Printf("Withdrawable Extra: %.09f GLF\n", poolsutil.ToFIL(info.WithdrawableExtraLockedFunds))
		}

		// TIER SWITCH TIMING
		penaltyWindow, _, err := PoolsSDK.Query().SPPlusTierSwitchPenaltyInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		windowStart, windowEnd, days, hours := getTierSwitchWindow(info, penaltyWindow)

		if info.Tier > 0 {
			fmt.Printf("\n── Tier Switch Info ──\n")
			windowStartFormatted := windowStart.UTC().Format("January 2 2006 15:04")
			fmt.Printf("Activated: %v\n", windowStartFormatted)
			if windowEnd.After(time.Now()) {
				windowEndFormatted := windowEnd.UTC().Format("January 2 2006 15:04")
				fmt.Printf("Free downgrade: %v UTC (%dd %dh)\n", windowEndFormatted, days, hours)
			} else {
				fmt.Printf("Free downgrade: Available now\n")
			}
		}

		filVaultBalance, err := PoolsSDK.Query().SPPlusFILVaultBalance(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		// CASHBACK INFORMATION
		fmt.Printf("\n── Cash Back Status ──\n")
		cashbackBasis, _ := info.PersonalCashBackPercent.Float64()
		fmt.Printf("Earned: %.09f FIL\n", poolsutil.ToFIL(info.FilCashbackEarned))
		fmt.Printf("Vault Balance: %.09f GLF\n", poolsutil.ToFIL(info.GLFVaultBalance))
		fmt.Printf("Cash Back Percentage: %.02f%%\n", cashbackBasis/100.00)

		// CONVERSION RATES
		fmt.Printf("\n── Conversion Rates ──\n")
		fmt.Printf("Base Rate: 1 FIL = %.09f GLF\n", poolsutil.ToFIL(info.BaseConversionRateFILtoGLF))

		// Get tier information to calculate tier premium rate
		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		if info.Tier > 0 && int(info.Tier) <= len(tierInfos) {
			tierInfo := tierInfos[info.Tier]

			// Calculate tier premium conversion rate using WAD math (18 decimals)
			// This matches the SpPlus contract logic: conversionRateWithPremium = filToGlf.rawMulWad(tierInfo.cashBackPremium)
			conversionRateWithPremium := poolsutil.MulWad(info.BaseConversionRateFILtoGLF, tierInfo.CashBackPremium)

			premium := new(big.Float).Mul(
				new(big.Float).Sub(poolsutil.ToFIL(tierInfo.CashBackPremium), big.NewFloat(1)),
				big.NewFloat(100),
			)

			fmt.Printf("Tier Rate: 1 FIL = %.09f GLF (+%.02f%%)\n",
				poolsutil.ToFIL(conversionRateWithPremium), premium)
		}

		fmt.Printf("\n── Cash back program vault balance ──\n")
		fmt.Printf("Total FIL available for cash back program: %.09f FIL\n", poolsutil.ToFIL(filVaultBalance))

	},
}

func init() {
	plusCmd.AddCommand(plusInfoCmd)
	plusInfoCmd.Flags().Int64("token-id", 0, "Token ID to query (defaults to your stored token)")
}
