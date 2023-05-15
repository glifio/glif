/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glifio/go-pools/util"
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

		agentAddrEthType, err := ethtypes.ParseEthAddress(agentAddr.String())
		if err != nil {
			log.Fatal(err)
		}

		agentAddrDel, err := agentAddrEthType.ToFilecoinAddress()
		if err != nil {
			log.Fatal(err)
		}

		agentID, err := getAgentID(cmd)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Fetching stats for %s", agentAddr.String())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		query := PoolsSDK.Query()

		version, err := query.AgentVersion(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		assets, err := query.AgentLiquidAssets(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		assetsFIL, _ := util.ToFIL(assets).Float64()

		lvl, cap, err := query.InfPoolGetAgentLvl(cmd.Context(), agentID)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		generateHeader("BASIC INFO")
		fmt.Printf("Agent Address: %s\n", agentAddr.String())
		fmt.Printf("Agent Address (del): %s\n", agentAddrDel.String())
		fmt.Printf("Agent ID: %s\n", agentID)
		fmt.Printf("Agent's lvl is %s and can borrow %.03f FIL\n", lvl.String(), cap)
		fmt.Printf("Agent Version: %v\n", version)

		generateHeader("AGENT ASSETS")
		fmt.Printf("%f FIL\n", assetsFIL)

		s.Start()

		account, err := query.InfPoolGetAccount(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatalf("Failed to get iFIL balance %s", err)
		}

		defaultEpoch, err := query.DefaultEpoch(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		amountOwed, gcred, err := query.AgentOwes(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		amountOwedFIL, _ := util.ToFIL(amountOwed).Float64()

		filPrincipal := util.ToFIL(account.Principal)
		generateHeader("INFINITY POOL ACCOUNT")

		principal, _ := filPrincipal.Float64()

		if principal == 0 {
			fmt.Println("No account exists with the Infinity Pool")
			return
		} else {
			defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch)
			epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid)
			fmt.Printf("You currently owe: %.08f FIL on %.02f FIL borrowed\n", amountOwedFIL, principal)
			fmt.Printf("Your current GCRED score is: %s\n", gcred)
			fmt.Printf("Your account must make a payment to-current within the next: %s (by epoch # %s)\n", formatSinceDuration(epochsPaidTime, defaultEpochTime), defaultEpoch)
			fmt.Println()

			fmt.Printf("Your account with the Infinity Pool opened at: %s\n", util.EpochHeightToTimestamp(account.StartEpoch).Format(time.RFC3339))
		}
		fmt.Println()
	},
}

func formatSinceDuration(t1 time.Time, t2 time.Time) string {
	d := t2.Sub(t1).Round(time.Minute)

	var parts []string

	weeks := int(d.Hours()) / (24 * 7)
	d -= time.Duration(weeks) * 7 * 24 * time.Hour
	if weeks > 1 {
		parts = append(parts, fmt.Sprintf("%d weeks", weeks))
	} else if weeks == 1 {
		parts = append(parts, fmt.Sprintf("%d week", weeks))
	}

	days := int(d.Hours()) / 24
	d -= time.Duration(days) * 24 * time.Hour
	if days > 1 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	} else if days == 1 {
		parts = append(parts, fmt.Sprintf("%d day", days))
	}

	h := d / time.Hour
	d -= h * time.Hour
	parts = append(parts, fmt.Sprintf("%02d hours", h))

	m := d / time.Minute
	parts = append(parts, fmt.Sprintf("and %02d minutes", m))

	return strings.Join(parts, " ")
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
	agentInfoCmd.Flags().String("agent-addr", "", "Agent address")
	agentInfoCmd.Flags().String("agent-id", "", "AgentID")
}
