package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View actions that the agent has taken",
	Run: func(cmd *cobra.Command, args []string) {
		// read journal json file
		// parse njson into chronologically ordered slice
		// make sure that the top level data structure makes filtering by event type is simple for future --filter feature
		evts, err := journal.ReadEvents()
		if err != nil {
			logFatal(err)
		}

		for _, e := range evts {
			fmt.Println(e)
		}
	},
}

func init() {
	agentCmd.AddCommand(historyCmd)
}
