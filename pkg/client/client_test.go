package client

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"liarslie/pkg/agents"
	"os"
	"strings"
	"testing"
)

type ConfigClientTestSuite struct {
	suite.Suite
}

func TestClientNetworkTestSuite(t *testing.T) {
	agents.ConfigFileName = "unittestNetwork.config"
	suite.Run(t, &ConfigClientTestSuite{})
	os.Remove("testwrite")
}

func (s *ConfigClientTestSuite) TestSpawnLiars() {
	pool := agents.AgentsRegistry{}
	spawnLiars(pool, 100, 0.5, 100)
	s.Equal(len(pool), 50)

	spawnLiars(pool, 100, 0.33, 100)
	s.Equal(len(pool), 83)
}

func (s *ConfigClientTestSuite) TestSpawnHonestAgents() {
	pool := agents.AgentsRegistry{}
	spawnHonestAgents(pool, 2, 100, 0.6)
	s.Equal(len(pool), 40)

	spawnHonestAgents(pool, 2, 100, 0.2)
	s.Equal(len(pool), 120)
}

func (s *ConfigClientTestSuite) TestLaunchAgents() {
	pool, err := LaunchAgents(5, 100, 100, 0.4)
	s.Equal(len(*pool), 100)
	s.Nil(err)
	// Err case
	pool, err = LaunchAgents(5, 100, 100, -0.4)
	s.NotNil(err)
	s.EqualError(err, "Invalid liar ratio, must be between 0 and 1.")
	s.Nil(pool)
	pool, err = LaunchAgents(5, 100, 100, 1.5)
	s.NotNil(err)
	s.EqualError(err, "Invalid liar ratio, must be between 0 and 1.")
	s.Nil(pool)
}

func (s *ConfigClientTestSuite) TestWriteConfigFile() {
	agents.ConfigFileName = "testwrite"
	id := uuid.New()
	pool := agents.AgentsRegistry{
		id: &agents.AgentChannel{
			ID: id,
		},
	}
	err := writeConfigFile(&pool)
	s.Nil(err)
	f, err := os.Open(agents.ConfigFileName)
	s.Nil(err)
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()
	line := fileScanner.Text()
	s.Equal(id.String(), line)
	os.Remove("testwrite")
}

func (s *ConfigClientTestSuite) TestCleanFile() {
	agents.ConfigFileName = "testwrite"
	id := uuid.New()
	agent := agents.NewHonestAgent(id, 55)
	chIn := agent.Start()
	pool := agents.AgentsRegistry{
		id: &agents.AgentChannel{
			ID:   id,
			ChIn: chIn,
		},
	}
	err := writeConfigFile(&pool)
	s.Nil(err)
	cleanFile()

	f, err := os.Open(agents.ConfigFileName)
	s.Nil(err)
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()
	line := fileScanner.Text()
	s.Equal("", line)

	f.Close()
	os.Remove(agents.ConfigFileName)

}

func (s *ConfigClientTestSuite) TestAddNetworkPeersToAgents() {
	agents.ConfigFileName = "testwrite"
	id := uuid.New()
	agent := agents.NewHonestAgent(id, 44)
	chIn := agent.Start()
	pool := agents.AgentsRegistry{
		id: &agents.AgentChannel{
			ID:   id,
			ChIn: chIn,
		},
	}
	agents.SetAgentsNetwork(pool)
	err := writeConfigFile(&pool)

	s.Nil(err)
	errAdd := addNetworkPeersToAgents(&pool)
	s.Nil(errAdd)
	result := agents.GetAgentsNetwork()
	s.Equal(pool, result)
	// Error case
	agents.ConfigFileName = "non existent file"
	errAdd = addNetworkPeersToAgents(&pool)
	s.NotNil(errAdd)
	s.EqualError(errAdd, "open non existent file: no such file or directory")
}

func (s *ConfigClientTestSuite) TestCreateNetwork() {
	fmt.Print("QWEQWEQWEWQE")
	pool, err := CreateNetwork(1, 10, 26, 0.5)
	s.Nil(err)
	s.Equal(len(pool), 26)
	fmt.Print("2222222222")
	// Error case invalid liar ratio
	pool, err = CreateNetwork(1, 10, 26, -0.5)
	s.NotNil(err)
	s.Nil(pool)
	s.EqualError(err, "Invalid liar ratio, must be between 0 and 1.")
	os.Remove(agents.ConfigFileName)
}

type MockReader struct {
	*bufio.Reader
	ReadStringFunc func(delim byte) (string, error)
}

func (r *MockReader) ReadString(delim byte) (string, error) {
	if r.ReadStringFunc != nil {
		return r.ReadStringFunc(delim)
	}
	return r.Reader.ReadString(delim)
}

func (s *ConfigClientTestSuite) TestStartClient() {
	cmd := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	cmd2 := cobra.Command{}
	cmd.AddCommand(&cmd2)
	mr := &MockReader{
		Reader: bufio.NewReader(strings.NewReader("exit")),
		ReadStringFunc: func(delim byte) (string, error) {
			return "exit", nil
		},
	}
	client := StartClient(&cmd2, 1, 100, 26, 0.5, mr)
	s.Equal(client, CurrentClient)
	s.Equal(len(client.AgentsFullNetwork), 26)
}

func (s *ConfigClientTestSuite) TestStopClient() {
	id := uuid.New()
	agent := agents.NewHonestAgent(id, 1)
	pool := agents.AgentsRegistry{
		id: &agents.AgentChannel{
			ID:     id,
			Online: true,
			ChIn:   agent.Start(),
		},
	}
	client := LiarsLieClient{
		AgentsFullNetwork: pool,
	}
	StopClient(client)
	for _, ag := range client.AgentsFullNetwork {
		s.False(ag.Online)
	}
}
