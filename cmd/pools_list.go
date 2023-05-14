/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var poolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of active Pools",
	Run: func(cmd *cobra.Command, args []string) {
		poolsList, err := PoolsSDK.Query().ListPools(cmd.Context())
		if err != nil {
			log.Fatalf("Failed to get list of active pools: %s", err)
		}

		poolsStr := util.StringifyArg(poolsList)

		fmt.Printf("Pools: %s", poolsStr)
	},
}

func init() {
	poolsCmd.AddCommand(poolsListCmd)
}
