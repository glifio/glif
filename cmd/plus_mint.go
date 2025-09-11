package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
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

		mintPrice, err := PoolsSDK.Query().SPPlusMintPrice(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		if dueNow {
			fmt.Printf("%0.f\n", poolsutil.ToFIL(mintPrice))
			return
		}

		fmt.Printf("Mint Price: %.09f GLF\n", poolsutil.ToFIL(mintPrice))

		err = checkGlfPlusBalanceAndAllowance(mintPrice)
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

		tx, err := PoolsSDK.Act().SPPlusMint(ctx, auth)
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		receipt, err := PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		// grab the token ID from the receipt's logs
		tokenID, err := PoolsSDK.Query().SPPlusTokenIDFromRcpt(cmd.Context(), receipt)
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
	plusMintCmd.Flags().BoolVar(&dueNow, "due-now", false, "Print amount of GLF tokens required to mint")
}
