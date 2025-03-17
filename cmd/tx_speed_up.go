/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var txSpeedUpCmd = &cobra.Command{
	Use:   "speed-up <tx hash>",
	Short: "Replaces a Eth transaction in the mempool with a higher premium",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		replaceTx(cmd, args, false)
	},
}

func init() {
	txCmd.AddCommand(txSpeedUpCmd)
}
