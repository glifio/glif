package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusFundGLFVaultCmd = &cobra.Command{
	Use:   "fund-glf-vault <amount>",
	Short: "Deposit GLF tokens into vault for cashback",
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

		tx, err := PoolsSDK.Act().SPPlusFundGLFVault(ctx, auth, big.NewInt(tokenID), amount)
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
	plusCmd.AddCommand(plusFundGLFVaultCmd)
}
