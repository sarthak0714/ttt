package actor

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/sarthak0714/ttt/internal/game"
	"github.com/sarthak0714/ttt/internal/message"
)

func TestActorModel(t *testing.T) {
	gameActor := NewGameActor()
	go gameActor.Run()

	var wg sync.WaitGroup

	// Create games
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			gameID := fmt.Sprintf("game%d", i)
			responseChan := make(chan interface{})
			gameActor.Send(message.CreateGame{
				GameID:   gameID,
				Player1:  game.Player{ID: fmt.Sprintf("player%d-1", i), Symbol: "X"},
				Player2:  game.Player{ID: fmt.Sprintf("player%d-2", i), Symbol: "O"},
				Response: responseChan,
			})

			createdGame := <-responseChan
			if createdGame.(*game.Game).ID != gameID {
				t.Errorf("expected game ID to be '%s', got %s", gameID, createdGame.(*game.Game).ID)
			}
		}(i)
	}

	wg.Wait()
	log.Printf("Created 10 games in %v\n", time.Since(start))

	// Make moves
	start = time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			gameID := fmt.Sprintf("game%d", i)
			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					moveResponseChan := make(chan interface{})
					playerID := fmt.Sprintf("player%d-%d", i, (j+k)%2+1)
					gameActor.Send(message.MakeMove{
						GameID:   gameID,
						PlayerID: playerID,
						Row:      j,
						Col:      k,
						Response: moveResponseChan,
					})

					err := <-moveResponseChan
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	log.Printf("Completed moves for 10 games in %v\n", time.Since(start))
}

func BenchmarkActorModel(b *testing.B) {
	gameActor := NewGameActor()
	go gameActor.Run()

	var totalGames, totalMoves int
	var mutex sync.Mutex // Mutex for safe concurrent updates

	// Use time.Now() to measure total benchmark duration
	startTime := time.Now()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var localGames, localMoves int

			var wg sync.WaitGroup

			// Create games
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					gameID := fmt.Sprintf("game%d", i)
					responseChan := make(chan interface{})
					gameActor.Send(message.CreateGame{
						GameID:   gameID,
						Player1:  game.Player{ID: fmt.Sprintf("player%d-1", i), Symbol: "X"},
						Player2:  game.Player{ID: fmt.Sprintf("player%d-2", i), Symbol: "O"},
						Response: responseChan,
					})

					<-responseChan
				}(i)
			}

			wg.Wait()

			localGames += 100

			// Make moves
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					gameID := fmt.Sprintf("game%d", i)
					for j := 0; j < 3; j++ {
						for k := 0; k < 3; k++ {
							moveResponseChan := make(chan interface{})
							playerID := fmt.Sprintf("player%d-%d", i, (j+k)%2+1)
							gameActor.Send(message.MakeMove{
								GameID:   gameID,
								PlayerID: playerID,
								Row:      j,
								Col:      k,
								Response: moveResponseChan,
							})

							<-moveResponseChan
						}
					}
				}(i)
			}

			wg.Wait()

			localMoves += 300

			// Update totals
			mutex.Lock()
			totalGames += localGames
			totalMoves += localMoves
			mutex.Unlock()
		}
	})

	// Calculate total time for all iterations
	totalTime := time.Since(startTime)

	log.Printf("Benchmark results:\n")
	log.Printf("%d games and %d tested in %v:\n", totalGames, totalMoves, totalTime)

}
