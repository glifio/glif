/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func panicIfKeyExists(key util.KeyType, addr common.Address, err error) {
	if err != nil {
		log.Fatal(err)
	}

	if !util.IsZeroAddress(addr) {
		log.Fatalf("Key already exists for %s", key)
	}
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a set of keys",
	Long:  `Creates an owner, an operator, and a requester key and stores the values in $HOME/.config/glif/keys.toml. Note that the owner and requester keys are only applicable to Agents, the operator key is the primary key for interacting with smart contracts.`,
	Run: func(cmd *cobra.Command, args []string) {
		ks := util.KeyStore()

		ownerAddr, ownerDelAddr, err := ks.GetAddrs(util.OwnerKey)
		panicIfKeyExists(util.OwnerKey, ownerAddr, err)

		operatorAddr, operatorDelAddr, err := ks.GetAddrs(util.OperatorKey)
		panicIfKeyExists(util.OperatorKey, operatorAddr, err)

		requestAddr, requestDelAddr, err := ks.GetAddrs(util.RequestKey)
		panicIfKeyExists(util.RequestKey, requestAddr, err)

		// Create the Ethereum private key
		ownerPrivateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		operatorPrivateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		requestPrivateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		if err := ks.SetKey(util.OwnerKey, ownerPrivateKey); err != nil {
			log.Fatal(err)
		}

		if err := ks.SetKey(util.OperatorKey, operatorPrivateKey); err != nil {
			log.Fatal(err)
		}

		if err := ks.SetKey(util.RequestKey, requestPrivateKey); err != nil {
			log.Fatal(err)
		}

		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}

		ownerAddr, ownerDelAddr, err = ks.GetAddrs(util.OwnerKey)
		if err != nil {
			log.Fatal(err)
		}
		operatorAddr, operatorDelAddr, err = ks.GetAddrs(util.OperatorKey)
		if err != nil {
			log.Fatal(err)
		}
		requestAddr, requestDelAddr, err = ks.GetAddrs(util.RequestKey)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Owner address: %s (ETH), %s (FIL)\n", ownerAddr, ownerDelAddr)
		log.Printf("Operator address: %s (ETH), %s (FIL)\n", operatorAddr, operatorDelAddr)
		log.Printf("Request key: %s (ETH), %s (FIL)\n", requestAddr, requestDelAddr)
	},
}

func init() {
	walletCmd.AddCommand(newCmd)

	//TODO: add flags that allow for specific keys to be generated
}
