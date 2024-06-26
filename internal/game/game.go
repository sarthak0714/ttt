package game

import (
	"errors"
)

var ErrGameNotFound = errors.New("game not found")

type Player struct {
	ID     string
	Symbol string
}

type Game struct {
	ID      string
	Board   Board
	Players [2]Player
	Turn    int
}

func NewGame(id string, player1, player2 Player) *Game {
	return &Game{
		ID:      id,
		Board:   NewBoard(),
		Players: [2]Player{player1, player2},
		Turn:    0,
	}
}

func (g *Game) MakeMove(playerID string, row, col int) error {
	if g.Players[g.Turn].ID != playerID {
		return errors.New("not your turn")
	}
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return errors.New("invalid move")
	}
	if g.Board[row][col] != "" {
		return errors.New("cell already occupied")
	}

	g.Board[row][col] = g.Players[g.Turn].Symbol
	g.Turn = 1 - g.Turn
	return nil
}

func (g *Game) Status() string {
	if winner := g.Board.Winner(); winner != "" {
		return "winner:" + winner
	}
	if g.Board.IsFull() {
		return "draw"
	}
	return "ongoing"
}
