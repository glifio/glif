/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var payCmd = &cobra.Command{
	Use: "pay",
}

func init() {
	agentCmd.AddCommand(payCmd)
}
