package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/glif-confidential/vc"
)

var ADOClient struct {
	Borrow      func(context.Context, common.Address, *big.Int) (*vc.SignedCredential, error)
	Pay         func(context.Context, common.Address, *big.Int) (*vc.SignedCredential, error)
	Withdraw    func(context.Context, common.Address, *big.Int) (*vc.SignedCredential, error)
	PushFunds   func(context.Context, common.Address, *big.Int, address.Address) (*vc.SignedCredential, error)
	PullFunds   func(context.Context, common.Address, *big.Int, address.Address) (*vc.SignedCredential, error)
	AddMiner    func(context.Context, common.Address, address.Address) (*vc.SignedCredential, error)
	RemoveMiner func(context.Context, common.Address, address.Address) (*vc.SignedCredential, error)
}

func NewADOClient(ctx context.Context, rpcurl string) (jsonrpc.ClientCloser, error) {
	return jsonrpc.NewClient(ctx, rpcurl, "Agent", &ADOClient, nil)
}
