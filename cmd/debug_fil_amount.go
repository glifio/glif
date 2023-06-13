/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"

	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var debugFilAmtCmd = &cobra.Command{
	Use:   "debug-fil-amt <amount>",
	Short: "Fetches the Agent ID (uses the address in agent.toml by default)",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		amount := args[0]
		fmt.Println("amount args[0]: ", amount)
		attofil, err := parseFILAmount(amount)
		if err != nil {
			logFatal(err)
		}
		fmt.Println("attofil: ", attofil)
		fmt.Println("fil: ", denoms.ToFIL(attofil))
		fmt.Println("Debug log for float:")
		fmt.Println("Prec: ", new(big.Float).Prec())
		fmt.Println("MinPrec: ", new(big.Float).MinPrec())
	},
}

func init() {
	rootCmd.AddCommand(debugFilAmtCmd)
}
