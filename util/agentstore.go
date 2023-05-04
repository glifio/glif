package util

var agentStore *Storage

func AgentStore() *Storage {
	return agentStore
}

func NewAgentStore(filename string) error {
	s, err := NewStorage(filename)
	if err != nil {
		return err
	}

	agentStore = s

	return nil
}
