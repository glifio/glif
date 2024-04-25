package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var wFILAllowanceCmd = &cobra.Command{
	Use:   "allowance [holder] [spender]",
	Short: "Get the wFIL balance of an address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		holderStr := args[0]
		spenderStr := args[1]
		fmt.Printf("Checking wFIL allowance of spender %s on holder %s...\n", spenderStr, holderStr)

		holder, err := AddressOrAccountNameToEVM(cmd.Context(), holderStr)
		if err != nil {
			logFatalf("Failed to parse address %s\n", err)
		}

		spender, err := AddressOrAccountNameToEVM(cmd.Context(), spenderStr)
		if err != nil {
			logFatalf("Failed to parse address %s\n", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		allowance, err := PoolsSDK.Query().WFILAllowance(cmd.Context(), holder, spender)
		if err != nil {
			logFatalf("Failed to get wFIL allowance %s\n", err)
		}

		s.Stop()

		fmt.Printf("wFIL allowance of spender %s on holder %s is: %.09f FIL\n", util.TruncateAddr(spenderStr), util.TruncateAddr(holderStr), allowance)
	},
}

func init() {
	wFILCmd.AddCommand(wFILAllowanceCmd)
}
