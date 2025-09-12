package cmd

import (
	"fmt"
	"math/big"

	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/util"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func percBigInt(num *big.Int) *big.Float {
	return new(big.Float).Mul(big.NewFloat(100), util.ToFIL(num))
}

func dtlToLeverage(dtl *big.Int) *big.Int {
	one := big.NewInt(1e18)
	denom := new(big.Int).Sub(one, dtl)
	leverage := new(big.Int).Quo(one, denom)
	return leverage
}

func premium(cashbackPremium *big.Int) *big.Float {
	return new(big.Float).Mul(big.NewFloat(100), util.ToFIL(cashbackPremium))
}

var tiersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all GLIF+ card tiers and their benefits",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// Get tier info from on-chain
		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		// Create and populate the table
		tbl := table.New("Tier", "Activation Amount ($GLF)", "Max DTL/Leverage", "Cashback Exchange Premium")

		// Add default (no card) row first
		defaultDTL := constants.MAX_BORROW_DTL
		percDefaultDTL := percBigInt(defaultDTL)
		defaultLeverage := dtlToLeverage(defaultDTL)
		tbl.AddRow("None (Default)", "0 GLF", fmt.Sprintf("%.01f%%/%vx", percDefaultDTL, defaultLeverage), "-")

		// Add each tier (skip index 0 which is "Inactive")
		for tier := uint8(1); tier < uint8(len(tierInfos)); tier++ {
			tierInfo := tierInfos[tier]

			// Get tier name
			name := tierName(tier)

			// Get GLF requirement
			glfRequirement := util.ToFIL(tierInfo.TokenLockAmount)
			glfReqStr := fmt.Sprintf("%.0f GLF", glfRequirement)

			dtl := percBigInt(tierInfo.DebtToLiquidationValue)
			leverage := dtlToLeverage(tierInfo.DebtToLiquidationValue)
			leverageStr := fmt.Sprintf("%.01f%%/%vx", dtl, leverage)

			premium := new(big.Float).Mul(
				new(big.Float).Sub(util.ToFIL(tierInfo.CashBackPremium), big.NewFloat(1)),
				big.NewFloat(100),
			)

			cashbackStr := fmt.Sprintf("+%.02f%% premium", premium)

			tbl.AddRow(name, glfReqStr, leverageStr, cashbackStr)
		}

		fmt.Println()
		tbl.AddRow("", "", "", "")
		tbl.Print()
	},
}

func init() {
	plusTiersCmd.AddCommand(tiersListCmd)
}
