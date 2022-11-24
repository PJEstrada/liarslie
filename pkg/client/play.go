package client

import (
	"fmt"
	"liarslie/pkg/agents"
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
func FindMajorityValue(values []int, numAgents int) int {
	valuesCount := map[int]int{}
	for _, val := range values {
		valuesCount[val] += 1
	}
	percentages := map[int]float64{}
	max := -1.0
	maxVal := -1
	for key, val := range valuesCount {
		percentages[key] += float64(float64(val) / float64(numAgents))
		if percentages[key] > float64(max) {
			max = percentages[key]
			maxVal = key
		}
	}
	return maxVal
}

// PlayStandard excutes the game by querying all agents on the network and determining real value V.
func PlayStandard(agentsRegistry agents.AgentsRegistry) {
	chOut := make(chan agents.MessageGetValueResult)
	wg := new(sync.WaitGroup)
	wg.Add(len(agentsRegistry))
	for key, agent := range agentsRegistry {
		msg := agents.MessageGetValue{
			ID: key,
		}

		go agent.GetValue(&msg, chOut)
	}
	var values []int
	go func() {
		for msg := range chOut {
			values = append(values, msg.Value)
			wg.Done()
		}
	}()
	if waitTimeout(wg, time.Second*5) {
		fmt.Println("Timed out waiting for agents.")
		return
	} else {

		fmt.Println("values", values)
		maxVal := FindMajorityValue(values, len(agentsRegistry))
		fmt.Println("The network value is: ", maxVal)
		close(chOut)
	}

}

func PlayExpert() {

}
