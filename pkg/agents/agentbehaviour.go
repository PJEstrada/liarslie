package agents

type AgentBehaviour interface {
	GetValue(msg *MessageGetValue) *MessageGetValueResult
	StartProcessing()
	IsOnline() bool
	SetOnline(v bool)
}
