/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func panicIfKeyExists(key util.KeyType) {
	as := util.AccountsStore()
	_, _, err := as.GetAddrs(string(key))
	if err == nil {
		logFatal("owner account already created")
	} else {
		var e *util.ErrKeyNotFound
		if !errors.As(err, &e) {
			logFatal(err)
		}
	}
}

// createAgentAccountsCmd represents the create-agent-accounts command
var createAgentAccountsCmd = &cobra.Command{
	Use:   "create-agent-accounts",
	Short: "Create a set of accounts for an agent",
	Long: `Create 3 new accounts and store them in the keystore. The following accounts will be created:
	  "owner" - a privileged account with full admin permissions for an agent (passphrase protected)
		"operator" - a sub-account with reduced permissions to perform routine transactions (eg. payments),
		             passphrase protection is optional
		"requester" - used for requesting credentials from the "Agent Data Oracle" (no passphrase)
	`,
	Run: func(cmd *cobra.Command, args []string) {

		panicIfKeyExists(util.OwnerKey)
		panicIfKeyExists(util.OperatorKey)
		panicIfKeyExists(util.RequestKey)

		ownerPassphrase, envSet := os.LookupEnv("GLIF_OWNER_PASSPHRASE")
		if !envSet {
			prompt := &survey.Password{
				Message: "Please type a passphrase to encrypt your owner private key",
			}
			survey.AskOne(prompt, &ownerPassphrase)
			var confirmPassphrase string
			confirmPrompt := &survey.Password{
				Message: "Confirm passphrase",
			}
			survey.AskOne(confirmPrompt, &confirmPassphrase)
			if ownerPassphrase != confirmPassphrase {
				logFatal("Aborting. Passphrase confirmation did not match.")
			}
		}

		ks := util.KeyStore()

		owner, err := ks.NewAccount(ownerPassphrase)
		if err != nil {
			logFatal(err)
		}

		operatorPassphrase := os.Getenv("GLIF_OPERATOR_PASSPHRASE")
		operator, err := ks.NewAccount(operatorPassphrase)
		if err != nil {
			logFatal(err)
		}

		requester, err := ks.NewAccount("")
		if err != nil {
			logFatal(err)
		}

		as := util.AccountsStore()

		as.Set(string(util.OwnerKey), owner.Address.String())
		as.Set(string(util.OperatorKey), operator.Address.String())
		as.Set(string(util.RequestKey), requester.Address.String())

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		ownerAddr, ownerDelAddr, err := as.GetAddrs(string(util.OwnerKey))
		if err != nil {
			logFatal(err)
		}
		operatorAddr, operatorDelAddr, err := as.GetAddrs(string(util.OperatorKey))
		if err != nil {
			logFatal(err)
		}
		requestAddr, requestDelAddr, err := as.GetAddrs(string(util.RequestKey))
		if err != nil {
			logFatal(err)
		}

		log.Printf("Owner address: %s (ETH), %s (FIL)\n", ownerAddr, ownerDelAddr)
		log.Printf("Operator address: %s (ETH), %s (FIL)\n", operatorAddr, operatorDelAddr)
		log.Printf("Request key: %s (ETH), %s (FIL)\n", requestAddr, requestDelAddr)
		log.Println()
		log.Println("Please make sure to fund your Owner Address with FIL before creating an Agent")

		bs := util.BackupsStore()
		bs.Invalidate()
	},
}

func init() {
	walletCmd.AddCommand(createAgentAccountsCmd)
}
