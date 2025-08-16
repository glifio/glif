package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var plusFundGLFVaultCmd = &cobra.Command{
	Use:   "fund-glf-vault <amount>",
	Short: "Deposit GLF tokens into vault for cashback",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		agentStore := util.AgentStore()

		tokenIDStr, err := agentStore.Get("plus-token-id")
		if err != nil && err.Error() != "key not found: plus-token-id" {
			logFatal(err)
		}

		if tokenIDStr == "" {
			logFatal("GLIF Card not minted yet.")
		}

		tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
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

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		tx, err := PoolsSDK.Act().PlusFundGLFVault(ctx, auth, big.NewInt(tokenID), amount)
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
