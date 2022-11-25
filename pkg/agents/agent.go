package agents

import (
	"github.com/google/uuid"
	"liarslie/pkg/consensus"
)

type Agent struct {
	ID     uuid.UUID
	value  int
	Online bool
}

// FindMajorityValue finds the number that is repeated the most among the list of integers in values param.
func FindMajorityValue(values []int) int {
	valuesCount := map[int]int{}
	for _, val := range values {
		valuesCount[val] += 1
	}
	percentages := map[int]float64{}
	max := -1.0
	maxVal := -1
	for key, val := range valuesCount {
		percentages[key] += float64(float64(val) / float64(len(values)))
		if percentages[key] > float64(max) {
			max = percentages[key]
			maxVal = key
		}
	}
	return maxVal
}

func NewHonestAgent(ID uuid.UUID, Value int) Agent {
	return Agent{
		ID:    ID,
		value: Value,
	}
}

func (a *Agent) GetPeers() AgentsRegistry {
	result := ReadConfigFile()
	return result

}
func (a *Agent) GetID() uuid.UUID {
	return a.ID
}

func (a *Agent) GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, withPeers bool) {
	if withPeers {
		a.getValueExpert(msg, chOut)

	} else {
		chOut <- MessageGetValueResult{
			ID:      msg.ID,
			AgentID: a.ID,
			Value:   a.value,
		}
	}

}
func (a *Agent) getValueExpert(msg *MessageGetValue, chOut chan MessageGetValueResult) {
	agentsNet := a.GetPeers()
	chAgents := make(chan MessageGetValueResult)
	msg.KnownValues[a.ID] = a.value
	for _, agent := range agentsNet {
		if _, ok := msg.KnownValues[a.ID]; !ok {
			agent.GetValue(msg, chAgents, true)
		}
	}
	var values []int
	for msgResponse := range chAgents {
		values = append(values, msgResponse.Value)
		msg.KnownValues[msgResponse.AgentID] = msgResponse.Value
	}
	maxVal := consensus.FindMajorityValue(values)
	chOut <- MessageGetValueResult{
		ID:      msg.ID,
		AgentID: a.ID,
		Value:   maxVal,
	}
}

func (a *Agent) SetOnline(v bool) {
	a.Online = v
}

func (a *Agent) IsOnline() bool {
	return a.Online
}
