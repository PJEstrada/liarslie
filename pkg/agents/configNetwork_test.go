package agents

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigNetworkTestSuite struct {
	suite.Suite
}

func TestConfigNetworkTestSuite(t *testing.T) {
	ConfigFileName = "unittestNetwork.config"
	suite.Run(t, &AgentTestSuite{})
	ConfigFileName = "app.config"
}

func (s *AgentTestSuite) TestReadConfigFile() {
	f, err := os.Create(ConfigFileName)
	if err != nil {
		s.Fail(err.Error())
	}
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New()}
	network := AgentsRegistry{
		ids[0]: &AgentChannel{ID: ids[0], Online: true, ChIn: make(chan LiarsLieMessageRequest)},
		ids[1]: &AgentChannel{ID: ids[1], Online: true, ChIn: make(chan LiarsLieMessageRequest)},
		ids[2]: &AgentChannel{ID: ids[2], Online: true, ChIn: make(chan LiarsLieMessageRequest)},
		ids[3]: &AgentChannel{ID: ids[3], Online: true, ChIn: make(chan LiarsLieMessageRequest)},
	}
	SetAgentsNetwork(network)
	for _, id := range ids {
		f.WriteString(fmt.Sprintf("%s\n", id.String()))
	}
	res, err := ReadConfigFile()
	s.Equal(len(res), len(network))
	s.Nil(err)
	for _, id := range ids {
		s.Equal(res[id].ID, id)
	}
	f.Close()
	os.Remove(ConfigFileName)

	// Test error non existent file
	ConfigFileName = "nofile"
	res, err = ReadConfigFile()
	s.NotNil(err)
	s.EqualError(err, "open nofile: no such file or directory")
}
