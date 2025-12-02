package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusFILVaultBalanceCmd = &cobra.Command{
	Use:   "fil-vault-balance",
	Short: "Get the FIL vault balance",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		balance, err := PoolsSDK.Query().SPPlusFILVaultBalance(ctx, nil)
		if err != nil {
			logFatalf("Failed to get FIL vault balance %s", err)
		}

		s.Stop()

		fmt.Printf(" %.09f FIL\n", poolsutil.ToFIL(balance))
	},
}

func init() {
	plusCashBackCmd.AddCommand(plusFILVaultBalanceCmd)
}
