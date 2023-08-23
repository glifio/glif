package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/glifio/go-pools/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentAutopilotInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print info about autopilot payment cycle",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		frequency := viper.GetFloat64("autopilot.frequency")
		// recalculate the due payment logic that exists in autopilot but then calculate
		// the difference between current point and when the payment will be due.

		agent, err := getAgentAddressWithFlags(cmd)
		if err != nil {
			log.Println(err)
		}

		account, err := PoolsSDK.Query().InfPoolGetAccount(ctx, agent, nil)
		if err != nil {
			log.Println(err)
		}

		chainHeadHeight, err := PoolsSDK.Query().ChainHeight(cmd.Context())
		if err != nil {
			log.Println(err)
		}

		// calculate epoch frequency
		epochFreq := big.NewFloat(float64(frequency * constants.EpochsInDay))

		dueEpoch := new(big.Int).Sub(new(big.Int).SetUint64(chainHeadHeight.Uint64()), account.EpochsPaid)

		epochFreqInt64, _ := epochFreq.Int64()
		epochFreqInt := big.NewInt(epochFreqInt64)

		if dueEpoch.Cmp(epochFreqInt) >= 0 {
			fmt.Println("based on the configured frequenc, a payment is due now")
		} else {
			dueIn := new(big.Int).Sub(epochFreqInt, dueEpoch)
			dueInFloat := new(big.Float).SetInt(dueIn)
			dueInTime := new(big.Float).Quo(dueInFloat, big.NewFloat(constants.EpochsInMinute))
			fmt.Printf("Next payment is due in: %0.1f mintues\n", dueInTime)
		}
	},
}

func init() {
	agentAutopilotCmd.AddCommand(agentAutopilotInfoCmd)
	agentAutopilotInfoCmd.Flags().String("agent-addr", "", "Agent address")
	agentAutopilotInfoCmd.Flags().String("pool-name", "infinity-pool", "name of the pool to make a payment")
}
