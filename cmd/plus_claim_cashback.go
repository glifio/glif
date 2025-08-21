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

var plusClaimCashBackCmd = &cobra.Command{
	Use:   "claim-cashback <receiver address>",
	Short: "Transfer earned FIL cashback to receiver address",
	Args:  cobra.ExactArgs(1),
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

		receiver, err := AddressOrAccountNameToEVM(cmd.Context(), args[0])
		if err != nil {
			logFatal(err)
		}

		tx, err := PoolsSDK.Act().PlusClaimCashBack(ctx, auth, big.NewInt(tokenID), receiver)
		if err != nil {
			logFatalf("Failed to claim cashback %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to claim cashback %s", err)
		}

		s.Stop()

		fmt.Println("FIL cashback successfully claimed.")
	},
}

func init() {
	plusCmd.AddCommand(plusClaimCashBackCmd)
	plusClaimCashBackCmd.Flags().String("from", "owner", "account to use")
}
