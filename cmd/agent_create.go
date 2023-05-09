/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  `Spins up a new Agent contract through the Agent Factory, passing the owner, operator, and requestor addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		ks := util.KeyStore()
		as := util.AgentStore()

		// Check if an agent already exists
		addressStr, err := as.Get("address")
		if err != nil && err != util.KeyNotFoundErr {
			log.Fatal(err)
		}
		if addressStr != "" {
			log.Fatalf("Agent already exists: %s", addressStr)
		}

		ownerAddr, _, err := ks.GetAddrs(util.OwnerKey)
		if err != nil {
			log.Fatal(err)
		}

		operatorAddr, _, err := ks.GetAddrs(util.OperatorKey)
		if err != nil {
			log.Fatal(err)
		}

		requestAddr, _, err := ks.GetAddrs(util.RequestKey)
		if err != nil {
			log.Fatal(err)
		}

		if util.IsZeroAddress(ownerAddr) || util.IsZeroAddress(operatorAddr) || util.IsZeroAddress(requestAddr) {
			log.Fatal("Keys not found. Please check your `keys.toml` file")
		}

		pk, err := ks.GetPrivate(util.OwnerKey)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Creating agent, owner %s, operator %s, request %s", ownerAddr, operatorAddr, requestAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		// submit the agent create transaction
		tx, err := fevm.Connection().AgentCreate(cmd.Context(), pk, ownerAddr, operatorAddr, requestAddr)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Printf("Agent create transaction submitted: %s\n", tx.Hash())
		fmt.Println("Waiting for confirmation...")

		s.Start()
		// transaction landed on chain or errored
		receipt, err := fevm.WaitReturnReceipt(tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		if receipt == nil {
			log.Fatal("Failed to get receipt")
		}

		// grab the ID and the address of the agent from the receipt's logs
		id, addr, err := fevm.Connection().AgentAddrID(cmd.Context(), receipt)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		as.Set("id", id.String())
		as.Set("address", addr.String())
		as.Set("tx", tx.Hash().String())
	},
}

func init() {
	agentCmd.AddCommand(createCmd)

	createCmd.Flags().String("ownerfile", "", "Owner eth address")
	createCmd.Flags().String("operatorfile", "", "Repayment eth address")
	createCmd.Flags().String("deployerfile", "", "Deployer eth address")
}
