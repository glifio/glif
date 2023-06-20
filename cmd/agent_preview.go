/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/filecoin-project/go-address"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/rpc"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use: "preview",
}

func init() {
	agentCmd.AddCommand(previewCmd)
}

func previewAction(cmd *cobra.Command, args []string, action constants.Method) {
	agentAddr, err := getAgentAddress(cmd)
	if err != nil {
		logFatal(err)
	}

	var minerAddr address.Address
	var amount *big.Int

	switch action {
	case constants.MethodAddMiner, constants.MethodRemoveMiner:
		minerAddr, err = address.NewFromString(args[0])
		if err != nil {
			logFatal(err)
		}
		amount = big.NewInt(0)

	case constants.MethodBorrow, constants.MethodPay, constants.MethodWithdraw:
		amount, err = parseFILAmount(args[0])
		if err != nil {
			logFatal(err)
		}
		minerAddr = address.Undef

	default:
		err = fmt.Errorf("unsupported preview action")
		logFatal(err)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	closer, err := PoolsSDK.Extern().ConnectAdoClient(cmd.Context())
	if err != nil {
		logFatal(err)
	}
	defer closer()

	agentDataBefore, err := rpc.ADOClient.AgentData(cmd.Context(), agentAddr)
	if err != nil {
		logFatal(err)
	}

	agentDataAfter, err := rpc.ADOClient.PreviewAction(cmd.Context(), agentAddr, minerAddr, amount, constants.MethodAddMiner)
	if err != nil {
		logFatal(err)
	}

	rateAfter, err := PoolsSDK.Query().InfPoolRateFromGCRED(cmd.Context(), agentDataAfter.Gcred)
	if err != nil {
		logFatal(err)
	}

	wpr := new(big.Float).Mul(rateAfter, big.NewFloat(constants.EpochsInWeek))
	wprFloat, _ := wpr.Float64()
	aprFloat, _ := new(big.Float).Mul(rateAfter, big.NewFloat(constants.EpochsInYear)).Float64()

	weeklyPmt := new(big.Float).Mul(new(big.Float).SetInt(agentDataAfter.Principal), wpr)
	weeklyPmt.Quo(weeklyPmt, big.NewFloat(1e18))

	generateHeader(fmt.Sprintf("PREVIEW %s", strings.ToUpper(string(action))))
	fmt.Printf("Total borrowed before/after: %0.09f => %0.09f\n", util.ToFIL(agentDataBefore.Principal), util.ToFIL(agentDataAfter.Principal))
	fmt.Printf("GCRED before/after: %s => %s\n", agentDataBefore.Gcred, agentDataAfter.Gcred)
	fmt.Printf("The weekly/annual fee rate: %.03f%% / %.03f%%\n", wprFloat*100, aprFloat*100)
	fmt.Printf("Your weekly min payment will be: %.06f FIL\n", weeklyPmt)
}
