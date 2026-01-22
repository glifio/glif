/*
Copyright © 2026 Glif LTD
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var plusChangeOwnerForAgentCmd = &cobra.Command{
	Use:   "transfer-owner",
	Short: "Transfers the ownership of a GLIF Card to a new Agent owner",
	Long: `Transfers the ownership of a GLIF Card to a new Agent owner.
This command is useful when transferring a GLIF Card to a new owner because the Agent's owner address changed.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		agentAddr, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		fmt.Println("Transferring ownership of GLIF Card to new Agent owner...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().SPPlusChangeOwnerForAgent(ctx, auth, agentAddr)
		if err != nil {
			logFatalf("Failed to change owner for agent: %s", err)
		}

		s.Stop()

		fmt.Printf("Submitted transaction, confirming...: %s\n", tx.Hash().Hex())

		s.Start()
		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to confirm transaction: %s", err)
		}
		s.Stop()
		fmt.Println("Successfully changed Card owner to new Agent owner")
	},
}

func init() {
	plusCmd.AddCommand(plusChangeOwnerForAgentCmd)
}
