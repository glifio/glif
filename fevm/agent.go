package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	//TODO: watch for event and return agent id and address

	afFilterer, err := abigen.NewAgentfactoryFilterer(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	aIDs := []*big.Int{}
	agents := []common.Address{owner}
	operators := []common.Address{operator}

	iter, err := afFilterer.FilterCreateAgent(&bind.FilterOpts{}, aIDs, agents, operators)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	var agentID *big.Int

	for iter.Next() {
		agent := iter.Event.Agent
		agentID = iter.Event.AgentID

		if agent == owner && agentID != nil {
			break
		}
	}

	return big.NewInt(0), common.Address{}, tx, nil
}
