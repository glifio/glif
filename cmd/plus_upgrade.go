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

var plusUpgradeCmd = &cobra.Command{
	Use:   "upgrade <new tier: bronze, silver or gold>",
	Short: "Upgrade to a higher tier",
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

		if tier <= info.Tier {
			err := fmt.Errorf("new tier must be higher than current tier: %s", tierName(info.Tier))
			logFatal(err)
		}

		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().PlusUpgrade(ctx, auth, big.NewInt(tokenID), tier)
		if err != nil {
			logFatalf("Failed to upgrade tier %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to upgrade tier %s", err)
		}

		s.Stop()

		fmt.Println("Tier successfully upgraded.")
	},
}

func init() {
	plusCmd.AddCommand(plusUpgradeCmd)
	plusUpgradeCmd.Flags().String("from", "owner", "account to use")
}
