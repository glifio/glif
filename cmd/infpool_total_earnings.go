/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/constants"
	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var infpoolTotalEarnings = &cobra.Command{
	Use:   "total-earnings",
	Short: "Returns the amount of FIL earned by the pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Querying the amount of FIL earned by the Infinity Pool")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		iFILPrice, err := PoolsSDK.Query().IFILPrice(cmd.Context())
		if err != nil {
			logFatalf("Failed to get iFIL price %s", err)
		}

		iFILSupply, err := PoolsSDK.Query().IFILSupply(cmd.Context())
		if err != nil {
			logFatalf("Failed to get iFIL supply %s", err)
		}

		totalAssets := new(big.Int).Mul(iFILPrice, iFILSupply)
		totalAssets = totalAssets.Div(totalAssets, constants.WAD)

		totalEarnings := new(big.Int).Sub(totalAssets, iFILSupply)

		s.Stop()

		fmt.Printf("Infinity Pool earnings: %.04f FIL\n", denoms.ToFIL(totalEarnings))
	},
}

func init() {
	infinitypoolCmd.AddCommand(infpoolTotalEarnings)
}
