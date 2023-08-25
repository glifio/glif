/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/glifio/cli/util"
	walletutils "github.com/glifio/go-wallet-utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Glif agent",
	Long:  `Spins up a new Agent contract through the Agent Factory, passing the owner, operator, and requestor addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()
		agentStore := util.AgentStore()
		ks := util.KeyStore()
		backends := []accounts.Backend{}
		backends = append(backends, ks)
		manager := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, backends...)

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
		passphrase, envSet := os.LookupEnv("GLIF_OWNER_PASSPHRASE")
		if !envSet {
			prompt := &survey.Password{
				Message: "Owner key passphrase",
			}
			survey.AskOne(prompt, &passphrase)
		}
		wallet, err := manager.Find(account)
		if err != nil {
			logFatal(err)
		}

		if util.IsZeroAddress(ownerAddr) || util.IsZeroAddress(operatorAddr) || util.IsZeroAddress(requestAddr) {
			logFatal("Keys not found. Please check your `keys.toml` file")
		}

		fmt.Printf("Creating agent, owner %s, operator %s, request %s\n", ownerAddr, operatorAddr, requestAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		auth, err := walletutils.NewEthWalletTransactor(wallet, &account, passphrase, big.NewInt(chainID))
		if err != nil {
			logFatal(err)
		}

		// submit the agent create transaction
		tx, err := PoolsSDK.Act().AgentCreate(
			cmd.Context(),
			auth,
			ownerAddr,
			operatorAddr,
			requestAddr,
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

		agentStore.Set("id", id.String())
		agentStore.Set("address", addr.String())
		agentStore.Set("tx", tx.Hash().String())
	},
}

func init() {
	agentCmd.AddCommand(createCmd)
}
