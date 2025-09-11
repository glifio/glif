package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var acceptPenalty bool

var plusDowngradeCmd = &cobra.Command{
	Use:   "downgrade <new tier: bronze or silver>",
	Short: "Downgrade to a lower tier",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		agentStore := util.AgentStore()

		tokenIDStr, err := agentStore.Get("plus-token-id")
		if err != nil && err.Error() != "key not found: plus-token-id" {
			logFatal(err)
		}

		if tokenIDStr == "" {
			logFatal("GLIF Card not minted yet.")
		}

		tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
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
		fmt.Printf("GLF lock amount for %s tier: %.0f GLF\n", tierName(info.Tier), poolsutil.ToFIL(oldLockAmount))
		fmt.Printf("GLF lock amount for %s tier: %.0f GLF\n", tierName(tier), poolsutil.ToFIL(newLockAmount))
		refundGlf := new(big.Int).Sub(oldLockAmount, newLockAmount)

		penaltyWindow, penaltyFee, err := PoolsSDK.Query().SPPlusTierSwitchPenaltyInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		windowStart, windowEnd, days, hours := getTierSwitchWindow(info, penaltyWindow)

		if refundGlf.Sign() == 1 && windowEnd.After(time.Now()) {
			fmt.Printf("Last tier switch timestamp: %v\n", windowStart.UTC())
			fmt.Printf("Free downgrade after %v (%d days, %d hours)\n", windowEnd.UTC(), days, hours)
			penaltyAmount := new(big.Int).Div(
				new(big.Int).Mul(refundGlf, penaltyFee),
				big.NewInt(10000))
			fmt.Printf("Penalty fee: %.09f GLF\n", poolsutil.ToFIL(penaltyAmount))
			expectedRefund := new(big.Int).Sub(refundGlf, penaltyAmount)
			fmt.Printf("Refund with penalty: %.09f GLF\n", poolsutil.ToFIL(expectedRefund))
			if !acceptPenalty {
				logFatal("Re-run with --accept-penalty flag to pay penalty and proceed with downgrade")
			}
		} else {
			downgradeAmount := new(big.Int).Sub(oldLockAmount, newLockAmount)
			fmt.Printf("GLF returned to owner after downgrade: %.0f GLF\n", poolsutil.ToFIL(downgradeAmount))
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
	plusCmd.AddCommand(plusDowngradeCmd)
	plusDowngradeCmd.Flags().BoolVar(&acceptPenalty, "accept-penalty", false, "Pay penalty for early downgrade")
}
