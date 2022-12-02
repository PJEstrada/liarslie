package agents

import "github.com/google/uuid"

type AgentChannel struct {
	ChIn   chan LiarsLieMessageRequest
	ID     uuid.UUID
	Online bool
}
type AgentsRegistry map[uuid.UUID]*AgentChannel

var agentsNetwork AgentsRegistry

func SetAgentsNetwork(reg AgentsRegistry) {
	agentsNetwork = reg
}

func GetAgentsNetwork() AgentsRegistry {
	return agentsNetwork
}
