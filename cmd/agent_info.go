/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var agentInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get the info associated with your Agent",
	Run: func(cmd *cobra.Command, args []string) {
		lapi, closer, err := PoolsSDK.Extern().ConnectLotusClient()
		if err != nil {
			log.Fatal(err)
		}
		defer closer()

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

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		query := PoolsSDK.Query()

		agentID, _, _, _, agentAdmin, err := basicInfo(cmd.Context(), agentAddr, agentAddrDel, lapi, s)
		if err != nil {
			log.Fatal(err)
		}

		err = econInfo(cmd.Context(), agentAddr, agentID, lapi, s)

		lvl, cap, err := query.InfPoolGetAgentLvl(cmd.Context(), agentID)
		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		fmt.Printf("Agent's lvl is %s and can borrow %.03f FIL\n", lvl.String(), cap)

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

		s.Stop()

		amountOwedFIL, _ := util.ToFIL(amountOwed).Float64()

		filPrincipal := util.ToFIL(account.Principal)
		generateHeader("INFINITY POOL ACCOUNT")

		principal, _ := filPrincipal.Float64()

		if principal == 0 {
			fmt.Println("No account exists with the Infinity Pool")
		} else {
			defaultEpochTime := util.EpochHeightToTimestamp(defaultEpoch)
			epochsPaidTime := util.EpochHeightToTimestamp(account.EpochsPaid)
			fmt.Println("Your account with the Infinity Pool is open", defaultEpoch, account.EpochsPaid)
			fmt.Printf("You currently owe: %.08f FIL on %.02f FIL borrowed\n", amountOwedFIL, principal)
			fmt.Printf("Your current GCRED score is: %s\n", gcred)
			fmt.Printf("Your account must make a payment to-current within the next: %s (by epoch # %s)\n", formatSinceDuration(defaultEpochTime, epochsPaidTime), defaultEpoch)
			fmt.Println()

			fmt.Printf("Your account with the Infinity Pool opened at: %s\n", util.EpochHeightToTimestamp(account.StartEpoch).Format(time.RFC3339))
		}

		defaulted, err := query.AgentDefaulted(cmd.Context(), agentAddr)
		if err != nil {
			s.Stop()
			log.Fatal(err)
		}

		faultySectorStart, err := query.AgentFaultyEpochStart(cmd.Context(), agentAddr)
		if err != nil {
			log.Fatal(err)
		}

		generateHeader("HEALTH")
		fmt.Printf("Agent's administrator: %s\n", agentAdmin)
		fmt.Printf("Agent in default: %t\n\n", defaulted)
		if faultySectorStart.Cmp(big.NewInt(0)) == 0 {
			fmt.Printf("Status healthy 🟢\n")
		} else {
			chainHeight, err := query.ChainHeight(cmd.Context())
			if err != nil {
				s.Stop()
				log.Fatal(err)
			}

			consecutiveFaultEpochTolerance, err := query.MaxConsecutiveFaultEpochs(cmd.Context())
			if err != nil {
				s.Stop()
				log.Fatal(err)
			}

			consecutiveFaultEpochs := new(big.Int).Sub(chainHeight, faultySectorStart)

			liableForFaultySectorDefault := consecutiveFaultEpochs.Cmp(consecutiveFaultEpochTolerance) >= 0

			if liableForFaultySectorDefault {
				fmt.Printf("🔴 Status unhealthy - you are at risk of liquidation due to consecutive faulty sectors 🔴\n")
				fmt.Printf("Faulty sector start epoch: %v", faultySectorStart)
			} else {
				epochsBeforeZeroTolerance := new(big.Int).Sub(consecutiveFaultEpochTolerance, consecutiveFaultEpochs)
				fmt.Printf("🟡 Status unhealthy - you are approaching risk of liquidation due to consecutive faulty sectors 🟡\n")
				fmt.Printf("- With %v more consecutive faulty sectors, you will be at risk of liquidation\n", epochsBeforeZeroTolerance)
			}
		}
		fmt.Println()
	},
}

func basicInfo(ctx context.Context, agent common.Address, agentDel address.Address, lapi *api.FullNodeStruct, s *spinner.Spinner) (
	agentID *big.Int,
	agentFILIDAddr address.Address,
	agVersion uint8,
	ntwVersion uint8,
	agentAdmin common.Address,
	err error,
) {
	query := PoolsSDK.Query()

	agentID, err = query.AgentID(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, common.Address{}, err
	}

	agentFILIDAddr, err = lapi.StateLookupID(ctx, agentDel, types.EmptyTSK)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, common.Address{}, err
	}

	agVersion, ntwVersion, err = query.AgentVersion(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, common.Address{}, err
	}

	agentAdmin, err = query.AgentAdministrator(ctx, agent)
	if err != nil {
		return common.Big0, address.Undef, 0, 0, common.Address{}, err
	}

	goodVersion := agVersion == ntwVersion

	s.Stop()
	generateHeader("BASIC INFO")
	fmt.Printf("Agent Address: %s\n", agent.String())
	fmt.Printf("Agent Address (del): %s\n", agentDel.String())
	fmt.Printf("Agent FIL ID Address: %s\n", agentFILIDAddr.String())
	fmt.Printf("Agent Pools Protocol ID: %s\n", agentID)
	if goodVersion {
		fmt.Printf("Agent Version: %v ✅ \n", agVersion)
	} else {
		fmt.Println("Agent requires upgrade, run `glif agent upgrade` to upgrade")
		fmt.Printf("Agent/Network version mismatch: %v/%v ❌ \n", agVersion, ntwVersion)
	}
	s.Start()

	return agentID, agentFILIDAddr, agVersion, ntwVersion, agentAdmin, nil
}

func econInfo(ctx context.Context, agent common.Address, agentID *big.Int, lapi *api.FullNodeStruct, s *spinner.Spinner) error {
	query := PoolsSDK.Query()

	assets, err := query.AgentLiquidAssets(ctx, agent)
	if err != nil {
		return err
	}

	assetsFIL, _ := util.ToFIL(assets).Float64()

	agentMiners, err := query.MinerRegistryAgentMinersList(ctx, agentID)
	if err != nil {
		return err
	}

	tasks := make([]util.TaskFunc, len(agentMiners))

	for i, minerAddr := range agentMiners {
		tasks[i] = func() (interface{}, error) {
			state, err := lapi.StateReadState(ctx, minerAddr, types.EmptyTSK)
			if err != nil {
				return nil, err
			}
			bal, ok := new(big.Int).SetString(state.Balance.String(), 10)
			if !ok {
				return nil, fmt.Errorf("failed to convert balance to big.Int")
			}

			return bal, nil
		}
	}

	bals, err := util.Multiread(tasks)
	if err != nil {
		return err
	}

	var totalMinerCollaterals = big.NewInt(0)
	for _, bal := range bals {
		totalMinerCollaterals.Add(totalMinerCollaterals, bal.(*big.Int))
	}

	fmt.Printf("Agent's liquid assets: %0.08f FIL\n", assetsFIL)
	fmt.Printf("Agent's pledged miner count: %v\n", len(agentMiners))
	return nil
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

func init() {
	agentCmd.AddCommand(agentInfoCmd)
	agentInfoCmd.Flags().String("agent-addr", "", "Agent address")
}
