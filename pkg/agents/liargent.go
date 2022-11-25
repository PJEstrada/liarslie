package agents

import (
	"github.com/google/uuid"
)

type LiarAgent struct {
	ID     uuid.UUID
	value  int
	Online bool
	Peers  AgentsRegistry
}

func NewLiarAgent(id uuid.UUID, value int) LiarAgent {
	return LiarAgent{
		ID:    id,
		value: value,
	}
}
func (a *LiarAgent) GetPeers() AgentsRegistry {
	return AgentsRegistry{}
}

func (a *LiarAgent) GetID() uuid.UUID {
	return a.ID
}

func (a *LiarAgent) GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, withPeers bool) {
	chOut <- MessageGetValueResult{
		ID:      msg.ID,
		AgentID: a.ID,
		Value:   a.value,
	}
}

func (a *LiarAgent) GetValueExpert(msg *MessageGetValue, chOut chan MessageGetValueResult) {
	chOut <- MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
	}
}

func (a *LiarAgent) SetOnline(v bool) {
	a.Online = v
}

func (a *LiarAgent) IsOnline() bool {
	return a.Online
}
