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

var agentInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get the info associated with your Agent",
	Run: func(cmd *cobra.Command, args []string) {
		agentAddr, err := getAgentAddress(cmd)
		if err != nil {
			log.Fatal(err)
		}

		agentID, err := getAgentID(cmd)
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		log.Printf("Fetching stats for %s", agentAddr.String())

		conn := fevm.Connection()

		assets, err := conn.AgentAssets(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		assetsFIL, _ := util.ToFIL(assets).Float64()

		s.Stop()

		generateHeader("AGENT ASSETS")
		fmt.Printf("%f FIL\n", assetsFIL)

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
		generateHeader("INFINITY POOL ACCOUNT")
		fmt.Printf("Account opened at epoch # %s\n", account.StartEpoch.String())
		fmt.Printf("Outstanding principal: %s\n", filPrincipal.String())
		fmt.Printf("Account owes %s epoch payments", new(big.Int).Sub(new(big.Int).SetUint64(uint64(chainHead.Height())), account.EpochsPaid))
		fmt.Printf("Account is paid up to epoch # %s\n", account.EpochsPaid.String())
		fmt.Printf("Account in default? %v\n", account.Defaulted)
	},
}

const headerWidth = 60

func generateHeader(title string) {
	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", title)
}

// var agentInfoCmd = &cobra.Command{
// 	Use:   "stats",
// 	Short: "Get the stats associated with your Agent",
// 	Run: func(cmd *cobra.Command, args []string) {

// 		defaultBlock := 1000
// 		currentBlock := 500000
// 		paidBlock := 20000

// 		lineLength := 50
// 		percentagePaid := float64(paidBlock-defaultBlock) / float64(currentBlock-defaultBlock)

// 		paidPosition := int(float64(lineLength) * percentagePaid)
// 		line := ""
// 		labelsTop := ""
// 		labelsBottom := ""

// 		for i := 0; i < lineLength; i++ {
// 			if i == 0 {
// 				line += "─"
// 				labelsTop += "default"
// 				labelsBottom += strconv.Itoa(defaultBlock)
// 			} else if i == lineLength-1 {
// 				line += "─"
// 				labelsTop += " current"
// 				labelsBottom += strconv.Itoa(currentBlock)
// 			} else if i == paidPosition {
// 				line += "⦿"
// 				labelsTop += "account paid"
// 				labelsBottom += strconv.Itoa(paidBlock)
// 			} else {
// 				line += "─"
// 				labelsTop += " "
// 				labelsBottom += " "
// 			}
// 		}

// 		fmt.Println(labelsTop)
// 		fmt.Println(line)
// 		fmt.Println(labelsBottom)
// 	},
// }

func init() {
	agentCmd.AddCommand(agentInfoCmd)
	agentInfoCmd.Flags().String("address", "", "Agent address")
	agentInfoCmd.Flags().String("agent-id", "", "AgentID")
}
