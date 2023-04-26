package fevm

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NonceCache struct {
	nonces map[common.Address]*big.Int
	rpcUrl string
	mutex sync.Mutex
}

var nonceCache *NonceCache

func Nonce() *NonceCache {
	return nonceCache
}

func (c *FEVMConnection) InitNonceCache() {
	nonceCache = &NonceCache{}
	nonceCache.rpcUrl = c.RpcURL
}

func (n *NonceCache) BumpNonce(address common.Address, nonceOverride uint64) (*big.Int, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	oldNonce, exists := n.nonces[address]

	if (!exists) {
		var nonce uint64
		if nonceOverride != 0 {
			nonce = nonceOverride
		} else {
			client, err := ethclient.Dial(n.rpcUrl)
			if err != nil {
				return nil, err
			}
			defer client.Close()

			startNonce, err := client.NonceAt(context.Background(), address, nil)
			if err != nil {
				return nil, err
			}
			nonce = startNonce
		}

		oldNonce = big.NewInt(int64(nonce))
	}

	// Bump the nonce in the nonce cache map
	newNonce := new(big.Int).Add(oldNonce, big.NewInt(1))
	n.nonces[address] = newNonce

	return oldNonce, nil
}

