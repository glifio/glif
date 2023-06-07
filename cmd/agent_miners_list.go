/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var minersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of miners owned by this Agent",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			logFatal(err)
		}

		list, err := PoolsSDK.Query().AgentMiners(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Agent's miners: %s", util.StringifyArg(list))
	},
}

func init() {
	minersCmd.AddCommand(minersListCmd)
	minersListCmd.Flags().String("agent-addr", "", "Agent address")
}
