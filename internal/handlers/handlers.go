package handlers

import (
	"fmt"
	"net/http"

	"github.com/sarthak0714/ttt/internal/actor"
	"github.com/sarthak0714/ttt/internal/game"
	"github.com/sarthak0714/ttt/internal/message"
	"github.com/sarthak0714/ttt/templates"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Handler struct {
	gameActor *actor.GameActor
}

func NewHandler(gameActor *actor.GameActor) *Handler {
	return &Handler{gameActor: gameActor}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}

func (h *Handler) Game(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")

	responseChan := make(chan interface{})
	h.gameActor.Send(message.GetGame{GameID: gameID, Response: responseChan})

	response := <-responseChan
	game, ok := response.(*game.Game)
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	state := templates.GameState{
		ID:     game.ID,
		Board:  game.Board,
		Turn:   game.Players[game.Turn].ID,
		Status: game.Status(),
	}

	templates.Game(state).Render(r.Context(), w)
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	player1ID := r.FormValue("player1")
	player2ID := r.FormValue("player2")
	gameID := "game-" + player1ID + "-" + player2ID

	player1 := game.Player{ID: player1ID, Symbol: "X"}
	player2 := game.Player{ID: player2ID, Symbol: "O"}

	responseChan := make(chan interface{})
	h.gameActor.Send(message.CreateGame{
		GameID:   gameID,
		Player1:  player1,
		Player2:  player2,
		Response: responseChan,
	})

	<-responseChan
	http.Redirect(w, r, "/game/"+gameID, http.StatusSeeOther)
}

func (h *Handler) MakeMove(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	playerID := r.FormValue("player_id")
	row := r.FormValue("row")
	col := r.FormValue("col")

	var rowInt, colInt int
	_, err := fmt.Sscanf(row+","+col, "%d,%d", &rowInt, &colInt)
	if err != nil {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	responseChan := make(chan interface{})
	h.gameActor.Send(message.MakeMove{
		GameID:   gameID,
		PlayerID: playerID,
		Row:      rowInt,
		Col:      colInt,
		Response: responseChan,
	})

	response := <-responseChan
	if err, ok := response.(error); ok {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	h.Game(w, r)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	updates := make(chan game.Game)
	h.gameActor.Send(message.Subscribe{GameID: gameID, Subscriber: updates})
	defer h.gameActor.Send(message.Unsubscribe{GameID: gameID, Subscriber: updates})

	for game := range updates {
		state := templates.GameState{
			ID:     game.ID,
			Board:  game.Board,
			Turn:   game.Players[game.Turn].ID,
			Status: game.Status(),
		}
		if err := conn.WriteJSON(state); err != nil {
			break
		}
	}
}
