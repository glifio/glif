/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Fetches the Agent ID (uses the address in agent.toml by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		as := util.AgentStore()

		address := cmd.Flag("address").Value.String()

		fmt.Println(address)

		if address == "" {
			// Check if an agent already exists
			cachedAddr, err := as.Get("address")
			if err != nil {
				log.Fatal(err)
			}

			address = cachedAddr

			if address == "" {
				log.Fatalf("Did you forget to create your agent or specify an address? Try `glif agent id --address <address>`")
			}

		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		addr := common.HexToAddress(address)

		log.Printf("Fetching agent ID for %s", util.TruncateAddr(addr.String()))

		id, err := fevm.Connection().AgentID(cmd.Context(), addr)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		log.Printf("Agent %s ID: %s", util.TruncateAddr(addr.String()), id)
	},
}

func init() {
	agentCmd.AddCommand(idCmd)

	idCmd.Flags().String("address", "", "Agent address")
}
