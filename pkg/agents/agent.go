package agents

import (
	"github.com/google/uuid"
)

type Agent struct {
	ID      uuid.UUID
	value   int
	conn    chan LiarsLieMessageRequest
	connOut chan LiarsLieMessageResult
	Online  bool
}

func NewHonestAgent(ID uuid.UUID, Value int) Agent {
	return Agent{
		ID:    ID,
		value: Value,
	}
}
func (a *Agent) GetValue(msg *MessageGetValue) *MessageGetValueResult {
	return &MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
	}
}

func (a *Agent) SetOnline(v bool) {
	a.Online = v
}
func (a *Agent) StartProcessing() {
	for msg := range a.conn {
		if msg.MessageGetValue != nil {
			result := a.GetValue(msg.MessageGetValue)
			a.connOut <- LiarsLieMessageResult{
				MessageGetValueResult: result,
			}
		}
	}
}

func (a *Agent) IsOnline() bool {
	return a.Online
}
