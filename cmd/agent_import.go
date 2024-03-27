/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/cli/util"
	"github.com/spf13/cobra"
)

var agentImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a Glif agent <agent-addr>",
	Long:  `Imports the Agent's ID and address in the agent.toml file to remove the need for passing --agent-addr flags.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		agentAddr, err := AddressOrAccountNameToEVM(cmd.Context(), args[0])
		if err != nil {
			logFatal(err)
		}

		exists, err := PoolsSDK.Query().AgentIsValid(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		if !exists {
			logFatal("Agent not found")
		}

		agentStore := util.AgentStore()

		// Check if an agent already exists
		addressStr, err := agentStore.Get("address")
		if err != nil {
			var e *util.ErrKeyNotFound
			if !errors.As(err, &e) {
				logFatal(err)
			}
		}
		if addressStr != "" {
			logFatalf("Agent already exists: %s", addressStr)
		}

		id, err := PoolsSDK.Query().AgentID(cmd.Context(), agentAddr)
		if err != nil {
			logFatal(err)
		}

		err = agentStore.Set("id", id.String())
		if err != nil {
			logFatal(err)
		}

		err = agentStore.Set("address", agentAddr.String())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		fmt.Printf("Successfully imported agent %s (%v)\n", agentAddr.String(), id)
	},
}

func init() {
	agentCmd.AddCommand(agentImportCmd)
}
