/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"log"

	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// labelAccountCmd represents the label-account command
var labelAccountCmd = &cobra.Command{
	Use:   "label-account <name> <address>",
	Short: "Label an account with a human readable name",
	Long:  "Labeling an account creates a read-only alias for an account's address.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AccountsStore()

		name := args[0]
		_, _, err := as.GetAddrs(name)
		if err == nil {
			logFatalf("%s account already created\n", name)
		} else {
			var e *util.ErrKeyNotFound
			if !errors.As(err, &e) {
				logFatal(err)
			}
		}

		addr, err := AddressOrAccountNameToEVM(cmd.Context(), args[1])
		if err != nil {
			logFatal(err)
		}

		as.Set(name, addr.Hex())

		if err := viper.WriteConfig(); err != nil {
			logFatal(err)
		}

		if addr.Hex() != args[1] {
			log.Printf("Converting %s into its EVM representation: %s\n", args[1], addr.Hex())
		}

		log.Printf("Successfully added new read-only account to wallet - %s\n", addr.Hex())
	},
}

func init() {
	walletCmd.AddCommand(labelAccountCmd)
}
