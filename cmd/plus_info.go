package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints information about the GLIF Card",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		agentStore := util.AgentStore()

		tokenIDStr, err := agentStore.Get("plus-token-id")
		if err != nil && err.Error() != "key not found: plus-token-id" {
			logFatal(err)
		}

		if tokenIDStr == "" {
			fmt.Println("GLIF Card not minted yet.")
			return
		}

		tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
		if err != nil {
			logFatal(err)
		}

		info, err := PoolsSDK.Query().PlusInfo(ctx, big.NewInt(tokenID), nil)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("GLIF Card Token ID: %d\n\n", tokenID)

		fmt.Printf("Tier: %s\n", tierName(info.Tier))
		fmt.Printf("Locked Amount: %.9f GLF\n", poolsutil.ToFIL(info.TierLockAmount))
		if info.WithdrawableExtraLockedFunds.Sign() == 1 {
			fmt.Printf("Withdrawable Extra Locked Funds: %.9f GLF\n", poolsutil.ToFIL(info.WithdrawableExtraLockedFunds))
		}

		penaltyWindow, penaltyFee, err := PoolsSDK.Query().PlusTierSwitchPenaltyInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		windowStartSecs, err := strconv.ParseInt(info.LastTierSwitchTimestamp.String(), 10, 64)
		if err != nil {
			logFatal(err)
		}
		windowStart := time.Unix(windowStartSecs, 0)

		windowEnd := windowStart.Add(time.Duration(time.Duration(penaltyWindow.Int64()) * time.Second))
		hoursLeft := time.Until(windowEnd) / time.Hour
		days := int(hoursLeft) / 24
		hours := int(hoursLeft) % 24

		if info.Tier > 0 {
			fmt.Printf("Last tier switch timestamp: %v\n", windowStart.UTC())
			if windowEnd.After(time.Now()) {
				fmt.Printf("Free downgrade after %v (%d days, %d hours)\n", windowEnd.UTC(), days, hours)
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
		fmt.Printf("Base Conversion Rate: 1 FIL = %.9f GLF\n", poolsutil.ToFIL(info.BaseConversionRateFILtoGLF))
	},
}

func init() {
	plusCmd.AddCommand(plusInfoCmd)
}
