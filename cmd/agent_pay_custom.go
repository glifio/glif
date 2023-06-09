/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payCustomCmd = &cobra.Command{
	Use:   "custom <amount> [flags]",
	Short: "Pay down a custom amount of FIL",
	Args:  cobra.ExactArgs(1),
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		payAmt, err := pay(cmd, args, "custom")
		if err != nil {
			logFatal(err)
		}
		fmt.Printf("Successfully paid %s FIL", util.ToFIL(payAmt).String())
	},
}

func init() {
	payCmd.AddCommand(payCustomCmd)
	payCustomCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payCustomCmd.Flags().String("from", "", "address to send the transaction from")
}
