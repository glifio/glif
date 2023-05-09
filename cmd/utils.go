package cmd

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
)

func ParseAddress(ctx context.Context, addr string) (common.Address, error) {
	lapi, closer, err := fevm.Connection().ConnectLotusClient()
	if err != nil {
		return common.Address{}, err
	}
	defer closer()

	return parseAddress(ctx, addr, lapi)
}

func parseAddress(ctx context.Context, addr string, lapi lotusapi.FullNode) (common.Address, error) {
	if strings.HasPrefix(addr, "0x") {
		return common.HexToAddress(addr), nil
	}
	// user passed f1, f2, f3, or f4
	filAddr, err := address.NewFromString(addr)

	if err != nil {
		return common.Address{}, err
	}

	if filAddr.Protocol() != address.ID && filAddr.Protocol() != address.Delegated {
		filAddr, err = lapi.StateLookupID(ctx, filAddr, types.EmptyTSK)
		if err != nil {
			return common.Address{}, err
		}
	}

	ethAddr, err := ethtypes.EthAddressFromFilecoinAddress(filAddr)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(ethAddr.String()), nil
}

func commonSetupOwnerCall() (common.Address, *ecdsa.PrivateKey, error) {
	as := util.AgentStore()
	ks := util.KeyStore()
	// Check if an agent already exists
	agentAddrStr, err := as.Get("address")
	if err != nil {
		return common.Address{}, nil, err
	}

	if agentAddrStr == "" {
		return common.Address{}, nil, errors.New("No agent found. Did you forget to create one?")
	}

	agentAddr := common.HexToAddress(agentAddrStr)

	pk, err := ks.GetPrivate(util.OwnerKey)
	if err != nil {
		return common.Address{}, nil, err
	}

	if pk == nil {
		return common.Address{}, nil, errors.New("Owner key not found. Please check your `keys.toml` file. Only an Agent's owner can add a miner to it")
	}

	return agentAddr, pk, nil
}

func commonOwnerOrOperatorSetup(cmd *cobra.Command) (common.Address, *ecdsa.PrivateKey, error) {
	as := util.AgentStore()
	ks := util.KeyStore()

	opEvm, opFevm, err := ks.GetAddrs(util.OperatorKey)
	if err != nil {
		return common.Address{}, nil, err
	}
	owEvm, owFevm, err := ks.GetAddrs(util.OwnerKey)
	if err != nil {
		return common.Address{}, nil, err
	}

	var pk *ecdsa.PrivateKey
	// if no flag was passed, we just use the operator address by default
	from := cmd.Flag("from").Value.String()
	if from == "" {
		pk, err = ks.GetPrivate(util.OperatorKey)
	} else if from == opEvm.String() || from == opFevm.String() {
		pk, err = ks.GetPrivate(util.OperatorKey)
	} else if from == owEvm.String() || from == owFevm.String() {
		pk, err = ks.GetPrivate(util.OwnerKey)
	} else {
		return common.Address{}, nil, errors.New("invalid from address")
	}

	if err != nil {
		return common.Address{}, nil, err
	}

	agentAddrStr, err := as.Get("address")
	if err != nil {
		return common.Address{}, nil, err
	}

	if agentAddrStr == "" {
		return common.Address{}, nil, errors.New("No agent found. Did you forget to create one?")
	}

	agentAddr := common.HexToAddress(agentAddrStr)

	return agentAddr, pk, nil
}

type PoolType uint64

const (
	InfinityPool PoolType = iota
)

var poolNames = map[string]PoolType{
	"infinity": InfinityPool,
}

func parsePoolType(pool string) (*big.Int, error) {
	if pool == "" {
		return big.NewInt(int64(InfinityPool)), nil
	}

	poolType, ok := poolNames[pool]
	if !ok {
		return nil, errors.New("invalid pool")
	}

	return big.NewInt(int64(poolType)), nil
}

func parseFILAmount(amount string) (*big.Int, error) {
	amt, ok := new(big.Float).SetString(amount)
	if !ok {
		return nil, errors.New("invalid amount")
	}

	return util.ToAtto(amt), nil
}
