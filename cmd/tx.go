/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Commands to manage transactions in the mempool",
}

func init() {
	rootCmd.AddCommand(txCmd)
}
