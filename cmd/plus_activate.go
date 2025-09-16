package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	poolsutil "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusActivateCmd = &cobra.Command{
	Use:   "activate <tier: bronze, silver or gold>",
	Short: "Activates an already minted GLIF Card with an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		tokenID, err := getPlusTokenID()
		if err != nil {
			logFatal(err)
		}

		tier, err := parseTierName(args[0])
		if err != nil {
			logFatal(err)
		}

		tierInfos, err := PoolsSDK.Query().SPPlusTierInfo(ctx, nil)
		if err != nil {
			logFatal(err)
		}
		lockAmount := tierInfos[tier].TokenLockAmount

		fmt.Printf("GLF lock amount for tier: %.09f GLF\n", poolsutil.ToFIL(lockAmount))

		err = checkGlfPlusBalanceAndAllowance(lockAmount)
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

		tx, err := PoolsSDK.Act().SPPlusActivate(ctx, auth, agentAddr, big.NewInt(tokenID), tier)
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
