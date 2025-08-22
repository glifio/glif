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

var plusDowngradeCmd = &cobra.Command{
	Use:   "downgrade <new tier: bronze, silver or gold>",
	Short: "Downgrade to a lower tier",
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

		tier, err := parseTierName(args[0])
		if err != nil {
			logFatal(err)
		}

		info, err := PoolsSDK.Query().PlusInfo(ctx, big.NewInt(tokenID))
		if err != nil {
			logFatal(err)
		}

		if tier >= info.Tier {
			err := fmt.Errorf("new tier must be lower than current tier: %s", tierName(info.Tier))
			logFatal(err)
		}

		agentAddr, auth, _, requesterKey, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().PlusDowngrade(ctx, auth, big.NewInt(tokenID), tier, agentAddr, requesterKey)
		if err != nil {
			logFatalf("Failed to downgrade tier %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to downgrade tier %s", err)
		}

		s.Stop()

		fmt.Println("Tier successfully downgraded.")
	},
}

func init() {
	plusCmd.AddCommand(plusDowngradeCmd)
}
