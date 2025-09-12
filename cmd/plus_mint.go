package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusMintCmd = &cobra.Command{
	Use:   "mint [tier: bronze, silver or gold] [--fund-glf-vault <amount>] [--personal-cashback-percent <percent>]",
	Short: "Mints a GLIF Card and optionally activates it with an agent",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		err := ensureNoPlusToken()
		if err != nil {
			logFatal(err)
		}

		mintPrice, err := PoolsSDK.Query().SPPlusMintPrice(ctx, nil)
		if err != nil {
			logFatal(err)
		}

		var tier uint8
		lockAmount := big.NewInt(0)
		fundAmount := big.NewInt(0)
		var cashBackPercent int64
		if len(args) == 1 {
			tier, err = parseTierName(args[0])
			if err != nil {
				logFatal(err)
			}

			tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
			if err != nil {
				logFatal(err)
			}
			lockAmount = tierInfos[tier].TokenLockAmount

			fundAmountStr, err := cmd.Flags().GetString("fund-glf-vault")
			if err != nil {
				log.Fatal(err)
			}
			fundAmount, err = parseFILAmount(fundAmountStr)
			if err != nil {
				logFatalf("Failed to parse amount %s", err)
			}

			cashBackPercentFloat, err := cmd.Flags().GetFloat64("personal-cashback-percent")
			if err != nil {
				logFatal(err)
			}
			cashBackPercent = int64(cashBackPercentFloat * 100.00)

		}

		combinedAmount := new(big.Int).Add(mintPrice, lockAmount)
		combinedAmount = new(big.Int).Add(combinedAmount, fundAmount)

		if dueNow {
			fmt.Printf("%.09f\n", poolsutil.ToFIL(combinedAmount))
			return
		}

		fmt.Printf("Mint Price: %.09f GLF\n", poolsutil.ToFIL(mintPrice))
		if len(args) == 1 {
			fmt.Printf("GLF lock amount for tier: %.09f GLF\n", poolsutil.ToFIL(lockAmount))
			fmt.Printf("GLF vault fund amount: %.09f GLF\n", poolsutil.ToFIL(fundAmount))
			fmt.Printf("Mint + Lock Amount: %.09f GLF\n", poolsutil.ToFIL(combinedAmount))
		}

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

		var tx *types.Transaction
		if len(args) == 0 { // mint
			tx, err = PoolsSDK.Act().SPPlusMint(ctx, auth)
			if err != nil {
				logFatalf("Failed to mint GLIF Plus NFT %s", err)
			}
		} else if fundAmount.Sign() == 0 {
			tx, err = PoolsSDK.Act().SPPlusMintAndActivate(ctx, auth, agentAddr, tier)
			if err != nil {
				logFatalf("Failed to mint and activate GLIF Plus NFT %s", err)
			}
		} else {
			tx, err = PoolsSDK.Act().SPPlusMintActivateAndFund(ctx, auth, big.NewInt(cashBackPercent), agentAddr, tier, fundAmount)
			if err != nil {
				logFatalf("Failed to mint and activate GLIF Plus NFT %s", err)
			}
		}

		receipt, err := PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint and/or activate GLIF Plus NFT %s", err)
		}

		// grab the token ID from the receipt's logs
		tokenID, err := PoolsSDK.Query().SPPlusTokenIDFromRcpt(cmd.Context(), receipt)
		if err != nil {
			logFatalf("pools sdk: query: token id from receipt: %s", err)
		}

		s.Stop()

		util.AgentStore().Set("plus-token-id", tokenID.String())

		if len(args) == 0 {
			fmt.Printf("GLIF Plus NFT minted: %s\n", tokenID.String())
		} else {
			fmt.Printf("GLIF Plus NFT minted and activated: %s\n", tokenID.String())
		}
	},
}

func init() {
	plusCmd.AddCommand(plusMintCmd)
	plusMintCmd.Flags().BoolVar(&dueNow, "due-now", false, "Print amount of GLF tokens required to mint and activate")
	plusMintCmd.PersistentFlags().String("fund-glf-vault", "0", "Amount of GLF to fund vault with")
	plusMintCmd.PersistentFlags().Float64("personal-cashback-percent", 5.0, "Set personal cashback percent")
}
