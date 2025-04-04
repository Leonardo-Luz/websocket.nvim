package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var mutex = &sync.Mutex{}

// Upgrades connection from http to ws
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	port      string
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	lines     []string
	adminCode string
}

// Creates a new instance of Server
func NewServer(port string, lines []string, adminCode string) *Server {
	return &Server{
		port:      port,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 100),
		lines:     lines,
		adminCode: adminCode,
	}
}

const ROLE_CODE = "EWFSDNASKDJNQQWEO"
const JOIN_CODE = "Ef232wefeEFAwdEFF"

func (s *Server) readLoop(ws *websocket.Conn) {
	for {
		//on message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			// fmt.Println("read error: ", err)
			continue
		}

		if string(msg) == (ROLE_CODE + "role") {
			ws.WriteMessage(websocket.TextMessage, []byte(ROLE_CODE+"user:"+strings.Join(s.lines, JOIN_CODE)))

			// if len(s.clients) == 1 {
			// 	ws.WriteMessage(websocket.TextMessage, []byte(ROLE_CODE+"admin:"+s.adminCode))
			// } else {
			// 	ws.WriteMessage(websocket.TextMessage, []byte(ROLE_CODE+"user:"+strings.Join(s.lines, JOIN_CODE)))
			// }
		} else {
			s.broadcast <- msg
		}
	}
}

func (s *Server) MessageHandler() {
	for {
		message := <-s.broadcast

		mutex.Lock()
		for client := range s.clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
		mutex.Unlock()
	}
}

func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// fmt.Println("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()

	server.clients[conn] = true

	server.readLoop(conn)
}

func generateAdminCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil // Use URL-safe encoding
}

// Starts the WebSocket Server
func StartServer(lines []string, port string) {
	admin_code, err := generateAdminCode(32)
	if err != nil {
		return
	}

	server := NewServer(port, lines, admin_code)

	// fmt.Println("Starting the WebSocket server... port: " + server.port)

	go func() {
		http.HandleFunc("/ws", server.WebsocketHandler)

		go server.MessageHandler()

		err := http.ListenAndServe(":"+server.port, nil)
		if err != nil {
			// fmt.Sprintf("Error starting server: %v", err)
		}
	}()
}
