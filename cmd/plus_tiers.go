/*
Copyright © 2025 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var plusTiersCmd = &cobra.Command{
	Use:   "tiers",
	Short: "Commands for Card tiers",
}

func init() {
	plusCmd.AddCommand(plusTiersCmd)
}
