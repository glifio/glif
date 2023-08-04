/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"

	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates keys from keys.toml to the new encrypted keystore",
	Run: func(cmd *cobra.Command, args []string) {
		err := migrateLegacyKeys()
		if err != nil {
			logFatal(err)
		}
		fmt.Println("Keys successfully migrated to the encrypted keystore.")
		fmt.Printf("For increased security, please delete the legacy %s/keys.toml file with the cleartext private keys.\n", cfgDir)
		fmt.Println("The legacy keys.toml file is no longer needed, unless you are just testing and plan to downgrade.")
	},
}

func init() {
	walletCmd.AddCommand(migrateCmd)
}

func migrateLegacyKeys() error {
	if err := migrateLegacyKey(util.OwnerKey); err != nil {
		return err
	}
	if err := migrateLegacyKey(util.OperatorKey); err != nil {
		return err
	}
	if err := migrateLegacyKey(util.RequestKey); err != nil {
		return err
	}
	return nil
}

func migrateLegacyKey(key util.KeyType) error {
	ksLegacy := util.KeyStoreLegacy()
	ks := util.KeyStore()
	as := util.AgentStore()

	pkStr, err := ksLegacy.Get(string(key))
	if err != nil {
		return err
	}
	if pkStr != "" {
		pk, err := ksLegacy.GetPrivate(key)
		if err != nil {
			return err
		}
		account, err := ks.ImportECDSA(pk, "")
		if err != nil {
			return err
		}
		as.Set(string(key), account.Address.String())
		fmt.Printf("Migrated %s key to encrypted key store with empty passphrase.\n", key)
	}
	return nil
}
