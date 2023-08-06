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
	"github.com/spf13/cobra"
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
	generateHeader("BASIC INFO")
	fmt.Printf("Agent Address: %s\n", agent.String())
	fmt.Printf("Agent Address (del): %s\n", agentDel.String())
	fmt.Printf("Agent FIL ID Address: %s\n", agentFILIDAddr.String())
	fmt.Printf("Agent Owner: %s\n", owner.String())
	fmt.Printf("Agent Pools Protocol ID: %s\n", agentID)
	if goodVersion {
		fmt.Printf("Agent Version: %v ✅ \n", agVersion)
	} else {
		fmt.Println("Agent requires upgrade, run `glif agent upgrade` to upgrade")
		fmt.Printf("Agent/Network version mismatch: %v/%v ❌ \n", agVersion, ntwVersion)
	}
	fmt.Printf("Agent's pledged miner count: %v\n", len(agentMiners))
	s.Start()

	return agentID, agentFILIDAddr, agVersion, ntwVersion, nil
}

func econInfo(ctx context.Context, agent common.Address, agentID *big.Int, lapi *api.FullNodeStruct, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	assets, err := query.AgentLiquidAssets(ctx, agent, nil)
	if err != nil {
		return err
	}

	assetsFIL, _ := util.ToFIL(assets).Float64()

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
	apr.Quo(apr, big.NewFloat(1e36))

	weeklyPmt := new(big.Float).Mul(new(big.Float).SetInt(agentData.Principal), wpr)
	weeklyPmt.Quo(weeklyPmt, big.NewFloat(1e54))

	weekOneDeadline := new(big.Int).Add(defaultEpoch, big.NewInt(constants.EpochsInWeek*2))

	weekOneDeadlineTime := util.EpochHeightToTimestamp(weekOneDeadline, query.ChainID())
	defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch, query.ChainID())
	epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid, query.ChainID())

	s.Stop()

	generateHeader("ECON INFO")
	if lvl.Cmp(big.NewInt(0)) == 0 && chainID == constants.MainnetChainID {
		fmt.Println("Please open up a request to borrow on GitHub: https://github.com/glifio/infinity-pool-gov/issues/new?assignees=Schwartz10&labels=Entry+request+opened&projects=&template=infinity-pool-entry-request.md&title=%5BENTRY+REQUEST%5D")
	} else {
		if agentData.Principal.Cmp(big.NewInt(0)) == 0 {
			fmt.Println("Total borrowed: 0 FIL")
		} else {
			fmt.Printf("Total borrowed: %0.09f FIL\n", util.ToFIL(account.Principal))
			fmt.Printf("You currently owe: %.09f FIL\n", util.ToFIL(amountOwed))
			fmt.Printf("Current borrow APR: %.03f%%\n", apr.Mul(apr, big.NewFloat(100)))
			fmt.Printf("Your weekly payment: %0.09f FIL\n", weeklyPmt)

			// check to see we're still in good standing wrt making our weekly payment
			if account.EpochsPaid.Cmp(weekOneDeadline) == 1 {
				fmt.Printf("Your account owes its weekly payment (`to-current`) within the next: %s (by epoch # %s)\n", formatSinceDuration(weekOneDeadlineTime, epochsPaidTime), weekOneDeadline)
			} else {
				fmt.Printf("🔴 Overdue weekly payment 🔴\n")
				fmt.Printf("Your account *must* make a payment to-current within the next: %s (by epoch # %s)\n", formatSinceDuration(defaultEpochTime, epochsPaidTime), defaultEpoch)
			}
			fmt.Printf("Agent's quota is %.03f FIL\n", cap)

			fmt.Println()

			fmt.Printf("Agent's liquid assets (liquid FIL on your Agent): %0.08f FIL\n", assetsFIL)
			fmt.Printf("Agent's total assets (includes Miner's balances): %0.08f FIL\n", util.ToFIL(agentData.AgentValue))

			equity := new(big.Int).Sub(agentData.AgentValue, agentData.Principal)
			dte := econ.DebtToEquityRatio(agentData.Principal, equity)
			fmt.Printf("Agent's equity: %0.08f FIL - debt-to-equity: %0.03f%% (must stay below 100%%)\n", util.ToFIL(equity), dte.Mul(dte, big.NewFloat(100)))

			liquidationVal := new(big.Int).Div(agentData.AgentValue, big.NewInt(2))
			ltlv := econ.LoanToCollateralRatio(agentData.Principal, liquidationVal)
			fmt.Printf("Agent's liquidation value: %0.08f FIL - loan-to-liquidation %0.03f%% (must stay below 100%%)\n", util.ToFIL(liquidationVal), ltlv.Mul(ltlv, big.NewFloat(100)))

			dailyFees := rate.Mul(rate, big.NewInt(constants.EpochsInDay))
			dailyFees.Mul(dailyFees, agentData.Principal)
			dailyFees.Div(dailyFees, constants.WAD)

			dti := new(big.Int).Div(dailyFees, agentData.ExpectedDailyRewards)
			dtiFloat := new(big.Float).Mul(new(big.Float).SetInt(dti), big.NewFloat(100))
			dtiFloat.Quo(dtiFloat, big.NewFloat(1e18))

			weeklyEarnings := new(big.Int).Mul(agentData.ExpectedDailyRewards, big.NewInt(constants.EpochsInWeek))

			fmt.Printf("Agent's expected weekly earnings: %0.08f FIL - debt-to-income %0.03f%% (must stay below 25%%)\n", util.ToFIL(weeklyEarnings), dtiFloat)
			fmt.Println()
		}

		printWithBoldPreface("Agent's max borrow:", fmt.Sprintf("%0.09f FIL", util.ToFIL(maxBorrow)))
		printWithBoldPreface("Agent's max withdraw:", fmt.Sprintf("%0.09f FIL", util.ToFIL(maxWithdraw)))
	}

	s.Start()

	return nil
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
	fmt.Printf("Agent's administrator: %s\n", agentAdmin)
	fmt.Printf("Agent in default: %t\n\n", defaulted)
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

const headerWidth = 60

func generateHeader(title string) {
	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", title)
}

func init() {
	agentCmd.AddCommand(agentInfoCmd)
	agentInfoCmd.Flags().String("agent-addr", "", "Agent address")
}
