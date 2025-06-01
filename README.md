## Minesweeper Arena

A browser-based, real-time multiplayer Minesweeper game built with Go (backend) and vanilla JavaScript (frontend). Players are matched randomly to compete on a shared board and race to reveal more safe cells than their opponent.

### Features

- Real-time 1v1 multiplayer gameplay
- Shared minefield logic for consistent state
- WebSocket-based client-server communication
- Matchmaking queue system
- Flag placement and mine detection
- Sound effects for game events (start, win, loss)
- Responsive UI

### Run

```sh
git clone https://github.com/yourname/minesweeper-multiplayer.git
cd minesweeper-multiplayer
go run .
```

The server will start at [http://localhost:8080](http://localhost:8080), unless configured otherwise.

### The Application

The application exposes a single "Start" button that connects the client to the server using a non-blocking WebSocket. The client enters a matchmaking queue and is paired with another client when available. The game starts immediately once a match is formed.

All interactions are routed through the server. The client sends requests to reveal cells, and the server responds with updated state. Game logic resides entirely on the server, which enforces rules and prevents client-side tampering.

The client sees its own board on the left with full visibility. The opponent's board is shown on the right but only displays interacted tiles, not their values. This allows tracking opponent progress without revealing their game state.

#### Game Termination Rules

- A player wins by revealing all safe cells first.
- If one player hits a mine and the other has revealed more safe cells, the latter wins immediately.
- If both players hit a mine, the player with more revealed safe cells wins. If tied, the faster player wins.

The game ends immediately when any of the above conditions are met, and both boards are fully revealed.

### Limitation of the current gameplay

The current gameplay is highly unforgiving. Because fully randomized MineSweeper contains too many undecideable positions, early and late stages often contain undecidable positions. This leads to many matches ending quickly based on chance. A potential improvement would be introducing non-terminal outcomes when a mine is triggered, to improve playability.

### Missing features

Currently, the application is missing user-reliant features that are not yet implemented. These are
- Ladder
- Recent matches

### Author

Created for Czech Technical University, Faculty of Electrical Engineering  
Course: Vývoj klientských aplikací v JavaScriptu (B0B39KAJ)  
2025 — Marek Jagoš
