package main

import (
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

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	log.Println("Creating new web socket")
	player := &Player{
		conn: conn,
		send: make(chan []byte, 10),
	}

	go writeLoop(player)
	matchPlayer(player)
}

func writeLoop(p *Player) {
	for msg := range p.send {
		err := p.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func matchPlayer(p *Player) {
	queueLock.Lock()
	defer queueLock.Unlock()

	if len(matchQueue) > 0 {
		opponent := matchQueue[0]
		matchQueue = matchQueue[1:]

		log.Println("Starting multiplayer game")
		startMultiplayerGame(p, opponent)
	} else {
		matchQueue = append(matchQueue, p)
		log.Println("Player added to match queue")
	}
}

func startMultiplayerGame(p1, p2 *Player) {
	game := NewMultiplayerGame(p1, p2)
	game.Start()
}
