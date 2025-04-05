package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

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
	adminConn *websocket.Conn
}

// Creates a new instance of Server
func NewServer(port string, lines []string, adminCode string) *Server {
	return &Server{
		port:      port,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 100),
		lines:     lines,
		adminCode: adminCode,
		adminConn: nil,
	}
}

const ROLE_CODE = "EWFSDNASKDJNQQWEO"
const JOIN_CODE = "Ef232wefeEFAwdEFF"
const UPDATE_CODE = "FSADnkj34sd1QQW3O"
const GET_LINES_CODE = "wfeFJEWO23ASD12oi"

func (s *Server) readLoop(conn *websocket.Conn) {
	for {
		//on message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("read error: ", err)
			continue
		}

		s.OnMessage(msg, conn)
	}
}

func (s *Server) OnMessage(msg []byte, conn *websocket.Conn) {
	re := regexp.MustCompile(`^wfeFJEWO23ASD12oilines\[(.*)\]`)
	matches := re.FindStringSubmatch(string(msg))

	if len(matches) > 0 {
		linesStr := matches[1]
		lines := strings.Split(linesStr, "Ef232wefeEFAwdEFF")

		s.lines = lines
	}

	if string(msg) == (ROLE_CODE + "role") {
		if len(s.clients) == 1 {
			s.adminConn = conn
			conn.WriteMessage(websocket.TextMessage, []byte(ROLE_CODE+"admin:"+s.adminCode))
		} else {
			s.adminConn.WriteMessage(websocket.TextMessage, []byte(UPDATE_CODE))

			time.Sleep(200 * time.Millisecond)

			conn.WriteMessage(websocket.TextMessage, []byte(ROLE_CODE+"user:"+strings.Join(s.lines, JOIN_CODE)))
		}
	} else {
		s.broadcast <- msg
	}
}

func (s *Server) MessageHandler() {
	for {
		message := <-s.broadcast

		mutex.Lock()
		for client := range s.clients {
			s.onBroadcast(message, client)
		}
		mutex.Unlock()
	}
}

func (s *Server) onBroadcast(message []byte, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, message)
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// fmt.Println("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()

	s.clients[conn] = true

	// onConnection()

	s.readLoop(conn)
}

// func (s *Server) onConnection(){}

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
