package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

const (
	PlayerOne = 0
	PlayerTwo = 1
	NoPlayer  = 2
)

type Player struct {
	conn     *websocket.Conn
	send     chan []byte
	closed   bool
	game     *MultiplayerGame
	playerId int
}

type MultiplayerGame struct {
	board   *Board
	games   [2]*Game
	players [2]*Player
	startX  int
	startY  int
	winner  int
}

type StartMessage struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Mines  int `json:"mines"`
	StartX int `json:"startX"`
	StartY int `json:"startY"`
}

type ClientMove struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Flag bool `json:"flag"`
}

type MoveResponse struct {
	CellUpdates []CellUpdate `json:"cellUpdates"`
	YourMove    bool         `json:"yourMove"`
	Win         bool         `json:"win"`
	Lose        bool         `json:"lose"`
}

func NewMultiplayerGame(p1 *Player, p2 *Player) *MultiplayerGame {
	board := NewBoard(9, 9, 10)
	board.Print()
	startY, startX := FindRandomSafeStartField(board)

	return &MultiplayerGame{
		board: board,
		games: [2]*Game{
			NewGame(board),
			NewGame(board),
		},
		players: [2]*Player{p1, p2},
		startX:  startX,
		startY:  startY,
		winner:  NoPlayer,
	}
}

func FindRandomSafeStartField(b *Board) (int, int) {
	candidates := make([][2]int, 0)
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if b.At(x, y) == 0 {
				candidates = append(candidates, [2]int{x, y})
			}
		}
	}

	if len(candidates) == 0 {
		for y := 0; y < b.Height; y++ {
			for x := 0; x < b.Width; x++ {
				if b.At(x, y) != MINE {
					candidates = append(candidates, [2]int{x, y})
				}
			}
		}
	}

	start := candidates[rand.Intn(len(candidates))]
	return start[0], start[1]
}

func (mg *MultiplayerGame) Start() {
	startMessage := StartMessage{
		Width:  mg.board.Width,
		Height: mg.board.Height,
		Mines:  mg.board.NumMines,
		StartX: mg.startX,
		StartY: mg.startY,
	}

	msg, _ := json.Marshal(startMessage)
	mg.players[0].send <- msg
	mg.players[1].send <- msg
}

func (mg *MultiplayerGame) ProcessPlayerMove(playerIndex int, cmove ClientMove) {
	fmt.Printf("Player %d sent move: %+v\n", playerIndex, cmove)

	cellUpdates := mg.games[playerIndex].processMove(cmove.X, cmove.Y, cmove.Flag)

	maskedCellUpdates := make([]CellUpdate, len(cellUpdates))
	for i, cell := range cellUpdates {
		value := uint8(0)
		if cell.Value == WIPE {
			value = WIPE
		}
		maskedCellUpdates[i] = CellUpdate{
			X:     cell.X,
			Y:     cell.Y,
			Value: value,
		}
	}

	if mg.games[playerIndex].finished && mg.games[playerIndex].won {
		mg.winner = playerIndex
	} else if mg.games[playerIndex].finished && mg.games[1-playerIndex].finished {
		if mg.games[playerIndex].revealed < mg.games[1-playerIndex].revealed {
			mg.winner = 1 - playerIndex
		} else {
			mg.winner = playerIndex
		}
	}

	respOpponent := MoveResponse{
		CellUpdates: maskedCellUpdates,
		YourMove:    false,
		Win:         1-playerIndex == mg.winner,
		Lose:        playerIndex == mg.winner,
	}
	mg.players[1-playerIndex].send <- toJson(respOpponent)

	if mg.games[playerIndex].finished || mg.winner != NoPlayer {
		mg.games[playerIndex].revealBoard(&cellUpdates)
	}

	respPlayer := MoveResponse{
		CellUpdates: cellUpdates,
		YourMove:    true,
		Win:         playerIndex == mg.winner,
		Lose:        1-playerIndex == mg.winner,
	}
	mg.players[playerIndex].send <- toJson(respPlayer)
}

func (mg *MultiplayerGame) HandleDisconnect(playerIndex int) {
	fmt.Printf("Player %d disconnected\n", playerIndex)

	if mg.winner == NoPlayer {
		fmt.Printf("Game has no winner. Declearing the other player winner\n")
		mg.winner = 1 - playerIndex

		resp := MoveResponse{
			CellUpdates: []CellUpdate{},
			YourMove:    true,
			Win:         true,
			Lose:        false,
		}
		mg.players[1-playerIndex].send <- toJson(resp)
	}
}

func toJson(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
