package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var wFILBalanceOfCmd = &cobra.Command{
	Use:   "balance-of [address]",
	Short: "Get the wFIL balance of an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		fmt.Printf("Checking wFIL balance of %s...\n", strAddr)

		addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		bal, err := PoolsSDK.Query().WFILBalanceOf(cmd.Context(), addr)
		if err != nil {
			logFatalf("Failed to get wFIL balance %s", err)
		}

		balFIL, _ := bal.Float64()

		s.Stop()

		fmt.Printf("wFIL balance of %s is %.09f\n", strAddr, balFIL)
	},
}

func init() {
	wFILCmd.AddCommand(wFILBalanceOfCmd)
}
