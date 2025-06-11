## Minesweeper Arena

A browser-based, real-time multiplayer Minesweeper game built with Go (backend) and vanilla JavaScript (frontend). Players are matched randomly to compete on a shared board and race to finish first.

The application currently runs at http://18.193.94.76/ and is in experiemental stage of development.

### Features

- Real-time 1v1 multiplayer gameplay
- Shared minefield logic for consistent state
- WebSocket-based client-server communication
- Matchmaking queue system
- Flag placement and mine detection
- Custom graphics, including logo and minesweeper icons
- Sound effects for game events (start, win, loss)
- Responsive UI

### Run

```sh
git clone https://github.com/yourname/minesweeper-multiplayer.git
cd minesweeper-multiplayer
go run .
```

The server will start at [http://localhost:80](http://localhost:80), unless configured otherwise.

### The Application

The center of the application is a single button that connects the client to the server using a non-blocking WebSocket. The client enters a matchmaking queue and is paired with another client when available. The clients enter a game once paired.

All interactions are routed through the server. The client sends requests to reveal cells, and the server responds with updated state. Game logic resides entirely on the server, which enforces rules and prevents any client-side tampering.

The client sees its own board on the left with full visibility. The opponent's board is shown on the right but only displays interacted tiles, not their values. This allows tracking opponent progress without revealing information about the board layout.

### Client Interaction & Gameplay

Clients are presented with a single button in the center of the screen:  
<img src="docs/image-6.png" width="300"/>

By pressing it, a socket is created and they join a queue:  
<img src="docs/image-1.png" width="300"/>

When paired with another client, they are shown a responsive double-board layout. The player's board is always on the left. The cross field marks the starting position. The player is not required to start on the cross but risks revealing a bomb on their first turn.  
<img src="docs/image-9.png" width="300"/>

Clients reveal fields on their own board while observing their opponent's moves.  
<img src="docs/image-13.png" width="300"/>

If a client reveals a mine, their board is fully revealed and they must wait for the other client to finish.  
<img src="docs/image-12.png" width="250"/>

Solving the full board results in an immediate victory:  
<img src="docs/image-4.png" width="400"/>

And a loss for the opponent:  
<img src="docs/image-5.png" width="400"/>



#### Game Termination Rules

- A player wins by revealing all safe cells first.
- If one player hits a mine, he waits for the other player to finish.
- If both players hit a mine, the player with more revealed safe cells wins. If tied, the faster player wins.

The game ends immediately when any of the above conditions are met.

### Missing planned features

Currently, the application is missing user-reliant features which is not yet implemented. The planned features were:
- Ladder, present in aside section of the application, and through a menu item
- Recent matches, present in aside section of the application, and through a profile

### Author

Created for Czech Technical University, Faculty of Electrical Engineering  
Course: Vývoj klientských aplikací v JavaScriptu (B0B39KAJ)  
2025 — Marek Jagoš
