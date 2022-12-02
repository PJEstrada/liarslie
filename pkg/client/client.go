package client

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"io"
	"liarslie/pkg/agents"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var CurrentClient LiarsLieClient

type LiarsLieClient struct {
	AgentsFullNetwork agents.AgentsRegistry
}

// WaitTimeoutGetResult waits for a waitgroup to execute before the given timeout
func WaitTimeoutGetResult(wg *sync.WaitGroup, timeout time.Duration, c chan agents.LiarsLieMessageResult) *agents.LiarsLieMessageResult {
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

// spawnLiars spawns the given number of liar agents respecting the given ratio
func spawnLiars(pool agents.AgentsRegistry, maxValue int, liarRatio float32, numAgents int) {
	numLiars := int(math.Round(float64(float32(numAgents) * liarRatio)))
	for i := 0; i < numLiars; i++ {
		id := uuid.New()
		value := rand.Intn(maxValue + 1)
		agent := agents.NewLiarAgent(id, value)
		channel := agent.Start()
		pool[id] = &agents.AgentChannel{ChIn: channel, ID: id, Online: true}
	}
}

// spawnHonestAgents spawns the given number of honest agents respecting the given ratio
func spawnHonestAgents(pool agents.AgentsRegistry, value int, numAgents int, liarRatio float32) {
	numHonest := int(math.Round(float64((1 - liarRatio) * float32(numAgents))))
	for i := 0; i < numHonest; i++ {
		id := uuid.New()
		agent := agents.NewHonestAgent(id, value)
		channel := agent.Start()
		pool[id] = &agents.AgentChannel{ChIn: channel, ID: id, Online: true}
	}
}

// LaunchAgents creates liar and honest agents based on liar ratio and number given.
func LaunchAgents(value int, maxValue int, numAgents int, liarRatio float32) (*agents.AgentsRegistry, error) {
	if liarRatio < 0 || liarRatio > 1 {
		return nil, errors.New("Invalid liar ratio, must be between 0 and 1.")
	}
	pool := agents.AgentsRegistry{}

	spawnHonestAgents(pool, value, numAgents, liarRatio)
	spawnLiars(pool, maxValue, liarRatio, numAgents)

	return &pool, nil
}

// writeConfigFile writes agents IDs on a config file
func writeConfigFile(pool *agents.AgentsRegistry) error {
	f, err := os.Create(agents.ConfigFileName)
	if err != nil {
		return err
	}

	defer f.Close()
	for key, _ := range *pool {
		_, err2 := f.WriteString(fmt.Sprintf("%s\n", key))
		if err2 != nil {
			return err2
		}
	}
	return nil
}

// cleanFile cleans the text on the app.config file
func cleanFile() {
	f, err := os.Create(agents.ConfigFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	f.WriteString("")

}

// addNetworkPeersToAgents initializes peer network registry on the existing agents.
func addNetworkPeersToAgents(pool *agents.AgentsRegistry) error {
	// Populate Peers
	agentsData, err := agents.ReadConfigFile()
	if err != nil {
		return err
	}
	for _, agent := range *pool {
		wg := new(sync.WaitGroup)
		chout := make(chan agents.LiarsLieMessageResult)
		wg.Add(1)
		agent.ChIn <- agents.LiarsLieMessageRequest{
			MessageSetPeers: &agents.MessageSetPeers{
				Peers: agentsData,
			},
			ChOut: chout,
		}
		WaitTimeoutGetResult(wg, time.Second*3, chout)
	}
	return nil
}

type StringReader interface {
	io.Reader
	ReadString(delim byte) (string, error)
}

// displayShell shows the command prompt on terminal and starts waiting for user input.
func displayShell(rootCmd *cobra.Command, reader StringReader) {
	run := true
	for run {
		fmt.Print("liarslie>>")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}
		text = strings.Trim(text, "\r\n")

		if strings.Compare(text, "exit") == 0 {
			run = false
			log.Println("Exiting....")
		}

		cmdPieces := strings.Split(text, " ")

		command, args, err := rootCmd.Parent().Find(cmdPieces)
		if err != nil || command == nil {
			log.Printf("Unknown Command to execute : %s\n", text)
			continue
		}
		command.ParseFlags(args)
		command.Run(command, args)
		//command.Execute()
	}
}

// CreateNetwork launches agents and writes config file
func CreateNetwork(value int, maxValue int, numAgents int, liarRatio float32) (agents.AgentsRegistry, error) {
	pool, err := LaunchAgents(value, maxValue, numAgents, liarRatio)
	if err != nil {
		log.Print(fmt.Sprintf("Error launching agents: %s", err.Error()))
		return nil, err
	}
	err = writeConfigFile(pool)
	if err != nil {
		log.Fatal(err)
	}
	agents.SetAgentsNetwork(*pool)
	err = addNetworkPeersToAgents(pool)
	if err != nil {
		return nil, err
	}
	return *pool, nil
}

// StartClient initializes agents and starts command prompt.
func StartClient(rootCmd *cobra.Command, value int, maxValue int, numAgents int, liarRatio float32, reader StringReader) LiarsLieClient {

	pool, err := CreateNetwork(value, maxValue, numAgents, liarRatio)
	if err != nil {
		log.Fatal(err.Error())
	}
	CurrentClient = LiarsLieClient{
		AgentsFullNetwork: pool,
	}

	displayShell(rootCmd, reader)
	return CurrentClient
}

// StopClient stops all agents and exits program.
func StopClient(client LiarsLieClient) {

	for _, val := range client.AgentsFullNetwork {
		wg := new(sync.WaitGroup)
		chout := make(chan agents.LiarsLieMessageResult)
		wg.Add(1)
		msg := agents.LiarsLieMessageRequest{
			MessageSetOnline: &agents.MessageSetOnline{
				Online: false,
			},
			ChOut: chout,
		}
		val.ChIn <- msg
		msgResponse := WaitTimeoutGetResult(wg, time.Second*3, chout)
		if msgResponse != nil && msgResponse.MessageSetOnlineResult != nil {
			val.Online = false
		}
	}
	cleanFile()
}
