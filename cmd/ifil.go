/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var iFILCmd = &cobra.Command{
	Use:   "ifil",
	Short: "Commands for interacting with the Infinity Pool Liquid Staking Token (iFIL)",
}

var iFILBalanceOfCmd = &cobra.Command{
  Use:   "balance-of [address]",
  Short: "Get the iFIL balance of an address",
  Args:  cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    strAddr := args[0]
    fmt.Printf("Checking iFIL balance of %s...", strAddr)

    lapi, closer, err := fevm.Connection().ConnectLotusClient()
    if err != nil {
      log.Fatalf("Failed to instantiate lotus client %s", err)
    }
    defer closer()

    addr, err := ParseAddress(cmd.Context(), strAddr, lapi)
    if err != nil {
      log.Fatalf("Failed to parse address %s", err)
    }

    bal, err := fevm.Connection().IFILBalanceOf(addr)
    if err != nil {
      log.Fatalf("Failed to get iFIL balance %s", err)
    }

    fmt.Printf("iFIL balance of %s is %s", strAddr, bal.String())
  },
}

func init() {
	rootCmd.AddCommand(iFILCmd)
  iFILCmd.AddCommand(iFILBalanceOfCmd)
}
