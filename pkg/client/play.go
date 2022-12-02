package client

import (
	"errors"
	"fmt"
	"liarslie/pkg/agents"
	"liarslie/pkg/consensus"
	"math/rand"
	"sync"
	"time"
)

// queryAgents queries the given set of agents to get the real value of the network v.
func queryAgents(agentsRegistry agents.AgentsRegistry, useConsensusRatio bool, queryOthers bool, consensusRatio float32) int {
	var values []int
	for key, agent := range agentsRegistry {
		chOut := make(chan agents.LiarsLieMessageResult)
		wg := new(sync.WaitGroup)
		wg.Add(len(agentsRegistry))
		msg := agents.LiarsLieMessageRequest{
			MessageGetValue: &agents.MessageGetValue{
				ID: key,
			},
			ChOut: chOut,
		}
		agent.ChIn <- msg
		msgResponse := WaitTimeoutGetResult(wg, time.Second*5, chOut)
		if msgResponse == nil {
			fmt.Println("Timed out waiting for agents.")
			return -1
		}
		if msgResponse.MessageGetValueResult != nil {
			values = append(values, msgResponse.MessageGetValueResult.Value)
		}

	}
	var maxVal int
	if useConsensusRatio {
		maxVal = consensus.FindMajorityValue(values)
	} else {
		maxVal = consensus.FindMajorityValuePercent(values, consensusRatio)
	}
	fmt.Printf("Queried %d Agents \n", len(agentsRegistry))
	fmt.Println("The network value is: ", maxVal)
	return maxVal

}

// getOnlineAgents filters to only online agents respecting the given liar ratio.
func getOnlineAgents(agentsRegistry agents.AgentsRegistry) []*agents.AgentChannel {
	result := []*agents.AgentChannel{}
	for _, agent := range agentsRegistry {
		if agent.Online {
			result = append(result, agent)
		}
	}
	return result
}

// PlayStandard executes the game by querying all agents on the network and determining real value V.
func PlayStandard(agentsRegistry agents.AgentsRegistry) int {
	onlineAgents := getOnlineAgents(agentsRegistry)
	onlineRegistry := agents.AgentsRegistry{}
	for _, agent := range onlineAgents {
		onlineRegistry[agent.ID] = agent
	}
	val := queryAgents(onlineRegistry, true, false, -1.0) // liarRatio not used
	return val
}

// PlayExpert executes the game by querying the given subset of agents and given liarRatio
func PlayExpert(agentsRegistry agents.AgentsRegistry, numAgents int, liarRatio float32) (int, error) {
	onlineAgents := getOnlineAgents(agentsRegistry)
	if len(onlineAgents) < numAgents {
		err := errors.New(fmt.Sprintf("Not enough online agents to play in expert mode. Want %d have %d", numAgents, len(onlineAgents)))
		return 0, err
	}
	indexesToQuery := rand.Perm(numAgents)
	onlineRegistry := agents.AgentsRegistry{}
	for _, index := range indexesToQuery {
		agent := onlineAgents[index]
		onlineRegistry[agent.ID] = agent
	}
	val := queryAgents(onlineRegistry, false, true, liarRatio)
	return val, nil
}
