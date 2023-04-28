package fevm

import (
	"context"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/go-address"
	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/cli/rpc"
	"github.com/spf13/viper"
)

func (c *FEVMConnection) MinerAdd(ctx context.Context, agentAddr common.Address, minerAddr address.Address) (*types.Transaction, error) {
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

	sc, err := rpc.ADOClient.AddMiner(ctx, agentAddr, minerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// record the block height
	blockHeight, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	minerCaller, err := abigen.NewMinerregistryTransactor(minerAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{minerAddr}

	tx, err := WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, minerCaller.AddMiner, "Miner Add")
	if err != nil {
		return nil, err
	}

	mrFilterer, err := abigen.NewMinerregistryFilterer(minerAddr, client)
	if err != nil {
		return nil, err
	}

	// wait for the miner to be added
	iter, err := mrFilterer.FilterAddMiner(&bind.FilterOpts{Start: blockHeight, End: nil}, []common.Address{minerAddr})
	if err != nil {
		return nil, err
	}

	for iter.Next() {

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
		return nil, common.Address{}, nil, err
	}

	minerCaller, err := abigen.NewMinerregistryTransactor(minerAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{minerAddr}

	return WriteTx(ctx, &ecdsa.PrivateKey{}, client, args, minerCaller.RemoveMiner, "Miner Remove")
}
