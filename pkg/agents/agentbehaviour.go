package agents

import "github.com/google/uuid"

type AgentBehaviour interface {
	GetID() uuid.UUID
	GetPeers() AgentsRegistry
	GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult, withPeers bool)
	IsOnline() bool
	SetOnline(v bool)
}
