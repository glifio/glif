package fevm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) AgentCreate(ctx context.Context, deployerPk *ecdsa.PrivateKey, owner common.Address, operator common.Address, request common.Address) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentFactoryTransactor, err := abigen.NewAgentFactoryTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{owner, operator, request}

	return WriteTx(ctx, deployerPk, client, args, agentFactoryTransactor.Create, "Agent Create")
}

func (c *FEVMConnection) AgentFilter(ctx context.Context, receipt *types.Receipt) (*big.Int, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentABI, err := abigen.AgentFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	agentFactoryFilterer, err := abigen.NewAgentFactoryFilterer(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, err
	}

	var agentID *big.Int

	for _, l := range receipt.Logs {
		event, err := agentABI.EventByID(l.Topics[0])
		if err != nil {
			return nil, err
		}
		if event.Name == "CreateAgent" {
			createAgentEvent, err := agentFactoryFilterer.ParseCreateAgent(*l)
			if err != nil {
				return nil, err
			}
			agentID = createAgentEvent.AgentID
		}
	}

	return agentID, nil
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
