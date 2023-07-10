package cmd

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentAutopilotCmd = &cobra.Command{
	Use:   "autopilot",
	Short: "Background service that automatically repays FIL to pools",
	Long:  `Background service that automatically repays FIL to pools.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			}
		}()

		defer journal.Close()

		if cmd.Flag("logfile") != nil && cmd.Flag("logfile").Changed {
			file, err := os.OpenFile(cmd.Flag("logfile").Value.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			log.SetOutput(file)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		ctx := cmd.Context()

		log.Println("Starting autopilot...")

		log.Println(viper.GetString("daemon.rpc-url"))
		for {
			select {
			case <-sigs:
				log.Println("Shutting down...")
				Exit(0)
			default:
				log.Println("Checking for payments...")
				// each loop retrieve config values aka hot-reload
				paymentType, err := ParsePaymentType(viper.GetString("autopilot.payment-type"))
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println("Payment type: ", paymentType)
				payargs := []string{}
				if paymentType == Principal || paymentType == Custom {
					amount := viper.GetInt64("autopilot.amount")
					payargs = append(payargs, fmt.Sprintf("%d", amount))
				}

				//TODO: maybe change frequency to max debt or max epoch difference
				frequency := viper.GetFloat64("autopilot.frequency")

				log.Println("frequency: ", frequency)

				// goto can't jump over variable declarations
				var epochFreq *big.Float
				var dueEpoch *big.Int
				var epochFreqInt64 int64
				var epochFreqInt *big.Int
				var chainHeadHeight *big.Int
				var account abigen.Account

				agent, err := getAgentAddress(cmd)
				if err != nil {
					log.Println(err)
					goto SLEEP
				}

				account, err = PoolsSDK.Query().InfPoolGetAccount(ctx, agent, nil)
				if err != nil {
					log.Println(err)
					goto SLEEP
				}
				if account == (abigen.Account{}) {
					log.Println("failed to get infinity pool account, check evm api provider status")
					goto SLEEP
				}

				chainHeadHeight, err = PoolsSDK.Query().ChainHeight(cmd.Context())
				if err != nil {
					log.Println(err)
					goto SLEEP
				}
				if chainHeadHeight == nil {
					log.Println("failed to get chainheight, check lotus api provider status")
					goto SLEEP
				}

				// calculate epoch frequency
				epochFreq = big.NewFloat(float64(frequency * constants.EpochsInDay))

				dueEpoch = new(big.Int).Sub(new(big.Int).SetUint64(chainHeadHeight.Uint64()), account.EpochsPaid)

				epochFreqInt64, _ = epochFreq.Int64()
				epochFreqInt = big.NewInt(epochFreqInt64)

				// check if payment is due
				// if so, make payment
				if dueEpoch.Cmp(epochFreqInt) >= 0 {
					switch paymentType {
					case Principal:
						_, err := pay(cmd, payargs, paymentType, true)
						if err != nil {
							log.Println(err)
						}

					case ToCurrent:
						_, err := pay(cmd, payargs, paymentType, true)
						if err != nil {
							log.Println(err)
						}

					case Custom:
						_, err := pay(cmd, payargs, paymentType, true)
						if err != nil {
							log.Println(err)
						}

					default:
						log.Println("Invalid payment type")
					}
				}
			SLEEP:
				select {
				case <-time.After(30 * time.Minute):
					continue
				case <-sigs:
					log.Println("Shutting down...")
					Exit(0)
				}
			}
		}
	},
}

func init() {
	agentCmd.AddCommand(agentAutopilotCmd)
	agentAutopilotCmd.Flags().String("agent-addr", "", "Agent address")
	agentAutopilotCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
	agentAutopilotCmd.Flags().String("from", "", "address to send the transaction from")
	agentAutopilotCmd.Flags().String("logfile", "", "Logfile path, if empty autopilot logs to stderr")
}
