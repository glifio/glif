/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var minersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of miners owned by this Agent",
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AgentStore()

		var agentIDStr string
		if cmd.Flag("agent-id") != nil && cmd.Flag("agent-id").Changed {
			agentIDStr = cmd.Flag("agent-id").Value.String()
		} else {
			// Check if an agent already exists
			cachedID, err := as.Get("id")
			if err != nil {
				log.Fatal(err)
			}

			agentIDStr = cachedID
		}

		agentID, ok := new(big.Int).SetString(agentIDStr, 10)
		if !ok {
			log.Fatal("could not convert agent id %s to big.Int", agentIDStr)
		}

		list, err := fevm.Connection().MinersList(cmd.Context(), agentID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Agent's miners: %s", fevm.StringifyArg(list))
	},
}

func init() {
	minersCmd.AddCommand(minersListCmd)
	minersListCmd.Flags().String("agent-id", "", "ID of the Agent")
}
