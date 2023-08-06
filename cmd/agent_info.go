/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/econ"
	"github.com/glifio/go-pools/rpc"
	"github.com/glifio/go-pools/util"
	"github.com/glifio/go-pools/vc"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

var agentInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get the info associated with your Agent",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			logFatal(err)
		}
		defer closer()

		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			logFatal(err)
		}

		agentAddrEthType, err := ethtypes.ParseEthAddress(agentAddr.String())
		if err != nil {
			logFatal(err)
		}

		agentAddrDel, err := agentAddrEthType.ToFilecoinAddress()
		if err != nil {
			logFatal(err)
		}

		agentID, _, _, _, err := basicInfo(cmd.Context(), agentAddr, agentAddrDel, lapi, s)
		if err != nil {
			logFatal(err)
		}

		err = econInfo(cmd.Context(), agentAddr, agentID, lapi, s)
		if err != nil {
			logFatal(err)
		}

		err = agentHealth(cmd.Context(), agentAddr, s)
		if err != nil {
			logFatal(err)
		}
	},
}

func basicInfo(ctx context.Context, agent common.Address, agentDel address.Address, lapi *api.FullNodeStruct, s *spinner.Spinner) (
	agentID *big.Int,
	agentFILIDAddr address.Address,
	agVersion uint8,
	ntwVersion uint8,
	err error,
) {
	query := PoolsSDK.Query()

	agentID, err = query.AgentID(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, err
	}

	agentFILIDAddr, err = lapi.StateLookupID(ctx, agentDel, types.EmptyTSK)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, err
	}

	agVersion, ntwVersion, err = query.AgentVersion(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, err
	}

	owner, err := query.AgentOwner(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, err
	}

	goodVersion := agVersion == ntwVersion

	agentMiners, err := query.MinerRegistryAgentMinersList(ctx, agentID, nil)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, err
	}

	s.Stop()

	versionCopy := fmt.Sprintf("%v ✅", agVersion)
	if !goodVersion {
		versionCopy = fmt.Sprintf("Please upgrade Agent ❌. Your version: %v, latest: %v", agVersion, ntwVersion)
	}

	basicInfoKeys := []string{
		"Agent 0x Addr",
		"Agent f4 Addr",
		"Agent f0 Addr",
		"Agent GLIF ID",
		"Agent Owner",
		"Agent Miners",
		"Version",
	}

	basicInfoValues := []string{
		agent.String(),
		agentDel.String(),
		agentFILIDAddr.String(),
		agentID.String(),
		owner.String(),
		fmt.Sprintf("%v", len(agentMiners)),
		versionCopy,
	}

	generateHeader("BASIC INFO")
	printTable(basicInfoKeys, basicInfoValues)

	s.Start()

	return agentID, agentFILIDAddr, agVersion, ntwVersion, nil
}

func econInfo(ctx context.Context, agent common.Address, agentID *big.Int, lapi *api.FullNodeStruct, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	assets, err := query.AgentLiquidAssets(ctx, agent, nil)
	if err != nil {
		return err
	}

	adoCloser, err := PoolsSDK.Extern().ConnectAdoClient(ctx)
	if err != nil {
		return err
	}
	defer adoCloser()

	agentData, err := rpc.ADOClient.AgentData(context.Background(), agent)
	if err != nil {
		return err
	}

	maxBorrow, err := PoolsSDK.Query().InfPoolAgentMaxBorrow(ctx, agent, agentData)
	if err != nil {
		return err
	}

	maxWithdraw, err := PoolsSDK.Query().InfPoolAgentMaxWithdraw(ctx, agent, agentData)
	if err != nil {
		return err
	}

	lvl, cap, err := query.InfPoolGetAgentLvl(ctx, agentID)
	if err != nil {
		return err
	}

	defaultEpoch, err := query.DefaultEpoch(ctx)
	if err != nil {
		return err
	}

	account, err := query.InfPoolGetAccount(ctx, agent, nil)
	if err != nil {
		return err
	}

	amountOwed, _, err := query.AgentOwes(ctx, agent)
	if err != nil {
		return err
	}

	nullCred, err := vc.NullishVerifiableCredential(*agentData)
	if err != nil {
		return err
	}

	rate, err := query.InfPoolGetRate(ctx, *nullCred)
	if err != nil {
		return err
	}

	wpr := new(big.Float).Mul(new(big.Float).SetInt(rate), big.NewFloat(constants.EpochsInWeek))

	apr := new(big.Float).Mul(new(big.Float).SetInt(rate), big.NewFloat(constants.EpochsInYear))
	apr.Quo(apr, big.NewFloat(1e34))

	weeklyEarnings := new(big.Int).Mul(agentData.ExpectedDailyRewards, big.NewInt(7))

	weeklyPmt := new(big.Float).Mul(new(big.Float).SetInt(agentData.Principal), wpr)
	weeklyPmt.Quo(weeklyPmt, big.NewFloat(1e54))

	weekOneDeadline := new(big.Int).Add(defaultEpoch, big.NewInt(constants.EpochsInWeek*2))

	weekOneDeadlineTime := util.EpochHeightToTimestamp(weekOneDeadline, query.ChainID())
	defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch, query.ChainID())
	epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid, query.ChainID())

	equity := new(big.Int).Sub(agentData.AgentValue, agentData.Principal)
	dte := econ.DebtToEquityRatio(agentData.Principal, equity)

	dailyFees := rate.Mul(rate, big.NewInt(constants.EpochsInDay))
	dailyFees.Mul(dailyFees, agentData.Principal)
	dailyFees.Div(dailyFees, constants.WAD)

	s.Stop()

	generateHeader("ECON INFO")
	if lvl.Cmp(big.NewInt(0)) == 0 && chainID == constants.MainnetChainID {
		fmt.Println()
		fmt.Println(chalk.Bold.TextStyle("Please open up a request for quota on GitHub: https://tinyurl.com/glif-entry-request"))
	}
	if agentData.Principal.Cmp(big.NewInt(0)) == 0 {
		nothingBorrowedKeys := []string{
			"Total borrowed",
			"Agent's liquid FIL",
			"Agent's total FIL",
			"Agent's equity",
			"Agent's expected weekly earnings",
			"Agent's quota",
		}

		nothingBorrowedValues := []string{
			"0 FIL",
			fmt.Sprintf("%0.08f FIL", util.ToFIL(assets)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(agentData.AgentValue)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(equity)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(weeklyEarnings)),
			fmt.Sprintf("%.03f FIL", cap),
		}
		printTable(nothingBorrowedKeys, nothingBorrowedValues)
	} else {
		somethingBorrowedKeys := []string{
			"Total borrowed",
			"You current owe",
			"Current borrow APR",
			"Your weekly payment",
			"Quota",
		}

		somethingBorrowedValues := []string{
			fmt.Sprintf("%0.09f FIL", util.ToFIL(account.Principal)),
			fmt.Sprintf("%0.09f FIL", util.ToFIL(amountOwed)),
			fmt.Sprintf("%.03f%%", apr),
			fmt.Sprintf("%0.09f FIL", weeklyPmt),
			fmt.Sprintf("%.03f FIL", cap),
		}
		printTable(somethingBorrowedKeys, somethingBorrowedValues)

		dti := new(big.Int).Div(dailyFees, agentData.ExpectedDailyRewards)
		dtiFloat := new(big.Float).Mul(new(big.Float).SetInt(dti), big.NewFloat(100))
		dtiFloat.Quo(dtiFloat, big.NewFloat(1e18))

		coreEconKeys := []string{
			"Agent's liquid FIL",
			"Agent's total FIL",
			"Agent's equity",
			"Agent's expected weekly earnings",
			"Agent's debt-to-equity (DTE)",
			"Agent's debt-to-income (DTI)",
		}

		coreEconValues := []string{
			fmt.Sprintf("%0.08f FIL", util.ToFIL(assets)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(agentData.AgentValue)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(equity)),
			fmt.Sprintf("%0.08f FIL", util.ToFIL(weeklyEarnings)),
			fmt.Sprintf("%0.03f%% (must stay below 100%%)", dte.Mul(dte, big.NewFloat(100))),
			fmt.Sprintf("%0.03f%% (must stay below 25%%)", dtiFloat),
		}

		printTable(coreEconKeys, coreEconValues)
	}

	printTable([]string{
		"Agent's max borrow",
		"Agent's max withdraw",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(maxBorrow)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(maxWithdraw)),
	})

	// check to see we're still in good standing wrt making our weekly payment
	fmt.Println()
	if account.Principal.Cmp(big.NewInt(0)) > 0 {
		if account.EpochsPaid.Cmp(weekOneDeadline) == 1 {
			fmt.Printf("Your account owes its weekly payment (`to-current`) within the next: %s (by epoch # %s)\n", formatSinceDuration(weekOneDeadlineTime, epochsPaidTime), weekOneDeadline)
		} else {
			fmt.Printf("🔴 Overdue weekly payment 🔴\n")
			fmt.Printf("Your account *must* make a payment to-current within the next: %s (by epoch # %s)\n", formatSinceDuration(defaultEpochTime, epochsPaidTime), defaultEpoch)
		}
	}

	s.Start()

	return nil
}

func printTable(keys []string, values []string) {
	// here we hacky get the same width for all separate tables in the info command by making the first row have a long width
	tbl := table.New("                                    ", "")

	for i, k := range keys {
		tbl.AddRow(k, values[i])
	}

	tbl.Print()
}

func agentHealth(ctx context.Context, agent common.Address, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	agentAdmin, err := query.AgentAdministrator(ctx, agent)
	if err != nil {
		return err
	}

	defaulted, err := query.AgentDefaulted(ctx, agent)
	if err != nil {
		return err
	}

	faultySectorStart, err := query.AgentFaultyEpochStart(ctx, agent)
	if err != nil {
		return err
	}

	s.Stop()

	generateHeader("HEALTH")
	printTable([]string{
		"Agent's administrator",
		"Agent in default",
	}, []string{
		agentAdmin.String(),
		fmt.Sprintf("%t", defaulted),
	})

	fmt.Println()

	if faultySectorStart.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("Status healthy 🟢\n")
	} else {
		chainHeight, err := query.ChainHeight(ctx)
		if err != nil {
			return err
		}

		consecutiveFaultEpochTolerance, err := query.MaxConsecutiveFaultEpochs(ctx)
		if err != nil {
			return err
		}

		consecutiveFaultEpochs := new(big.Int).Sub(chainHeight, faultySectorStart)

		liableForFaultySectorDefault := consecutiveFaultEpochs.Cmp(consecutiveFaultEpochTolerance) >= 0

		if liableForFaultySectorDefault {
			fmt.Printf("🔴 Status unhealthy - you are at risk of liquidation due to consecutive faulty sectors 🔴\n")
			fmt.Printf("Faulty sector start epoch: %v\n", faultySectorStart)
		} else {
			epochsBeforeZeroTolerance := new(big.Int).Sub(consecutiveFaultEpochTolerance, consecutiveFaultEpochs)
			fmt.Printf("🟡 Status unhealthy - you are approaching risk of liquidation due to consecutive faulty sectors 🟡\n")
			fmt.Printf("- With %v more consecutive faulty sectors, you will be at risk of liquidation\n", epochsBeforeZeroTolerance)
		}
	}
	fmt.Println()

	return nil
}

func formatSinceDuration(t1 time.Time, t2 time.Time) string {
	d := t2.Sub(t1).Round(time.Minute)

	var parts []string

	weeks := int(d.Hours()) / (24 * 7)
	d -= time.Duration(weeks) * 7 * 24 * time.Hour
	if weeks > 1 {
		parts = append(parts, fmt.Sprintf("%d weeks", weeks))
	} else if weeks == 1 {
		parts = append(parts, fmt.Sprintf("%d week", weeks))
	}

	days := int(d.Hours()) / 24
	d -= time.Duration(days) * 24 * time.Hour
	if days > 1 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	} else if days == 1 {
		parts = append(parts, fmt.Sprintf("%d day", days))
	}

	h := d / time.Hour
	d -= h * time.Hour
	parts = append(parts, fmt.Sprintf("%02d hours", h))

	m := d / time.Minute
	parts = append(parts, fmt.Sprintf("and %02d minutes", m))

	return strings.Join(parts, " ")
}

func generateHeader(title string) {
	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", chalk.Underline.TextStyle(title))
}

func init() {
	agentCmd.AddCommand(agentInfoCmd)
	agentInfoCmd.Flags().String("agent-addr", "", "Agent address")
}
