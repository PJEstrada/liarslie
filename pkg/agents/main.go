package agents

import "github.com/google/uuid"

type AgentsRegistry map[uuid.UUID]AgentBehaviour

var agentsNetwork AgentsRegistry

func SetAgentsNetwork(reg AgentsRegistry) {
	agentsNetwork = reg
}

func GetAgentsNetwork() AgentsRegistry {
	return agentsNetwork
}
