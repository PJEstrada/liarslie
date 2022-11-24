package agents

import (
	"github.com/google/uuid"
	"sync"
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

func (a *LiarAgent) GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, wg *sync.WaitGroup) {
	chOut <- MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
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
