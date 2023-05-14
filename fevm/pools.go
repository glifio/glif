package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) PoolDepositFIL(ctx context.Context, poolAddr common.Address, receiver common.Address, amount *big.Int, pk *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	poolTransactor, err := abigen.NewInfinityPoolTransactor(poolAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{receiver}

	return WriteTx(ctx, pk, client, amount, args, poolTransactor.Deposit0, "Deposit FIL")
}
