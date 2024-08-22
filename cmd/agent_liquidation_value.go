/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var liquidationValueCmd = &cobra.Command{
	Use:   "liquidation-value",
	Short: "Fetches the Agent's liquidation value",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		log.Printf("Fetching liquidation value for %s", util.TruncateAddr(agentAddr.String()))
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		// agentCollateralStats, err := PoolsSDK.Query().AgentCollateralStatsQuick(cmd.Context(), agentAddr)
		// if err != nil {
		// 	logFatal(err)
		// }
		s.Stop()

		// ats := agentCollateralStats.Summarize()

		minersKeys := []string{
			"Miner liquidation values",
		}

		minersValues := []string{
			"",
		}

		// for _, minerCollateral := range agentCollateralStats.MinersTerminationStats {
		// 	minersKeys = append(minersKeys, fmt.Sprintf("%s", minerCollateral.Address))
		// 	// here we instantiate a PreviewAgentTerminationSummary type to reuse its liquidation value and recovery rate funcs
		// 	ts := terminate.PreviewAgentTerminationSummary{
		// 		TerminationPenalty: minerCollateral.TerminationPenalty,
		// 		InitialPledge:      minerCollateral.Pledged,
		// 		VestingBalance:     minerCollateral.Vesting,
		// 		MinersAvailableBal: minerCollateral.Available,
		// 		AgentAvailableBal:  big.NewInt(0),
		// 	}

		// 	minersValues = append(minersValues, fmt.Sprintf("%0.04f FIL (%0.02f%%)", util.ToFIL(ts.LiquidationValue()), bigIntAttoToPercent(ts.RecoveryRate())))
		// }

		// agentCollateralStatsKeys := []string{
		// 	"Agent liquidation value",
		// }

		// agentCollateralStatsVals := []string{
		// 	fmt.Sprintf("%0.03f FIL (%0.02f%% recovery)", util.ToFIL(ats.LiquidationValue()), bigIntAttoToPercent(ats.RecoveryRate())),
		// }

		// agentLiquidFILKey := []string{
		// 	"Agent's liquid FIL",
		// }

		// agentLiquidFILValue := []string{
		// 	fmt.Sprintf("%0.04f FIL", util.ToFIL(ats.AgentAvailableBal)),
		// }

		// printTable(agentCollateralStatsKeys, agentCollateralStatsVals)
		// printTable(agentLiquidFILKey, agentLiquidFILValue)
		printTable(minersKeys, minersValues)
		fmt.Println()
	},
}

func init() {
	agentCmd.AddCommand(liquidationValueCmd)
	liquidationValueCmd.Flags().String("agent-addr", "", "Agent address")
}
