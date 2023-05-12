/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var poolsCmd = &cobra.Command{
	Use:   "pools",
	Short: "Commands for interacting with the GLIF Pools Protocol",
}

func init() {
	rootCmd.AddCommand(poolsCmd)
}
