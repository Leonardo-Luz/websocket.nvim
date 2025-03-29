package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leonardo-luz/websocket.nvim/client"
	"github.com/leonardo-luz/websocket.nvim/server"
	"github.com/neovim/go-client/nvim"
)

var _client *client.Client
var _server *server.Server

var n_client *nvim.Nvim

func nvimClient() *nvim.Nvim {
	client, err := nvim.New(os.Stdin, os.Stdout, os.Stdout, log.Printf)
	if err != nil {
		log.Fatalf("Failed to create a Nvim Client...")
	}

	return client
}

func getNvimClient() *nvim.Nvim {
	return n_client
}

func nvimPrint(message string) error {
	return getNvimClient().WriteOut(message)
}

func setClient(host string, port string) {
	_client = client.NewBasicClient(host, port)
}

func setServer(host string, port string) {
	_server = server.NewServer(host, port, func(client_addr string, client_message []byte) {
		nvimPrint(fmt.Sprintf("Client: %s message: %s\n", client_addr, client_message))
	}, func(client_addr string) {
		nvimPrint((fmt.Sprintf("Client %s connected\n", client_addr)))
	}, func(client_addr string) {
		nvimPrint((fmt.Sprintf("Client %s Disconnected\n", client_addr)))
	}, func(addr string) {
		nvimPrint((fmt.Sprintf("Server listenning at: %s\n", addr)))
	}, func() {
		nvimPrint((fmt.Sprintf("Server Closed\n")))
	}, func(errorMsg string) {
		nvimPrint((fmt.Sprintf("Server Error: %s\n", errorMsg)))
	})

	_server.StartServer()
}

func main() {
	n_client = nvimClient()

	n_client.RegisterHandler("newWsServer", setServer)
	n_client.RegisterHandler("setWsServerOnMessage", _client.SetOnMessage)
	n_client.RegisterHandler("startWsServer", _server.StartServer)
	n_client.RegisterHandler("sendMessage", _server.SendMessage)

	n_client.RegisterHandler("newWsClient", setClient)
	n_client.RegisterHandler("setWsClientOnMessage", _client.SetOnMessage)
	n_client.RegisterHandler("connectToWsServer", _client.ConnectToServer)
	n_client.RegisterHandler("sendMessageToServer", _client.SendMessage)

	if err := n_client.Serve(); err != nil {
		log.Fatal(err)
	}
}

// func stringsToBytes(lines []string) [][]byte {
// 	bytes := make([][]byte, len(lines))
// 	for i, line := range lines {
// 		bytes[i] = []byte(line)
// 	}
// 	return bytes
// }
//
// // / live plugin
// func set_lines(client *nvim.Nvim, bufnum int, start int, end int, lines []string) {
// 	if err := client.SetBufferLines(nvim.Buffer(bufnum), start, end, false, stringsToBytes(lines)); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// client, err := nvim.New(os.Stdin, os.Stdout, os.Stdout, log.Printf)
// if err != nil {
// 	log.Fatalf("Failed to create a Nvim Client...")
// }
//
// http.HandleFunc("ws", websocketHandler)
//
// client.RegisterHandler("hello", set_lines)
//
// if err := client.Serve(); err != nil {
// 	log.Fatal(err)
// }

// func stringsToBytes(lines []string) [][]byte {
// 	bytes := make([][]byte, len(lines))
// 	for i, line := range lines {
// 		bytes[i] = []byte(line)
// 	}
// 	return bytes
// }
//
// // / live plugin
// func set_lines(client *nvim.Nvim, bufnum int, start int, end int, lines []string) {
// 	if err := client.SetBufferLines(nvim.Buffer(bufnum), start, end, false, stringsToBytes(lines)); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// client, err := nvim.New(os.Stdin, os.Stdout, os.Stdout, log.Printf)
// if err != nil {
// 	log.Fatalf("Failed to create a Nvim Client...")
// }
//
// http.HandleFunc("ws", websocketHandler)
//
// client.RegisterHandler("hello", set_lines)
//
// if err := client.Serve(); err != nil {
// 	log.Fatal(err)
// }
