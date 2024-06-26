package main

import (
	"log"
	"net/http"

	"github.com/sarthak0714/ttt/internal/actor"
	"github.com/sarthak0714/ttt/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	gameActor := actor.NewGameActor()
	go gameActor.Run()

	h := handlers.NewHandler(gameActor)

	r.Get("/", h.Index)
	r.Get("/game/{id}", h.Game)
	r.Post("/game", h.CreateGame)
	r.Post("/game/{id}/move", h.MakeMove)
	r.Get("/ws/{id}", h.WebSocket)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
