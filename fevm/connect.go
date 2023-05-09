package fevm

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/spf13/viper"
)

type FEVMConnection struct {
	// the lotus RPC url contains a token in it - see https://lotus.filecoin.io/storage-providers/setup/initialize/ for more info
	LotusRpcUrl       string
	EthRpcUrl         string
	ChainID           int64
	RouterAddress     common.Address
	IFILAddr          common.Address
	InfinityPoolAddr  common.Address
	AgentFactoryAddr  common.Address
	PoolRegistryAddr  common.Address
	MinerRegistryAddr common.Address
}

var connection *FEVMConnection

func Connection() *FEVMConnection {
	return connection
}

// NewConnectParams instantiates an FEVMConnection struct singleton
func InitFEVMConnection(ctx context.Context) error {
	ethRpcUrl := viper.GetString("daemon.rpc-url")

	client, err := ethclient.Dial(ethRpcUrl)
	if err != nil {
		return err
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	// https://api.node.glif.io/rpc/v1
	connection = &FEVMConnection{
		LotusRpcUrl:       ethRpcUrl,
		EthRpcUrl:         ethRpcUrl,
		ChainID:           chainID.Int64(),
		RouterAddress:     common.HexToAddress(viper.GetString("routes.router")),
		IFILAddr:          common.HexToAddress(viper.GetString("routes.ifil")),
		InfinityPoolAddr:  common.HexToAddress(viper.GetString("routes.infinity-pool")),
		AgentFactoryAddr:  common.HexToAddress(viper.GetString("routes.agent-factory")),
		PoolRegistryAddr:  common.HexToAddress(viper.GetString("routes.pool-registry")),
		MinerRegistryAddr: common.HexToAddress(viper.GetString("routes.miner-registry")),
	}

	return nil
}

func (c *FEVMConnection) ConnectEthClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.EthRpcUrl)
}

func (c *FEVMConnection) ConnectLotusClient() (*lotusapi.FullNodeStruct, jsonrpc.ClientCloser, error) {
	head := http.Header{}

	if viper.GetString("daemon.token") != "" {
		head.Add("Authorization", "Bearer "+viper.GetString("lotus.token"))
	}

	lapi := &lotusapi.FullNodeStruct{}

	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		c.LotusRpcUrl,
		"Filecoin",
		lotusapi.GetInternalStructs(lapi),
		head,
	)

	if err != nil {
		return nil, nil, err
	}

	return lapi, closer, nil
}
