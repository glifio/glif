/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/glif/v2/util"
	"github.com/spf13/cobra"
)

var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Fetches the Agent ID (uses the address in agent.toml by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			logFatal(err)
		}

		log.Printf("Fetching agent ID for %s", util.TruncateAddr(agentAddr.String()))
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		id, err := PoolsSDK.Query().AgentID(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		log.Printf("Agent %s ID: %s\n", util.TruncateAddr(agentAddr.String()), id)
	},
}

func init() {
	agentCmd.AddCommand(idCmd)
	idCmd.Flags().String("agent-addr", "", "Agent address")
}
