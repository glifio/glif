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

	tx, err := WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.AddMiner, "Agent Add Miner")
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

	tx, err := WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.RemoveMiner, "Agent Remove Miner")
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

func (c *FEVMConnection) ChangeWorker(
	ctx context.Context,
	agentAddr common.Address,
	minerAddr address.Address,
	workerAddr address.Address,
	controlAddrs []address.Address,
	pk *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentTransactor, err := abigen.NewAgentTransactor(agentAddr, client)
	if err != nil {
		return nil, err
	}

	// convert miner address to ID address
	minerID, err := address.IDFromAddress(minerAddr)
	if err != nil {
		return nil, err
	}

	// convert worker address to ID address
	workerID, err := address.IDFromAddress(workerAddr)
	if err != nil {
		return nil, err
	}

	// convert control addresses to ID addresses
	var controlIDs []uint64
	for _, controlAddr := range controlAddrs {
		controlID, err := address.IDFromAddress(controlAddr)
		if err != nil {
			return nil, err
		}
		controlIDs = append(controlIDs, controlID)
	}

	args := []interface{}{minerID, workerID, controlIDs}

	tx, err := WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.ChangeMinerWorker, "Agent Change Miner Worker")
	if err != nil {
		return nil, err
	}

	return tx, nil
}
