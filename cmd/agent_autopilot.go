package cmd

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Println("Starting autopilot...")

		for {
			select {
			case <-sigs:
				log.Println("Shutting down...")
				Exit(0)
			default:
				log.Println("Checking for payments...")
				// each loop retrieve config values aka hot-reload
				paymentType := viper.GetString("autopilot.payment-type")
				log.Println("Payment type: ", paymentType)
				if paymentType == "principal" || paymentType == "custom" {
					amount := viper.GetInt64("autopilot.amount")
					args = append(args, fmt.Sprintf("%d", amount))
				}

				//TODO: maybe change frequency to max debt or max epoch difference
				frequency := viper.GetFloat64("autopilot.frequency")

				agent, err := getAgentAddress(cmd)
				if err != nil {
					log.Println(err)
				}

				account, err := query.InfPoolGetAccount(ctx, agent)
				if err != nil {
					log.Println(err)
				}

				chainHeadHeight, err := PoolsSDK.Query().ChainHeight(cmd.Context())
				if err != nil {
					logFatal(err)
				}

				// calculate epoch frequency
				// frequency is in days so multiply by 24 hours, 60 minutes, 60 seconds to get seconds
				// divide by 30 seconds to get epochs
				epochFreq := big.NewFloat(float64(frequency * 24 * 60 * 60 / 30))
				log.Println("Payment Frequency: ", epochFreq, " epochs")

				dueEpoch := new(big.Int).Sub(new(big.Int).SetUint64(chainHeadHeight.Uint64()), account.EpochsPaid)
				log.Println("Last Payment Made: ", dueEpoch, " epochs ago")

				epochFreqInt64, _ := epochFreq.Int64()
				epochFreqInt := big.NewInt(epochFreqInt64)
				log.Println("Should a Payment be made: ", dueEpoch.Cmp(epochFreqInt))

				// check if payment is due
				// if so, make payment
				if dueEpoch.Cmp(epochFreqInt) >= 0 {
					switch paymentType {
					case "principle":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							log.Println(err)
						}

					case "to-current":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							log.Println(err)
						}

					case "custom":
						_, err := pay(cmd, args, paymentType)
						if err != nil {
							log.Println(err)
						}

					default:
						log.Println("Invalid payment type")
					}
				}
				log.Println("reseting args...")

				// reset args in case of hot-reload change to the amount config value
				args = []string{}

				time.Sleep(10 * time.Second)
			}
		}
	},
}

func init() {
	agentCmd.AddCommand(agentAutopilotCmd)
	agentAutopilotCmd.Flags().String("agent-addr", "", "Agent address")
	agentAutopilotCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	agentAutopilotCmd.Flags().String("from", "", "address to send the transaction from")
}
