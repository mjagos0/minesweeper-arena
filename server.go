// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// type Player struct {
// 	conn *websocket.Conn
// 	send chan []byte
// }

// type GameInit struct {
// 	Type   string `json:"type"`
// 	Width  int    `json:"width"`
// 	Height int    `json:"height"`
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// var (
// 	matchQueue []*Player
// 	queueLock  sync.Mutex
// )

// func main() {
// 	// Serve static files (HTML, JS, CSS)
// 	http.Handle("/", http.FileServer(http.Dir("static")))

// 	// WebSocket handler
// 	http.HandleFunc("/ws", handleWebSocket)

// 	fmt.Println("Server listening on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade HTTP to WebSocket
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade failed:", err)
// 		return
// 	}

// 	// Create player object with buffered send channel
// 	player := &Player{
// 		conn: conn,
// 		send: make(chan []byte, 10), // buffered to avoid blocking
// 	}

// 	// Start the write loop immediately
// 	go writeLoop(player)

// 	// Try to match the player
// 	matchPlayer(player)

// 	// Optional: Keep read loop to detect disconnects, etc.
// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			break // connection closed
// 		}
// 	}

// 	conn.Close()
// }

// func writeLoop(p *Player) {
// 	for msg := range p.send {
// 		err := p.conn.WriteMessage(websocket.TextMessage, msg)
// 		if err != nil {
// 			log.Println("Write error:", err)
// 			break
// 		}
// 	}
// }

// func matchPlayer(p *Player) {
// 	queueLock.Lock()
// 	defer queueLock.Unlock()

// 	if len(matchQueue) > 0 {
// 		opponent := matchQueue[0]
// 		matchQueue = matchQueue[1:]

// 		startGame(p, opponent)
// 	} else {
// 		matchQueue = append(matchQueue, p)
// 	}
// }

// func startGame(p1, p2 *Player) {
// 	gameData := GameInit{
// 		Type:   "game-start",
// 		Width:  10,
// 		Height: 10,
// 	}
// 	jsonData, err := json.Marshal(gameData)
// 	if err != nil {
// 		log.Println("JSON marshal failed:", err)
// 		return
// 	}

// 	// Send to both players
// 	p1.send <- jsonData
// 	p2.send <- jsonData
// }
