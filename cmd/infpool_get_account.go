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

var getAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Gets the details associated with an active account borrowing from the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var agentIDStr string

		if cmd.Flag("agent-id") != nil && cmd.Flag("agent-id").Changed {
			agentIDStr = cmd.Flag("agent-id").Value.String()
		} else {
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

		fmt.Printf("Querying the Account of AgentID %s", agentID.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		conn := fevm.Connection()

		account, err := conn.PoolGetAccount(cmd.Context(), conn.InfinityPoolAddr, agentID)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		lapi, closer, err := conn.ConnectLotusClient()
		if err != nil {
			log.Fatal(err)
		}
		defer closer()

		chainHead, err := lapi.ChainHead(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		filPrincipal := util.ToFIL(account.Principal)

		log.Printf("Account opened at epoch # %s", account.StartEpoch.String())
		log.Printf("Outstanding principal: %s", filPrincipal.String())
		log.Printf("Account owes %s epoch payments", new(big.Int).Sub(new(big.Int).SetUint64(uint64(chainHead.Height())), account.EpochsPaid))
		log.Printf("Account is paid up to epoch # %s", account.EpochsPaid.String())
		log.Printf("Account in default? %v", account.Defaulted)

	},
}

func init() {
	infinitypoolCmd.AddCommand(getAccountCmd)
	getAccountCmd.Flags().String("account-id", "", "ID of the Agent")
}
