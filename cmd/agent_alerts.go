package cmd

import (
	"fmt"

	"github.com/glifio/cli/journal/alerting"
	"github.com/spf13/cobra"
)

var agentAlertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Manage Autopilot alerts",
}

var agentAlertsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List currently raised alerts",
	Run: func(cmd *cobra.Command, args []string) {
		// alrts, err := alerts.GetAlertsFromFile()
		alrts := alerts.GetAlerts()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

		raised := []alerting.Alert{}
		for _, alrt := range alrts {
			if alerts.IsRaised(alrt.Type) {
				raised = append(raised, alrt)
			}
		}

		for _, r := range raised {
			fmt.Println(r)
		}
	},
}

func init() {
	agentCmd.AddCommand(agentAlertsCmd)
	agentAlertsCmd.AddCommand(agentAlertsListCmd)
}
