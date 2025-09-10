package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusMintAndActivateCmd = &cobra.Command{
	Use:   "mint-and-activate <tier: bronze, silver or gold>",
	Short: "Mints a GLIF Card and activates it with an agent",
	Args:  cobra.ExactArgs(1),
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

		tier, err := parseTierName(args[0])
		if err != nil {
			logFatal(err)
		}

		mintPrice, err := PoolsSDK.Query().SPPlusMintPrice(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Mint Price: %.0f GLF\n", poolsutil.ToFIL(mintPrice))

		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}
		lockAmount := tierInfos[tier].TokenLockAmount

		fmt.Printf("GLF lock amount for tier: %.0f GLF\n", poolsutil.ToFIL(lockAmount))

		combinedAmount := new(big.Int).Add(mintPrice, lockAmount)
		fmt.Printf("Mint + Lock Amount: %.0f GLF\n", poolsutil.ToFIL(combinedAmount))

		err = checkGlfPlusBalanceAndAllowance(combinedAmount)
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

		tx, err := PoolsSDK.Act().SPPlusMintAndActivate(ctx, auth, agentAddr, tier)
		if err != nil {
			logFatalf("Failed to mint and activate GLIF Plus NFT %s", err)
		}

		receipt, err := PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint and activate GLIF Plus NFT %s", err)
		}

		// grab the token ID from the receipt's logs
		tokenID, err := PoolsSDK.Query().SPPlusTokenIDFromRcpt(cmd.Context(), receipt)
		if err != nil {
			logFatalf("pools sdk: query: token id from receipt: %s", err)
		}

		s.Stop()

		agentStore.Set("plus-token-id", tokenID.String())

		fmt.Printf("GLIF Plus NFT minted and activated: %s\n", tokenID.String())
	},
}

func init() {
	plusCmd.AddCommand(plusMintAndActivateCmd)
}
