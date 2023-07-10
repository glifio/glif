package cmd

import (
	"context"

	"github.com/filecoin-project/go-address"
	actorstypes "github.com/filecoin-project/go-state-types/actors"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type MockFullNodeAPI struct {
	api.FullNode
}

var idStr = "f01234"
var maskedIDStr = "0xFF000000000000000000000000000000000004d2"

func (m *MockFullNodeAPI) StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return address.NewFromString(idStr)
}

var evmActorID = "f05678"
var ethAccountID = "f09876"

func (m *MockFullNodeAPI) StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) {
	code := cid.Cid{}
	if actor.String() == evmActorID {
		code, _ = actors.GetActorCodeID(actorstypes.Version(actors.LatestVersion), manifest.EvmKey)
	} else if actor.String() == ethAccountID {
		code, _ = actors.GetActorCodeID(actorstypes.Version(actors.LatestVersion), manifest.EthAccountKey)
	}

	return &types.Actor{
		Code:    code,
		Head:    cid.Cid{},
		Nonce:   0,
		Balance: types.NewInt(0),
	}, nil
}
