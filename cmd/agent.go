/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Commands for interacting with the Glif Agent",
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
