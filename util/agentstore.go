package util

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
