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
	Use:   "activate",
	Short: "Activates an already minted GLIF Card with an agent",
	Args:  cobra.NoArgs,
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
		tx, err := PoolsSDK.Act().PlusActivate(ctx, auth, beneficiary, big.NewInt(tokenID), tier)
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
	plusActivateCmd.Flags().String("from", "owner", "account to activate GLIF Card from")
}
