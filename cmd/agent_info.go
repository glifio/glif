/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"context"
	"fmt"
	"math/big"
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
	"github.com/glifio/go-pools/util"
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

		_, _, _, _, afi, tokenID, tier, tierInfos, err := basicInfo(cmd.Context(), agentAddr, agentAddrDel, lapi, s)
		if err != nil {
			logFatal(err)
		}

		maxDTL := getDTLForTier(tier, tierInfos)

		err = econInfo(cmd.Context(), agentAddr, afi, maxDTL, s)
		if err != nil {
			logFatal(err)
		}

		err = agentHealth(cmd.Context(), agentAddr, afi, maxDTL, s)
		if err != nil {
			logFatal(err)
		}

		err = plusCardInfo(cmd.Context(), tokenID, tier, tierInfos, s)
		if err != nil {
			logFatal(err)
		}
	},
}

func getDTLForTier(tier uint8, tierInfos []abigen.TierInfo) *big.Int {
	if tier == 0 {
		return constants.MAX_BORROW_DTL
	}
	if int(tier) <= len(tierInfos) {
		tierInfo := tierInfos[tier]
		return tierInfo.DebtToLiquidationValue
	}

	return constants.MAX_BORROW_DTL
}

func basicInfo(ctx context.Context, agent common.Address, agentDel address.Address, lapi *api.FullNodeStruct, s *spinner.Spinner) (
	agentID *big.Int,
	agentFILIDAddr address.Address,
	agVersion uint8,
	ntwVersion uint8,
	afi *econ.AgentFi,
	tokenID *big.Int,
	tier uint8,
	tierInfos []abigen.TierInfo,
	err error,
) {
	query := PoolsSDK.Query()

	agentID, err = query.AgentID(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, nil, common.Big0, 0, nil, err
	}

	tasks := []util.TaskFunc{
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
			return query.MinerRegistryAgentMinersList(ctx, agentID, nil)
		},
		func() (interface{}, error) {
			return econ.GetAgentFiFromAPI(agent, PoolsSDK.Extern().GetEventsURL())
		},
		func() (interface{}, error) {
			return query.SPPlusAgentIdToTokenId(ctx, agentID, nil)
		},
		func() (interface{}, error) {
			return query.SPPlusTierFromAgentAddress(ctx, agent, nil)
		},
		func() (interface{}, error) {
			return query.SPPlusTierInfo(ctx, nil)
		},
	}
	results, err := util.Multiread(tasks)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, nil, common.Big0, 0, nil, err
	}

	agentFILIDAddr = results[0].(address.Address)
	versionResults := results[1].([]interface{})
	agVersion = versionResults[0].(uint8)
	ntwVersion = versionResults[1].(uint8)
	owner := results[2].(common.Address)
	operator := results[3].(common.Address)
	requester := results[4].(common.Address)
	agentMiners := results[5].([]address.Address)
	afi = results[6].(*econ.AgentFi)
	tokenID = results[7].(*big.Int)
	tier = results[8].(uint8)
	tierInfos = results[9].([]abigen.TierInfo)

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
		fmt.Sprintf("%v", len(agentMiners)),
		versionCopy,
	}

	generateHeader("BASIC INFO")
	printTable(basicInfoKeys, basicInfoValues)

	s.Start()

	return agentID, agentFILIDAddr, agVersion, ntwVersion, afi, tokenID, tier, tierInfos, nil
}

func econInfo(ctx context.Context, agent common.Address, afi *econ.AgentFi, maxDTL *big.Int, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	tasks := []util.TaskFunc{
		func() (interface{}, error) {
			return query.InfPoolGetRate(ctx)
		},
		func() (interface{}, error) {
			return query.AgentLiquidAssets(ctx, agent, nil)
		},
	}

	results, err := util.Multiread(tasks)
	if err != nil {
		return err
	}
	rate := results[0].(*big.Int)
	agentLiquidAssets := results[1].(*big.Int)

	apr := new(big.Float).Mul(new(big.Float).SetInt(rate), big.NewFloat(constants.EpochsInYear))
	apr.Quo(apr, big.NewFloat(1e34))

	s.Stop()

	generateHeader("ECON INFO")

	printTable([]string{
		"Liquidation value",
		"Total debt",
		"Debt-to-liquidation ratio (DTL)",
		"Max DTL",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LiquidationValue())),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.Debt())),
		fmt.Sprintf("%0.02f%%", new(big.Float).Mul(afi.DTL(), big.NewFloat(100))),
		fmt.Sprintf("%0.02f%%", new(big.Float).Mul(big.NewFloat(100), util.ToFIL(maxDTL))),
	})

	printTable([]string{
		"Max borrow to seal",
		"Max borrow to withdraw",
		"Available to withdraw",
	}, []string{
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.BorrowLimit(maxDTL))),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.MaxBorrowAndWithdraw(maxDTL))),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.WithdrawLimit(maxDTL))),
	})

	printTable([]string{
		"Assets Breakdown",
		"Total assets",
		"Liquid assets",
		"Agent balance",
	}, []string{
		"",
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.Balance)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.AvailableBalance)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(agentLiquidAssets)),
	})

	printTable([]string{
		"Liabilities Breakdown",
		"Total borrowed",
		"Interest owed",
		"APR",
	}, []string{
		"",
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.Principal)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.Interest)),
		fmt.Sprintf("%.02f%%", apr),
	})

	printTable([]string{
		"Liquidation Value Breakdown",
		"Available balance",
		"Initial pledge",
		"Locked rewards",
		"Termination fee",
		"Total liquidation value",
	}, []string{
		"",
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.AvailableBalance)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.InitialPledge)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LockedRewards)),
		fmt.Sprintf("-%0.09f FIL", util.ToFIL(afi.TerminationFee)),
		fmt.Sprintf("%0.09f FIL", util.ToFIL(afi.LiquidationValue())),
	})

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

func agentHealth(ctx context.Context, agent common.Address, afi *econ.AgentFi, maxDTL *big.Int, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	tasks := []util.TaskFunc{
		func() (interface{}, error) {
			return query.AgentAdministrator(ctx, agent)
		},

		func() (interface{}, error) {
			return query.AgentDefaulted(ctx, agent)
		},
	}

	results, err := util.Multiread(tasks)
	if err != nil {
		return err
	}

	agentAdmin := results[0].(common.Address)
	defaulted := results[1].(bool)

	s.Stop()

	generateHeader("HEALTH")
	fmt.Println()

	if defaulted {
		fmt.Println("Agent is in default")
		return nil
	}

	if agentAdmin != common.HexToAddress("") {
		fmt.Println("Agent is on administration")
	}

	if afi.LiquidationValue().Sign() == 0 {
		fmt.Println("Agent is inactive")
		return nil
	}

	overLTV := util.DivWad(afi.Debt(), afi.LiquidationValue()).Cmp(maxDTL) > 0

	if overLTV {
		fmt.Println(chalk.Bold.TextStyle("Status unhealthy 🔴"))
		fmt.Println("Agent is over the debt to liquidation value borrowing limit. Pay down principal or increase your collateral to avoid liquidation.")
		fmt.Println()
	} else {
		fmt.Printf("Status healthy 🟢\n")
	}

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

func plusCardInfo(ctx context.Context, tokenID *big.Int, tier uint8, tierInfos []abigen.TierInfo, s *spinner.Spinner) error {
	s.Stop()

	generateHeader("GLIF+ CARD")

	if tokenID.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("No GLIF+ Card minted for Agent")
		return nil
	}
	if tier == 0 {
		fmt.Println("Agent's GLIF+ Card is inactive")
		return nil
	}

	// Basic card info
	printTable([]string{
		"Card ID",
		"Tier",
	}, []string{
		tokenID.String(),
		tierName(tier),
	})

	// Display tier-specific information
	if int(tier) <= len(tierInfos) && tier > 0 {
		tierInfo := tierInfos[tier]

		info, err := PoolsSDK.Query().SPPlusInfo(ctx, tokenID, nil)
		if err != nil {
			logFatal(err)
		}

		conversionRateWithPremium := util.MulWad(info.BaseConversionRateFILtoGLF, tierInfo.CashBackPremium)

		premium := new(big.Float).Mul(
			new(big.Float).Sub(util.ToFIL(tierInfo.CashBackPremium), big.NewFloat(1)),
			big.NewFloat(100),
		)

		printTable([]string{
			"Tier Benefits",
			"Max Debt-to-Liquidation Ratio",
			"Cash Back Exchange Rate",
		}, []string{
			"",
			fmt.Sprintf("%.2f%%", new(big.Float).Mul(big.NewFloat(100), util.ToFIL(tierInfo.DebtToLiquidationValue))),
			fmt.Sprintf("1 FIL = %.09f GLF (+%.02f%%)", util.ToFIL(conversionRateWithPremium), premium),
		})
	}

	s.Start()
	return nil
}

func generateHeader(title string) {
	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", chalk.Underline.TextStyle(title))
}

func init() {
	agentCmd.AddCommand(agentInfoCmd)
	agentInfoCmd.Flags().String("agent-addr", "", "Agent address")
}
