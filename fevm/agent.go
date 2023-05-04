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

func (c *FEVMConnection) AgentAddrID(ctx context.Context, receipt *types.Receipt) (*big.Int, common.Address, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, common.Address{}, err
	}
	defer client.Close()

	agentFactoryABI, err := abigen.AgentFactoryMetaData.GetAbi()
	if err != nil {
		return nil, common.Address{}, err
	}

	agentFactoryFilterer, err := abigen.NewAgentFactoryFilterer(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, err
	}

	var agentID *big.Int
	var agentAddr common.Address

	for _, l := range receipt.Logs {
		event, err := agentFactoryABI.EventByID(l.Topics[0])
		if err != nil {
			return nil, common.Address{}, err
		}
		if event.Name == "CreateAgent" {
			createAgentEvent, err := agentFactoryFilterer.ParseCreateAgent(*l)
			if err != nil {
				return nil, common.Address{}, err
			}
			agentAddr = createAgentEvent.Agent
			agentID = createAgentEvent.AgentID
		}
	}

	return agentID, agentAddr, nil
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
