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

		collateralStatsKeys := []string{
			"Liquidation value",
		}

		lv, ok := new(big.Int).SetString(agentCollateralStats.LiquidationValue, 10)
		if !ok {
			logFatal(err)
		}

		collateralStatsVals := []string{
			fmt.Sprintf("%0.09f FIL", util.ToFIL(lv)),
		}

		s.Stop()

		printTable(collateralStatsKeys, collateralStatsVals)
		fmt.Println()
	},
}

func init() {
	agentCmd.AddCommand(liquidationValueCmd)
	liquidationValueCmd.Flags().String("agent-addr", "", "Agent address")
}
