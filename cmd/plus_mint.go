package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusMintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mints a GLIF Card",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
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

		fmt.Printf("GLIF Plus NFT minted: %s\n", tokenID.String())
	},
}

func init() {
	plusCmd.AddCommand(plusMintCmd)
	plusMintCmd.Flags().String("from", "owner", "account to mint GLIF Card from")
}
