//go:build advanced
// +build advanced

/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Commands for controlling the Agent's ownership and health status",
}

func init() {
	agentCmd.AddCommand(adminCmd)
}
