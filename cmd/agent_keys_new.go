/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a set of keys for the Agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: create .glif home directory, if it doesn't exist

		// Create the Ethereum private key
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		//TODO: store private key in $HOME/.glif/token

		// privateKeyBytes := crypto.FromECDSA(privateKey)
		// fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
		fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	},
}

func init() {
	keysCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
