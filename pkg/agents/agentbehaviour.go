package agents

type AgentBehaviour interface {
	GetValue(msg *MessageGetValue, chOut chan MessageGetValueResult)
	IsOnline() bool
	SetOnline(v bool)
}
