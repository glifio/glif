/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payPrincipalCmd = &cobra.Command{
	Use:   "principal <amount> [flags]",
	Short: "Pay down an amount of principal (will also pay fees if any are owed)",
	Long:  "<amount> is the amount of principal to pay down, in FIL. Any fees owed will be paid off as well in order to make the principal payment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payAmt, err := pay(cmd, args, "principal", false)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Successfully paid %s FIL", util.ToFIL(payAmt).String())
	},
}

func init() {
	payCmd.AddCommand(payPrincipalCmd)
	payPrincipalCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payPrincipalCmd.Flags().String("from", "", "address to send the transaction from")
}
