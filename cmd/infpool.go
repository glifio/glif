/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var infinitypoolCmd = &cobra.Command{
	Use:   "infinity-pool",
	Short: "Commands for interacting with the Infinity Pool",
}

func init() {
	rootCmd.AddCommand(infinitypoolCmd)
}
