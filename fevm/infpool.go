package fevm

import (
	"context"
	"math/big"

	abigen "github.com/glif-confidential/abigen/bindings"
	"github.com/glif-confidential/cli/util"
)

func (c *FEVMConnection) AgentLevel(ctx context.Context, agentID *big.Int) (*big.Int, float64, error) {
	client, err := c.ConnectEthClient()
	if err != nil {
		return nil, 0, err
	}
	defer client.Close()

	infpool, err := abigen.NewInfinityPoolCaller(c.InfinityPoolAddr, client)
	if err != nil {
		return nil, 0, err
	}

	rateModuleAddr, err := infpool.RateModule(nil)
	if err != nil {
		return nil, 0, err
	}

	rateModule, err := abigen.NewRateModuleCaller(rateModuleAddr, client)
	if err != nil {
		return nil, 0, err
	}

	lvl, err := rateModule.AccountLevel(nil, agentID)
	if err != nil {
		return nil, 0, err
	}

	cap, err := rateModule.Levels(nil, lvl)
	if err != nil {
		return nil, 0, err
	}

	capInFIL, _ := util.ToFIL(cap).Float64()

	return lvl, capInFIL, nil
}
