package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Clients connected to the server
var clients = make(map[*websocket.Conn]bool)

// Messages broadcasted
var broadcast = make(chan []byte, 100) // Buffered channel

// Protect clients map
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
	host         string
	port         string
	addr         string // host:port
	OnMessage    func(client_addr string, client_message []byte)
	OnConnect    func(client_addr string)
	OnDisconnect func(client_addr string)
	OnOpen       func(addr string)
	OnClose      func()
	OnError      func(errorMsg string)
}

func (server *Server) GetAddr() string {
	return server.addr
}
func (server *Server) SetOnMessage(onMessage func(client_addr string, client_message []byte)) {
	server.OnMessage = onMessage
}
func (server *Server) SetOnConnect(onConnect func(client_addr string)) {
	server.OnConnect = onConnect
}
func (server *Server) SetOnDisconnect(onDisconnect func(client_addr string)) {
	server.OnDisconnect = onDisconnect
}
func (server *Server) SetOnOpen(onOpen func(addr string)) {
	server.OnOpen = onOpen
}
func (server *Server) SetOnClose(onClose func()) {
	server.OnClose = onClose
}
func (server *Server) SetOnError(onError func(errorMsg string)) {
	server.OnError = onError
}

// Wait for a client to send a message
func (server *Server) ServerListener(conn *websocket.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			server.OnError(fmt.Sprintf("Error reading message: %v", err))
			break
		}

		server.OnMessage(conn.LocalAddr().String(), message)
		broadcast <- message
	}
}

// Upgrades the connection and starts listening
func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.OnError("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	go server.ServerListener(conn)
}

// Sends received messages to all connected clients
func (server *Server) MessageHandler() {
	for {
		message := <-broadcast

		mutex.Lock()
		for client := range clients {
			server.SendMessage(client, message)
		}
		mutex.Unlock()
	}
}

// Send message to connected client
func (server *Server) SendMessage(client *websocket.Conn, message []byte) {
	err := client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error sending message to client %s: %v", client.LocalAddr().String(), err)
		client.Close()
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		server.OnDisconnect(client.LocalAddr().String())
	}
}

// Starts the WebSocket Server
func (server *Server) StartServer() {
	log.Println("Starting the WebSocket server...") // Log when server starts
	server.OnOpen(server.GetAddr())

	http.HandleFunc("/ws", server.WebsocketHandler)

	go server.MessageHandler() // Ensure this is in a goroutine to not block

	err := http.ListenAndServe(":"+server.port, nil)
	if err != nil {
		server.OnError(fmt.Sprintf("Error starting server: %v", err))
	}
}

// Creates a new instance of Server
func NewServer(host string, port string,
	OnMessage func(client_addr string, client_message []byte),
	OnConnect func(client_addr string),
	OnDisconnect func(client_addr string),
	OnOpen func(addr string),
	OnClose func(),
	OnError func(errorMsg string),
) *Server {
	return &Server{
		host:         host,
		port:         port,
		addr:         fmt.Sprintf("%s:%s", host, port),
		OnMessage:    OnMessage,
		OnConnect:    OnConnect,
		OnDisconnect: OnDisconnect,
		OnOpen:       OnOpen,
		OnClose:      OnClose,
		OnError:      OnError,
	}
}

// Creates a new instance of Server with default methods
func NewBasicServer(host string, port string) *Server {
	return &Server{
		host: host,
		port: port,
		addr: fmt.Sprintf("%s:%s", host, port),
		OnMessage: func(client_addr string, client_message []byte) {
			log.Printf("received message: %s from %s", client_message, client_addr)
		},
		OnConnect: func(client_addr string) {
			log.Printf("Client %s connected to the server", client_addr)
		},
		OnDisconnect: func(client_addr string) {
			log.Printf("Client disconnected: %s", client_addr)
		},
		OnOpen: func(addr string) {
			log.Printf("Listening at %s", addr)
		},
		OnClose: func() {
			log.Printf("Server Closed")
		},
		OnError: func(errorMsg string) {
			log.Printf(errorMsg)
		},
	}
}
