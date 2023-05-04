/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage Glif Agent keys",
	Long:  ``,
}

func init() {
	agentCmd.AddCommand(keysCmd)
}
