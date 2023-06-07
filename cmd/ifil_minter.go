package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var iFILMinterCmd = &cobra.Command{
	Use:   "minter",
	Short: "Get the contract address that can mint iFIL tokens",
	Run: func(cmd *cobra.Command, args []string) {
		minter, err := PoolsSDK.Query().IFILMinter(cmd.Context())
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		fmt.Printf("iFIL Minter addr: %s\n", minter)
	},
}

func init() {
	iFILCmd.AddCommand(iFILMinterCmd)
}
