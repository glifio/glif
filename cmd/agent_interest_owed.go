package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var agentInterestOwedCmd = &cobra.Command{
	Use:   "interest-owed",
	Short: "Get the total amount of interest owed by the agent",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Getting agent interest owed...")

		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		assets, err := PoolsSDK.Query().AgentInterestOwed(cmd.Context(), agentAddr, nil)
		if err != nil {
			logFatalf("Failed to get agent interest owed %s", err)
		}

		s.Stop()

		fmt.Printf("Agent %s owes %.04f FIL in interest\n", agentAddr, util.ToFIL(assets))
	},
}

func init() {
	agentCmd.AddCommand(agentInterestOwedCmd)
	agentInterestOwedCmd.Flags().String("agent-addr", "", "Agent address")
}
