package fevm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abigen "github.com/glif-confidential/abigen/bindings"
)

func (c *FEVMConnection) IFILBalanceOf(address common.Address) (*big.Int, error) {
  client, err := c.ConnectEthClient()
  if err != nil {
    return nil, err
  }
  defer client.Close()

  poolTokenCaller, err := abigen.NewPooltokenCaller(address, client)
  if err != nil {
    return nil, err
  }

  return poolTokenCaller.BalanceOf(nil, c.IFILAddr)
}
