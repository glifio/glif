/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/econ"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
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

		miners, baseFis, err := econ.GetBaseFisFromAPI(agentAddr, PoolsSDK.Extern().GetEventsURL())
		if err != nil {
			logFatal(err)
		}
		s.Stop()

		minersKeys := []string{
			"Miner liquidation values",
		}

		minersValues := []string{
			"",
		}

		for i, miner := range miners {
			baseFi := baseFis[i]
			minersKeys = append(minersKeys, miner.String())
			minersValues = append(minersValues, fmt.Sprintf("%0.04f FIL (%0.02f%%)", util.ToFIL(baseFi.LiquidationValue()), baseFi.RecoveryRate()*100))
		}

		agentCollateralStatsKeys := []string{
			"Agent liquidation value breakdown",
			"Liquid FIL",
			"Pledged",
			"Vesting",
			"Termination penalty",
			"Total",
		}

		afi, err := econ.GetAgentFiFromAPI(agentAddr, PoolsSDK.Extern().GetEventsURL())
		if err != nil {
			logFatal(err)
		}

		agentCollateralStatsVals := []string{
			"",
			fmt.Sprintf("%0.04f FIL", util.ToFIL(afi.AvailableBalance)),
			fmt.Sprintf("%0.04f FIL", util.ToFIL(afi.InitialPledge)),
			fmt.Sprintf("%0.04f FIL", util.ToFIL(afi.LockedRewards)),
			fmt.Sprintf("-%0.04f FIL", util.ToFIL(afi.TerminationFee)),
			fmt.Sprintf(chalk.Bold.TextStyle("%0.03f FIL (%0.02f%% recovery rate)"), util.ToFIL(afi.LiquidationValue()), afi.RecoveryRate()*100),
		}

		printTable(agentCollateralStatsKeys, agentCollateralStatsVals)
		printTable(minersKeys, minersValues)
		fmt.Println()
	},
}

func init() {
	agentCmd.AddCommand(liquidationValueCmd)
	liquidationValueCmd.Flags().String("agent-addr", "", "Agent address")
}
