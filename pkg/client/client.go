package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"liarslie/pkg/agents"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
)

var CurrentClient LiarsLieClient

type LiarsLieClient struct {
	AgentsFullNetwork agents.AgentsRegistry
}

// spawnLiars spawns the given number of liar agents respecting the given ratio
func spawnLiars(pool agents.AgentsRegistry, maxValue int, liarRatio float32, numAgents int) {
	numLiars := int(math.Round(float64(float32(numAgents) * liarRatio)))
	fmt.Println("Num liars: ", numLiars)
	for i := 0; i < numLiars; i++ {
		id := uuid.New()
		value := rand.Intn(maxValue + 1)
		agent := agents.NewLiarAgent(id, value)
		pool[id] = &agent
	}
}

// spawnHonestAgents spawns the given number of honest agents respecting the given ratio
func spawnHonestAgents(pool agents.AgentsRegistry, value int, numAgents int, liarRatio float32) {
	numHonest := int(math.Round(float64((1 - liarRatio) * float32(numAgents))))
	fmt.Println("Num honest: ", numHonest)
	for i := 0; i < numHonest; i++ {
		id := uuid.New()
		agent := agents.NewHonestAgent(id, value)
		pool[id] = &agent
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
func writeConfigFile(pool *agents.AgentsRegistry) {
	f, err := os.Create("app.config")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	for key, _ := range *pool {
		_, err2 := f.WriteString(fmt.Sprintf("%s\n", key))
		if err2 != nil {
			log.Fatal(err)
		}
	}
}

// cleanFile cleans the text on the app.config file
func cleanFile() {
	f, err := os.Create("app.config")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	f.WriteString("")

}

// addNetworkPeersToAgents initializes peer network registry on the existing agents.
func addNetworkPeersToAgents(pool *agents.AgentsRegistry) {
	// Populate Peers
	for _, agent := range *pool {
		agent.GetPeers()
	}
}

// displayShell shows the command prompt on terminal and starts waiting for user input.
func displayShell(rootCmd *cobra.Command) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("liarslie>>")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}
		text = strings.Trim(text, "\r\n")

		if strings.Compare(text, "exit") == 0 {

			log.Println("Exiting....")
			os.Exit(0)
		}

		cmdPieces := strings.Split(text, " ")

		command, args, err := rootCmd.Parent().Find(cmdPieces)
		if err != nil {
			log.Printf("Unknown Command to execute : %s\n", text)
			continue
		}
		command.ParseFlags(args)
		command.Run(command, args)
		command.Execute()
	}
}

// StartClient initializes agents and starts command prompt.
func StartClient(rootCmd *cobra.Command, value int, maxValue int, numAgents int, liarRatio float32) {

	pool, err := LaunchAgents(value, maxValue, numAgents, liarRatio)
	for _, agent := range *pool {
		agent.SetOnline(true)
	}
	writeConfigFile(pool)
	if err != nil {
		fmt.Sprintf("Error launching agents: %s", err.Error())
		os.Exit(1)
	}

	agents.SetAgentsNetwork(*pool)
	addNetworkPeersToAgents(pool)
	CurrentClient = LiarsLieClient{
		AgentsFullNetwork: *pool,
	}
	displayShell(rootCmd)
}

// StopClient stops all agents and exits program.
func StopClient() {

	for _, val := range CurrentClient.AgentsFullNetwork {
		val.SetOnline(false)
	}
	cleanFile()
	os.Exit(0)
}
