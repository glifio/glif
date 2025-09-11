package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plusApproveSpendCmd = &cobra.Command{
	Use:   "approve-spend <amount>",
	Short: "Set allowance for transfer of GLF tokens to SPPlus contract from owner",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		client, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to connect client %s", err)
		}
		defer client.Close()

		_, auth, _, _, err := commonSetupOwnerCall(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		glfAddr := PoolsSDK.Query().GLF()
		plusAddr := PoolsSDK.Query().SPPlus()

		poolTokenTransactor, err := abigen.NewPoolTokenTransactor(glfAddr, client)
		if err != nil {
			logFatalf("Failed to get GLF transactor %s", err)
		}

		tx, err := poolTokenTransactor.Approve(auth, plusAddr, amount)
		if err != nil {
			logFatalf("Failed to approve GLF spend %s", err)
		}

		tx, err = util.TxPostProcess(tx, err)
		if err != nil {
			logFatalf("Failed to approve GLF spend %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to approve GLF spend %s", err)
		}

		s.Stop()

		fmt.Println("GLF spend allowance set for SPPlus.")
	},
}

func init() {
	plusCmd.AddCommand(plusApproveSpendCmd)
}
