/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var agentLvlCmd = &cobra.Command{
	Use:   "get-agent-level",
	Short: "Gets the level of the Agent within the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentIDStr := cmd.Flag("agent-id").Value.String()
		if agentIDStr == "" {
			as := util.AgentStore()
			storedAgent, err := as.Get("id")
			if err != nil {
				log.Fatal(err)
			}

			agentIDStr = storedAgent
		}

		agentID := new(big.Int)
		if _, ok := agentID.SetString(agentIDStr, 10); !ok {
			log.Fatalf("could not convert agent id %s to big.Int", agentIDStr)
		}

		fmt.Println("Querying the level of AgentID %s", agentID.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		conn := fevm.Connection()

		lvl, borrowCap, err := conn.AgentLevel(cmd.Context(), agentID)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		s.Stop()

		fmt.Printf("Agent's lvl is %s and can borrow %.03f FIL", lvl.String(), borrowCap)
	},
}

func init() {
	infinitypoolCmd.AddCommand(agentLvlCmd)
	agentLvlCmd.Flags().String("agent-id", "", "ID of the agent")
}
