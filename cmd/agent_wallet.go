/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage Glif wallets",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(walletCmd)
}
