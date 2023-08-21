package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	walletutils "github.com/glifio/go-wallet-utils"
	"github.com/spf13/cobra"
)

var refreshRoutesCmd = &cobra.Command{
	Use:   "refresh-routes",
	Short: "Update cached routes on your Agent",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		from := cmd.Flag("from").Value.String()
		agentAddr, senderWallet, senderAccount, senderPassphrase, _, err := commonOwnerOrOperatorSetup(ctx, from)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		auth, err := walletutils.NewEthWalletTransactor(senderWallet, &senderAccount, senderPassphrase, big.NewInt(chainID))
		if err != nil {
			logFatal(err)
		}

		tx, err := PoolsSDK.Act().AgentRefreshRoutes(ctx, auth, agentAddr)
		if err != nil {
			logFatalf("Failed to refresh routes %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
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
