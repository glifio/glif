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
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/econ"
	"github.com/glifio/go-pools/terminate"
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

		agentAddr, err := getAgentAddressWithFlags(cmd)
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

		agentID, _, _, _, afi, err := basicInfo(cmd.Context(), agentAddr, agentAddrDel, lapi, s)
		if err != nil {
			logFatal(err)
		}

		err = econInfo(cmd.Context(), agentAddr, agentID, afi, lapi, s)
		if err != nil {
			logFatal(err)
		}

		// err = agentHealth(cmd.Context(), agentAddr, agentData, ats, s)
		// if err != nil {
		// 	logFatal(err)
		// }
	},
}

func basicInfo(ctx context.Context, agent common.Address, agentDel address.Address, lapi *api.FullNodeStruct, s *spinner.Spinner) (
	agentID *big.Int,
	agentFILIDAddr address.Address,
	agVersion uint8,
	ntwVersion uint8,
	afi *econ.AgentFi,
	err error,
) {
	query := PoolsSDK.Query()

	tasks := []util.TaskFunc{
		func() (interface{}, error) {
			return query.AgentID(ctx, agent)
		},
		func() (interface{}, error) {
			return lapi.StateLookupID(ctx, agentDel, types.EmptyTSK)
		},
		func() (interface{}, error) {
			ver, ntwVer, err := query.AgentVersion(ctx, agent)
			return []interface{}{ver, ntwVer}, err
		},
		func() (interface{}, error) {
			return query.AgentOwner(ctx, agent)
		},
		func() (interface{}, error) {
			return query.AgentOperator(ctx, agent)
		},
		func() (interface{}, error) {
			return query.AgentRequester(ctx, agent)
		},
		func() (interface{}, error) {
			agentID, err := query.AgentID(ctx, agent)
			if err != nil {
				return nil, err
			}
			return query.MinerRegistryAgentMinersList(ctx, agentID, nil)
		},
		func() (interface{}, error) {
			return query.AgentLiquidAssets(ctx, agent, nil)
		},
		func() (interface{}, error) {
			return fetchAgentEconFromAPI(agent)
		},
	}
	results, err := util.Multiread(tasks)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, nil, err
	}

	agentID = results[0].(*big.Int)
	agentFILIDAddr = results[1].(address.Address)
	versionResults := results[2].([]interface{})
	agVersion = versionResults[0].(uint8)
	ntwVersion = versionResults[1].(uint8)
	owner := results[3].(common.Address)
	operator := results[4].(common.Address)
	requester := results[5].(common.Address)
	agentMiners := results[6].([]address.Address)
	agentLiquidFIL := results[7].(*big.Int)
	afi = results[8].(*econ.AgentFi)

	goodVersion := agVersion == ntwVersion

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
		"Agent Operator",
		"Agent ADO Requester",
		"Agent Balance",
		"Agent Total Value",
		"Agent Miners",
		"Version",
	}

	basicInfoValues := []string{
		agent.String(),
		agentDel.String(),
		agentFILIDAddr.String(),
		agentID.String(),
		owner.String(),
		operator.String(),
		requester.String(),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(agentLiquidFIL)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.Balance)),
		fmt.Sprintf("%v", len(agentMiners)),
		versionCopy,
	}

	generateHeader("BASIC INFO")
	printTable(basicInfoKeys, basicInfoValues)

	s.Start()

	return agentID, agentFILIDAddr, agVersion, ntwVersion, afi, nil
}

func econInfo(ctx context.Context, agent common.Address, agentID *big.Int, afi *econ.AgentFi, lapi *api.FullNodeStruct, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	tasks := []util.TaskFunc{
		func() (interface{}, error) {
			return query.AgentPrincipal(ctx, agent, nil)
		},
		func() (interface{}, error) {
			amountOwed, err := query.AgentInterestOwed(ctx, agent, nil)
			if err != nil {
				return nil, err
			}
			return amountOwed, nil
		},
		func() (interface{}, error) {
			return query.InfPoolGetRate(ctx)
		},
	}

	results, err := util.Multiread(tasks)
	if err != nil {
		return err
	}
	principal := results[0].(*big.Int)
	interestOwed := results[1].(*big.Int)
	totalDebt := new(big.Int).Add(principal, interestOwed)
	rate := results[2].(*big.Int)

	apr := new(big.Float).Mul(new(big.Float).SetInt(rate), big.NewFloat(constants.EpochsInYear))
	apr.Quo(apr, big.NewFloat(1e34))

	s.Stop()

	generateHeader("ECON INFO")

	printTable([]string{
		"Liquidation Value",
		"Total Debt",
		"Debt to liquidation value % (DTL)",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LiquidationValue())),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(totalDebt)),
		fmt.Sprintf("%0.02f%%", afi.DTL()),
	})

	printTable([]string{
		"Max borrow to seal",
		"Max borrow to withdraw",
		"Current borrow limit",
		"Current withdraw limit",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.MaxBorrowAndSeal())),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.MaxBorrowAndWithdraw())),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.BorrowLimit())),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.WithdrawLimit())),
	})

	printTable([]string{
		"Total borrowed",
		"Interest owed",
		"APR",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(principal)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(interestOwed)),
		fmt.Sprintf("%.02f%%", apr),
	})

	printTable([]string{
		"Liquidation Value Breakdown",
		"Available balance",
		"Intial Pledge",
		"Locked Rewards",
		"Total Liquidation Value",
	}, []string{
		"",
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.AvailableBalance)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.InitialPledge)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LockedRewards)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LiquidationValue())),
	})

	s.Start()

	return nil
}

func bigIntAttoToPercent(atto *big.Int) *big.Float {
	return new(big.Float).Mul(util.ToFIL(atto), big.NewFloat(100))
}

func printTable(keys []string, values []string) {
	// here we hacky get the same width for all separate tables in the info command by making the first row have a long width
	tbl := table.New("                                    ", "")

	for i, k := range keys {
		tbl.AddRow(k, values[i])
	}

	tbl.Print()
}

func agentHealth(ctx context.Context, agent common.Address, agentData *vc.AgentData, ats terminate.PreviewAgentTerminationSummary, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	tasks := []util.TaskFunc{
		func() (interface{}, error) {
			return query.AgentAdministrator(ctx, agent)
		},

		func() (interface{}, error) {
			return query.AgentDefaulted(ctx, agent)
		},
		func() (interface{}, error) {
			return big.NewInt(0), nil
			// return query.AgentFaultyEpochStart(ctx, agent)
		},
		func() (interface{}, error) {
			// return query.DefaultEpoch(ctx)
			return big.NewInt(0), nil
		},
		func() (interface{}, error) {
			return query.InfPoolGetAccount(ctx, agent, nil)
		},
	}

	results, err := util.Multiread(tasks)
	if err != nil {
		return err
	}

	agentAdmin := results[0].(common.Address)
	defaulted := results[1].(bool)
	// faultySectorStart := results[2].(*big.Int)
	defaultEpoch := results[3].(*big.Int)
	account := results[4].(abigen.Account)

	// overLTV := ats.LTV(agentData.Principal).Cmp(constants.MAX_LTV) > 0

	weekOneDeadline := new(big.Int).Add(defaultEpoch, big.NewInt(constants.EpochsInWeek*2))

	// weekOneDeadlineTime := util.EpochHeightToTimestamp(weekOneDeadline, query.ChainID())
	defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch, query.ChainID())
	epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid, query.ChainID())

	s.Stop()

	generateHeader("HEALTH")
	fmt.Println()
	// check to see we're still in good standing wrt making our weekly payment
	owesPmt := account.Principal.Cmp(big.NewInt(0)) > 0
	badPmtStatus := owesPmt && account.EpochsPaid.Cmp(weekOneDeadline) < 1
	// badFaultStatus := faultySectorStart.Cmp(big.NewInt(0)) > 0

	faultRatio := big.NewFloat(0)

	// check to see if we have faulty sectors (regardless of the Agent's state)
	// pendingBadFaultStatus := false
	// if agentData.LiveSectors.Int64() > 0 {
	// 	faultRatio = new(big.Float).Quo(new(big.Float).SetInt(agentData.FaultySectors), new(big.Float).SetInt(agentData.LiveSectors))
	// 	// faulty sectors exist over the limit
	// 	if faultRatio.Cmp(constants.FAULTY_SECTOR_TOLERANCE) > 0 {
	// 		pendingBadFaultStatus = true
	// 	}
	// }

	// convert faults into percentage for logging
	faultRatio = faultRatio.Mul(faultRatio, big.NewFloat(100))
	// convert limit into percentage for logging
	// limit := new(big.Float).Mul(constants.FAULTY_SECTOR_TOLERANCE, big.NewFloat(100))

	// if !badPmtStatus && !badFaultStatus && !pendingBadFaultStatus && !overLTV {
	// 	fmt.Printf("Status healthy 🟢\n")
	// 	if owesPmt {
	// 		fmt.Printf("Your account owes its weekly payment (`to-current`) within the next: %s (by epoch # %s)\n", formatSinceDuration(weekOneDeadlineTime, epochsPaidTime), weekOneDeadline)
	// 	}
	// } else {
	// 	fmt.Println(chalk.Bold.TextStyle("Status unhealthy 🔴"))
	// }

	// if overLTV {
	// 	fmt.Printf("WARNING: Your Agent is over the LTV limit of %0.00f%%\n", bigIntAttoToPercent(constants.MAX_LTV))
	// 	fmt.Printf("Your Agent must pay down its debt or increase its collateral to avoid liquidation\n")
	// 	fmt.Printf("Contact the GLIF team as soon as possible\n")
	// }

	if badPmtStatus {
		fmt.Println("You are late on your weekly payment")
		fmt.Printf("Your account *must* make a payment to-current within the next: %s (by epoch # %s)\n", formatSinceDuration(defaultEpochTime, epochsPaidTime), defaultEpoch)
	}

	// since we have to report faulty sectors when the Agent is overLTV, we only display this message if the Agent is not overLTV AND has faulty sectors
	// if badFaultStatus && !overLTV {
	// 	chainHeight, err := query.ChainHeight(ctx)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	consecutiveFaultEpochTolerance, err := query.MaxConsecutiveFaultEpochs(ctx)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	consecutiveFaultEpochs := new(big.Int).Sub(chainHeight, faultySectorStart)

	// 	liableForFaultySectorDefault := consecutiveFaultEpochs.Cmp(consecutiveFaultEpochTolerance) >= 0

	// 	if liableForFaultySectorDefault {
	// 		fmt.Printf("You are at risk of liquidation due to consecutive faulty sectors - recover your sectors as soon as possible\n")
	// 		fmt.Printf("Faulty sector start epoch: %v\n", faultySectorStart)
	// 	} else {
	// 		epochsBeforeZeroTolerance := new(big.Int).Sub(consecutiveFaultEpochTolerance, consecutiveFaultEpochs)
	// 		fmt.Printf("WARNING: You are approaching risk of liquidation due to consecutive faulty sectors\n")
	// 		fmt.Printf("With %v more consecutive epochs of faulty sectors, you will be at risk of liquidation\n", epochsBeforeZeroTolerance)
	// 	}
	// } else if pendingBadFaultStatus {
	// 	fmt.Printf("WARNING: Your Agent has one or more miners with faulty sectors - recover your sectors as soon as possible\n")
	// 	fmt.Printf("Faulty sector ratio: %.02f%%\n", faultRatio)
	// 	fmt.Printf("Faulty sector ratio limit: %v%%\n", limit.String())
	// }

	printTable([]string{
		"Agent's administrator",
		"Agent in default",
	}, []string{
		agentAdmin.String(),
		fmt.Sprintf("%t", defaulted),
	})
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
