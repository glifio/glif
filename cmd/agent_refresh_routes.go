package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var refreshRoutesCmd = &cobra.Command{
	Use:   "refresh-routes",
	Short: "Update cached routes on your Agent",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, senderWallet, senderAccount, senderPassphrase, _, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().AgentRefreshRoutes(cmd.Context(), agentAddr, senderWallet, senderAccount, senderPassphrase)
		if err != nil {
			logFatalf("Failed to refresh routes %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatalf("Failed to refresh routes %s", err)
		}

		s.Stop()

		fmt.Printf("Routes refreshed!\n")
	},
}

func init() {
	agentCmd.AddCommand(refreshRoutesCmd)
	refreshRoutesCmd.Flags().String("from", "", "address to send the transaction from")
}
