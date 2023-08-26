/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createAccountCmd represents the create-account command
var createAccountCmd = &cobra.Command{
	Use:   "create-account [account-name]",
	Short: "Create a single named account",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()

		var name string
		if len(args) == 1 {
			name = args[0]
		} else {
			name = "default"
		}

		_, err := as.Get(name)
		if err != util.ErrKeyNotFound {
			logFatalf("Account %s already exists", name)
		}

		fmt.Println("Creating account:", name)

		var passphrase string
		prompt := &survey.Password{
			Message: "Please type a passphrase to encrypt your private key",
		}
		survey.AskOne(prompt, &passphrase)
		var confirmPassphrase string
		confirmPrompt := &survey.Password{
			Message: "Confirm passphrase",
		}
		survey.AskOne(confirmPrompt, &confirmPassphrase)
		if passphrase != confirmPassphrase {
			logFatal("Aborting. Passphrase confirmation did not match.")
		}

		ks := util.KeyStore()

		account, err := ks.NewAccount(passphrase)
		if err != nil {
			logFatal(err)
		}

		as.Set(name, account.Address.String())

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		accountAddr, accountDelAddr, err := as.GetAddrs(name)
		if err != nil {
			logFatal(err)
		}

		log.Printf("%s address: %s (ETH), %s (FIL)\n", name, accountAddr, accountDelAddr)
	},
}

func init() {
	walletCmd.AddCommand(createAccountCmd)
}
