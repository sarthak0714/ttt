package message

import "github.com/sarthak0714/ttt/internal/game"

type Message interface {
	isMessage()
}

type CreateGame struct {
	GameID   string
	Player1  game.Player
	Player2  game.Player
	Response chan<- interface{}
}

type GetGame struct {
	GameID   string
	Response chan<- interface{}
}

type MakeMove struct {
	GameID   string
	PlayerID string
	Row      int
	Col      int
	Response chan<- interface{}
}

type Subscribe struct {
	GameID     string
	Subscriber chan<- game.Game
}

type Unsubscribe struct {
	GameID     string
	Subscriber chan<- game.Game
}

func (CreateGame) isMessage()  {}
func (GetGame) isMessage()     {}
func (MakeMove) isMessage()    {}
func (Subscribe) isMessage()   {}
func (Unsubscribe) isMessage() {}
