package cmd

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var plusInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints information about the GLIF Card",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		agentStore := util.AgentStore()

		tokenIDStr, err := agentStore.Get("plus-token-id")
		if err != nil && err.Error() != "key not found: plus-token-id" {
			logFatal(err)
		}

		if tokenIDStr == "" {
			fmt.Println("GLIF Card not minted yet.")
			return
		}

		tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
		if err != nil {
			logFatal(err)
		}

		info, err := PoolsSDK.Query().PlusInfo(ctx, big.NewInt(tokenID), nil)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("GLIF Card Token ID: %d\n", tokenID)
		fmt.Printf("Tier: %s\n", tierName(info.Tier))
		fmt.Printf("Info: %+v\n", info)
	},
}

func init() {
	plusCmd.AddCommand(plusInfoCmd)
}
