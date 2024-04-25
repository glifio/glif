/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmAccountCmd represents the import-account command
var rmAccountCmd = &cobra.Command{
	Use:   "remove-account [account-name]",
	Short: "Remove an account and its private key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()

		reallyDo, err := cmd.Flags().GetBool("really-do-it")
		if err != nil {
			logFatal(err)
		}

		if !reallyDo {
			logFatal("DANGEROUS COMMAND - are you really trying to export a raw private key from your wallet? If so, you must pass --really-do-it to complete the export")
		}

		name := strings.ToLower(args[0])
		addrToRemove, err := as.Get(name)

		var e *util.ErrKeyNotFound
		if errors.As(err, &e) {
			logFatalf("Account %s not found", name)
		} else {
			log.Printf("Removing account: %s, %s\n", name, addrToRemove)
		}

		var passphrase string
		var message = "Passphrase for account (or hit enter for no passphrase)"
		prompt := &survey.Password{Message: message}
		survey.AskOne(prompt, &passphrase)

		ks := util.KeyStore()

		account, err := ks.Find(accounts.Account{Address: common.HexToAddress(addrToRemove)})
		if err != nil {
			logFatal(err)
		}

		if err := ks.Delete(account, passphrase); err != nil {
			logFatal(err)
		}

		if err := as.Delete(name); err != nil {
			logFatal(err)
		}

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		log.Printf("Account %s removed successfully\n", name)
	},
}

func init() {
	walletCmd.AddCommand(rmAccountCmd)
	rmAccountCmd.Flags().Bool("really-do-it", false, "really remove the account")
}
