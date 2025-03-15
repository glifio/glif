package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	token "github.com/glifio/go-pools/token/data"
	"github.com/spf13/cobra"
)

var getPriceCmd = &cobra.Command{
	Use:   "price",
	Short: "Get the current price of $GLF in FIL from Sushi V3 on FEVM",
	Run: func(cmd *cobra.Command, args []string) {
		if PoolsSDK.Query().ChainID().Cmp(big.NewInt(constants.MainnetChainID)) != 0 {
			logFatalf("Sushi is only available on Filecoin Mainnet")
		}

		ctx := cmd.Context()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		client, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to connect to Ethereum client: %s", err)
		}
		defer client.Close()

		// Connect to the Uniswap V3 Pool contract
		pool, err := abigen.NewUniswapV3PoolCaller(deploy.SushiGLFWFILPool, client)
		if err != nil {
			logFatalf("Failed to instantiate pool caller %s", err)
		}

		// Get the current price sqrt from slot0
		slot0, err := pool.Slot0(&bind.CallOpts{Context: ctx})
		if err != nil {
			logFatalf("Failed to get slot0 data %s", err)
		}

		s.Stop()

		fmt.Printf("Current price of GLF/FIL: 1 GLIF ≈ %0.05f FIL\n", token.GLFToFIL(slot0.SqrtPriceX96))
		fmt.Printf("Current price of FIL/GLF: 1 FIL ≈ %0.05f GLIF\n", token.FILToGLF(slot0.SqrtPriceX96))
	},
}

func init() {
	glifCmd.AddCommand(getPriceCmd)
}
