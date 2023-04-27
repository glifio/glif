package fevm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) PoolsList() ([]common.Address, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	poolRegCaller, err := abigen.NewPoolregistryCaller(c.PoolRegistryAddr, client)
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
