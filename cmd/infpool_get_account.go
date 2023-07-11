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
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Gets the details associated with an active account borrowing from the Infinity Pool",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Querying the Account of agent %s", agentAddr.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		account, err := PoolsSDK.Query().InfPoolGetAccount(cmd.Context(), agentAddr, nil)
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		chainHeadHeight, err := PoolsSDK.Query().ChainHeight(cmd.Context())
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		filPrincipal := util.ToFIL(account.Principal)

		log.Printf("Account opened at epoch # %s", account.StartEpoch.String())
		log.Printf("Outstanding principal: %0.09f", filPrincipal)
		log.Printf("Account owes %s epoch payments", new(big.Int).Sub(new(big.Int).SetUint64(chainHeadHeight.Uint64()), account.EpochsPaid))
		log.Printf("Account is paid up to epoch # %s", account.EpochsPaid.String())
		log.Printf("Account in default? %v", account.Defaulted)

	},
}

func init() {
	infinitypoolCmd.AddCommand(getAccountCmd)
	getAccountCmd.Flags().String("agent-addr", "", "Address of the Agent")
}
