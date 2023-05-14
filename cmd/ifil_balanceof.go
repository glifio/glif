package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var iFILBalanceOfCmd = &cobra.Command{
	Use:   "balance-of [address]",
	Short: "Get the iFIL balance of an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		fmt.Printf("Checking iFIL balance of %s...", strAddr)

		addr, err := ParseAddress(cmd.Context(), strAddr)
		if err != nil {
			log.Fatalf("Failed to parse address %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		bal, err := PoolsSDK.Query().IFILBalanceOf(cmd.Context(), addr)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		balFIL, _ := bal.Float64()

		s.Stop()

		fmt.Printf("iFIL balance of %s is %.09f", strAddr, balFIL)
	},
}

func init() {
	iFILCmd.AddCommand(iFILBalanceOfCmd)
}
