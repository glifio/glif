/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var wFILCmd = &cobra.Command{
	Use:        "wfil",
	Short:      "Commands for interacting with Wrapped Filecoin tokens",
	Deprecated: "wFIL command palette has been moved under the `tokens` commands. These wFIL commands will be moved in the next major release.",
}

func init() {
	rootCmd.AddCommand(wFILCmd)
}
