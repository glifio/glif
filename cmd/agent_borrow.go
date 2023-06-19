/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/cli/events"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/util"
	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var borrowPreview bool

// borrowCmd represents the borrow command
var borrowCmd = &cobra.Command{
	Use:   "borrow <amount> [flags]",
	Short: "Borrow FIL from a Pool",
	Long:  "Borrow FIL from a Pool. If you do not pass a `pool-name` flag, the default pool is the Infinity Pool.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if borrowPreview {
			previewAction(cmd, args, constants.MethodBorrow)
			return
		}

		agentAddr, ownerKey, requesterKey, err := commonSetupOwnerCall()
		if err != nil {
			logFatal(err)
		}

		amount, err := parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}

		if amount.Cmp(util.WAD) == -1 {
			logFatal("Borrow amount must be greater than 1 FIL")
		}

		poolName := cmd.Flag("pool-name").Value.String()

		poolID, err := parsePoolType(poolName)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Borrowing %v FIL from the %s into agent %s\n", amount, poolName, agentAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		borrowevt := journal.RegisterEventType("agent", "borrow")
		evt := &events.AgentBorrow{
			AgentID: agentAddr.String(),
			PoolID:  poolID.String(),
			Amount:  amount.String(),
		}
		defer journal.Close()
		defer journal.RecordEvent(borrowevt, func() interface{} { return evt })

		// prevent attempting borrow action if the agent would get a PayUp error
		// a PayUp is triggered in the smart contract if
		// the agent epochs paid + max owed tolerance < current block height
		bh, err := PoolsSDK.Query().ChainHeight(cmd.Context())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		account, err := PoolsSDK.Query().AgentAccount(cmd.Context(), agentAddr, big.NewInt(0))
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		tolerance, err := PoolsSDK.Query().InfPoolMaxEpochsOwedTolerance(cmd.Context(), agentAddr)
		agentState := big.NewInt(0).Add(account.EpochsPaid, tolerance)

		if agentState.Cmp(bh) < 0 {
			evt.Error = fmt.Errorf("Agent needs to pay current debt before borrow more funds").Error()
			logFatal(err)
		}

		tx, err := PoolsSDK.Act().AgentBorrow(cmd.Context(), agentAddr, poolID, amount, ownerKey, requesterKey)
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}
		evt.Tx = tx.Hash().String()

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			evt.Error = err.Error()
			logFatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully borrowed %0.08f FIL\n", denoms.ToFIL(amount))
	},
}

func init() {
	agentCmd.AddCommand(borrowCmd)
	borrowCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to borrow from")
	borrowCmd.Flags().BoolVar(&borrowPreview, "preview", false, "preview the financial outcome of a borrow action")
}
