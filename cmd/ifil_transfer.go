package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var iFILTransferCmd = &cobra.Command{
	Use:   "transfer [to] [amount]",
	Short: "Transfer iFIL to another address",
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
		fmt.Printf("Transferring %s iFIL balance to %s...\n", strAmt, strAddr)

		addr, err := AddressOrAccountNameToEVM(ctx, strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		amt := big.NewInt(0)
		amt, ok := amt.SetString(strAmt, 10)
		if !ok {
			logFatalf("Failed to parse amount %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		tx, err := PoolsSDK.Act().IFILTransfer(ctx, auth, addr, amt)
		if err != nil {
			logFatalf("Failed to transfer iFIL %s", err)
		}

		eapi, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to instantiate eth client %s", err)
		}
		defer eapi.Close()

		_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
		if err != nil {
			logFatalf("Failed to transfer iFIL %s", err)
		}

		s.Stop()

		fmt.Printf("iFIL sent!\n")
	},
}

func init() {
	iFILCmd.AddCommand(iFILTransferCmd)
	iFILTransferCmd.Flags().String("from", "default", "account to transfer iFIL from")
}
