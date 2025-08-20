package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var plusMintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mints a GLIF Card",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		agentStore := util.AgentStore()

		oldTokenID, err := agentStore.Get("plus-token-id")
		if err != nil && err.Error() != "key not found: plus-token-id" {
			logFatal(err)
		}

		if oldTokenID != "" {
			logFatal("GLIF Card already minted.")
		}

		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().PlusMint(ctx, auth)
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		receipt, err := PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		// grab the token ID from the receipt's logs
		tokenID, err := PoolsSDK.Query().PlusTokenIDFromRcpt(cmd.Context(), receipt)
		if err != nil {
			logFatalf("pools sdk: query: token id from receipt: %s", err)
		}

		s.Stop()

		agentStore.Set("plus-token-id", tokenID.String())

		fmt.Printf("GLIF Plus NFT minted: %s\n", tokenID.String())
	},
}

func init() {
	plusCmd.AddCommand(plusMintCmd)
	plusMintCmd.Flags().String("from", "owner", "account to mint GLIF Card from")
}
