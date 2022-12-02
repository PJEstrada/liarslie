package agents

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type AgentTestSuite struct {
	suite.Suite
}

func TestAgentTestSuite(t *testing.T) {
	ConfigFileName = "unittest.config"
	suite.Run(t, &AgentTestSuite{})
	ConfigFileName = "app.config"
}

func (s *AgentTestSuite) TestNewAgent() {
	id := uuid.New()
	value := 25
	agent := NewHonestAgent(id, value)
	s.Equal(value, agent.value)
	s.Equal(id, agent.GetID())
}

func (s *AgentTestSuite) TestSetValue() {
	id := uuid.New()
	value := 25
	newVal := 7
	agent := NewHonestAgent(id, value)
	agent.Start()
	chout := make(chan LiarsLieMessageResult)
	go agent.SetValue(newVal, chout)
	result := <-chout
	s.NotNil(result)
	s.NotNil(result.MessageSetValueResult)
	s.Equal(result.MessageSetValueResult.Value, newVal)
	s.Equal(agent.value, newVal)
}

func (s *AgentTestSuite) TestWaitTimeoutGetValue() {
	wg := new(sync.WaitGroup)
	ch := make(chan LiarsLieMessageResult)
	agent := NewHonestAgent(uuid.New(), 20)
	// No timeout
	go func() {
		ch <- LiarsLieMessageResult{
			MessageGetValueResult: &MessageGetValueResult{
				Value: 20,
				ID:    uuid.New(),
			},
		}
	}()
	msgValue := agent.WaitTimeoutGetValue(wg, 2*time.Second, ch)
	s.NotNil(msgValue)
	s.NotNil(msgValue.MessageGetValueResult)
	s.Equal(msgValue.MessageGetValueResult.Value, 20)
	// With timeout
	go func() {
		time.Sleep(5 * time.Second)
		ch <- LiarsLieMessageResult{
			MessageGetValueResult: &MessageGetValueResult{
				Value: 20,
				ID:    uuid.New(),
			},
		}
	}()
	msgValue = agent.WaitTimeoutGetValue(wg, 1*time.Second, ch)
	s.Nil(msgValue)
}

func (s *AgentTestSuite) TestSetPeers() {
	agent := NewHonestAgent(uuid.New(), 20)
	chIn := agent.Start()
	peers := AgentsRegistry{}
	chOut := make(chan LiarsLieMessageResult)
	msg := LiarsLieMessageRequest{
		MessageSetPeers: &MessageSetPeers{
			Peers: peers,
		},
		ChOut: chOut,
	}
	wg := new(sync.WaitGroup)

	chIn <- msg
	msgResult := agent.WaitTimeoutGetValue(wg, time.Second*3, chOut)
	s.NotNil(msgResult)
	s.NotNil(msgResult.MessageSetPeersResult)
	res, err := agent.GetPeers()
	s.Equal(res, peers)
	s.Nil(err)
	s.Equal(*agent.peers, peers)
}
func (s *AgentTestSuite) TestGetPeers() {
	agent := NewHonestAgent(uuid.New(), 20)
	peers := AgentsRegistry{}
	agent.peers = &peers
	peers2, err := agent.GetPeers()
	if err != nil {
		s.Fail(err.Error())
	}
	s.Equal(peers2, peers)
}

func (s *AgentTestSuite) TestGetValue() {
	agent := NewHonestAgent(uuid.New(), 20)
	agent.peers = &AgentsRegistry{}

	msg := LiarsLieMessageRequest{
		MessageGetValue: &MessageGetValue{ID: uuid.New()},
	}
	chout := make(chan LiarsLieMessageResult)
	go agent.GetValue(msg.MessageGetValue, chout)
	res := <-chout
	s.Equal(res.MessageGetValueResult.Value, 20)

	// With peers
	chout2 := make(chan LiarsLieMessageResult)
	go agent.GetValue(msg.MessageGetValue, chout2)
	res = <-chout2
	s.Equal(res.MessageGetValueResult.Value, 20)
}

func (s *AgentTestSuite) TestGetValueForMsg() {
	agent := NewHonestAgent(uuid.New(), 20)
	id := uuid.New()
	id2 := uuid.New()
	agent.setValueForMessage(id, 75)
	value, found := agent.getValueForMsg(id)
	s.Equal(value, 75)
	s.True(found)

	value, found = agent.getValueForMsg(id2)
	s.Equal(value, 0)
	s.False(found)
}

func (s *AgentTestSuite) TestSetOnline() {
	agent := NewHonestAgent(uuid.New(), 20)
	chIn := agent.Start()
	chout := make(chan LiarsLieMessageResult)
	msg := LiarsLieMessageRequest{
		MessageSetOnline: &MessageSetOnline{
			Online: true,
		},
		ChOut: chout,
	}

	wg := new(sync.WaitGroup)
	chIn <- msg
	msgResponse := agent.WaitTimeoutGetValue(wg, time.Second*3, chout)
	s.NotNil(msgResponse)
	s.NotNil(msgResponse.MessageSetOnlineResult)
	s.True(msgResponse.MessageSetOnlineResult.Online)
	s.True(agent.IsOnline())

	// False case
	chout = make(chan LiarsLieMessageResult)
	wg = new(sync.WaitGroup)
	msg = LiarsLieMessageRequest{
		MessageSetOnline: &MessageSetOnline{
			Online: false,
		},
		ChOut: chout,
	}
	chIn <- msg
	msgResponse = agent.WaitTimeoutGetValue(wg, time.Second*3, chout)
	s.NotNil(msgResponse)
	s.NotNil(msgResponse.MessageSetOnlineResult)
	s.False(msgResponse.MessageSetOnlineResult.Online)
	s.False(agent.IsOnline())
}

func (s *AgentTestSuite) TestCopyMap() {
	testVal := map[uuid.UUID]int{
		uuid.New(): 1,
		uuid.New(): 2,
		uuid.New(): 3,
	}
	res := copyMap(testVal)
	s.Equal(res, testVal)
}

func (s *AgentTestSuite) TestQueryPeers() {
	agent := NewHonestAgent(uuid.New(), 20)
	agentsPeers := []Agent{
		NewHonestAgent(uuid.New(), 20),
		NewHonestAgent(uuid.New(), 20),
		NewHonestAgent(uuid.New(), 20),
		NewHonestAgent(uuid.New(), 20),
		NewHonestAgent(uuid.New(), 20),
	}
	net := AgentsRegistry{
		agentsPeers[0].ID: &AgentChannel{ID: agentsPeers[0].ID, Online: true, ChIn: agentsPeers[0].Start()},
		agentsPeers[1].ID: &AgentChannel{ID: agentsPeers[1].ID, Online: true, ChIn: agentsPeers[1].Start()},
		agentsPeers[2].ID: &AgentChannel{ID: agentsPeers[2].ID, Online: true, ChIn: agentsPeers[2].Start()},
		agentsPeers[3].ID: &AgentChannel{ID: agentsPeers[3].ID, Online: true, ChIn: agentsPeers[3].Start()},
	}
	result := agent.queryPeers(net, MessageGetValue{ID: uuid.New()})
	s.Equal(5, len(result))
	for _, val := range result {
		s.Equal(val, 20)
	}
}
