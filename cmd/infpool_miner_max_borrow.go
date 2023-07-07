/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/glifio/go-pools/mstat"
	psdk "github.com/glifio/go-pools/sdk"
	denoms "github.com/glifio/go-pools/util"
	"github.com/glifio/go-pools/vc"
	"github.com/spf13/cobra"
)

var infpoolMinerQuote = &cobra.Command{
	Use:   "miner-quote <miner-addr>",
	Short: "Returns the amount of FIL the miner can borrow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Generating a quote for %s\n", minerAddr)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		lCli, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		tipset, err := lCli.ChainHead(cmd.Context())
		if err != nil {
			logFatal(err)
		}

		minerstat, err := mstat.ComputeMinerStats(cmd.Context(), minerAddr, tipset, lCli)
		if err != nil {
			logFatal(err)
		}

		edr := new(big.Int).Add(minerstat.ExpectedDailyReward, new(big.Int).Div(minerstat.VestingFunds, big.NewInt(180)))

		agentData := &vc.AgentData{
			AgentValue:                  minerstat.Balance,
			CollateralValue:             new(big.Int).Div(minerstat.Balance, big.NewInt(2)),
			ExpectedDailyFaultPenalties: minerstat.PenaltyFaultPerDay,
			ExpectedDailyRewards:        edr,
			Gcred:                       big.NewInt(100),
			QaPower:                     minerstat.QualityAdjPower,
			Principal:                   common.Big0,
			FaultySectors:               minerstat.FaultySectors,
			LiveSectors:                 minerstat.LiveSectors,
			GreenScore:                  common.Big0,
		}

		nullCred, err := vc.NullishVerifiableCredential(*agentData)
		if err != nil {
			logFatal(err)
		}

		rate, err := PoolsSDK.Query().InfPoolGetRate(cmd.Context(), *nullCred)
		if err != nil {
			logFatal(err)
		}

		borrowNowMax := psdk.MaxBorrowFromAgentData(agentData, rate)
		if err != nil {
			logFatal(err)
		}

		fmt.Println()
		fmt.Println()

		s.Stop()

		fmt.Printf("Answer: %.04f FIL\n", denoms.ToFIL(borrowNowMax))
	},
}

func init() {
	infinitypoolCmd.AddCommand(infpoolMinerQuote)
}
