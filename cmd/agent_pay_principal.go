/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"

	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payPrincipalPreview bool

var payPrincipalCmd = &cobra.Command{
	Use:   "principal <amount> [flags]",
	Short: "Pay down an amount of principal (will also pay fees if any are owed)",
	Long:  "<amount> is the amount of principal to pay down, in FIL. Any fees owed will be paid off as well in order to make the principal payment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if payPrincipalPreview {
			agentAddr, err := getAgentAddressWithFlags(cmd)
			if err != nil {
				logFatal(err)
			}
			amount, err := parseFILAmount(args[0])
			if err != nil {
				logFatal(err)
			}

			amountOwed, err := PoolsSDK.Query().AgentInterestOwed(cmd.Context(), agentAddr, nil)
			if err != nil {
				logFatal(err)
			}

			payAmt := new(big.Int).Add(amount, amountOwed)
			args = append(args, util.ToFIL(payAmt).String())
			previewAction(cmd, args, constants.MethodPay)
			return
		}
		payAmt, err := pay(cmd, args, Principal)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Successfully paid %s FIL\n", util.ToFIL(payAmt).String())
	},
}

func init() {
	payCmd.AddCommand(payPrincipalCmd)
	payPrincipalCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payPrincipalCmd.Flags().String("from", "", "address to send the transaction from")
	payPrincipalCmd.Flags().BoolVar(&payPrincipalPreview, "preview", false, "preview financial outcome of pay principal action")
}
