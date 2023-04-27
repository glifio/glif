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

// agentCmd represents the agent command
var poolsCmd = &cobra.Command{
	Use:   "pools",
	Short: "Commands for interacting with the GLIF Pools Protocol",
}

var poolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of active Pools",
	Run: func(cmd *cobra.Command, args []string) {
		poolsList, err := fevm.Connection().PoolsList()
		if err != nil {
			log.Fatalf("Failed to get list of active pools: %s", err)
		}

		poolsStr := fevm.StringifyArg(poolsList)

		fmt.Printf("Pools: %s", poolsStr)
	},
}

func init() {
	rootCmd.AddCommand(poolsCmd)
	poolsCmd.AddCommand(poolsListCmd)
}
