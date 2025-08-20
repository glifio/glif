package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusMintAndActivateCmd = &cobra.Command{
	Use:   "mint-and-activate",
	Short: "Mints a GLIF Card and activates it with an agent",
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

		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		/*
			enum Tier {
					Inactive,
					Bronze,
					Silver,
					Gold
			}
		*/

		// var tier uint8 = 0 // Inactive
		var tier uint8 = 1 // Bronze

		fmt.Printf("auth.From %v\n", auth.From)
		fmt.Printf("agentAddr %v\n", agentAddr)
		fmt.Printf("tier %v\n", tier)
		// beneficiary := common.Address{}
		beneficiary := agentAddr
		fmt.Printf("beneficiary %v\n", beneficiary)
		tx, err := PoolsSDK.Act().PlusMintAndActivate(ctx, auth, beneficiary, tier)
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
	plusCmd.AddCommand(plusMintAndActivateCmd)
	plusMintAndActivateCmd.Flags().String("from", "owner", "account to mint GLIF Card from")
}
