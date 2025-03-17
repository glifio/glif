package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"

	"github.com/glifio/go-pools/token"
	"github.com/spf13/cobra"
)

var airdropCmd = &cobra.Command{
	Use:   "airdrop",
	Short: "Airdrop related commands",
}

var checkEligibilityCmd = &cobra.Command{
	Use:   "check-eligibility [address]",
	Short: "Check airdrop eligibility for an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		fmt.Printf("Checking airdrop eligibility for %s...\n", strAddr)

		addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		isTestDrop := PoolsSDK.Query().ChainID().Int64() != constants.MainnetChainID
		amount, claimer, err := token.CheckAirdropEligibility(addr, isTestDrop)
		if err != nil {
			logFatalf("Failed to check airdrop eligibility %s", err)
		}

		// if the claimer is not the same as the address, then this is an agent, its claimer is its owner
		isAgent := claimer != addr
		if isAgent {
			fmt.Println("This address is an agent, its claimer is its owner")
		}
		fmt.Printf("Amount: %0.02f GLF, can be claimed by: %s\n", amount, claimer.Hex())
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim [address]",
	Short: "Claim airdrop for an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		addressToClaimOnBehalf := addr

		// first check if this is an agent address
		isTestDrop := PoolsSDK.Query().ChainID().Int64() != constants.MainnetChainID
		agentToOwnerMap, err := token.ReadAgentOwnerMap(isTestDrop)
		if err != nil {
			logFatalf("Failed to read agent owner map %s", err)
		}

		owner, ok := agentToOwnerMap[addr]
		if ok {
			addressToClaimOnBehalf = owner
			fmt.Println("This is an Agent address - you are claiming with your Agent's Owner wallet: ", addressToClaimOnBehalf.Hex())
		}

		// generic account setup
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd.Context(), from)
		if err != nil {
			logFatal(err)
		}

		if addressToClaimOnBehalf != auth.From {
			logFatal("Invalid 'from' address - you cannot claim on behalf of someone else")
		}

		fmt.Printf("Claiming airdrop for %s from %s...\n", strAddr, auth.From.Hex())

		delegatee, err := interactiveClaimExp(cmd.Context(), auth.From)
		if err != nil {
			logFatal(err)
		}

		fmt.Printf("You have selected to delegate your vote to: %s\n", delegatee.Hex())

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		airdropInstance, err := abigen.NewIHedgeyAirdropTransactor(PoolsSDK.Query().DelegatedClaimsCampaigns(), ethClient)
		if err != nil {
			logFatalf("Failed to create IHedgeyAirdrop instance %s", err)
		}

		mt := &token.MerkleTree{}
		mt, err = mt.ReadFromJSON(isTestDrop)
		if err != nil {
			logFatal(err)
		}

		proof, err := mt.GetProofForAddr(addressToClaimOnBehalf)
		if err != nil {
			logFatal(err)
		}

		value, err := mt.GetLeafValueForAddr(addressToClaimOnBehalf)
		if err != nil {
			logFatal(err)
		}

		tx, err := airdropInstance.ClaimAndDelegate(auth, mt.ID(), proof, value, delegatee, abigen.IHedgeyAirdropSignatureParams{
			V:      0,
			R:      [32]byte{},
			S:      [32]byte{},
			Nonce:  big.NewInt(0),
			Expiry: big.NewInt(0),
		})
		if err != nil {
			logFatalf("Failed to claim airdrop %s", err)
		}

		s.Stop()

		fmt.Printf("Claim transaction submitted: %s\n", tx.Hash().Hex())

		s.Start()
		_, err = PoolsSDK.Query().StateWaitReceipt(cmd.Context(), tx.Hash())
		if err != nil {
			logFatalf("Failed to claim airdrop %s", err)
		}
		s.Stop()

		fmt.Printf("Airdrop claimed successfully.\n")
	},
}

func init() {
	rootCmd.AddCommand(airdropCmd)
	airdropCmd.AddCommand(checkEligibilityCmd)
	airdropCmd.AddCommand(claimCmd)
	claimCmd.Flags().String("from", "", "address of the owner or operator of the agent")
}
