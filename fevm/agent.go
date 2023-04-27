package fevm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) AgentCreate(owner common.Address, operator common.Address) (*big.Int, common.Address, *types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, common.Address{}, err
	}
	defer client.Close()

	agentFactoryTransactor, err := abigen.NewAgentfactoryTransactor(c.AgentFactoryAddr, client)
	if err != nil {
		return nil, common.Address{}, err
	}

	args := []interface{}{owner, operator}

	tx, err := WriteTx(&c.ctx, &c.privateKey, client, args, agentFactoryTransactor.Create, "Agent Create")
	if err != nil {
		return nil, common.Address{}, err
	}

	return big.NewInt(0), common.Address{}, tx, nil
}
