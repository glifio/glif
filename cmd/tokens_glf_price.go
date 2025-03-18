package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	token "github.com/glifio/go-pools/token"
	"github.com/glifio/go-pools/util"
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

		fmt.Printf("Current price of GLF/FIL: 1 GLIF ≈ %0.08f FIL\n", token.GLFToFIL(slot0.SqrtPriceX96))
		fmt.Printf("Current price of FIL/GLF: 1 FIL ≈ %0.08f GLIF\n", token.FILToGLF(slot0.SqrtPriceX96))
	},
}

type QuotePath string

const (
	QuotePathFILGLF QuotePath = "fil:glf"
	QuotePathGLFFIL QuotePath = "glf:fil"
)

var quoteCmd = &cobra.Command{
	Use:   "quote <path> <amount>",
	Short: "Get the amount of token1 that would be received for swapping a `amount` of token0 from Sushi V3. Path is either: fil:glf or glf:fil",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if PoolsSDK.Query().ChainID().Cmp(big.NewInt(constants.MainnetChainID)) != 0 {
			logFatalf("Sushi is only available on Filecoin Mainnet")
		}

		var tokenIn, tokenOut common.Address
		path := QuotePath(args[0])
		switch path {
		case QuotePathFILGLF:
			tokenIn = PoolsSDK.Query().WFIL()
			tokenOut = PoolsSDK.Query().GLF()
		case QuotePathGLFFIL:
			tokenIn = PoolsSDK.Query().GLF()
			tokenOut = PoolsSDK.Query().WFIL()
		}

		amount, err := parseFILAmount(args[1])
		if err != nil {
			logFatalf("Failed to parse amount %s", err)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		client, err := PoolsSDK.Extern().ConnectEthClient()
		if err != nil {
			logFatalf("Failed to connect to Ethereum client: %s", err)
		}
		defer client.Close()

		quoteParams := abigen.IQuoterV2QuoteExactInputSingleParams{
			TokenIn:           tokenIn,
			TokenOut:          tokenOut,
			AmountIn:          amount,
			Fee:               big.NewInt(3000),
			SqrtPriceLimitX96: big.NewInt(0),
		}

		quoterABI, err := abigen.QuoterV2MetaData.GetAbi()
		if err != nil {
			logFatalf("Failed to get quoter ABI %s", err)
		}

		calldata, err := quoterABI.Pack("quoteExactInputSingle", quoteParams)
		if err != nil {
			logFatalf("Failed to pack quoteExactInputSingle %s", err)
		}

		// Step 5: Execute the eth_call
		callMsg := ethereum.CallMsg{
			To:   &deploy.SushiQuoterV2,
			Data: calldata,
		}

		result, err := client.CallContract(context.Background(), callMsg, nil)
		if err != nil {
			log.Fatalf("Failed to call contract: %v", err)
		}

		// Step 6: Decode the return value
		outputs, err := quoterABI.Unpack("quoteExactInputSingle", result)
		if err != nil {
			log.Fatalf("Failed to unpack return value: %v", err)
		}

		// Extract the estimated output amount
		// uint256 amountOut,
		// uint160 sqrtPriceX96After,
		// uint32 initializedTicksCrossed,
		// uint256 gasEstimate
		amountOut := outputs[0].(*big.Int)
		sqrtPriceX96After := outputs[1].(*big.Int)

		s.Stop()

		switch path {
		case QuotePathFILGLF:
			fmt.Printf("for %0.04f FIL, you would receive approximately %0.06f GLF\n", util.ToFIL(amount), util.ToFIL(amountOut))
		case QuotePathGLFFIL:
			fmt.Printf("for %0.04f GLF, you would receive approximately %0.06f FIL\n", util.ToFIL(amount), util.ToFIL(amountOut))
		}
		fmt.Printf("the price after the swap would be %0.06f GLF/FIL\n", token.GLFToFIL(sqrtPriceX96After))
	},
}

func init() {
	glifCmd.AddCommand(getPriceCmd)
	glifCmd.AddCommand(quoteCmd)
}
