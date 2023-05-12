/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Fetches the Agent ID (uses the address in agent.toml by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		log.Printf("Fetching agent ID for %s", util.TruncateAddr(agentAddr.String()))

		id, err := fevm.Connection().AgentID(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		log.Printf("Agent %s ID: %s", util.TruncateAddr(agentAddr.String()), id)
	},
}

func init() {
	agentCmd.AddCommand(idCmd)
	idCmd.Flags().String("address", "", "Agent address")
}
