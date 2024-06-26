package actor

import (
	"github.com/sarthak0714/ttt/internal/game"
	"github.com/sarthak0714/ttt/internal/message"
)

type GameActor struct {
	games       map[string]*game.Game
	inbox       chan message.Message
	subscribers map[string][]chan<- game.Game
}

func NewGameActor() *GameActor {
	return &GameActor{
		games:       make(map[string]*game.Game),
		inbox:       make(chan message.Message, 100),
		subscribers: make(map[string][]chan<- game.Game),
	}
}

func (ga *GameActor) Run() {
	for msg := range ga.inbox {
		ga.handleMessage(msg)
	}
}

func (ga *GameActor) handleMessage(msg message.Message) {
	switch m := msg.(type) {
	case message.CreateGame:
		ga.createGame(m)
	case message.GetGame:
		ga.getGame(m)
	case message.MakeMove:
		ga.makeMove(m)
	case message.Subscribe:
		ga.subscribe(m)
	case message.Unsubscribe:
		ga.unsubscribe(m)
	}
}

func (ga *GameActor) createGame(m message.CreateGame) {
	newGame := game.NewGame(m.GameID, m.Player1, m.Player2)
	ga.games[m.GameID] = newGame
	m.Response <- newGame
}

func (ga *GameActor) getGame(m message.GetGame) {
	g, ok := ga.games[m.GameID]
	if !ok {
		m.Response <- nil
		return
	}
	m.Response <- g
}

func (ga *GameActor) makeMove(m message.MakeMove) {
	g, ok := ga.games[m.GameID]
	if !ok {
		m.Response <- game.ErrGameNotFound
		return
	}

	err := g.MakeMove(m.PlayerID, m.Row, m.Col)
	m.Response <- err

	if err == nil {
		ga.notifySubscribers(m.GameID)
	}
}

func (ga *GameActor) subscribe(m message.Subscribe) {
	ga.subscribers[m.GameID] = append(ga.subscribers[m.GameID], m.Subscriber)
}

func (ga *GameActor) unsubscribe(m message.Unsubscribe) {
	subs := ga.subscribers[m.GameID]
	for i, sub := range subs {
		if sub == m.Subscriber {
			ga.subscribers[m.GameID] = append(subs[:i], subs[i+1:]...)
			close(sub)
			break
		}
	}
}

func (ga *GameActor) notifySubscribers(gameID string) {
	g := ga.games[gameID]
	for _, sub := range ga.subscribers[gameID] {
		sub <- *g
	}
}

func (ga *GameActor) Send(msg message.Message) {
	ga.inbox <- msg
}
