/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
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

		addr, err := ParseAddress(cmd.Context(), strAddr)
		if err != nil {
			log.Fatalf("Failed to parse address %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		bal, err := fevm.Connection().IFILBalanceOf(addr)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		fmt.Printf("iFIL balance of %s is %s", strAddr, bal.String())
	},
}

var iFILTransferCmd = &cobra.Command{
	Use:   "transfer [to] [amount]",
	Short: "Transfer iFIL to another address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		strAmt := args[1]
		fmt.Printf("Transferring %s iFIL balance to %s...", strAmt, strAddr)

		addr, err := ParseAddress(cmd.Context(), strAddr)
		if err != nil {
			log.Fatalf("Failed to parse address %s", err)
		}

		amt := big.NewInt(0)
		amt, ok := amt.SetString(strAmt, 10)
		if !ok {
			log.Fatalf("Failed to parse amount %s", err)
		}

		tx, err := fevm.Connection().IFILTransfer(cmd.Context(), addr, amt)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		if tx == nil {
			log.Fatal("Failed to transfer iFIL")
		}

		eapi, err := fevm.Connection().ConnectEthClient()
		if err != nil {
			log.Fatalf("Failed to instantiate eth client %s", err)
		}
		defer eapi.Close()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		fevm.WaitForReceipt(tx.Hash())

		s.Stop()

		fmt.Printf("iFIL sent!")
	},
}

var iFILApproveCmd = &cobra.Command{
	Use:   "approve <spender> <allowance>",
	Short: "Approve another address to spend your iFIL",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		strAmt := args[1]
		fmt.Printf("Approving %s to spend %s of your iFIL balance...", strAddr, strAmt)

		addr, err := ParseAddress(cmd.Context(), strAddr)
		if err != nil {
			log.Fatalf("Failed to parse address %s", err)
		}

		amt := big.NewInt(0)
		amt, ok := amt.SetString(strAmt, 10)
		if !ok {
			log.Fatalf("Failed to parse amount %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := fevm.Connection().IFILApprove(cmd.Context(), addr, amt)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		if tx == nil {
			log.Fatal("Failed to transfer iFIL")
		}

		fevm.WaitForReceipt(tx.Hash())

		s.Stop()

		fmt.Printf("iFIL sent!")
	},
}

var iFILPriceCmd = &cobra.Command{
	Use:   "price",
	Short: "Get the iFIL price, denominated in FIL",
	Long:  "Get the iFIL price, denominated in FIL. The number returned is the amount of FIL that 1 iFIL is worth.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Checking iFIL prices...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		price, err := fevm.Connection().IFILPrice()
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		fmt.Printf("1 iFIL is worth %s FIL", price.String())
	},
}

func init() {
	rootCmd.AddCommand(iFILCmd)
	iFILCmd.AddCommand(iFILBalanceOfCmd)
	iFILCmd.AddCommand(iFILTransferCmd)
	iFILCmd.AddCommand(iFILApproveCmd)
}
