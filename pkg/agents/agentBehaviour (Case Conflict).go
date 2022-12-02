package agents

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type AgentBehaviour interface {
	GetID() uuid.UUID
	GetPeers(forceRead bool) (AgentsRegistry, error)
	SetPeers(AgentsRegistry)
	GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, withPeers bool) error
	WaitTimeoutGetValue(wg *sync.WaitGroup, timeout time.Duration, c chan MessageGetValueResult) *MessageGetValueResult
	IsOnline() bool
	SetOnline(v bool)
	SetValue(v int)
}
