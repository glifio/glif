/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  `Spins up a new Agent contract through the Agent Factory, passing the owner, operator, and requestor addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AgentStore()
		ks := util.KeyStore()
		wallet := ks.Wallets()[0]

		// Check if an agent already exists
		addressStr, err := as.Get("address")
		if err != nil && err != util.ErrKeyNotFound {
			logFatal(err)
		}
		if addressStr != "" {
			logFatalf("Agent already exists: %s", addressStr)
		}

		ownerAddr, _, err := as.GetAddrs(util.OwnerKey)
		if err != nil {
			logFatal(err)
		}

		operatorAddr, _, err := as.GetAddrs(util.OperatorKey)
		if err != nil {
			logFatal(err)
		}

		requestAddr, _, err := as.GetAddrs(util.RequestKey)
		if err != nil {
			logFatal(err)
		}

		account := accounts.Account{Address: ownerAddr}
		passphrase := ""

		if util.IsZeroAddress(ownerAddr) || util.IsZeroAddress(operatorAddr) || util.IsZeroAddress(requestAddr) {
			logFatal("Keys not found. Please check your `keys.toml` file")
		}

		fmt.Printf("Creating agent, owner %s, operator %s, request %s", ownerAddr, operatorAddr, requestAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		// submit the agent create transaction
		tx, err := PoolsSDK.Act().AgentCreate(
			cmd.Context(),
			ownerAddr,
			operatorAddr,
			requestAddr,
			wallet,
			account,
			passphrase,
		)
		if err != nil {
			logFatalf("pools sdk: agent create: %s", err)
		}

		s.Stop()

		fmt.Printf("Agent create transaction submitted: %s\n", tx.Hash())
		fmt.Println("Waiting for confirmation...")

		s.Start()
		// transaction landed on chain or errored
		receipt, err := PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatalf("pools sdk: query: state wait receipt: %s", err)
		}

		// grab the ID and the address of the agent from the receipt's logs
		addr, id, err := PoolsSDK.Query().AgentAddrIDFromRcpt(cmd.Context(), receipt)
		if err != nil {
			logFatalf("pools sdk: query: agent addr id from receipt: %s", err)
		}

		s.Stop()

		fmt.Printf("Agent created: %s\n", addr.String())
		fmt.Printf("Agent ID: %s\n", id.String())

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
