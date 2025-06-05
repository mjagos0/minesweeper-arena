package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	matchQueue []*Player
	queueLock  sync.Mutex
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Server listening on 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	player := &Player{
		conn: conn,
		send: make(chan []byte, 10),
	}

	go writeLoop(player)
	go readLoop(player)

	joinQueue(player)
}

func writeLoop(p *Player) {
	for msg := range p.send {
		err := p.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
	p.conn.Close()
	p.closed = true
}

func readLoop(p *Player) {
	defer func() {
		p.conn.Close()
		p.closed = true

		if p.game != nil {
			p.game.HandleDisconnect(p.playerId)
		}
	}()

	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("Read error):", err)
			break
		}
		processMsg(p, msg)
	}
}

func joinQueue(p *Player) {
	queueLock.Lock()
	defer queueLock.Unlock()

	for len(matchQueue) > 0 {
		opponent := matchQueue[0]
		matchQueue = matchQueue[1:]

		if opponent.closed {
			continue
		} else {
			log.Println("Starting multiplayer game")
			startMultiplayerGame(p, opponent)
			return
		}
	}

	log.Println("Player added to match queue")
	matchQueue = append(matchQueue, p)
}

func processMsg(p *Player, msg []byte) {
	var move ClientMove
	if err := json.Unmarshal(msg, &move); err != nil {
		log.Println("Invalid move JSON:", err)
		return
	}

	if p.game != nil {
		p.game.ProcessPlayerMove(p.playerId, move)
	}
}

func startMultiplayerGame(p1, p2 *Player) {
	game := NewMultiplayerGame(p1, p2)
	p1.game = game
	p2.game = game
	p1.playerId = 0
	p2.playerId = 1

	game.Start()
}
