/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var iFILCmd = &cobra.Command{
	Use:   "ifil",
	Short: "Commands for interacting with the Infinity Pool Liquid Staking Token (iFIL)",
}

func init() {
	rootCmd.AddCommand(iFILCmd)
}
