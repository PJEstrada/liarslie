package client

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"liarslie/pkg/agents"
	"os"
	"testing"
)

type ConfigExtendTestSuite struct {
	suite.Suite
}

func TestExtendTestSuite(t *testing.T) {
	suite.Run(t, &ConfigExtendTestSuite{})
}

func (s *ConfigExtendTestSuite) TestExtendNetwork() {
	id := uuid.New()
	id2 := uuid.New()
	agent1 := agents.NewHonestAgent(id, 5)
	agent2 := agents.NewHonestAgent(id2, 5)
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
	//agents.SetAgentsNetwork(originalNetwork)
	//CurrentClient.AgentsFullNetwork = originalNetwork
	//extendedPool, err := ExtendNetwork(70, 100, 12, 0.5)
	//s.Equal(len(extendedPool), 12)
	//s.Nil(err)

	//// Existing Pool
	writeConfigFile(&originalNetwork)
	agents.SetAgentsNetwork(originalNetwork)

	CurrentClient.AgentsFullNetwork = originalNetwork
	extendedPool, err := ExtendNetwork(70, 100, 12, 0.5)
	s.Equal(len(extendedPool), 12)
	s.Nil(err)
	os.Remove(agents.ConfigFileName)

	//// Existing Pool Error
	//agents.SetAgentsNetwork(originalNetwork)
	//writeConfigFile(&originalNetwork)
	//CurrentClient.AgentsFullNetwork = originalNetwork
	//extendedPool, err = ExtendNetwork(70, 100, 1, 0.5)
	//s.Nil(extendedPool)
	//s.NotNil(err)
	//os.Remove(agents.ConfigFileName)
}
