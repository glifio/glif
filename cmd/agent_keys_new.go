/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a set of keys for the Agent",
	Long:  `Creates an owner and an operator key and stores the values in $HOME/.config/glif/keys.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		ownerKey := viper.GetString("keys.owner")
		operatorKey := viper.GetString("keys.operator")
		if ownerKey != "" || operatorKey != "" {
			log.Fatal("Keys already exists")
		}

		// Create the Ethereum private key
		ownerPrivateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		operatorPrivateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		KeyStorage.Set("owner", hexutil.Encode(crypto.FromECDSA(ownerPrivateKey))[2:])
		KeyStorage.Set("operator", hexutil.Encode(crypto.FromECDSA(operatorPrivateKey))[2:])
		// viper.Set("keys.owner", hexutil.Encode(crypto.FromECDSA(ownerPrivateKey))[2:])
		// viper.Set("keys.operator", hexutil.Encode(crypto.FromECDSA(operatorPrivateKey))[2:])

		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}

		ownerAddr, ownerDelAddr, err := deriveAddrFromPk(ownerPrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		operatorAddr, operatorDelAddr, err := deriveAddrFromPk(operatorPrivateKey)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Owner address: %s (ETH), %s (FIL)", ownerAddr, ownerDelAddr)
		log.Printf("Operator address: %s (ETH), %s (FIL)", operatorAddr, operatorDelAddr)
	},
}

func init() {
	keysCmd.AddCommand(newCmd)
}
