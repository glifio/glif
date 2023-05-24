/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var payPrincipalCmd = &cobra.Command{
	Use:   "principal <amount> [flags]",
	Short: "Pay down an amount of principal (will also pay fees if any are owed)",
	Long:  "<amount> is the amount of principal to pay down, in FIL. Any fees owed will be paid off as well in order to make the principal payment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, senderKey, requesterKey, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		amountOwed, _, err := PoolsSDK.Query().AgentOwes(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		payAmt := new(big.Int).Add(amount, amountOwed)

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Paying %s FIL to the %s", util.ToFIL(amountOwed).String(), poolName)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().AgentPay(cmd.Context(), agentAddr, poolID, payAmt, senderKey, requesterKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		paidAmount, err := parseFILAmount(args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully paid %s FIL", util.ToFIL(paidAmount).String())
	},
}

func init() {
	payCmd.AddCommand(payPrincipalCmd)
	payPrincipalCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	payPrincipalCmd.Flags().String("from", "", "address to send the transaction from")
}
