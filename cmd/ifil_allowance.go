package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

var iFILAllowanceCmd = &cobra.Command{
	Use:   "allowance [owner] [spender]",
	Short: "Get the iFIL balance of an address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		owner := args[0]
		spender := args[1]
		fmt.Printf("Checking iFIL allowance of spender: %s on behalf of owner: %s ...", spender, owner)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()
		defer s.Stop()
		ownerAddr, err := AddressOrAccountNameToEVM(cmd.Context(), owner)
		if err != nil {
			logFatalf("Failed to parse owner address %s", err)
		}

		spenderAddr, err := AddressOrAccountNameToEVM(cmd.Context(), spender)
		if err != nil {
			logFatalf("Failed to parse spender address %s", err)
		}

		client, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}
		defer client.Close()

		poolToken, err := abigen.NewPoolTokenCaller(PoolsSDK.Query().IFIL(), client)
		if err != nil {
			logFatalf("Failed to get iFIL balance %s", err)
		}

		allow, err := poolToken.Allowance(&bind.CallOpts{}, ownerAddr, spenderAddr)
		if err != nil {
			logFatalf("Failed to get iFIL allowance %s", err)
		}

		s.Stop()

		fmt.Printf("iFIL allowance for spender: %s on behalf of owner: %s is %.09f\n", spender, owner, util.ToFIL(allow))

	},
}

func init() {
	iFILCmd.AddCommand(iFILAllowanceCmd)
}
