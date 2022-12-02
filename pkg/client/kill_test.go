package client

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"liarslie/pkg/agents"
	"testing"
)

type ConfigKillTestSuite struct {
	suite.Suite
}

func TestKillTestSuite(t *testing.T) {
	suite.Run(t, &ConfigExtendTestSuite{})
}

func (s *ConfigExtendTestSuite) TestKill() {
	id := uuid.New()
	id2 := uuid.New()
	agent1 := agents.NewHonestAgent(id, 3)
	agent2 := agents.NewHonestAgent(id2, 3)
	originalNetwork := agents.AgentsRegistry{
		id: &agents.AgentChannel{
			ID:     id,
			Online: true,
			ChIn:   agent1.Start(),
		},
		id2: &agents.AgentChannel{
			ID:     id2,
			Online: true,
			ChIn:   agent2.Start(),
		},
	}
	CurrentClient.AgentsFullNetwork = originalNetwork
	KillAgent(id)
	s.False(originalNetwork[id].Online)
	s.True(originalNetwork[id2].Online)

	KillAgent(id2)
	s.False(originalNetwork[id].Online)
	s.False(originalNetwork[id2].Online)
}
