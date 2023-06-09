package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentAutopilotCmd = &cobra.Command{
	Use:   "autopilot",
	Short: "Background service that automatically repays FIL to pools",
	Long:  `Background service that automatically repays FIL to pools.`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		ctx := cmd.Context()

		query := PoolsSDK.Query()

		for {
			select {
			case <-sigs:
				fmt.Println("Shutting down...")
				Exit(0)
			default:
				// each loop retrieve config values aka hot-reload
				paymentType := viper.GetString("autopilot.payment-type")
				if paymentType == "principle" || paymentType == "custom" {
					amount := viper.GetInt64("autopilot.amount")
					args = append(args, fmt.Sprintf("%d", amount))
				}

				//TODO: maybe change frequency to max debt or max epoch difference
				frequency := viper.GetInt("autopilot.frequency")

				agent, err := getAgentAddress(cmd)
				if err != nil {
					logFatal(err)
				}

				defaultEpoch, err := query.DefaultEpoch(ctx)
				if err != nil {
					logFatal(err)
				}

				account, err := query.InfPoolGetAccount(ctx, agent)
				if err != nil {
					logFatal(err)
				}

				defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch)
				epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid)
				dueTime := int(epochsPaidTime.Sub(defaultEpochTime).Round(time.Minute).Hours()) / 24

				// check if payment is due
				// if so, make payment
				if dueTime >= frequency {
					switch paymentType {
					case "principle":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							logFatal(err)
						}

					case "to-current":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							logFatal(err)
						}

					case "custom":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							logFatal(err)
						}

					default:
						fmt.Println("Invalid payment type")
					}
				}

				// reset args in case of hot-reload change to the amount config value
				args = []string{}

				time.Sleep(10 * time.Second)
			}
		}
	},
}

func init() {
	agentCmd.AddCommand(agentAutopilotCmd)
}
