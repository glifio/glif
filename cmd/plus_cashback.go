/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var plusCashBackCmd = &cobra.Command{
	Use:   "cashback",
	Short: "Manage cash back operations",
}

func init() {
	plusCmd.AddCommand(plusCashBackCmd)
}
