/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payToCurrentCmd = &cobra.Command{
	Use:   "to-current [flags]",
	Short: "Make your account current",
	Long:  "Pays off all fees owed",
	Run: func(cmd *cobra.Command, args []string) {
		payAmt, err := pay(cmd, args, "to-current", false)
		if err != nil {
			logFatal(err)
		}
		fmt.Printf("Successfully paid %s FIL", util.ToFIL(payAmt).String())
	},
}

func init() {
	payCmd.AddCommand(payToCurrentCmd)
	payToCurrentCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payToCurrentCmd.Flags().String("from", "", "address to send the transaction from")
}
