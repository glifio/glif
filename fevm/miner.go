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

func (c *FEVMConnection) AddMiner(ctx context.Context, agentAddr common.Address, minerAddr address.Address) (*types.Transaction, error) {
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

	tx, err := WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, agentTransactor.AddMiner, "Agent Add Miner")
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Remove miner
func (c *FEVMConnection) MinerRemove(ctx context.Context, minerAddr common.Address) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// record the block height
	blockHeight, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("blockHeight: %d", blockHeight)

	minerCaller, err := abigen.NewMinerRegistryTransactor(minerAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{minerAddr}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, minerCaller.RemoveMiner, "Miner Remove")
}
