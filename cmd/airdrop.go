package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/glifio/go-pools/abigen"
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

		amount, claimer, err := token.CheckAirdropEligibility(addr)
		if err != nil {
			logFatalf("Failed to check airdrop eligibility %s", err)
		}

		// if the claimer is the same as the address, then this is an agent, its claimer is its owner
		isAgent := claimer == addr

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

		// generic account setup
		from := cmd.Flag("from").Value.String()
		auth, _, err := commonGenericAccountSetup(cmd.Context(), from)
		if err != nil {
			logFatal(err)
		}
		fmt.Printf("Claiming airdrop for %s from %s...", strAddr, from)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()

		addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		ethClient, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatal(err)
		}

		airdropInstance, err := abigen.NewIHedgeyAirdropTransactor(PoolsSDK.Query().DelegatedClaimsCampaigns(), ethClient)
		if err != nil {
			logFatalf("Failed to create IHedgeyAirdrop instance %s", err)
		}

		mt := &token.MerkleTree{}
		mt, err = mt.ReadFromJSON()
		if err != nil {
			logFatal(err)
		}

		proof, err := mt.GetProofForAddr(addr)
		if err != nil {
			logFatal(err)
		}

		value, err := mt.GetLeafValueForAddr(addr)
		if err != nil {
			logFatal(err)
		}

		tx, err := airdropInstance.Claim(auth, mt.ID(), proof, value)
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
	// airdropCmd.AddCommand(redeemCmd)
}
