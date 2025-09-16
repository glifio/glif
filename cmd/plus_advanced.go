/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var plusAdvancedCmd = &cobra.Command{
	Use:   "advanced",
	Short: "Manage advanced Card settings",
}

func init() {
	plusCmd.AddCommand(plusAdvancedCmd)
}
