package agents

import (
	"fmt"
	"github.com/google/uuid"
	"liarslie/pkg/consensus"
	"sync"
	"time"
)

type Agent struct {
	ID               uuid.UUID
	value            int
	Online           bool
	MsgNetworkValues sync.Map
	peers            *AgentsRegistry
}

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
		defer close(c)
		return true // timed out
	}
}

// FindMajorityValue finds the number that is repeated the most among the list of integers in values param.
func FindMajorityValue(values []int) int {
	valuesCount := map[int]int{}
	for _, val := range values {
		valuesCount[val] += 1
	}
	percentages := map[int]float64{}
	max := -1.0
	maxVal := -1
	for key, val := range valuesCount {
		percentages[key] += float64(float64(val) / float64(len(values)))
		if percentages[key] > float64(max) {
			max = percentages[key]
			maxVal = key
		}
	}
	return maxVal
}

func NewHonestAgent(ID uuid.UUID, Value int) Agent {
	return Agent{
		ID:               ID,
		value:            Value,
		MsgNetworkValues: sync.Map{},
	}
}

func (a *Agent) GetPeers() AgentsRegistry {
	if a.peers == nil {
		result := ReadConfigFile()
		return result

	} else {
		return *a.peers
	}

}
func (a *Agent) GetID() uuid.UUID {
	return a.ID
}

func (a *Agent) GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, withPeers bool) {
	if withPeers {
		a.getValueExpert(*msg, chOut)

	} else {
		chOut <- MessageGetValueResult{
			ID:      msg.ID,
			AgentID: a.ID,
			Value:   a.value,
		}
	}

}
func copyMap(m map[uuid.UUID]int) map[uuid.UUID]int {
	res := map[uuid.UUID]int{}
	for key, val := range m {
		res[key] = val
	}
	return res
}
func (a *Agent) getValueForMsg(msgID uuid.UUID) (int, bool) {
	val, ok := a.MsgNetworkValues.Load(msgID)
	var result int
	if ok {
		result = val.(int)
	}
	return result, ok
}
func (a *Agent) setValueForMessage(msgID uuid.UUID, value int) {
	a.MsgNetworkValues.Store(msgID, value)
}
func (a *Agent) getValueExpert(msg MessageGetValue, chOut chan MessageGetValueResult) {
	// Case of value already processed by Node
	existingValue, found := a.getValueForMsg(msg.ID)
	if found {
		fmt.Println("Existing Value on node", a.ID, existingValue)
		chOut <- MessageGetValueResult{
			ID:      msg.ID,
			AgentID: a.ID,
			Value:   existingValue,
		}
		return
	}
	agentsNet := a.GetPeers()
	chAgents := make(chan MessageGetValueResult)
	fmt.Println("Node value is", a.value, a.ID)
	wg := new(sync.WaitGroup)
	var values []int
	fmt.Println(fmt.Sprintf("GET VALUE EXPERT -------- %s", a.GetID().String()))
	msg.KnownValues[a.ID] = a.value
	fmt.Println("KNOWN", msg.KnownValues)
	numCalls := 0
	numReturns := 0
	for _, agent := range agentsNet {
		if _, ok := msg.KnownValues[agent.GetID()]; !ok {
			//fmt.Println(fmt.Sprintf("[%s] querying %s", a.ID.String(), agent.GetID().String()))
			wg.Add(1)
			numCalls += 1
			go agent.GetValue(&MessageGetValue{
				ID:          msg.ID,
				KnownValues: copyMap(msg.KnownValues),
			}, chAgents, true)
		} else {
			fmt.Println(fmt.Sprintf("[%s] using cached %d value from: %s", a.ID.String(), msg.KnownValues[agent.GetID()], agent.GetID().String()))
			values = append(values, msg.KnownValues[agent.GetID()])
		}
	}

	go func() {
		for msgResponse := range chAgents {
			numReturns += 1
			fmt.Println("ID - NUM CALLS - NUM RETURNS", a.ID, numCalls, numReturns)
			fmt.Println("msg response", msgResponse)
			values = append(values, msgResponse.Value)
			//msg.KnownValues[msgResponse.AgentID] = msgResponse.Value
			wg.Done()
		}
	}()
	var networkVal int
	if waitTimeout(wg, time.Second*2) {
		fmt.Println(fmt.Sprintf("[Agent %s] Timeout querying agents ", a.ID.String()))
		fmt.Println("values timeout", values)
		networkVal = consensus.FindMajorityValue(values)
	} else {

		networkVal = consensus.FindMajorityValue(values)
		fmt.Println(fmt.Sprintf("Finished: [Agent %s] Network Val: %d", a.ID.String(), networkVal))
	}
	a.setValueForMessage(msg.ID, networkVal)
	defer close(chAgents)
	chOut <- MessageGetValueResult{
		ID:      msg.ID,
		AgentID: a.ID,
		Value:   networkVal,
	}
}

func (a *Agent) SetOnline(v bool) {
	a.Online = v
}

func (a *Agent) IsOnline() bool {
	return a.Online
}
