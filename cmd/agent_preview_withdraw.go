package cmd

import (
	"github.com/glifio/go-pools/constants"
	"github.com/spf13/cobra"
)

var previewWithdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Preview the financial position of performing a withdraw action",
	Run: func(cmd *cobra.Command, args []string) {
		previewAction(cmd, args, constants.MethodWithdraw)
	},
}

func init() {
	previewCmd.AddCommand(previewWithdrawCmd)
}
