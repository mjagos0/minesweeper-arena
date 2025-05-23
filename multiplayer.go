package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	conn *websocket.Conn
	send chan []byte
}

type MultiplayerGame struct {
	games   [2]*Game
	players [2]*Player
	mu      sync.Mutex
}

func NewMultiplayerGame(p1 *Player, p2 *Player) *MultiplayerGame {
	board := NewBoard(9, 9, 10)
	return &MultiplayerGame{
		games: [2]*Game{
			NewGame(board),
			NewGame(board),
		},
		players: [2]*Player{p1, p2},
	}
}

func (mg *MultiplayerGame) Start() {
	startMessage := struct {
		Type   string `json:"type"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}{
		Type:   "start",
		Width:  mg.games[0].board.width,
		Height: mg.games[0].board.height,
	}
	msg, _ := json.Marshal(startMessage)
	mg.players[0].send <- msg
	mg.players[1].send <- msg

	for i := 0; i < 2; i++ {
		go mg.listen(i)
	}
}

func (mg *MultiplayerGame) listen(playerIndex int) {
	player := mg.players[playerIndex]
	defer player.conn.Close()

	for {
		var move struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
		err := player.conn.ReadJSON(&move)
		if err != nil {
			log.Printf("Error reading from player %d: %v\n", playerIndex, err)
			break
		}
		log.Printf("Player %d sent move: %+v\n", playerIndex, move)

		mg.mu.Lock()
		result := mg.games[playerIndex].Reveal(move.X, move.Y)
		log.Printf("Reveal result for player %d: %+v\n", playerIndex, result)

		// Actual result for initiating player
		response := struct {
			RevealResult
			YourMove bool `json:"yourMove"`
		}{
			RevealResult: result,
			YourMove:     true,
		}
		actualBytes, _ := json.Marshal(response)

		// Zeroed result for the opponent
		masked := make([]RevealedCell, len(result.Positions))
		for i, cell := range result.Positions {
			masked[i] = RevealedCell{
				X:     cell.X,
				Y:     cell.Y,
				Value: 0, // hide real value
			}
		}
		opponentResponse := struct {
			RevealResult
			YourMove bool `json:"yourMove"`
		}{
			RevealResult: RevealResult{
				Positions: masked,
				HitMine:   result.HitMine,
				Won:       result.Won,
			},
			YourMove: false,
		}
		maskedBytes, _ := json.Marshal(opponentResponse)
		mg.mu.Unlock()

		// Send to both players
		log.Printf("Sending actual result to player %d", playerIndex)
		mg.players[playerIndex].send <- actualBytes

		otherIndex := 1 - playerIndex
		log.Printf("Sending masked result to player %d", otherIndex)
		mg.players[otherIndex].send <- maskedBytes
	}
}
