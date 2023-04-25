/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// minerCmd represents the miner command
var minerCmd = &cobra.Command{
	Use: "miner",
}

func init() {
	agentCmd.AddCommand(minerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// minerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// minerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
