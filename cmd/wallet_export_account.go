/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

// exportAccountCmd represents the export-account command
var exportAccountCmd = &cobra.Command{
	Use:   "export-account [account-name]",
	Short: "(Dangerous) Export a single private key account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reallyDo, err := cmd.Flags().GetBool("really-do-it")
		if err != nil {
			logFatal(err)
		}

		lotus, err := cmd.Flags().GetBool("lotus")
		if err != nil {
			logFatal(err)
		}

		if !reallyDo {
			logFatal("DANGEROUS COMMAND - are you really trying to export a raw private key from your wallet? If so, you must pass --really-do-it to complete the export")
		}

		addrNameToExport := args[0]

		as := util.AccountsStore()
		addrStr, err := as.Get(addrNameToExport)
		if err != nil {
			logFatal(err)
		}

		ks := util.KeyStore()

		account, err := ks.Find(accounts.Account{Address: common.HexToAddress(addrStr)})
		if err != nil {
			logFatal(err)
		}

		var passphrase string
		var message = "Passphrase for account"
		err = ks.Unlock(account, "")
		if err != nil {
			prompt := &survey.Password{Message: message}
			survey.AskOne(prompt, &passphrase)
			if passphrase == "" {
				fmt.Println("Aborted")
				return
			}
		}

		keyJSON, err := ks.Export(account, passphrase, passphrase)
		if err != nil {
			logFatal(err)
		}
		if lotus {
			pk, err := keystore.DecryptKey(keyJSON, passphrase)
			if err != nil {
				logFatal(err)
			}
			keyBytes := crypto.FromECDSA(pk.PrivateKey)
			ki := types.KeyInfo{
				Type:       types.KTDelegated,
				PrivateKey: keyBytes,
			}
			b, err := json.Marshal(ki)
			if err != nil {
				logFatal(err)
			}
			fmt.Println(hex.EncodeToString(b))
		} else {
			fmt.Println(hex.EncodeToString(keyJSON))
		}
	},
}

func init() {
	walletCmd.AddCommand(exportAccountCmd)
	exportAccountCmd.Flags().Bool("really-do-it", false, "really export the account")
	exportAccountCmd.Flags().Bool("lotus", false, "export in format that Lotus can import")
}
