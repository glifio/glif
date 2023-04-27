/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Read in the owner and operator addresses
		// 2. Call AgentCreate, which gives you an address, agent ID, and a transaction hash
		// 3. Given the tx hash, WaitForReceipt(tx.Hash())
		// 4. Print the address, agent ID, and tx hash
		// 5. Write the address, agent ID, and tx hash to the config

		// check flags for custom owner/repayment private key files
		// otherwise load default private keys files

		// read private key files into memory

		// agentCreatEvent := contract_utils.AgentCreate()
		// if agentCreateEvent == nil {
		//		log.Fatal("failed to create agent")
		// }

		// agentID := agentCreateEvent.Agent

		// store the agentID into a file in the config folder

	},
}

func init() {
	agentCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().String("ownerfile", "", "Owner eth address")
	createCmd.Flags().String("repaymentfile", "", "Repayment eth address")
}
