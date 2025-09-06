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

var plusActivateCmd = &cobra.Command{
	Use:   "activate <tier: bronze, silver or gold>",
	Short: "Activates an already minted GLIF Card with an agent.",
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

		agentAddr, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().PlusActivate(ctx, auth, agentAddr, big.NewInt(tokenID), tier)
		if err != nil {
			logFatalf("Failed to activate GLIF Plus NFT %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to activate GLIF Plus NFT %s", err)
		}

		s.Stop()

		fmt.Println("GLIF Plus NFT activated.")
	},
}

func init() {
	plusCmd.AddCommand(plusActivateCmd)
}
