/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/terminate"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

func fetchAgentCollateralStats(agentID string) (*terminate.AgentCollateralStats, error) {
	url := fmt.Sprintf("https://events.glif.link/agent/%s/collateral-value", agentID)
	// Making an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return &terminate.AgentCollateralStats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &terminate.AgentCollateralStats{}, fmt.Errorf("Error fetching collateral stats. Status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &terminate.AgentCollateralStats{}, err
	}

	var response terminate.AgentCollateralStats
	if err := json.Unmarshal(body, &response); err != nil {
		return &terminate.AgentCollateralStats{}, err
	}

	return &response, nil
}

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

		agentID, err := PoolsSDK.Query().AgentID(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		agentCollateralStats, err := fetchAgentCollateralStats(agentID.String())
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

		totalMinerBal := big.NewInt(0)

		for _, minerCollateral := range agentCollateralStats.MinersTerminationStats {
			minersKeys = append(minersKeys, fmt.Sprintf("%s", minerCollateral.MinerAddr))
			minerBal, ok := new(big.Int).SetString(minerCollateral.MinerBal, 10)
			if !ok {
				logFatal(err)
			}

			totalMinerBal = totalMinerBal.Add(totalMinerBal, minerBal)

			minerTerminationPenalty, ok := new(big.Int).SetString(minerCollateral.TerminationPenalty, 10)
			if !ok {
				logFatal(err)
			}

			minerLiquidationValue := new(big.Int).Sub(minerBal, minerTerminationPenalty)
			// format recovery rate as a percentage
			recoveryRate := new(big.Int).Div(new(big.Int).Mul(minerLiquidationValue, big.NewInt(1e18)), minerBal)
			recoveryRate.Mul(recoveryRate, big.NewInt(100))

			minersValues = append(minersValues, fmt.Sprintf("%0.09f FIL (%0.02f%%)", util.ToFIL(minerLiquidationValue), util.ToFIL(recoveryRate)))
		}

		liquidFIL, ok := new(big.Int).SetString(agentCollateralStats.AgentLiquidCollateral, 10)
		if !ok {
			logFatal(err)
		}

		agentCollateralStatsKeys := []string{
			"Agent liquidation value",
		}

		lv, ok := new(big.Int).SetString(agentCollateralStats.LiquidationValue, 10)
		if !ok {
			logFatal(err)
		}

		recoveryRate := new(big.Int).Div(new(big.Int).Mul(lv, big.NewInt(1e18)), totalMinerBal)
		recoveryRate.Mul(recoveryRate, big.NewInt(100))

		agentCollateralStatsVals := []string{
			fmt.Sprintf("%0.03f FIL (%0.02f%% recovery)", util.ToFIL(lv), util.ToFIL(recoveryRate)),
		}

		agentLiquidFILKey := []string{
			"Agent's liquid FIL",
		}

		agentLiquidFILValue := []string{
			fmt.Sprintf("%0.04f FIL", util.ToFIL(liquidFIL)),
		}

		printTable(agentCollateralStatsKeys, agentCollateralStatsVals)
		printTable(agentLiquidFILKey, agentLiquidFILValue)
		printTable(minersKeys, minersValues)
		fmt.Println()
	},
}

func init() {
	agentCmd.AddCommand(liquidationValueCmd)
	liquidationValueCmd.Flags().String("agent-addr", "", "Agent address")
}
