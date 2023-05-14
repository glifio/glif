package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) IFILTransfer(ctx context.Context, toAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	poolTokenCaller, err := abigen.NewPoolTokenTransactor(c.IFILAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{toAddr, amount}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, common.Big0, args, poolTokenCaller.Transfer, "iFIL Transfer")
}

func (c *FEVMConnection) IFILApprove(ctx context.Context, spender common.Address, allowance *big.Int) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	poolTokenCaller, err := abigen.NewPoolTokenTransactor(c.IFILAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{spender, allowance}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, common.Big0, args, poolTokenCaller.Approve, "iFIL Approve")
}
