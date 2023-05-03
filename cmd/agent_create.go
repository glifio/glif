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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: handle deployer key

		// 1. Read in the owner and operator addresses
		ownerKey, err := KeyStorage.Get("owner")
		if err != nil {
			log.Fatal(err)
		}
		ownerAddr, _, err := deriveAddrFromPkString(ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		operatorKey, err := KeyStorage.Get("operator")
		if err != nil {
			log.Fatal(err)
		}
		operatorAddr, _, err := deriveAddrFromPkString(operatorKey)
		if err != nil {
			log.Fatal(err)
		}

		requestKey, err := KeyStorage.Get("request")
		if err != nil {
			log.Fatal(err)
		}
		requestAddr, _, err := deriveAddrFromPkString(requestKey)
		if err != nil {
			log.Fatal(err)
		}

		bh, err := fevm.Connection().BlockNumber()
		if err != nil {
			log.Fatal(err)
		}

		// 2. Call AgentCreate, which gives you an address, agent ID, and a transaction hash
		tx, err := fevm.Connection().AgentCreate(cmd.Context(), nil, ownerAddr, operatorAddr, requestAddr)
		if err != nil {
			log.Fatal(err)
		}

		// 3. Given the tx hash, WaitForReceipt(tx.Hash())
		receipt := fevm.WaitReturnReceipt(tx.Hash())
		if receipt == nil {
			log.Fatal("Failed to get receipt")
		}

		// 4. Call AgentFilter, which gives you the agent ID
		id, err := fevm.Connection().AgentFilter(cmd.Context(), receipt, bh)
		if err != nil {
			log.Fatal(err)
		}

		// 4. Print the address, agent ID, and tx hash
		// fmt.Printf("Agent address: %s\n", addr)
		fmt.Printf("Agent ID: %s\n", id)
		fmt.Printf("Tx hash: %s\n", tx.Hash())

		// 5. Write the address, agent ID, and tx hash to the config
		// AgentStorage.Set("agent.address", addr.String())
		AgentStorage.Set("agent.id", id.String())
		AgentStorage.Set("agent.tx", tx.Hash().String())
	},
}

func init() {
	agentCmd.AddCommand(createCmd)

	createCmd.Flags().String("ownerfile", "", "Owner eth address")
	createCmd.Flags().String("operatorfile", "", "Repayment eth address")
	createCmd.Flags().String("deployerfile", "", "Deployer eth address")
}
