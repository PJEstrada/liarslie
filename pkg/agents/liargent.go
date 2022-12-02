package agents

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type LiarAgent struct {
	ID        uuid.UUID
	value     int
	Online    bool
	Peers     AgentsRegistry
	channelIn chan LiarsLieMessageRequest
}

// SetValue sets the agent value (network reconfig only)
func (a *LiarAgent) SetValue(v int, chout chan LiarsLieMessageResult) {
	chout <- LiarsLieMessageResult{
		MessageSetValueResult: &MessageSetValueResult{
			Value: v,
			ID:    a.ID,
		},
	}
}

// WaitTimeoutGetValue waits for a waitgroup to execute before the given timeout
func (a *LiarAgent) WaitTimeoutGetValue(wg *sync.WaitGroup, timeout time.Duration, c chan MessageGetValueResult) *MessageGetValueResult {
	go func() {
		wg.Wait()
	}()
	select {
	case res := <-c:
		return &res // completed normally
	case <-time.After(timeout):
		return nil // timed out
	}
}

// Start initializes agent and starts receiving messages
func (a *LiarAgent) Start() chan LiarsLieMessageRequest {
	a.channelIn = make(chan LiarsLieMessageRequest)
	go func() {
		for {
			select {
			case msg := <-a.channelIn:
				if !a.IsOnline() {
					msg.ChOut <- LiarsLieMessageResult{}
				}
				if msg.MessageGetValue != nil {
					a.GetValue(msg.MessageGetValue, msg.ChOut)
				} else if msg.MessageStop != nil {

				} else if msg.MessageSetPeers != nil {
					a.SetPeers(msg.Peers, msg.ChOut)
				} else if msg.MessageSetOnline != nil {
					a.setOnline(msg.MessageSetOnline.Online, msg.ChOut)
				} else if msg.MessageSetValue != nil {
					a.SetValue(msg.MessageSetValue.Value, msg.ChOut)
				}
			}
		}
	}()
	return a.channelIn

}
func NewLiarAgent(id uuid.UUID, value int) LiarAgent {
	return LiarAgent{
		ID:     id,
		value:  value,
		Online: true,
	}
}
func (a *LiarAgent) GetPeers() (AgentsRegistry, error) {
	return AgentsRegistry{}, nil
}

func (a *LiarAgent) SetPeers(peers AgentsRegistry, chout chan LiarsLieMessageResult) {
	chout <- LiarsLieMessageResult{
		MessageSetPeersResult: &MessageSetPeersResult{
			Peers: peers,
		},
	}
	return
}

func (a *LiarAgent) GetID() uuid.UUID {
	return a.ID
}

func (a *LiarAgent) GetValue(msg *MessageGetValue, chOut chan LiarsLieMessageResult) error {
	// Returning the wrong value always. I'm such a liar.
	chOut <- LiarsLieMessageResult{
		MessageGetValueResult: &MessageGetValueResult{
			ID:      msg.ID,
			AgentID: a.ID,
			Value:   a.value,
		},
	}
	return nil
}

func (a *LiarAgent) GetValueExpert(msg *MessageGetValue, chOut chan MessageGetValueResult) {
	chOut <- MessageGetValueResult{
		ID:    msg.ID,
		Value: a.value,
	}
}

func (a *LiarAgent) IsOnline() bool {
	return a.Online
}

// SetOnline sets online value
func (a *LiarAgent) setOnline(v bool, chOut chan LiarsLieMessageResult) {
	a.Online = v
	chOut <- LiarsLieMessageResult{
		MessageSetOnlineResult: &MessageSetOnlineResult{
			Online: a.Online,
		},
	}
	defer close(chOut)
}
