package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) AgentCreate(ctx context.Context, deployerPk *ecdsa.PrivateKey, owner common.Address, operator common.Address) (*big.Int, common.Address, *types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, common.Address{}, nil, err
	}
	defer client.Close()

	agentFactoryTransactor, err := abigen.NewAgentfactoryTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	args := []interface{}{owner, operator}

	tx, err := WriteTx(ctx, deployerPk, client, args, agentFactoryTransactor.Create, "Agent Create")
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	return big.NewInt(0), common.Address{}, tx, nil
}
