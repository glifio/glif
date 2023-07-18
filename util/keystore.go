package util

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type KeyStorage struct {
	*Storage
}

var keyStore *keystore.KeyStore

func KeyStore() *keystore.KeyStore {
	return keyStore
}

func NewKeyStore(keydir string) {
	keyStore = keystore.NewKeyStore(
		keydir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
}
