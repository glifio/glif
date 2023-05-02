/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: handle deployer key

		// 1. Read in the owner and operator addresses
		ownerAddr, _, err := deriveAddrFromPkString(viper.GetString("keys.owner"))
		if err != nil {
			log.Fatal(err)
		}

		operatorAddr, _, err := deriveAddrFromPkString(viper.GetString("keys.operator"))
		if err != nil {
			log.Fatal(err)
		}

		requestAddr, _, err := deriveAddrFromPkString(viper.GetString("keys.request"))
		if err != nil {
			log.Fatal(err)
		}

		// 2. Call AgentCreate, which gives you an address, agent ID, and a transaction hash
		id, addr, tx, err := fevm.Connection().AgentCreate(cmd.Context(), nil, ownerAddr, operatorAddr, requestAddr)
		if err != nil {
			log.Fatal(err)
		}

		// 3. Given the tx hash, WaitForReceipt(tx.Hash())
		fevm.WaitForReceipt(tx.Hash())

		// 4. Print the address, agent ID, and tx hash
		fmt.Printf("Agent address: %s\n", addr)
		fmt.Printf("Agent ID: %s\n", id)
		fmt.Printf("Tx hash: %s\n", tx.Hash())

		// 5. Write the address, agent ID, and tx hash to the config
		viper.Set("agent.address", addr)
		viper.Set("agent.id", id)
		viper.Set("agent.tx", tx.Hash())
	},
}

func init() {
	agentCmd.AddCommand(createCmd)

	createCmd.Flags().String("ownerfile", "", "Owner eth address")
	createCmd.Flags().String("operatorfile", "", "Repayment eth address")
	createCmd.Flags().String("deployerfile", "", "Deployer eth address")
}
