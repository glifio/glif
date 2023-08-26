/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/filecoin-project/lotus/api"
	"github.com/glifio/cli/util"
	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

func printBalance(ctx context.Context, lapi *api.FullNodeStruct, as *util.AccountsStorage, name string) {
	_, addr, err := as.GetAddrs(name)
	if err != nil {
		fmt.Printf("%s balance: Error %v\n", name, err)
		return
	}

	bal, err := lapi.WalletBalance(ctx, addr)
	if err != nil {
		fmt.Printf("%s balance: Error %v\n", name, err)
		return
	}
	balance := denoms.ToFIL(bal.Int)
	bf64, _ := balance.Float64()
	fmt.Printf("%s balance: %.02f FIL\n", name, bf64)
}

// newCmd represents the new command
var balCmd = &cobra.Command{
	Use:   "balance",
	Short: "Gets the balances associated with your accounts",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		as := util.AccountsStore()

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatalf("Failed to instantiate eth client %s", err)
		}
		defer closer()

		owner, _ := as.Get(string(util.OwnerKey))
		operator, _ := as.Get(string(util.OperatorKey))
		if owner != "" || operator != "" {
			agentNames := []string{
				string(util.OwnerKey),
				string(util.OperatorKey),
			}
			fmt.Printf("Agent accounts:\n\n")
			for _, name := range agentNames {
				printBalance(ctx, lapi, as, name)
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
				printBalance(ctx, lapi, as, name)
			}
			fmt.Println()
		}
	},
}

func init() {
	walletCmd.AddCommand(balCmd)
}
