package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade your Agent to the latest version",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, ownerKey, err := commonSetupOwnerCall()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Upgrading agent")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		tx, err := fevm.Connection().AgentUpgrade(cmd.Context(), agentAddr, ownerKey)
		if err != nil {
			log.Fatal(err)
		}

		// transaction landed on chain or errored
		receipt, err := fevm.WaitReturnReceipt(tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		if receipt == nil {
			log.Fatal("Failed to get receipt")
		}

		if receipt.Status == 0 {
			log.Fatal("Transaction failed")
		}

		// grab the ID and the address of the agent from the receipt's logs
		addr, err := fevm.Connection().UpgradedAgentAddr(cmd.Context(), receipt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Agent upgraded, new Agent Address: ", addr.String())

		s.Stop()

		as := util.AgentStore()

		as.Set("address", addr.String())
		as.Set("tx", tx.Hash().String())

		fmt.Println("Agent upgraded")
	},
}

func init() {
	agentCmd.AddCommand(upgradeCmd)
}
