/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

var minersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of miners owned by this Agent",
	Run: func(cmd *cobra.Command, args []string) {
		agentID, err := getAgentID(cmd)
		if err != nil {
			log.Fatal(err)
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
