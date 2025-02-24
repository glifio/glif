package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/filecoin-project/lotus/chain/actors/builtin"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Airdrop plans related commands",
}

var getPlanCmd = &cobra.Command{
	Use:   "get [plan-id]",
	Short: "Get the airdrop plan for an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		fmt.Printf("Getting airdrop plan id %s...\n", planID)

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		planIDBig, ok := big.NewInt(0).SetString(planID, 10)
		if !ok {
			logFatalf("Failed to parse plan ID %s", planID)
		}

		caller, err := abigen.NewIHedgeyVoteTokenVestingPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}

		plan, err := caller.Plans(&bind.CallOpts{Context: cmd.Context()}, planIDBig)
		if err != nil {
			logFatal(err)
		}

		unixNow := big.NewInt(time.Now().Unix())

		balance, err := caller.PlanBalanceOf(&bind.CallOpts{Context: cmd.Context()}, planIDBig, unixNow, unixNow)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("available to claim: %s\n", util.ToFIL(balance.Balance))

		printVestingSchedule(planIDBig, &plan)
	},
}

var listPlansCmd = &cobra.Command{
	Use:   "list [address]",
	Short: "List all vesting plans for a given address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ethAddr, err := AddressOrAccountNameToEVM(cmd.Context(), args[0])
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		fmt.Printf("Getting vesting plans for %s...\n", ethAddr.Hex())

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		caller, err := abigen.NewIHedgeyVoteTokenVestingPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}

		balance, err := caller.BalanceOf(&bind.CallOpts{Context: cmd.Context()}, ethAddr)
		if err != nil {
			logFatal(err)
		}

		for i := big.NewInt(0); i.Cmp(balance) == -1; i.Add(i, big.NewInt(1)) {
			tokenId, err := caller.TokenOfOwnerByIndex(&bind.CallOpts{Context: cmd.Context()}, ethAddr, i)
			if err != nil {
				logFatal(err)
			}

			plan, err := caller.Plans(&bind.CallOpts{Context: cmd.Context()}, tokenId)
			if err != nil {
				logFatal(err)
			}

			printVestingSchedule(tokenId, &plan)
		}
	},
}

var redeemPlanCmd = &cobra.Command{
	Use:   "redeem [plan-id]",
	Short: "Redeem tokens from an airdrop plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		// generic account setup
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd.Context(), from)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Redeeming GLF tokens from airdrop plan %s...", planID)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		planIDBig, ok := big.NewInt(0).SetString(planID, 10)
		if !ok {
			logFatalf("Failed to parse plan ID %s", planID)
		}

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		delegatedClaimsTx, err := abigen.NewIHedgeyVoteTokenVestingPlanTransactor(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}

		tx, err := delegatedClaimsTx.RedeemPlans(auth, []*big.Int{planIDBig})
		if err != nil {
			logFatalf("Failed to redeem airdrop %s", err)
		}

		s.Stop()

		fmt.Printf("Confirming redeem transaction: %s...\n", tx.Hash().Hex())

		s.Start()
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatalf("Failed to redeem airdrop %s", err)
		}

		s.Stop()

		fmt.Printf("GLF tokens redeemed successfully\n")
	},
}

func printVestingSchedule(tokenID *big.Int, plan *abigen.IHedgeyVoteTokenVestingPlanPlan) {
	amountFIL := util.ToFIL(plan.Amount)
	fmt.Println("")
	fmt.Printf("plan ID: %s\n", tokenID.String())
	fmt.Printf("amount: %0.04f GLF\n", amountFIL)

	vestPerDay := big.NewInt(0).Mul(plan.Rate, big.NewInt(builtin.EpochsInDay))
	fmt.Printf("vesting rate: %0.04f GLF per day\n", util.ToFIL(vestPerDay))

	periods := big.NewInt(0).Div(plan.Amount, plan.Rate)
	elapsedSecondsUntilEnd := big.NewInt(0).Mul(periods, plan.Period)

	startDate := time.Unix(plan.Start.Int64(), 0)
	fmt.Printf("vesting start date: %s\n", startDate.Format(time.RFC1123))
	vestingEnd := time.Unix(plan.Start.Int64()+elapsedSecondsUntilEnd.Int64(), 0)
	fmt.Printf("vesting end date: %s\n", vestingEnd.Format(time.RFC1123))
}

func init() {
	airdropCmd.AddCommand(plansCmd)
	plansCmd.AddCommand(getPlanCmd)
	plansCmd.AddCommand(listPlansCmd)
	plansCmd.AddCommand(redeemPlanCmd)
	redeemPlanCmd.Flags().String("from", "", "address to redeem the tokens from")
}
