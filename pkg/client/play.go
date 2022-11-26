package client

import (
	"fmt"
	"github.com/google/uuid"
	"liarslie/pkg/agents"
	"liarslie/pkg/consensus"
	"math/rand"
	"os"
	"sync"
	"time"
)

// waitTimeout waits for a waitgroup to execute before the given timeout
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// queryAgents queries the given set of agents to get the real value of the network v.
func queryAgents(agentsRegistry agents.AgentsRegistry, queryAll bool, queryOthers bool, consensusRatio float32) int {
	chOut := make(chan agents.MessageGetValueResult)
	wg := new(sync.WaitGroup)
	wg.Add(len(agentsRegistry))
	for key, agent := range agentsRegistry {
		msg := agents.MessageGetValue{
			ID:          key,
			KnownValues: map[uuid.UUID]int{},
		}

		go agent.GetValue(&msg, chOut, queryOthers)
	}
	var values []int
	go func() {
		for msg := range chOut {
			values = append(values, msg.Value)
			wg.Done()
		}
	}()
	if waitTimeout(wg, time.Second*10) {
		fmt.Println("Timed out waiting for agents.")
		return -1
	} else {
		var maxVal int
		fmt.Println("values", values)
		if queryAll {
			maxVal = consensus.FindMajorityValue(values)
		} else {
			maxVal = consensus.FindMajorityValuePercent(values, consensusRatio)
		}

		fmt.Println("The network value is: ", maxVal)
		close(chOut)
		return maxVal
	}

}

// getOnlineAgents filters to only online agents respecting the given liar ratio.
func getOnlineAgents(agentsRegistry agents.AgentsRegistry) []agents.AgentBehaviour {
	result := []agents.AgentBehaviour{}
	for _, agent := range agentsRegistry {
		if agent.IsOnline() {
			result = append(result, agent)
		}
	}
	return result
}

// PlayStandard executes the game by querying all agents on the network and determining real value V.
func PlayStandard(agentsRegistry agents.AgentsRegistry) {
	queryAgents(agentsRegistry, true, false, -1.0) // liarRatio not used
}

// PlayExpert executes the game by querying the given subset of agents and given liarRatio
func PlayExpert(agentsRegistry agents.AgentsRegistry, numAgents int, liarRatio float32) {
	onlineAgents := getOnlineAgents(agentsRegistry)
	if len(onlineAgents) < numAgents {
		fmt.Println(fmt.Sprintf("Not enough online agents to play in expert mode. Want %d have %d", numAgents, len(onlineAgents)))
		os.Exit(1)
	}
	indexesToQuery := rand.Perm(numAgents)
	onlineRegistry := agents.AgentsRegistry{}
	for _, index := range indexesToQuery {
		agent := onlineAgents[index]
		onlineRegistry[agent.GetID()] = agent
	}
	queryAgents(onlineRegistry, false, true, liarRatio)
}
