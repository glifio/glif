/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var wFILCmd = &cobra.Command{
	Use:   "wfil",
	Short: "Commands for interacting with Wrapped Filecoin tokens",
}

func init() {
	rootCmd.AddCommand(wFILCmd)
}
