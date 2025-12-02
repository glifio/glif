package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusFundGLFVaultCmd = &cobra.Command{
	Use:   "fund <amount>",
	Short: "Deposit GLF tokens to use in the Card's cash back program",
	Long:  "Deposit GLF tokens to use in the Card's cash back program. The cash back program exchanges GLF tokens for 5% of every payment in FIL at a premium to the DEX price",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		tokenID, err := getPlusTokenID()
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		cashbackPercent, err := cmd.Flags().GetString("cashback-percent")
		if err != nil {
			logFatal(err)
		}
		var cashbackPercentBigInt *big.Int = nil
		if cashbackPercent != "" {
			cashbackPercentFloat, err := strconv.ParseFloat(cashbackPercent, 64)
			if err != nil {
				logFatal(err)
			}

			fmt.Printf("Setting cash back percent: %.02f%%\n", cashbackPercentFloat)

			cashbackPercentBigInt = big.NewInt(int64(cashbackPercentFloat * 100.00))
		}

		err = checkGlfPlusBalanceAndAllowance(amount)
		if err != nil {
			logFatal(err)
		}

		_, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().SPPlusFundGLFVault(ctx, auth, big.NewInt(tokenID), amount, cashbackPercentBigInt)
		if err != nil {
			logFatalf("Failed to fund GLF vault %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to fund GLF vault %s", err)
		}

		s.Stop()

		fmt.Println("GLF tokens transferred to vault.")
	},
}

func init() {
	plusCashBackCmd.AddCommand(plusFundGLFVaultCmd)
	plusFundGLFVaultCmd.Flags().String("cashback-percent", "", "Optional cash back percent to use for the GLF/FIL price")
}
