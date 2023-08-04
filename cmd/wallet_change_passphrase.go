/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

var changePassphraseCmd = &cobra.Command{
	Use:   "change-passphrase <address>",
	Short: "Change the passphrase for an encrypted key in the keystore",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ks := util.KeyStore()

		addr := common.HexToAddress(args[0])
		account := accounts.Account{Address: addr}

		oldPassphrase := ""
		err := ks.Unlock(account, "")
		if err != nil {
			prompt := &survey.Password{
				Message: "Old passphrase",
			}
			survey.AskOne(prompt, &oldPassphrase)
		}

		newPassphrase := ""
		prompt := &survey.Password{
			Message: "New passphrase",
		}
		survey.AskOne(prompt, &newPassphrase)

		err = ks.Update(account, oldPassphrase, newPassphrase)
		if err != nil {
			logFatal(err)
		}
		fmt.Println("Passphrase successfully changed.")
	},
}

func init() {
	walletCmd.AddCommand(changePassphraseCmd)
}
