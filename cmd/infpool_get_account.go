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
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Gets the details associated with an active account borrowing from the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Querying the Account of agent %s", agentAddr.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		account, err := PoolsSDK.Query().InfPoolGetAccount(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		chainHeadHeight, err := PoolsSDK.Query().ChainHeight(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		filPrincipal := util.ToFIL(account.Principal)

		log.Printf("Account opened at epoch # %s", account.StartEpoch.String())
		log.Printf("Outstanding principal: %s", filPrincipal.String())
		log.Printf("Account owes %s epoch payments", new(big.Int).Sub(new(big.Int).SetUint64(chainHeadHeight.Uint64()), account.EpochsPaid))
		log.Printf("Account is paid up to epoch # %s", account.EpochsPaid.String())
		log.Printf("Account in default? %v", account.Defaulted)

	},
}

func init() {
	infinitypoolCmd.AddCommand(getAccountCmd)
	getAccountCmd.Flags().String("agent-addr", "", "Address of the Agent")
}
