package cmd

import (
	"fmt"
	"log"
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
		_, pk, err := commonOwnerOrOperatorSetup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		strAddr := args[0]
		strAmt := args[1]
		fmt.Printf("Transferring %s iFIL balance to %s...", strAmt, strAddr)

		addr, err := ParseAddress(cmd.Context(), strAddr)
		if err != nil {
			log.Fatalf("Failed to parse address %s", err)
		}

		amt := big.NewInt(0)
		amt, ok := amt.SetString(strAmt, 10)
		if !ok {
			log.Fatalf("Failed to parse amount %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := PoolsSDK.Act().IFILTransfer(cmd.Context(), addr, amt, pk)
		if err != nil {
			log.Fatalf("Failed to transfer iFIL %s", err)
		}

		eapi, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			log.Fatalf("Failed to instantiate eth client %s", err)
		}
		defer eapi.Close()

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			log.Fatalf("Failed to transfer iFIL %s", err)
		}

		s.Stop()

		fmt.Printf("iFIL sent!")
	},
}

func init() {
	iFILCmd.AddCommand(iFILTransferCmd)
}
