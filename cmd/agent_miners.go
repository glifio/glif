/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// minerCmd represents the miner command
var minersCmd = &cobra.Command{
	Use: "miners",
}

func init() {
	agentCmd.AddCommand(minersCmd)
}
