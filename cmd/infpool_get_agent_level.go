/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var agentLvlCmd = &cobra.Command{
	Use:   "get-agent-level",
	Short: "Gets the level of the Agent within the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentID, err := getAgentID(cmd)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Querying the level of AgentID %s", agentID.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		lvl, borrowCap, err := PoolsSDK.Query().InfPoolGetAgentLvl(cmd.Context(), agentID)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		fmt.Printf("Agent's lvl is %s and can borrow %.03f FIL", lvl.String(), borrowCap)
	},
}

func init() {
	infinitypoolCmd.AddCommand(agentLvlCmd)
	agentLvlCmd.Flags().String("agent-id", "", "ID of the Agent")
}
