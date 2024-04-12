/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"encoding/hex"
	"strings"

	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

// importAccountCmd represents the import-account command
// the validateImportKeyParams and completeImport functions are defined in the importAccountRawCmd
var importAccountCmd = &cobra.Command{
	Use:   "import-account [account-name] [account-encrypted-key-json]",
	Short: "Import a single private key account, using an encrypted json key file",
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

		pkJSON, err := hex.DecodeString(args[1])
		if err != nil {
			logFatalf("Invalid private key JSON file hex string")
		}

		account, err := util.KeyStore().Import(pkJSON, passphrase, passphrase)
		if err != nil {
			logFatal(err)
		}

		if err := completeImport(account.Address, name, rename, addrToOverwrite, overwrite); err != nil {
			logFatal(err)
		}
	},
}

func init() {
	walletCmd.AddCommand(importAccountCmd)
	importAccountCmd.Flags().Bool("overwrite", false, "overwrite an existing account with the same name")
}
