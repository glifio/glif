/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var txCancelCmd = &cobra.Command{
	Use:   "cancel <tx hash or cid>",
	Short: "Replaces a transaction in the mempool with a dummy (zero value transfer to self)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		replaceTx(cmd, args, true)
	},
}

func init() {
	txCmd.AddCommand(txCancelCmd)
}
