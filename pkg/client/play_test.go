package client

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"liarslie/pkg/agents"
	"testing"
)

type ConfigPlayTestSuite struct {
	suite.Suite
}

func TestPlayTestSuite(t *testing.T) {
	suite.Run(t, &ConfigExtendTestSuite{})
}

func (s *ConfigExtendTestSuite) TestQueryAgents() {
	pool, err := CreateNetwork(7, 100, 10, 0.4)
	s.Nil(err)
	val := queryAgents(pool, true, true, 0.6)
	s.Equal(val, 7)

	// Query without consensus ratio
	val = queryAgents(pool, false, true, 0.35)
	s.Equal(val, 7)

	// Query without consensus ratio
	val = queryAgents(pool, false, true, 0.8)
	s.NotEqual(val, 7)
}

func (s *ConfigExtendTestSuite) TestGetOnlineAgents() {
	pool := agents.AgentsRegistry{
		uuid.New(): &agents.AgentChannel{
			Online: true,
		},
		uuid.New(): &agents.AgentChannel{
			Online: true,
		},
		uuid.New(): &agents.AgentChannel{
			Online: false,
		},
	}
	result := getOnlineAgents(pool)

	s.Equal(len(result), 2)
}

func (s *ConfigExtendTestSuite) TestPlayStandard() {
	pool, err := CreateNetwork(7, 100, 10, 0.4)
	s.Nil(err)
	result := PlayStandard(pool)

	s.Equal(result, 7)
}

func (s *ConfigExtendTestSuite) TestPlayExpert() {
	pool, err := CreateNetwork(7, 100, 10, 0.4)
	s.Nil(err)
	result, err := PlayExpert(pool, 6, 0.2)
	s.Nil(err)
	s.Equal(result, 7)
	// Fail case
	pool, err = CreateNetwork(7, 100, 10, 1)
	result, err = PlayExpert(pool, 2, 0.1)
	s.Nil(err)
	s.NotEqual(result, 7)

	// Error case
	pool, err = CreateNetwork(7, 100, 15, 1)
	result, err = PlayExpert(pool, 20, 0.1)
	s.NotNil(err)
	s.EqualError(err, "Not enough online agents to play in expert mode. Want 20 have 15")
}
