package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	host         string
	port         string
	addr         string
	OnMessage    func(client_addr string, server_message []byte)
	OnConnect    func(client_addr string, server_message []byte)
	OnDisconnect func(client_addr string, server_message []byte)
	OnOpen       func(client_addr string, server_message []byte)
	OnClose      func()
	OnError      func(client_addr string, server_message []byte)
}

func (client *Client) GetAddr() string {
	return client.addr
}
func (client *Client) SetOnMessage(onMessage func(client_addr string, server_message []byte)) {
	client.OnMessage = onMessage
}
func (client *Client) SetOnConnect(onConnect func(client_addr string, server_message []byte)) {
	client.OnConnect = onConnect
}
func (client *Client) SetOnDisconnect(onDisconnect func(client_addr string, server_message []byte)) {
	client.OnDisconnect = onDisconnect
}
func (client *Client) SetOnOpen(onOpen func(client_addr string, server_message []byte)) {
	client.OnOpen = onOpen
}
func (client *Client) SetOnClose(onClose func()) {
	client.OnClose = onClose
}
func (client *Client) SetOnError(onError func(client_addr string, server_message []byte)) {
	client.OnError = onError
}

func (client *Client) SendMessage(conn *websocket.Conn, messageClient string) {
	message := []byte("Hello, WebSocket Server!")
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Fatal("Failed to send message:", err)
	}
}

func (client *Client) handleMessages(conn *websocket.Conn, interrupt chan os.Signal) {
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			client.OnMessage(client.GetAddr(), message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(t.String())); err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("write close:", err)
				return
			}
			client.OnClose()

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (client *Client) ConnectToServer() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(client.GetAddr(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}

	go client.handleMessages(conn, interrupt)
}

func NewClient(host string, port string,
	OnMessage func(client_addr string, server_message []byte),
	OnConnect func(client_addr string, server_message []byte),
	OnDisconnect func(client_addr string, server_message []byte),
	OnOpen func(client_addr string, server_message []byte),
	OnClose func(),
	OnError func(client_addr string, server_message []byte),
) *Client {
	return &Client{
		host:         host,
		port:         port,
		addr:         fmt.Sprintf("ws://%s:%s/ws", host, port),
		OnMessage:    OnMessage,
		OnConnect:    OnConnect,
		OnDisconnect: OnDisconnect,
		OnOpen:       OnOpen,
		OnClose:      OnClose,
		OnError:      OnError,
	}
}

func NewBasicClient(host string, port string) *Client {
	return &Client{
		host: host,
		port: port,
		addr: fmt.Sprintf("ws://%s:%s/ws", host, port),
		OnMessage: func(client_addr string, server_message []byte) {
			log.Println("server message: ", server_message)
		},
		OnConnect: func(client_addr string, server_message []byte) {
			log.Println("server message: ", server_message)
		},
		OnDisconnect: func(client_addr string, server_message []byte) {
			log.Println("server message: ", server_message)
		},
		OnOpen: func(client_addr string, server_message []byte) {
			log.Println("server message: ", server_message)
		},
		OnClose: func() {
			log.Println("Disconnected from server...")
		},
		OnError: func(client_addr string, server_message []byte) {
			log.Println("server message: ", server_message)
		},
	}
}
