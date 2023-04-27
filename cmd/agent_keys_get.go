/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the addresses associated with your owner and operator keys",
	Run: func(cmd *cobra.Command, args []string) {
		ownerKey := viper.GetString("keys.owner")
		operatorKey := viper.GetString("keys.operator")
		if ownerKey == "" || operatorKey == "" {
			log.Fatal("Missing keys")
		}

		ownerPrivateKey, err := crypto.HexToECDSA(ownerKey)
		if err != nil {
			log.Fatal(err)
		}
		operatorPrivateKey, err := crypto.HexToECDSA(operatorKey)
		if err != nil {
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
	keysCmd.AddCommand(getCmd)
}
