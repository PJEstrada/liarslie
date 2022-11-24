package agents

import (
	"github.com/google/uuid"
	"liarslie/pkg/client"
)

type Agent struct {
	ID     uuid.UUID
	value  int
	Online bool
}

func (a *Agent) GetPeers() *AgentsRegistry {
	result := ReadConfigFile()
	return result

}

func NewHonestAgent(ID uuid.UUID, Value int) Agent {
	return Agent{
		ID:    ID,
		value: Value,
	}
}
func (a *Agent) GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult) {
	chOut <- MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
	}
}

func (a *Agent) GetValueExpert(msg *MessageGetValue, chOut chan MessageGetValueResult) {
	agentsNet := a.GetPeers()
	onlineAgents := 0
	chAgents := make(chan MessageGetValueResult)
	for _, agent := range *agentsNet {
		agent.GetValue(msg, chAgents)
		if agent.IsOnline() {
			onlineAgents += 1
		}
	}
	var values []int
	for msg := range chAgents {
		values = append(values, msg.Value)
	}
	maxVal := client.FindMajorityValue(values, onlineAgents)
	chOut <- MessageGetValueResult{
		ID:    msg.ID,
		Value: maxVal,
	}
}

func (a *Agent) SetOnline(v bool) {
	a.Online = v
}

func (a *Agent) IsOnline() bool {
	return a.Online
}
