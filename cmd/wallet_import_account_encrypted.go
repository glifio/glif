/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importAccountCmd represents the import-account command
var importAccountCmd = &cobra.Command{
	Use:   "import-account [account-name] [account-encrypted-key-json]",
	Short: "Import a single private key account, using an encrypted json key file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()

		overWrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			logFatal(err)
		}

		name := strings.ToLower(args[0])
		addrToOverwrite, err := as.Get(name)

		rename := fmt.Sprintf("%s-replaced-%s", name, time.Now().Format(time.RFC3339))

		var e *util.ErrKeyNotFound
		if !errors.As(err, &e) && !overWrite {
			logFatalf("Account %s already exists", name)
		} else if !errors.As(err, &e) {
			log.Printf("Warning: account '%s' already exists, renaming to '%s' and overriding with new '%s' key\n", name, rename, name)
		} else if overWrite {
			// here we dont actually have any keys to overwrite, so we set overWrite to false to avoid downstream issues
			overWrite = false
			log.Printf("Warning: no %s account to overwrite... importing %s\n", name, name)
		} else {
			log.Printf("Importing account: %s\n", name)
		}

		var passphrase string
		var message = "Passphrase for account (or hit enter for no passphrase)"
		prompt := &survey.Password{Message: message}
		survey.AskOne(prompt, &passphrase)

		re := regexp.MustCompile(`^[tf][0-9]`)
		if strings.HasPrefix(name, "0x") || re.MatchString(name) {
			logFatalf("Invalid name")
		}

		ks := util.KeyStore()

		pkJSON, err := hex.DecodeString(args[1])
		if err != nil {
			logFatalf("Invalid private key hex string")
		}

		account, err := ks.Import(pkJSON, passphrase, passphrase)
		if err != nil {
			logFatal(err)
		}

		// we rename the old named account to a new name so we dont lose a reference to this key
		if overWrite {
			as.Set(rename, addrToOverwrite)
		}

		as.Set(name, account.Address.String())

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		accountAddr, accountDelAddr, err := as.GetAddrs(name)
		if err != nil {
			logFatal(err)
		}

		bs := util.BackupsStore()
		bs.Invalidate()

		log.Printf("%s address: %s (ETH), %s (FIL) imported successfully\n", name, accountAddr, accountDelAddr)
	},
}

func init() {
	walletCmd.AddCommand(importAccountCmd)
	importAccountCmd.Flags().Bool("overwrite", false, "overwrite an existing account with the same name")
}
