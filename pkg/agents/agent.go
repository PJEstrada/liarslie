package agents

import (
	"fmt"
	"github.com/google/uuid"
	"liarslie/pkg/consensus"
	"sync"
	"time"
)

type Agent struct {
	ID               uuid.UUID
	value            int
	Online           bool
	MsgNetworkValues sync.Map
	peers            *AgentsRegistry
	channelIn        chan LiarsLieMessageRequest
}

func (a *Agent) Start() chan LiarsLieMessageRequest {
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

// SetValue sets the agent value (network reconfig only)
func (a *Agent) SetValue(v int, chout chan LiarsLieMessageResult) {
	a.value = v
	chout <- LiarsLieMessageResult{
		MessageSetValueResult: &MessageSetValueResult{
			Value: a.value,
			ID:    a.ID,
		},
	}
}

// WaitTimeoutGetValue waits for a waitgroup to execute before the given timeout
func (a *Agent) WaitTimeoutGetValue(wg *sync.WaitGroup, timeout time.Duration, c chan LiarsLieMessageResult) *LiarsLieMessageResult {
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

// NewHonestAgent creates a new agent object.
func NewHonestAgent(ID uuid.UUID, Value int) Agent {
	return Agent{
		ID:               ID,
		value:            Value,
		MsgNetworkValues: sync.Map{},
		Online:           true,
	}
}

// SetPeers sets the peer agents value
func (a *Agent) SetPeers(peers AgentsRegistry, chout chan LiarsLieMessageResult) {
	a.peers = &peers
	chout <- LiarsLieMessageResult{
		MessageSetPeersResult: &MessageSetPeersResult{
			Peers: peers,
		},
	}
}

// GetPeers gets the peer agents of the agent.
func (a *Agent) GetPeers() (AgentsRegistry, error) {
	return *a.peers, nil

}

// GetID gets the id of the agent.
func (a *Agent) GetID() uuid.UUID {
	return a.ID
}

// GetValue gets the value of the agent. Can be without asking other agents for their value or asking all peers for its value.
func (a *Agent) GetValue(msg *MessageGetValue, chOut chan LiarsLieMessageResult) error {
	if msg.WithPeers {
		err := a.getValueExpert(*msg, chOut)
		if err != nil {
			return err
		}

	} else {
		chOut <- LiarsLieMessageResult{
			MessageGetValueResult: &MessageGetValueResult{
				ID:      msg.ID,
				AgentID: a.ID,
				Value:   a.value,
			},
		}
	}
	// We are assuming this a unique channel for this sender, thus the sender can safely close it.
	defer close(chOut)
	return nil

}

// copyMap copies the given map
func copyMap(m map[uuid.UUID]int) map[uuid.UUID]int {
	res := map[uuid.UUID]int{}
	for key, val := range m {
		res[key] = val
	}
	return res
}

// getValueForMsg gets a cached value from the agent
func (a *Agent) getValueForMsg(msgID uuid.UUID) (int, bool) {
	val, ok := a.MsgNetworkValues.Load(msgID)
	var result int
	if ok {
		result = val.(int)
	}
	return result, ok
}

// setValueForMessage sets the cached value for a given message on the agent.
func (a *Agent) setValueForMessage(msgID uuid.UUID, value int) {
	a.MsgNetworkValues.Store(msgID, value)
}

// queryPeers asks all peer agents for its value and determines majority
func (a *Agent) queryPeers(agentsNet AgentsRegistry, msg MessageGetValue) []int {
	var values []int
	values = append(values, a.value)
	for _, agent := range agentsNet {
		fmt.Println("doing age", agent.ChIn)
		chOutPeer := make(chan LiarsLieMessageResult)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		agent.ChIn <- LiarsLieMessageRequest{
			MessageGetValue: &MessageGetValue{
				ID:        msg.ID,
				WithPeers: false,
			},
			ChOut: chOutPeer,
		}
		msgResponse := a.WaitTimeoutGetValue(wg, time.Second*3, chOutPeer)
		if msgResponse == nil {
			// Timeout case: no value apperessspnded.
			continue
		}
		// Successful Response from Agent
		if msgResponse.MessageGetValueResult != nil {
			values = append(values, msgResponse.MessageGetValueResult.Value)
		}

	}
	return values
}

// getValueExpert gets a value for playing expert mode.
func (a *Agent) getValueExpert(msg MessageGetValue, chOut chan LiarsLieMessageResult) error {
	// Case of value already processed by Node
	existingValue, found := a.getValueForMsg(msg.ID)
	if found {
		chOut <- LiarsLieMessageResult{
			MessageGetValueResult: &MessageGetValueResult{
				ID:      msg.ID,
				AgentID: a.ID,
				Value:   existingValue,
			},
		}
		return nil
	}
	// New message processing
	agentsNet, err := a.GetPeers()
	if err != nil {
		return err
	}

	// Query each peer agent for its value
	values := a.queryPeers(agentsNet, msg)

	// Determine real network value based on all peer values obtained.
	networkVal := consensus.FindMajorityValue(values)
	a.setValueForMessage(msg.ID, networkVal)
	// Set Latest Message Value
	a.value = networkVal
	// Send Response on Channel
	chOut <- LiarsLieMessageResult{
		MessageGetValueResult: &MessageGetValueResult{
			ID:      msg.ID,
			AgentID: a.ID,
			Value:   networkVal,
		},
	}

	return nil
}

// SetOnline sets online value
func (a *Agent) setOnline(v bool, chOut chan LiarsLieMessageResult) {
	a.Online = v
	chOut <- LiarsLieMessageResult{
		MessageSetOnlineResult: &MessageSetOnlineResult{
			Online: a.Online,
		},
	}
	defer close(chOut)
}

// IsOnline checks if agent is online
func (a *Agent) IsOnline() bool {
	return a.Online
}
