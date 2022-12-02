package agents

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

var ConfigFileName string = "app.config"

func ReadConfigFile() (AgentsRegistry, error) {
	readFile, err := os.Open(ConfigFileName)
	defer readFile.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
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
			return nil, err
		}
		network := GetAgentsNetwork()
		agent := network[id]
		if agent.Online {
			result[agent.ID] = agent
		}

	}
	return result, nil
}
