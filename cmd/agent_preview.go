/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use: "preview",
}

func init() {
	agentCmd.AddCommand(previewCmd)
}
