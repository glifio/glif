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

	// record the block height
	blockHeight, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	agentFactoryTransactor, err := abigen.NewAgentFactoryTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	args := []interface{}{owner, operator}

	tx, err := WriteTx(ctx, deployerPk, client, args, agentFactoryTransactor.Create, "Agent Create")
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	//TODO: watch for event and return agent id and address

	afFilterer, err := abigen.NewAgentFactoryFilterer(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	aIDs := []*big.Int{}
	agents := []common.Address{owner}
	operators := []common.Address{operator}

	//TODO: I don't think this filter works
	iter, err := afFilterer.FilterCreateAgent(&bind.FilterOpts{Start: blockHeight}, aIDs, agents, operators)
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

	return agentID, common.Address{}, tx, nil
}

func (c *FEVMConnection) AgentPullFunds(ctx context.Context, agentID *big.Int, amount *big.Int) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentTransactor, err := abigen.NewAgentTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{amount}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, agentTransactor.PullFunds, "Agent Pull Funds")
}

// AgentPushFunds pushes funds from the agent to a miner
func (c *FEVMConnection) AgentPushFunds(ctx context.Context, agentID *big.Int, miner common.Address, amount *big.Int) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentTransactor, err := abigen.NewAgentTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{miner, amount}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, agentTransactor.PushFunds, "Agent Push Funds")
}
