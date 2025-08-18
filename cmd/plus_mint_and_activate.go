package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
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

		personalCashBackPercent := big.NewInt(0)

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
		beneficiary := common.Address{}
		fmt.Printf("beneficiary %v\n", beneficiary)
		tx, err := PoolsSDK.Act().PlusMintAndActivate(ctx, auth, personalCashBackPercent, beneficiary, tier)
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint GLIF Plus NFT %s", err)
		}

		s.Stop()

		fmt.Printf("GLIF Plus NFT minted!\n")
	},
}

func init() {
	plusCmd.AddCommand(plusMintAndActivateCmd)
	plusMintAndActivateCmd.Flags().String("from", "default", "account to mint GLIF Card from")
}
