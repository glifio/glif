/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payToCurrentPreview bool

var payToCurrentCmd = &cobra.Command{
	Use:   "to-current [flags]",
	Short: "Make your account current",
	Long:  "Pays off all fees owed",
	Run: func(cmd *cobra.Command, args []string) {
		if payToCurrentPreview {
			previewAction(cmd, args, constants.MethodPay)
			return
		}

		payAmt, err := pay(cmd, args, ToCurrent)
		if err != nil {
			logFatal(err)
		}
		fmt.Printf("Successfully paid %s FIL\n", util.ToFIL(payAmt).String())
	},
}

func init() {
	payCmd.AddCommand(payToCurrentCmd)
	payToCurrentCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payToCurrentCmd.Flags().String("from", "", "address to send the transaction from")
	payToCurrentCmd.Flags().BoolVar(&payToCurrentPreview, "preview", false, "DEPRECATED: preview financial outcome of pay to-current action")
}
