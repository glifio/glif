/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var iFILCmd = &cobra.Command{
	Use:        "ifil",
	Short:      "Commands for interacting with the Infinity Pool Liquid Staking Token (iFIL)",
	Deprecated: "ifil command palette has been moved under the `tokens` commands. These ifil commands will be moved in the next major release.",
}

func init() {
	rootCmd.AddCommand(iFILCmd)
}
