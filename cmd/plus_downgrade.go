package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var acceptPenalty bool

var plusDowngradeCmd = &cobra.Command{
	Use:   "downgrade <new tier: inactive, bronze or silver>",
	Short: "Downgrade to a lower tier",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		tokenID, err := getPlusTokenID()
		if err != nil {
			logFatal(err)
		}

		tier, err := parseTierName(args[0])
		if err != nil {
			logFatal(err)
		}

		info, err := PoolsSDK.Query().SPPlusInfo(ctx, big.NewInt(tokenID), nil)
		if err != nil {
			logFatal(err)
		}

		if tier >= info.Tier {
			err := fmt.Errorf("new tier must be lower than current tier: %s", tierName(info.Tier))
			logFatal(err)
		}

		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}
		oldLockAmount := info.TierLockAmount
		newLockAmount := tierInfos[tier].TokenLockAmount

		err = printGlfOwnerBalance("GLF balance of owner before downgrade")
		if err != nil {
			logFatal(err)
		}
		fmt.Printf("GLF lock amount for %s tier: %.09f GLF\n", tierName(info.Tier), poolsutil.ToFIL(oldLockAmount))
		fmt.Printf("GLF lock amount for %s tier: %.09f GLF\n", tierName(tier), poolsutil.ToFIL(newLockAmount))
		refundGlf := new(big.Int).Sub(oldLockAmount, newLockAmount)

		penaltyWindow, penaltyFee, err := PoolsSDK.Query().SPPlusTierSwitchPenaltyInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		windowStart, windowEnd, days, hours := getTierSwitchWindow(info, penaltyWindow)

		if refundGlf.Sign() == 1 && windowEnd.After(time.Now()) {
			fmt.Println("Attempting to downgrade early...")
			windowStartFormatted := windowStart.UTC().Format("January 2 2006 15:04")
			fmt.Printf("Tier activation timestamp: %v\n", windowStartFormatted)
			windowEndFormatted := windowEnd.UTC().Format("January 2 2006 15:04")
			fmt.Printf("Free downgrade after %v UTC (%d days, %d hours)\n", windowEndFormatted, days, hours)
			penaltyAmount := new(big.Int).Div(
				new(big.Int).Mul(refundGlf, penaltyFee),
				big.NewInt(10000))
			fmt.Printf("Penalty fee: %.09f GLF\n", poolsutil.ToFIL(penaltyAmount))
			expectedRefund := new(big.Int).Sub(refundGlf, penaltyAmount)
			fmt.Printf("Refund with penalty: %.09f GLF\n", poolsutil.ToFIL(expectedRefund))
			if !acceptPenalty {
				logFatal("Re-run with --accept-penalty flag to pay penalty and proceed with early downgrade")
			}
		} else if refundGlf.Sign() == -1 {
			extraGlf := new(big.Int).Neg(refundGlf)
			fmt.Printf("GLF required to downgrade: %.09f GLF\n", poolsutil.ToFIL(extraGlf))

			err = checkGlfPlusBalanceAndAllowance(extraGlf)
			if err != nil {
				logFatal(err)
			}
		} else {
			downgradeAmount := new(big.Int).Sub(oldLockAmount, newLockAmount)
			fmt.Printf("GLF returned to owner after downgrade: %.09f GLF\n", poolsutil.ToFIL(downgradeAmount))
		}

		agentAddr, auth, _, requesterKey, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().SPPlusDowngrade(ctx, auth, big.NewInt(tokenID), tier, agentAddr, requesterKey)
		if err != nil {
			logFatalf("Failed to downgrade tier %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to downgrade tier %s", err)
		}

		s.Stop()

		fmt.Println("Tier successfully downgraded.")
		err = printGlfOwnerBalance("GLF balance of owner after downgrade")
		if err != nil {
			logFatal(err)
		}
	},
}

func init() {
	plusTiersCmd.AddCommand(plusDowngradeCmd)
	plusDowngradeCmd.Flags().BoolVar(&acceptPenalty, "accept-penalty", false, "Pay penalty for early downgrade")
}
