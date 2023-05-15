package util

var agentStore *Storage

func AgentStore() *Storage {
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

	agentStore = s

	return nil
}
