/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var plusCmd = &cobra.Command{
	Use:   "plus",
	Short: "Manage Glif Plus",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(plusCmd)
}
