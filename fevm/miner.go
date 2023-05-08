package fevm

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/go-address"
	ltypes "github.com/filecoin-project/lotus/chain/types"
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
	if err != nil {
		return nil, err
	}

	sc, err := rpc.ADOClient.AddMiner(ctx, agentAddr, minerAddr)
	if err != nil {
		return nil, err
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
	newOwnerAddr address.Address,
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

	sc, err := rpc.ADOClient.RemoveMiner(ctx, agentAddr, minerAddr)
	if err != nil {
		return nil, err
	}

	args := []interface{}{newOwnerAddr, sc}

	tx, err := WriteTx(ctx, pk, client, args, agentTransactor.RemoveMiner, "Agent Remove Miner")
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// converts a FIL address to an ID address
func (c *FEVMConnection) ToMinerID(miner string) (address.Address, error) {
	minerAddr, err := address.NewFromString(miner)
	if err != nil {
		return address.Undef, err
	}

	if minerAddr.Protocol() == address.ID {
		return minerAddr, nil
	}

	lapi, closer, err := c.ConnectLotusClient()
	if err != nil {
		return address.Undef, err
	}
	defer closer()

	idAddr, err := lapi.StateLookupID(context.Background(), minerAddr, ltypes.EmptyTSK)
	if err != nil {
		return address.Undef, err
	}

	return idAddr, nil
}
