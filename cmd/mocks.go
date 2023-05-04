package cmd

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type MockFullNodeAPI struct {
	api.FullNode
}

var idStr = "f01234"
var maskedIDStr = "0xFF000000000000000000000000000000000004d2"

func (m *MockFullNodeAPI) StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return address.NewFromString(idStr)
}
