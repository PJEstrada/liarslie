package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"liarslie/pkg/agents"
	"log"
	"math/rand"
	"os"
	"strings"
)

type LiarsLieClient struct {
	StartedAgents bool
	Agents        agents.AgentsPool
}

func spawnLiars(pool agents.AgentsPool, maxValue int, liarRatio float32) {
	numLiars := int(1 * liarRatio)
	for i := 0; i < numLiars; i++ {
		id := uuid.New()
		value := rand.Intn(maxValue + 1)
		agent := agents.NewLiarAgent(id, value)
		pool[id] = &agent
	}
}

func spawnHonestAgents(pool agents.AgentsPool, value int, numAgents int, liarRatio float32) {
	numHonest := int((1 - liarRatio) * float32(numAgents))
	for i := 0; i < numHonest; i++ {
		id := uuid.New()
		agent := agents.NewHonestAgent(id, value)
		pool[id] = &agent
	}
}

func LaunchAgents(value int, maxValue int, numAgents int, liarRatio float32) (*agents.AgentsPool, error) {
	if liarRatio < 0 || liarRatio > 1 {
		return nil, errors.New("Invalid liar ratio, must be between 0 and 1.")
	}
	pool := agents.AgentsPool{}

	spawnHonestAgents(pool, value, numAgents, liarRatio)
	spawnLiars(pool, maxValue, liarRatio)

	return &pool, nil
}
func writeConfigFile(pool *agents.AgentsPool) {
	f, err := os.Create("app.config")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	for key, agent := range *pool {
		_, err2 := f.WriteString(fmt.Sprintf("%s %v\n", key, agent.IsOnline()))
		if err2 != nil {
			log.Fatal(err)
		}
	}
}
func StartClient(rootCmd *cobra.Command, value int, maxValue int, numAgents int, liarRatio float32) {

	reader := bufio.NewReader(os.Stdin)
	pool, err := LaunchAgents(value, maxValue, numAgents, liarRatio)
	for _, agent := range *pool {
		agent.SetOnline(true)
		go agent.StartProcessing()
	}
	writeConfigFile(pool)
	if err != nil {
		fmt.Sprintf("Error launching agents: %s", err.Error())
		os.Exit(1)
	}
	_ = LiarsLieClient{
		Agents: *pool,
	}
	for {
		fmt.Print("liarslie>>")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\r\n")

		//split the text into pieces by spaces
		cmdPieces := strings.Split(text, " ")
		if len(cmdPieces) == 0 || (len(cmdPieces) == 1 && cmdPieces[0] == "") {
			continue
		}

		command, args, err := rootCmd.Find([]string{cmdPieces[0]})
		if err != nil || command == rootCmd {
			log.Printf("Unknown Command to execute : %s\n", text)
			continue
		}

		args = append(args, cmdPieces[1:]...)

		command.Run(command, args)
		command.Execute()
	}
}
