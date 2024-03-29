package agents

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

func ReadConfigFile() AgentsRegistry {
	readFile, err := os.Open("app.config")
	defer readFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}
	result := AgentsRegistry{}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		values := strings.Split(line, " ")
		id, err := uuid.Parse(values[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return nil
		}
		network := GetAgentsNetwork()
		agent := network[id]
		if agent.IsOnline() {
			result[agent.GetID()] = agent
		}

	}
	return result
}
