/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var minersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of miners owned by this Agent",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		list, err := PoolsSDK.Query().AgentMiners(cmd.Context(), agentAddr, nil)
		if err != nil {
			logFatal(err)
		}

		if len(list) == 0 {
			fmt.Printf("Agent has no miners\n")
			return
		}

		totalBal := big.NewInt(0)

		fmt.Printf("\033[1m%s\033[0m", "Agent's miners:\n")
		for _, miner := range list {
			bal, err := lapi.WalletBalance(cmd.Context(), miner)
			if err != nil {
				logFatal(err)
			}

			totalBal = new(big.Int).Add(totalBal, bal.Int)
			fmt.Printf("Miner %s - %0.09f FIL\n", miner, util.ToFIL(bal.Int))
		}
		fmt.Printf("\nTotal balance: %0.09f\n", util.ToFIL(totalBal))
	},
}

func init() {
	minersCmd.AddCommand(minersListCmd)
	minersListCmd.Flags().String("agent-addr", "", "Agent address")
}
