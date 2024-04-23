/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateImportKeyParams(name string, overwrite bool) (string, string, string, error) {
	as := util.AccountsStore()

	var passphrase string
	addrToOverwrite, err := as.Get(name)

	rename := fmt.Sprintf("%s-replaced-%s", name, time.Now().Format(time.RFC3339))

	var e *util.ErrKeyNotFound
	if !errors.As(err, &e) && !overwrite {
		return passphrase, addrToOverwrite, rename, errors.New("Account already exists")
	} else if !errors.As(err, &e) {
		log.Printf("Warning: account '%s' already exists, renaming to '%s' and overriding with new '%s' key\n", name, rename, name)
	} else if overwrite {
		// here we dont actually have any keys to overwrite, so we set overwrite to false to avoid downstream issues
		overwrite = false
		log.Printf("Warning: no %s account to overwrite... importing %s\n", name, name)
	} else {
		log.Printf("Importing account: %s\n", name)
	}

	var message = "Passphrase for account (or hit enter for no passphrase)"
	prompt := &survey.Password{Message: message}
	survey.AskOne(prompt, &passphrase)

	re := regexp.MustCompile(`^[tf][0-9]`)
	if strings.HasPrefix(name, "0x") || re.MatchString(name) {
		return passphrase, addrToOverwrite, rename, errors.New("Invalid name")
	}
	return passphrase, addrToOverwrite, rename, nil
}

func completeImport(address common.Address, name string, rename string, addrToOverwrite string, overwrite bool) error {
	as := util.AccountsStore()

	// we rename the old named account to a new name so we dont lose a reference to this key
	if overwrite {
		as.Set(rename, addrToOverwrite)
	}

	as.Set(name, address.String())

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	accountAddr, accountDelAddr, err := as.GetAddrs(name)
	if err != nil {
		return err
	}

	bs := util.BackupsStore()
	bs.Invalidate()

	log.Printf("%s address: %s (ETH), %s (FIL) imported successfully\n", name, accountAddr, accountDelAddr)
	return nil
}

// importAccountRawCmd represents the import-account command
var importAccountRawCmd = &cobra.Command{
	Use:   "import-account-raw [account-name] [account-private-key]",
	Short: "Import a single private key account, using a raw, unencrypted private key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			logFatal(err)
		}

		name := strings.ToLower(args[0])
		passphrase, addrToOverwrite, rename, err := validateImportKeyParams(name, overwrite)
		if err != nil {
			logFatal(err)
		}

		pk := args[1]
		pkECDSA, err := crypto.HexToECDSA(pk)
		if err != nil {
			logFatalf("Invalid private key")
		}

		account, err := util.KeyStore().ImportECDSA(pkECDSA, passphrase)
		if err != nil {
			logFatal(err)
		}

		if err := completeImport(account.Address, name, rename, addrToOverwrite, overwrite); err != nil {
			logFatal(err)
		}
	},
}

func init() {
	walletCmd.AddCommand(importAccountRawCmd)
	importAccountRawCmd.Flags().Bool("overwrite", false, "overwrite an existing account with the same name")
}
