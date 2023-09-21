//go:build advanced
// +build advanced

/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newRequestKeyCmd represents the new-key command
var newRequestKeyCmd = &cobra.Command{
	Use:   "new-key <name>",
	Short: "Create a new request key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]

		if keyName != string(util.OwnerKey) &&
			keyName != string(util.OperatorKey) &&
			keyName != string(util.RequestKey) {
			logFatalf("Invalid Agent key name passed %s. Key must be one of `owner`, `operator`, or `request`", keyName)
		}

		as := util.AccountsStore()

		fmt.Printf("Creating new %s key for Agent\n", keyName)

		passphrase, envSet := os.LookupEnv("GLIF_PASSPHRASE")
		// only prompt for passphrase if it's owner key
		if !envSet && keyName == string(util.OwnerKey) {
			prompt := &survey.Password{
				Message: "Please type a passphrase to encrypt your Agent's owner key",
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
		}

		key, err := as.Get(keyName)
		if err != util.ErrKeyNotFound {
			// rename the existing key
			newKeyName := fmt.Sprintf("%s-%s", keyName, time.Now().Format(time.RFC3339))
			as.Set(newKeyName, key)
			fmt.Printf("Renamed existing %s key to %s\n", keyName, newKeyName)
		}

		ks := util.KeyStore()

		account, err := ks.NewAccount(passphrase)
		if err != nil {
			logFatal(err)
		}

		as.Set(keyName, account.Address.String())

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		accountAddr, accountDelAddr, err := as.GetAddrs(keyName)
		if err != nil {
			logFatal(err)
		}

		bs := util.BackupsStore()
		bs.Invalidate()

		log.Printf("Created new %s key for Agent. Address: %s (ETH), %s (FIL)\n", keyName, accountAddr, accountDelAddr)
	},
}

func init() {
	adminCmd.AddCommand(newRequestKeyCmd)
}
