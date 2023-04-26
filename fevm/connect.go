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
	RpcURL             string
	Token              string
	ChainID            int64
	RouterAddress      common.Address
	IFILAddr           common.Address
	InfinityPoolAddr   common.Address
	AgentFactoryAddr   common.Address
}

var connection *FEVMConnection

func Connection() *FEVMConnection {
	return connection
}

// NewConnectParams instantiates an FEVMConnection struct singleton
func InitFEVMConnection(ctx context.Context) error {
	rpcUrl := viper.GetString("daemon.rpc-url")

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	connection = &FEVMConnection{
		RpcURL: viper.GetString("daemon.rpc-url"),
		ChainID: chainID.Int64(),
		RouterAddress: common.HexToAddress(viper.GetString("routes.router")),
		IFILAddr: common.HexToAddress(viper.GetString("routes.ifil")),
		InfinityPoolAddr: common.HexToAddress(viper.GetString("routes.infinity-pool")),
		AgentFactoryAddr: common.HexToAddress(viper.GetString("routes.agent-factory")),
	}

	return nil
}

func (c *FEVMConnection) ConnectEthClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.RpcURL)
}

func (c *FEVMConnection) ConnectLotusClient() (*lotusapi.FullNodeStruct, jsonrpc.ClientCloser, error) {
		head := http.Header{}
		// var api lotusapi.FullNodeStruct

		if c.Token != "" {
			head.Set("Authorization", "Bearer "+c.Token)
		}

		lapi := &lotusapi.FullNodeStruct{}

		closer, err := jsonrpc.NewMergeClient(
			context.Background(),
			c.RpcURL,
			"Filecoin",
			//[]interface{}{&api.Internal, &api.CommonStruct.Internal},
			lotusapi.GetInternalStructs(lapi),
			head,
		)

		if err != nil {
			return nil, nil, err
		}

		return lapi, closer, nil
}
