/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a set of keys for the Agent",
	Long:  `Creates an owner and an operator key and stores the values in $HOME/.config/glif/keys.toml`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create the Ethereum private key
		// ownerPrivateKey, err := crypto.GenerateKey()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// hexutil.Encode(crypto.FromECDSA(ownerPrivateKey))

		// operatorPrivateKey, err := crypto.GenerateKey()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// hexutil.Encode(crypto.FromECDSA(operatorPrivateKey))

		//TODO: store private key in $HOME/.config/glif/token.toml

		// privateKeyBytes := crypto.FromECDSA(privateKey)
		// fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

		// fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

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
