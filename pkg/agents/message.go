package agents

import "github.com/google/uuid"

type MessageGetValue struct {
	ID          uuid.UUID
	KnownValues map[uuid.UUID]int
}

type MessageStop struct {
	ID uuid.UUID
}
type LiarsLieMessageRequest struct {
	*MessageGetValue
	*MessageStop
}

type LiarsLieMessageResult struct {
	*MessageGetValueResult
	*MessageStopResult
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
