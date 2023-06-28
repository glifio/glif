package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var inpoolUtilizationRateCmd = &cobra.Command{
	Use:   "utilization-rate",
	Short: "Returns the percentage of FIL currently deployed from the Infinity Pool to Agents",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Querying the amount of FIL currently outstanding from the Infinity Pool")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		borrowed, err := PoolsSDK.Query().InfPoolTotalBorrowed(cmd.Context())
		if err != nil {
			logFatalf("Failed to get total borrowed %s", err)
		}
		assets, err := PoolsSDK.Query().InfPoolTotalAssets(cmd.Context())
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		q := big.NewFloat(0).Quo(borrowed, assets)
		per := big.NewFloat(0).Mul(q, big.NewFloat(100))

		fmt.Printf("Infinity Pool deployed: %.04f%%\n", per)
	},
}

func init() {
	infinitypoolCmd.AddCommand(inpoolUtilizationRateCmd)
}
