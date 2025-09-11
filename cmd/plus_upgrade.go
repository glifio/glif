package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusUpgradeCmd = &cobra.Command{
	Use:   "upgrade <new tier: bronze, silver or gold>",
	Short: "Upgrade to a higher tier",
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

		if tier <= info.Tier {
			err := fmt.Errorf("new tier must be higher than current tier: %s", tierName(info.Tier))
			logFatal(err)
		}

		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}
		oldLockAmount := info.TierLockAmount
		newLockAmount := tierInfos[tier].TokenLockAmount

		upgradeAmount := new(big.Int).Sub(newLockAmount, oldLockAmount)

		if dueNow {
			fmt.Printf("%.09f\n", poolsutil.ToFIL(upgradeAmount))
			return
		}

		fmt.Printf("GLF lock amount for %s tier: %.09f GLF\n", tierName(info.Tier), poolsutil.ToFIL(oldLockAmount))
		fmt.Printf("GLF lock amount for %s tier: %.09f GLF\n", tierName(tier), poolsutil.ToFIL(newLockAmount))
		fmt.Printf("GLF required to upgrade: %.09f GLF\n", poolsutil.ToFIL(upgradeAmount))

		err = checkGlfPlusBalanceAndAllowance(upgradeAmount)
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

		tx, err := PoolsSDK.Act().SPPlusUpgrade(ctx, auth, big.NewInt(tokenID), tier)
		if err != nil {
			logFatalf("Failed to upgrade tier %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to upgrade tier %s", err)
		}

		s.Stop()

		fmt.Println("Tier successfully upgraded.")
	},
}

func init() {
	plusCmd.AddCommand(plusUpgradeCmd)
	plusUpgradeCmd.Flags().BoolVar(&dueNow, "due-now", false, "Print amount of GLF tokens required to upgrade")
}
