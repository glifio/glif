package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/filecoin-project/lotus/chain/actors/builtin"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/token"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Airdrop plans related commands",
}

var getPlanCmd = &cobra.Command{
	Use:   "get <plan-id>",
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

		caller, err := abigen.NewIHedgeyVoteTokenLockupPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
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

		fmt.Printf("available to claim: %0.04f GLF\n", util.ToFIL(balance.Balance))

		printVestingSchedule(planIDBig, &plan, false)
	},
}

var listPlansCmd = &cobra.Command{
	Use:   "list <address>",
	Short: "List all vesting plans for a given address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ethAddr, err := AddressOrAccountNameToEVM(cmd.Context(), args[0])
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		fmt.Printf("Getting vesting plans for %s...\n", ethAddr.Hex())

		agentOwnerMap, err := token.ReadAgentOwnerMap(false)
		if err != nil {
			logFatal(err)
		}

		addr, ok := agentOwnerMap[ethAddr]
		if ok {
			fmt.Printf("This address is an agent, its claimer is its owner: %s\n", addr.Hex())
			ethAddr = addr
		}

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		caller, err := abigen.NewIHedgeyVoteTokenLockupPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}

		balance, err := caller.BalanceOf(&bind.CallOpts{Context: cmd.Context()}, ethAddr)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("Found %d vesting plans for %s\n", balance, ethAddr.Hex())

		for i := big.NewInt(0); i.Cmp(balance) == -1; i.Add(i, big.NewInt(1)) {
			tokenId, err := caller.TokenOfOwnerByIndex(&bind.CallOpts{Context: cmd.Context()}, ethAddr, i)
			if err != nil {
				logFatal(err)
			}

			plan, err := caller.Plans(&bind.CallOpts{Context: cmd.Context()}, tokenId)
			if err != nil {
				logFatal(err)
			}

			printVestingSchedule(tokenId, &plan, true)
		}
	},
}

var redeemPlanCmd = &cobra.Command{
	Use:   "redeem <plan-id>",
	Short: "Redeem tokens from an airdrop plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		// generic account setup
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		fmt.Println("Fetching the amount of GLF tokens available to redeem...")

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

		votingTokenLockupPlanCaller, err := abigen.NewIHedgeyVoteTokenLockupPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}
		unixNow := big.NewInt(time.Now().Unix())

		balance, err := votingTokenLockupPlanCaller.PlanBalanceOf(&bind.CallOpts{Context: cmd.Context()}, planIDBig, unixNow, unixNow)
		if err != nil {
			logFatal(err)
		}

		s.Stop()

		if balance.Balance.Cmp(big.NewInt(0)) == 0 {
			logFatalf("No tokens available to redeem")
		}

		fmt.Printf("Available to redeem: %0.06f GLF\n", util.ToFIL(balance.Balance))

		s.Start()

		votingTokenLockupPlanTxor, err := abigen.NewIHedgeyVoteTokenLockupPlanTransactor(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatal(err)
		}

		tx, err := votingTokenLockupPlanTxor.RedeemPlans(auth, []*big.Int{planIDBig})
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

		fmt.Printf("%0.06f GLF tokens redeemed successfully\n", util.ToFIL(balance.Balance))
	},
}

var getDelegateCmd = &cobra.Command{
	Use:   "get-delegate <plan-id>",
	Short: "Get the delegate address for GLF tokens locked in an airdrop plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		planIDBig, ok := big.NewInt(0).SetString(planID, 10)
		if !ok {
			logFatalf("Failed to parse plan ID %s", planID)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		client, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to connect to Ethereum client: %s", err)
		}
		defer client.Close()

		// Connect to the Hedgey NFT contract
		votingTokenLockupPlanCaller, err := abigen.NewIHedgeyVoteTokenLockupPlanCaller(PoolsSDK.Query().TokenNFTWrapper(), client)
		if err != nil {
			logFatalf("Failed to instantiate Hedgey NFT contract: %s", err)
		}

		// Get the vault address
		vaultAddr, err := votingTokenLockupPlanCaller.VotingVaults(&bind.CallOpts{Context: cmd.Context()}, planIDBig)
		if err != nil {
			logFatalf("Failed to get vault address: %s", err)
		}

		// Connect to the token ERC20 contract
		tokenContract, err := abigen.NewTokenCaller(PoolsSDK.Query().GLF(), client)
		if err != nil {
			logFatalf("Failed to instantiate token ERC20 contract: %s", err)
		}

		// Get the delegatee address
		delegatee, err := tokenContract.Delegates(&bind.CallOpts{Context: cmd.Context()}, vaultAddr)
		if err != nil {
			logFatalf("Failed to get delegatee address: %s", err)
		}

		s.Stop()

		fmt.Println(delegatee.Hex())
	},
}

var delegateCmd = &cobra.Command{
	Use:   "set-delegate <plan-id> <delegatee>",
	Short: "Delegate GLF tokens locked in an airdrop plan to a delegatee address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		delegatee := args[1]

		planIDBig, ok := big.NewInt(0).SetString(planID, 10)
		if !ok {
			logFatalf("Failed to parse plan ID %s", planID)
		}

		delegateeAddr, err := AddressOrAccountNameToEVM(cmd.Context(), delegatee)
		if err != nil {
			logFatalf("Failed to parse delegatee address %s", err)
		}

		// generic account setup
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd, from)
		if err != nil {
			logFatal(err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to connect to Ethereum ethClient: %s", err)
		}
		defer ethClient.Close()

		// Connect to the Hedgey NFT contract
		votingTokenLockupPlanTxor, err := abigen.NewIHedgeyVoteTokenLockupPlanTransactor(PoolsSDK.Query().TokenNFTWrapper(), ethClient)
		if err != nil {
			logFatalf("Failed to instantiate Hedgey NFT contract: %s", err)
		}

		tx, err := votingTokenLockupPlanTxor.Delegate(auth, planIDBig, delegateeAddr)
		if err != nil {
			logFatalf("Failed to delegate tokens: %s", err)
		}

		s.Stop()

		fmt.Printf("Confirming delegate transaction: %s...\n", tx.Hash().Hex())

		s.Start()

		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatalf("Failed to wait for transaction receipt: %s", err)
		}

		s.Stop()

		fmt.Printf("Airdrop plan %s delegated to %s\n", planID, delegateeAddr.Hex())
	},
}

func printVestingSchedule(tokenID *big.Int, plan *abigen.IHedgeyVoteTokenLockupPlanPlan, newLine bool) {
	amountFIL := util.ToFIL(plan.Amount)
	if newLine {
		fmt.Println("")
	}
	fmt.Printf("plan ID: %s\n", tokenID.String())
	fmt.Printf("amount: %0.04f GLF\n", amountFIL)

	vestPerDay := big.NewInt(0).Mul(plan.Rate, big.NewInt(builtin.EpochsInDay))
	fmt.Printf("vesting rate: %0.06f GLF per day\n", util.ToFIL(vestPerDay))

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
	plansCmd.AddCommand(getDelegateCmd)
	plansCmd.AddCommand(delegateCmd)
	redeemPlanCmd.Flags().String("from", "", "address to redeem the tokens from")
	delegateCmd.Flags().String("from", "", "address to delegate the tokens from")
}
