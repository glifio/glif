/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// borrowCmd represents the borrow command
var borrowCmd = &cobra.Command{
	Use:   "borrow",
	Short: "Borrow FIL from the Glif Infinity Pool",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("borrow called")
	},
}

func init() {
	agentCmd.AddCommand(borrowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// borrowCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// borrowCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
