package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var iFILApproveCmd = &cobra.Command{
	Use:   "approve <spender> <allowance>",
	Short: "Approve another address to spend your iFIL",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		strAddr := args[0]
		strAmt := args[1]
		fmt.Printf("Approving %s to spend %s of your iFIL balance...\n", strAddr, strAmt)

		addr, err := AddressOrAccountNameToEVM(ctx, strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		amount, err := parseFILAmount(strAmt)
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().IFILApprove(ctx, auth, addr, amount)
		if err != nil {
			logFatalf("Failed to approve iFIL %s", err)
		}

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to approve iFIL %s", err)
		}

		s.Stop()

		fmt.Printf("iFIL approved!\n")
	},
}

func init() {
	iFILCmd.AddCommand(iFILApproveCmd)
	iFILApproveCmd.Flags().String("from", "default", "account to approve iFIL transfer from")
}
