package agents

import (
	"github.com/google/uuid"
)

type LiarAgent struct {
	ID      uuid.UUID
	value   int
	conn    chan LiarsLieMessageRequest
	connOut chan LiarsLieMessageResult
	Online  bool
}

func NewLiarAgent(id uuid.UUID, value int) LiarAgent {
	return LiarAgent{
		ID:    id,
		value: value,
	}
}

func (a *LiarAgent) GetValue(msg *MessageGetValue) *MessageGetValueResult {
	return &MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
	}
}

func (a *LiarAgent) SetOnline(v bool) {
	a.Online = v
}

func (a *LiarAgent) StartProcessing() {
	a.Online = true
	for msg := range a.conn {
		if msg.MessageGetValue != nil {
			result := a.GetValue(msg.MessageGetValue)
			a.connOut <- LiarsLieMessageResult{
				MessageGetValueResult: result,
			}
		}
	}
}

func (a *LiarAgent) IsOnline() bool {
	return a.Online
}
