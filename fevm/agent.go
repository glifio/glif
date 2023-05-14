package fevm

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/filecoin-project/go-address"
	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/ado/constants"
	"github.com/glif-confidential/cli/rpc"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (c *FEVMConnection) AgentID(ctx context.Context, address common.Address) (*big.Int, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentCaller, err := abigen.NewAgentCaller(address, client)
	if err != nil {
		return nil, err
	}

	agentID, err := agentCaller.Id(nil)
	if err != nil {
		return nil, err
	}

	return agentID, nil
}

func (c *FEVMConnection) AgentOwner(ctx context.Context, address common.Address) (common.Address, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()

	agentCaller, err := abigen.NewAgentCaller(address, client)
	if err != nil {
		return common.Address{}, err
	}

	owner, err := agentCaller.Owner(nil)
	if err != nil {
		return common.Address{}, err
	}

	return owner, nil
}

func (c *FEVMConnection) AgentAssets(ctx context.Context, address common.Address) (*big.Int, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	agentCaller, err := abigen.NewAgentCaller(address, client)
	if err != nil {
		return nil, err
	}

	assets, err := agentCaller.LiquidAssets(nil)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

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

	return WriteTx(ctx, deployerPk, client, common.Big0, args, agentFactoryTransactor.Create, "Agent Create")
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

// AgentPullFunds pulls funds from the agent to a miner
func (c *FEVMConnection) AgentPullFunds(
	ctx context.Context,
	agentAddr common.Address,
	amount *big.Int,
	miner address.Address,
	pk *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	as := util.AgentStore()
	agentIDStr, err := as.Get("id")
	if err != nil {
		return nil, err
	}

	agentID, ok := new(big.Int).SetString(agentIDStr, 10)
	if !ok {
		return nil, errors.New("could not convert agent id to big.Int")
	}

	minerRegistryCaller, err := abigen.NewMinerRegistryCaller(c.MinerRegistryAddr, client)
	if err != nil {
		return nil, err
	}

	minerU64, err := address.IDFromAddress(miner)
	if err != nil {
		return nil, err
	}

	registered, err := minerRegistryCaller.MinerRegistered(nil, agentID, minerU64)
	if err != nil {
		return nil, err
	}

	if !registered {
		return nil, errors.New("Miner not registered with agent. Be sure to call `agent add-miner` first before pulling funds.")
	}

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	sc, err := rpc.ADOClient.PullFunds(ctx, agentAddr, amount, miner)
	if err != nil {
		return nil, err
	}

	agentTransactor, err := abigen.NewAgentTransactor(agentAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{sc}

	return WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.PullFunds, "Agent Pull Funds")
}

// AgentPushFunds pushes funds from the agent to a miner
func (c *FEVMConnection) AgentPushFunds(
	ctx context.Context,
	agentAddr common.Address,
	amount *big.Int,
	miner address.Address,
	pk *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	as := util.AgentStore()
	agentIDStr, err := as.Get("id")
	if err != nil {
		return nil, err
	}

	agentID, ok := new(big.Int).SetString(agentIDStr, 10)
	if !ok {
		return nil, errors.New("could not convert agent id to big.Int")
	}

	minerRegistryCaller, err := abigen.NewMinerRegistryCaller(c.MinerRegistryAddr, client)
	if err != nil {
		return nil, err
	}

	minerU64, err := address.IDFromAddress(miner)
	if err != nil {
		return nil, err
	}

	registered, err := minerRegistryCaller.MinerRegistered(nil, agentID, minerU64)
	if err != nil {
		return nil, err
	}

	if !registered {
		return nil, errors.New("Miner not registered with agent. Be sure to call `agent add-miner` first before pushing funds.")
	}

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	sc, err := rpc.ADOClient.PushFunds(ctx, agentAddr, amount, miner)
	if err != nil {
		return nil, err
	}

	agentTransactor, err := abigen.NewAgentTransactor(agentAddr, client)
	if err != nil {
		return nil, err
	}

	args := []interface{}{sc}

	return WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.PushFunds, "Agent Push Funds")
}

func (c *FEVMConnection) AgentBorrow(
	ctx context.Context,
	agentAddr common.Address,
	poolID *big.Int,
	amount *big.Int,
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

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	sc, err := rpc.ADOClient.Borrow(ctx, agentAddr, amount)
	if err != nil {
		log.Fatal(err)
	}

	args := []interface{}{poolID, sc}

	return WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.Borrow, "Agent Borrow")
}

func (c *FEVMConnection) AgentPay(
	ctx context.Context,
	agentAddr common.Address,
	poolID *big.Int,
	amount *big.Int,
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

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	sc, err := rpc.ADOClient.Pay(ctx, agentAddr, amount)
	if err != nil {
		log.Fatal(err)
	}

	args := []interface{}{poolID, sc}

	return WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.Pay, "Agent Pay")
}

func (c *FEVMConnection) AgentOwes(cmd *cobra.Command, agentAddr common.Address) (*big.Int, *big.Int, error) {
	closer, err := rpc.NewADOClient(cmd.Context(), viper.GetString("ado.address"))
	if err != nil {
		return nil, nil, err
	}
	defer closer()

	agentOwed, err := rpc.ADOClient.AmountOwed(cmd.Context(), agentAddr, constants.INFINITY_POOL_ID)
	if err != nil {
		return nil, nil, err
	}

	return agentOwed.AmountOwed, agentOwed.Gcred, nil
}

func (c *FEVMConnection) AgentWithdraw(
	ctx context.Context,
	agentAddr common.Address,
	receiver common.Address,
	amount *big.Int,
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

	closer, err := rpc.NewADOClient(ctx, viper.GetString("ado.address"))
	if err != nil {
		return nil, err
	}
	defer closer()

	sc, err := rpc.ADOClient.Withdraw(ctx, agentAddr, amount)
	if err != nil {
		log.Fatal(err)
	}

	args := []interface{}{receiver, sc}

	return WriteTx(ctx, pk, client, common.Big0, args, agentTransactor.Withdraw, "Agent Withdraw")
}
