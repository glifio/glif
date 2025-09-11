package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusMintActivateAndFundCmd = &cobra.Command{
	Use:   "mint-activate-and-fund <tier: bronze, silver or gold> <GLF vault fund amount> <personal cashback percent>",
	Short: "Mints a GLIF Card and activates it with an agent, funds GLF vault and sets personal cashback percent",
	Args:  cobra.ExactArgs(3),
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

		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}
		lockAmount := tierInfos[tier].TokenLockAmount

		fundAmount, err := parseFILAmount(args[1])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		cashBackPercentFloat, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			logFatal(err)
		}
		cashBackPercent := int64(cashBackPercentFloat * 100.00)

		combinedAmount := new(big.Int).Add(mintPrice, lockAmount)
		combinedAmount = new(big.Int).Add(combinedAmount, fundAmount)

		if dueNow {
			fmt.Printf("%0.09f\n", poolsutil.ToFIL(combinedAmount))
			return
		}

		fmt.Printf("Mint Price: %.09f GLF\n", poolsutil.ToFIL(mintPrice))
		fmt.Printf("GLF lock amount for tier: %.09f GLF\n", poolsutil.ToFIL(lockAmount))
		fmt.Printf("GLF vault fund amount: %.09f GLF\n", poolsutil.ToFIL(fundAmount))
		fmt.Printf("Mint + Lock + Fund Amount: %.09f GLF\n", poolsutil.ToFIL(combinedAmount))

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

		tx, err := PoolsSDK.Act().SPPlusMintActivateAndFund(ctx, auth, big.NewInt(cashBackPercent), agentAddr, tier, fundAmount)
		if err != nil {
			logFatalf("Failed to mint, activate and fund GLIF Plus NFT %s", err)
		}

		receipt, err := PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to mint, activate and fund GLIF Plus NFT %s", err)
		}

		// grab the token ID from the receipt's logs
		tokenID, err := PoolsSDK.Query().SPPlusTokenIDFromRcpt(cmd.Context(), receipt)
		if err != nil {
			logFatalf("pools sdk: query: token id from receipt: %s", err)
		}

		s.Stop()

		agentStore.Set("plus-token-id", tokenID.String())

		fmt.Printf("GLIF Plus NFT minted, activated and funded: %s\n", tokenID.String())
	},
}

func init() {
	plusCmd.AddCommand(plusMintActivateAndFundCmd)
	plusMintActivateAndFundCmd.Flags().BoolVar(&dueNow, "due-now", false, "Print amount of GLF tokens required to mint, activate and fund")
}
