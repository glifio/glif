/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the addresses associated with your accounts",
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()

		owner, _ := as.Get(string(util.OwnerKey))
		operator, _ := as.Get(string(util.OperatorKey))
		request, _ := as.Get(string(util.RequestKey))
		if owner != "" || operator != "" || request != "" {
			agentNames := []string{
				string(util.OwnerKey),
				string(util.OperatorKey),
				string(util.RequestKey),
			}
			fmt.Printf("Agent accounts:\n\n")
			for _, name := range agentNames {
				printAddresses(as, name)
			}
			fmt.Println()
		}

		allNames := as.AccountNames()
		names := make([]string, 0)
		for _, name := range allNames {
			if name == string(util.OwnerKey) ||
				name == string(util.OperatorKey) ||
				name == string(util.RequestKey) {
				continue
			}
			names = append(names, name)
		}

		if len(names) > 0 {
			fmt.Printf("Regular accounts:\n\n")
			for _, name := range names {
				printAddresses(as, name)
			}
			fmt.Println()
		}
	},
}

func printAddresses(as *util.AccountsStorage, name string) {
	evm, fevm, err := as.GetAddrs(name)
	if err != nil {
		if err == util.ErrKeyNotFound {
			return
		}
		logFatal(err)
	}
	fmt.Printf("%s: %s (EVM), %s (FIL)\n", name, evm, fevm)
}

func init() {
	walletCmd.AddCommand(listCmd)
}
