package util

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glif-confidential/cli/fevm"
)

var ZeroEvmAddr = common.Address{}

type KeyType string

const (
	OwnerKey    KeyType = "owner"
	OperatorKey KeyType = "operator"
	RequestKey  KeyType = "request"
)

type KeyStorage struct {
	*Storage
}

var keyStore *KeyStorage

func KeyStore() *KeyStorage {
	return keyStore
}

func NewKeyStore(filename string) error {
	s, err := NewStorage(filename)
	if err != nil {
		return err
	}

	keyStore = &KeyStorage{s}

	return nil
}

func (s *KeyStorage) GetPrivate(key KeyType) (*ecdsa.PrivateKey, error) {
	pk, ok := s.data[string(key)]
	if !ok {
		return &ecdsa.PrivateKey{}, fmt.Errorf("key not found: %s", key)
	}

	pkECDSA, err := crypto.HexToECDSA(pk)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}

	return pkECDSA, nil
}

// returns
func (s *KeyStorage) GetAddrs(key KeyType) (common.Address, address.Address, error) {
	pk, ok := s.data[string(key)]
	if !ok {
		return common.Address{}, address.Address{}, fmt.Errorf("key not found: %s", key)
	}

	pkECDSA, err := crypto.HexToECDSA(pk)
	if err != nil {
		return common.Address{}, address.Address{}, err
	}

	return DeriveAddrFromPk(pkECDSA)
}

func (s *KeyStorage) SetKey(key KeyType, pk *ecdsa.PrivateKey) error {
	pkStr := hexutil.Encode(crypto.FromECDSA(pk))[2:]
	err := s.Set(string(key), pkStr)

	return err
}

func DeriveAddrFromPkString(pk string) (common.Address, address.Address, error) {
	pkECDSA, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	return DeriveAddrFromPk(pkECDSA)
}

func DeriveAddrFromPk(pk *ecdsa.PrivateKey) (common.Address, address.Address, error) {
	evmAddr, err := fevm.DeriveAddressFromPk(pk)
	if err != nil {
		log.Fatal(err)
	}

	fevmAddr, err := ethtypes.ParseEthAddress(evmAddr.String())
	if err != nil {
		log.Fatal(err)
	}

	delegatedAddr, err := fevmAddr.ToFilecoinAddress()
	if err != nil {
		log.Fatal(err)
	}

	return evmAddr, delegatedAddr, nil
}

func DelegatedFromEthAddr(addr common.Address) (address.Address, error) {
	fevmAddr, err := ethtypes.ParseEthAddress(addr.String())
	if err != nil {
		return address.Address{}, err
	}

	return fevmAddr.ToFilecoinAddress()
}
