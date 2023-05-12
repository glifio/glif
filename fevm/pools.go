package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/cli/util"
)

func (c *FEVMConnection) PoolsList() ([]common.Address, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	poolRegCaller, err := abigen.NewPoolRegistryCaller(c.PoolRegistryAddr, client)
	if err != nil {
		return nil, err
	}

	numPools, err := poolRegCaller.AllPoolsLength(nil)
	if err != nil {
		return nil, err
	}

	var pools []common.Address

	for i := big.NewInt(0); i.Cmp(numPools) < 0; i.Add(i, big.NewInt(1)) {
		poolAddr, err := poolRegCaller.AllPools(nil, i)
		if err != nil {
			return nil, err
		}
		pools = append(pools, poolAddr)
	}

	return pools, nil
}

func (c *FEVMConnection) PoolGetAccount(ctx context.Context, poolAddr common.Address, agentID *big.Int) (abigen.Account, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return abigen.Account{}, err
	}
	defer client.Close()

	poolRegCaller, err := abigen.NewInfinityPoolCaller(poolAddr, client)
	if err != nil {
		return abigen.Account{}, err
	}

	id, err := poolRegCaller.Id(nil)
	if err != nil {
		return abigen.Account{}, err
	}

	routerCaller, err := abigen.NewRouterCaller(c.RouterAddr, client)
	if err != nil {
		return abigen.Account{}, err
	}

	return routerCaller.GetAccount(nil, agentID, id)
}

func (c *FEVMConnection) PoolAvailableLiquidity(ctx context.Context, poolAddr common.Address) (float64, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	poolCaller, err := abigen.NewInfinityPoolCaller(poolAddr, client)
	if err != nil {
		return 0, err
	}

	assets, err := poolCaller.TotalBorrowableAssets(nil)
	if err != nil {
		return 0, err
	}

	inFIL, _ := util.ToFIL(assets).Float64()

	return inFIL, nil
}

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
