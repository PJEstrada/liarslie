package client

import (
	"errors"
	"fmt"
	"liarslie/pkg/agents"
	"log"
	"os"
	"sync"
	"time"
)

func ExtendNetwork(value int, maxValue int, numAgents int, liarRatio float32) (agents.AgentsRegistry, error) {
	var pool agents.AgentsRegistry
	if _, err := os.Stat(agents.ConfigFileName); err == nil {
		pool, err = agents.ReadConfigFile()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		oldNumAgents := len(pool)
		if len(pool) > numAgents {
			err := errors.New(fmt.Sprintf("Cannot extend network. Current Size [%d] is less than current network size [%d].", oldNumAgents, numAgents))
			fmt.Println(err)
			return nil, err

		}
		// Set value for all agents (liars might ignore this...)
		for _, agent := range pool {
			chout := make(chan agents.LiarsLieMessageResult)
			wg := new(sync.WaitGroup)
			msg := agents.LiarsLieMessageRequest{
				MessageSetValue: &agents.MessageSetValue{
					Value: value,
				},
				ChOut: chout,
			}
			agent.ChIn <- msg
			msgResponse := WaitTimeoutGetResult(wg, time.Second*3, chout)
			if msgResponse == nil {
				return nil, errors.New("Cannot set value for agent. Communication timed out.")
			}
		}
		missingAgents := numAgents - oldNumAgents
		spawnLiars(pool, maxValue, liarRatio, missingAgents)
		spawnHonestAgents(pool, value, missingAgents, liarRatio)
		writeConfigFile(&pool)
		fmt.Println(fmt.Sprintf("Network reconfigured. Total Agents %d", len(pool)))
		agents.SetAgentsNetwork(pool)

	} else {
		pool, err = CreateNetwork(value, maxValue, numAgents, liarRatio)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

	}
	CurrentClient.AgentsFullNetwork = pool
	return pool, nil
}
