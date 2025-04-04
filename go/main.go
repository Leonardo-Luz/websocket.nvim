package main

import (
	"log"
	"os"

	"github.com/leonardo-luz/websocket.nvim/client"
	"github.com/leonardo-luz/websocket.nvim/server"
	"github.com/neovim/go-client/nvim"
)

func main() {
	_client, err := nvim.New(os.Stdin, os.Stdout, os.Stdout, log.Printf)
	if err != nil {
		log.Fatalf("Failed to create a Nvim Client...")
	}

	_client.RegisterHandler("startWsServer", server.StartServer)

	_client.RegisterHandler("startWsClient", client.StartClient)
	_client.RegisterHandler("writeWsClient", client.WriteLines)

	if err := _client.Serve(); err != nil {
		log.Fatal(err)
	}
}
