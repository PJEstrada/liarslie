package agents

import "github.com/google/uuid"

type LiarsLieMessageRequest struct {
	*MessageGetValue
	*MessageStop
	*MessageSetPeers
	*MessageSetOnline
	*MessageSetValue
	ChOut chan LiarsLieMessageResult
}

type LiarsLieMessageResult struct {
	*MessageGetValueResult
	*MessageStopResult
	*MessageSetOnlineResult
	*MessageSetPeersResult
	*MessageSetValueResult
}
type MessageSetValue struct {
	ID    uuid.UUID
	Value int
}

type MessageSetValueResult struct {
	ID    uuid.UUID
	Value int
}

type MessageGetValue struct {
	ID uuid.UUID

	WithPeers bool
}

type MessageSetPeers struct {
	Peers AgentsRegistry
}

type MessageSetOnline struct {
	Online bool
}

type MessageStop struct {
	ID uuid.UUID
}

type MessageSetPeersResult struct {
	Peers AgentsRegistry
}

type MessageGetValueResult struct {
	ID      uuid.UUID
	AgentID uuid.UUID
	Value   int
}

type MessageStopResult struct {
	ID      uuid.UUID
	AgentID uuid.UUID
}

type MessageSetOnlineResult struct {
	Online bool
	ID     uuid.UUID
}
