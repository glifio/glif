package fevm

import (
	"context"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/go-address"
	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/cli/rpc"
	"github.com/spf13/viper"
)

func (c *FEVMConnection) AddMiner(
	ctx context.Context,
	agentAddr common.Address,
	minerAddr address.Address,
	pk *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	agentTransactor, err := abigen.NewAgentTransactor(agentAddr, client)

	sc, err := rpc.ADOClient.AddMiner(ctx, agentAddr, minerAddr)
	if err != nil {
		log.Fatal(err)
	}

	args := []interface{}{sc}

	tx, err := WriteTx(ctx, pk, client, args, agentTransactor.AddMiner, "Agent Add Miner")
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Remove miner
func (c *FEVMConnection) RemoveMiner(
	ctx context.Context,
	agentAddr common.Address,
	minerAddr address.Address,
	recipientAddr address.Address,
	pk *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	agentTransactor, err := abigen.NewAgentTransactor(agentAddr, client)

	sc, err := rpc.ADOClient.AddMiner(ctx, agentAddr, minerAddr)
	if err != nil {
		log.Fatal(err)
	}

	args := []interface{}{recipientAddr, sc}

	tx, err := WriteTx(ctx, pk, client, args, agentTransactor.RemoveMiner, "Agent Remove Miner")
	if err != nil {
		return nil, err
	}

	return tx, nil
}
