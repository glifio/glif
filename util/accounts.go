package util

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
)

type AccountsStorage struct {
	*Storage
}

var accountsStore *AccountsStorage

func AccountsStore() *AccountsStorage {
	return accountsStore
}

func NewAccountsStore(filename string) error {
	accountsDefault := map[string]string{}

	s, err := NewStorage(filename, accountsDefault)
	if err != nil {
		return err
	}

	accountsStore = &AccountsStorage{s}

	return nil
}

func (a *AccountsStorage) GetAddrs(key KeyType) (common.Address, address.Address, error) {
	evmAddress := common.HexToAddress(a.data[string(key)])

	delegated, err := DelegatedFromEthAddr(evmAddress)
	if err != nil {
		return evmAddress, address.Address{}, err
	}

	return evmAddress, delegated, nil
}
