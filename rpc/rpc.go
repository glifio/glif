package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/vc"
)

var ADOClient struct {
	Borrow      func(context.Context, common.Address, *big.Int) (abigen.SignedCredential, error)
	Pay         func(context.Context, common.Address, *big.Int) (abigen.SignedCredential, error)
	Withdraw    func(context.Context, common.Address, *big.Int) (abigen.SignedCredential, error)
	PushFunds   func(context.Context, common.Address, *big.Int, address.Address) (abigen.SignedCredential, error)
	PullFunds   func(context.Context, common.Address, *big.Int, address.Address) (abigen.SignedCredential, error)
	AddMiner    func(context.Context, common.Address, address.Address) (abigen.SignedCredential, error)
	RemoveMiner func(context.Context, common.Address, address.Address) (abigen.SignedCredential, error)
}

func NewADOClient(ctx context.Context, rpcurl string) (jsonrpc.ClientCloser, error) {
	return jsonrpc.NewClient(ctx, rpcurl, "Mock", &ADOClient, nil)
}

var MinerInfoClient struct {
	Info func(context.Context, []address.Address) (vc.VerifiableCredential, error)
}

func NewMinerInfoClient(ctx context.Context, rpcurl string) (jsonrpc.ClientCloser, error) {
	return jsonrpc.NewClient(ctx, rpcurl, "Miner", &ADOClient, nil)
}
