package util

import (
	"fmt"
	"strconv"
)

const (
	OperatorKeyFunded KeyType = "opkeyf"
	OwnerKeyFunded    KeyType = "ownkeyf"
)

type AgentStorage struct {
	*Storage
}

var agentStore *AgentStorage

func AgentStore() *AgentStorage {
	return agentStore
}

func NewAgentStore(filename string) error {
	agentDefault := map[string]string{
		"id":      "",
		"address": "",
		"tx":      "",
	}

	s, err := NewStorage(filename, agentDefault)
	if err != nil {
		return err
	}

	agentStore = &AgentStorage{s}

	return nil
}

func (a *AgentStorage) IsFunded(keytype KeyType, key string) (bool, error) {
	switch keytype {
	case OperatorKeyFunded, OwnerKeyFunded:
		f, ok := a.data[mapkey(keytype, key)]
		if !ok {
			return false, fmt.Errorf("key not found: %s", key)
		}

		return strconv.ParseBool(f)
	default:
		return false, fmt.Errorf("not supported key type for funded operation")
	}
}

func (a *AgentStorage) SetFunded(keytype KeyType, key string, funded bool) error {
	switch keytype {
	case OperatorKeyFunded, OwnerKeyFunded:
		return a.Set(mapkey(keytype, key), strconv.FormatBool(funded))
	default:
		return fmt.Errorf("not supported key type for funded operation")
	}
}

func mapkey(keytype KeyType, key string) string {
	return fmt.Sprintf("%s-%s", string(keytype), key)
}
