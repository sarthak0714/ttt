package actor

import (
	"github.com/sarthak0714/ttt/internal/message"
)

type PlayerActor struct {
	ID    string
	inbox chan message.Message
}

func NewPlayerActor(id string) *PlayerActor {
	return &PlayerActor{
		ID:    id,
		inbox: make(chan message.Message, 10),
	}
}

func (pa *PlayerActor) Run() {
	for msg := range pa.inbox {
		pa.handleMessage(msg)
	}
}

func (pa *PlayerActor) handleMessage(msg message.Message) {
	// NO USE WILL ADD EMOTE SPAM LATER XDDDDD
}

func (pa *PlayerActor) Send(msg message.Message) {
	pa.inbox <- msg
}
