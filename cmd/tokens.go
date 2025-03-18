package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/go-pools/abigen"
	"github.com/glifio/go-pools/util"
	denoms "github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Commands for interacting with tokens",
}

var iFILNewCmd = &cobra.Command{
	Use:   "ifil",
	Short: "Commands for interacting with the Infinity Pool Liquid Staking Token (iFIL)",
}

var glifCmd = &cobra.Command{
	Use:   "glf",
	Short: "Commands for interacting with the GLIF token",
}

var wFILNewCmd = &cobra.Command{
	Use:   "wfil",
	Short: "Commands for interacting with the Wrapped Filecoin token",
}

// generic methods for ERC20 tokens
var allowanceFunc = func(cmd *cobra.Command, args []string) {
	token, tokenAddress := parseToken(cmd)

	owner := args[0]
	spender := args[1]
	fmt.Printf("Checking %s allowance of spender: %s on behalf of owner: %s ...", token, spender, owner)

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

	poolToken, err := abigen.NewPoolTokenCaller(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get iFIL balance %s", err)
	}

	allow, err := poolToken.Allowance(&bind.CallOpts{}, ownerAddr, spenderAddr)
	if err != nil {
		logFatalf("Failed to get iFIL allowance %s", err)
	}

	s.Stop()

	fmt.Printf("%s allowance for spender: %s on behalf of owner: %s is %.09f\n", token, spender, owner, util.ToFIL(allow))
}

var approveFunc = func(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	token, tokenAddress := parseToken(cmd)

	from := cmd.Flag("from").Value.String()

	client, err := PoolsSDK.Extern().ConnectEthClient()
	if err != nil {
		logFatalf("Failed to get iFIL balance %s", err)
	}
	defer client.Close()

	auth, _, err := commonGenericAccountSetup(cmd, from)
	if err != nil {
		logFatal(err)
	}

	strAddr := args[0]
	strAmt := args[1]
	fmt.Printf("Approving %s to spend %s of your %s balance...\n", strAddr, strAmt, token)

	addr, err := AddressOrAccountNameToEVM(ctx, strAddr)
	if err != nil {
		logFatalf("Failed to parse address %s", err)
	}

	amount, err := parseFILAmount(strAmt)
	if err != nil {
		logFatalf("Failed to parse amount %s", err)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	poolTokenTransactor, err := abigen.NewPoolTokenTransactor(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get %s transactor %s", token, err)
	}

	tx, err := poolTokenTransactor.Approve(auth, addr, amount)
	if err != nil {
		logFatalf("Failed to approve %s %s", token, err)
	}

	_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
	if err != nil {
		logFatalf("Failed to approve %s %s", token, err)
	}

	s.Stop()

	fmt.Printf("%s approved!\n", token)
}

var transferFunc = func(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	token, tokenAddress := parseToken(cmd)
	from := cmd.Flag("from").Value.String()
	auth, _, err := commonGenericAccountSetup(cmd, from)
	if err != nil {
		logFatal(err)
	}

	strAddr := args[0]
	strAmt := args[1]
	fmt.Printf("Transferring %s %s to %s...\n", strAmt, token, strAddr)

	addr, err := AddressOrAccountNameToEVM(ctx, strAddr)
	if err != nil {
		logFatalf("Failed to parse address %s", err)
	}

	amount, err := parseFILAmount(strAmt)
	if err != nil {
		logFatalf("Failed to parse amount %s", err)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	client, err := PoolsSDK.Extern().ConnectEthClient()
	if err != nil {
		logFatalf("Failed to get %s transactor %s", token, err)
	}
	defer client.Close()

	poolTokenTransactor, err := abigen.NewPoolTokenTransactor(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get %s transactor %s", token, err)
	}

	tx, err := poolTokenTransactor.Transfer(auth, addr, amount)
	if err != nil {
		logFatalf("Failed to transfer %s %s", token, err)
	}

	s.Stop()

	fmt.Printf("Confirming transfer transaction: %s...\n", tx.Hash().Hex())

	s.Start()

	_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
	if err != nil {
		logFatalf("Failed to transfer %s %s", token, err)
	}

	s.Stop()

	fmt.Printf("Successfully transferred %0.03f %s from %s to %s!\n", util.ToFIL(amount), token, from, strAddr)
}

var transferFromFunc = func(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	holder := args[0]
	to := args[1]
	strAmt := args[2]
	token, tokenAddress := parseToken(cmd)
	from := cmd.Flag("from").Value.String()
	auth, _, err := commonGenericAccountSetup(cmd, from)
	if err != nil {
		logFatal(err)
	}

	fmt.Printf("Transferring %s %s from %s to %s...\n", token, strAmt, from, to)

	fromAddr, err := AddressOrAccountNameToEVM(ctx, holder)
	if err != nil {
		logFatalf("Failed to parse from address %s", err)
	}

	toAddr, err := AddressOrAccountNameToEVM(ctx, to)
	if err != nil {
		logFatalf("Failed to parse to address %s", err)
	}

	amount, err := parseFILAmount(strAmt)
	if err != nil {
		logFatalf("Failed to parse amount %s", err)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	client, err := PoolsSDK.Extern().ConnectEthClient()
	if err != nil {
		logFatalf("Failed to get %s transactor %s", token, err)
	}
	defer client.Close()

	poolTokenTransactor, err := abigen.NewPoolTokenTransactor(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get %s transactor %s", token, err)
	}

	tx, err := poolTokenTransactor.TransferFrom(auth, fromAddr, toAddr, amount)
	if err != nil {
		logFatalf("Failed to transfer %s %s", token, err)
	}

	s.Stop()

	fmt.Printf("Confirming transfer from transaction: %s...\n", tx.Hash().Hex())

	s.Start()

	_, err = PoolsSDK.Query().StateWaitReceipt(ctx, tx.Hash())
	if err != nil {
		logFatalf("Failed to transfer from %s %s", token, err)
	}

	s.Stop()

	fmt.Printf("Successfully transferred %0.03f %s from %s to %s!\n", util.ToFIL(amount), token, from, to)
}

var balanceOfFunc = func(cmd *cobra.Command, args []string) {
	strAddr := args[0]
	token, tokenAddress := parseToken(cmd)

	fmt.Printf("Checking %s balance of %s...\n", strAddr, token)

	addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
	if err != nil {
		logFatalf("Failed to parse address %s", err)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	client, err := PoolsSDK.Extern().ConnectEthClient()
	if err != nil {
		logFatalf("Failed to get %s caller %s", token, err)
	}
	defer client.Close()

	poolTokenCaller, err := abigen.NewPoolTokenCaller(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get %s caller %s", token, err)
	}

	bal, err := poolTokenCaller.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		logFatalf("Failed to get %s balance %s", token, err)
	}

	s.Stop()

	fmt.Printf("%s balance of %s is %.09f\n", token, strAddr, util.ToFIL(bal))
}

var supplyFunc = func(cmd *cobra.Command, args []string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	token, tokenAddress := parseToken(cmd)

	client, err := PoolsSDK.Extern().ConnectEthClient()
	if err != nil {
		logFatalf("Failed to get %s caller %s", token, err)
	}
	defer client.Close()

	poolTokenCaller, err := abigen.NewPoolTokenCaller(tokenAddress, client)
	if err != nil {
		logFatalf("Failed to get %s caller %s", token, err)
	}

	supply, err := poolTokenCaller.TotalSupply(&bind.CallOpts{Context: cmd.Context()})
	if err != nil {
		logFatalf("Failed to get %s supply %s", token, err)
	}

	supplyFIL, _ := denoms.ToFIL(supply).Float64()

	s.Stop()

	fmt.Printf("%.09f %s\n", supplyFIL, token)
}

var allowanceCmd = cobra.Command{
	Use:   "allowance <owner> <spender>",
	Short: "Get the amount of tokens that `spender` is allowed to spend on behalf of `owner`",
	Args:  cobra.ExactArgs(2),
	Run:   allowanceFunc,
}

var approveCmd = cobra.Command{
	Use:   "approve <spender> <amount>",
	Short: "Approve a spender for a token",
	Args:  cobra.ExactArgs(2),
	Run:   approveFunc,
}

var transferCmd = cobra.Command{
	Use:   "transfer <recipient> <amount>",
	Short: "Transfer `amount` of tokens to the recipient address",
	Args:  cobra.ExactArgs(2),
	Run:   transferFunc,
}

var transferFromCmd = cobra.Command{
	Use:   "transfer-from <from> <to> <amount>",
	Short: "Transfer `amount` of tokens from the `from` address to the `to` address for the token",
	Args:  cobra.ExactArgs(3),
	Run:   transferFromFunc,
}

var balanceOfCmd = cobra.Command{
	Use:   "balance-of <address>",
	Short: "Get the token balance of an address",
	Args:  cobra.ExactArgs(1),
	Run:   balanceOfFunc,
}

var supplyCmd = cobra.Command{
	Use:   "supply",
	Short: "Get the supply of a token",
	Args:  cobra.NoArgs,
	Run:   supplyFunc,
}

// this allows us to effectively create the same methods for each ERC20 token but have the 'parent' be different so we can identify the correct token to use
func createCommand(cmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   cmd.Use,
		Short: cmd.Short,
		Args:  cmd.Args,
		Run:   cmd.Run,
	}
}

func parseToken(cmd *cobra.Command) (string, common.Address) {
	switch cmd.Parent().Use {
	case "ifil":
		return "iFIL", PoolsSDK.Query().IFIL()
	case "glf":
		return "GLF", PoolsSDK.Query().GLF()
	case "wfil":
		return "wFIL", PoolsSDK.Query().WFIL()
	}
	return "", common.HexToAddress("0x0000000000000000000000000000000000000000")
}

func init() {
	rootCmd.AddCommand(tokensCmd)
	tokensCmd.AddCommand(iFILNewCmd)
	tokensCmd.AddCommand(glifCmd)
	tokensCmd.AddCommand(wFILNewCmd)

	tokensCmd.PersistentFlags().String("from", "", "address to send the transaction from")

	wFILNewCmd.AddCommand(createCommand(&allowanceCmd))
	iFILNewCmd.AddCommand(createCommand(&allowanceCmd))
	glifCmd.AddCommand(createCommand(&allowanceCmd))

	wFILNewCmd.AddCommand(createCommand(&approveCmd))
	iFILNewCmd.AddCommand(createCommand(&approveCmd))
	glifCmd.AddCommand(createCommand(&approveCmd))

	wFILNewCmd.AddCommand(createCommand(&transferCmd))
	iFILNewCmd.AddCommand(createCommand(&transferCmd))
	glifCmd.AddCommand(createCommand(&transferCmd))

	wFILNewCmd.AddCommand(createCommand(&balanceOfCmd))
	iFILNewCmd.AddCommand(createCommand(&balanceOfCmd))
	glifCmd.AddCommand(createCommand(&balanceOfCmd))

	wFILNewCmd.AddCommand(createCommand(&supplyCmd))
	iFILNewCmd.AddCommand(createCommand(&supplyCmd))
	glifCmd.AddCommand(createCommand(&supplyCmd))

	wFILNewCmd.AddCommand(createCommand(&transferFromCmd))
	iFILNewCmd.AddCommand(createCommand(&transferFromCmd))
	glifCmd.AddCommand(createCommand(&transferFromCmd))
}
