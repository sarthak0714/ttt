
# Tic-Tac-Toe

This project is a web-based Tic-Tac-Toe game implemented using the actor model in Go. It leverages channels for communication between actors to ensure better encapsulation and scalability. The game supports creating new games, making moves, and subscribing to game updates via WebSockets.

## Project Structure

```
tictactoe/
├── server/
│   └── main.go
│
├── internal/
│   ├── actor/
│   │   ├── actor.go
│   │   └── actor.go
│   ├── game/
│   │   ├── game.go
│   │   └── board.go
│   ├── handlers/
│   │   └── handlers.go
│   └── message/
│       └── message.go
├── templates/
│   ├── base.templ
│   ├── index.templ
│   └── game.templ
├── static/
│   └── styles.css
├── go.mod
└── go.sum
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sarthak0714/tictactoe.git
   cd tictactoe
   ```

2. Install the `templ` cli:
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```

3. Install dependencies:
   ```bash
   make deps
   ```

4. Generate Go files from templ:
   ```bash
   make gen
   ```

5. Build Server:
   ```bash
   make build
   ```

6. Run Server (Running will also build and gen):
   ```bash
   make run
   ```

The server will start on port 8080. You can access the game at `http://localhost:8080`.

## Usage

### Creating a New Game

1. Go to the home page (`http://localhost:8080`).
2. Enter the IDs of two players and create a new game.
3. You will be redirected to the game page where you can make moves.

### Making Moves

1. On the game page, players can make moves by clicking on the board.
2. The game will indicate whose turn it is and display the game's status (ongoing, draw, or winner).

### WebSocket Updates

1. The game page uses WebSocket to subscribe to game updates.
2. When a move is made, all connected clients will be updated in real-time.

## Internal Structure

### Actors

- **GameActor**: Manages all games and handles game-related messages.
- **PlayerActor**: Manages player-specific messages (not extensively used in this implementation but can be expanded).

### Messages

Messages define the communication between actors:
- `CreateGame`
- `GetGame`
- `MakeMove`
- `Subscribe`
- `Unsubscribe`

### Handlers

- **handlers.go**: Defines HTTP handlers for serving the web pages and handling user actions (creating games, making moves, WebSocket connections).

### Game Logic

- **game.go**: Defines the `Game` and `Player` structs and the main game logic.
- **board.go**: Defines the `Board` struct and functions for checking the game status (winner, draw, etc.).

### Templates

- **base.templ**: Base HTML template.
- **index.templ**: Template for the home page.
- **game.templ**: Template for the game page.

### Static Files

- **styles.css**: Basic CSS for styling the web pages.

## Contributing

Feel free to submit issues and pull requests. Contributions are welcome!


