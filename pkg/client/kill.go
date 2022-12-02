package client

import (
	"fmt"
	"github.com/google/uuid"
	"liarslie/pkg/agents"
	"sync"
	"time"
)

func KillAgent(id uuid.UUID) {
	pool := CurrentClient.AgentsFullNetwork
	if _, ok := pool[id]; ok {
		wg := new(sync.WaitGroup)
		chout := make(chan agents.LiarsLieMessageResult)
		msg := agents.LiarsLieMessageRequest{
			MessageSetOnline: &agents.MessageSetOnline{
				Online: false,
			},
			ChOut: chout,
		}
		pool[id].ChIn <- msg
		msgResult := WaitTimeoutGetResult(wg, time.Second*3, chout)
		if msgResult.MessageSetOnlineResult != nil {
			pool[id].Online = false
		}

		fmt.Println(fmt.Sprintf("Agent %s killed", id.String()))
	}

}
