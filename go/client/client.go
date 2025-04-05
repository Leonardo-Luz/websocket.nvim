package client

import (
	"fmt"
	"regexp"
	"strings"

	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/neovim/go-client/nvim"
)

func receiveLines(client *nvim.Nvim, bufnum int, conn *websocket.Conn) {
	for {
		_, response, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
			os.Exit(1)
		}

		handleNewUserConn(response, client, bufnum)

		handleNewLine(response, client, bufnum)

		handleUpdateLines(response, client, bufnum, conn)
	}
}

func handleNewUserConn(response []byte, client *nvim.Nvim, bufnum int) {
	// Check if the message matches the pattern containing role, user, and lines
	re := regexp.MustCompile(`^EWFSDNASKDJNQQWEOuser:(.*)`)
	matches := re.FindStringSubmatch(string(response))

	if len(matches) > 0 {
		// Extract all lines for the user
		userLines := matches[1]

		// Split the userLines string into individual lines
		arrayUserlines := strings.Split(userLines, "Ef232wefeEFAwdEFF")

		var lines [][]byte

		// Convert each line into []byte and append to the lines array
		for _, userLine := range arrayUserlines {
			lines = append(lines, []byte(userLine))
		}

		// Update the entire buffer with the user's lines
		if err := client.SetBufferLines(nvim.Buffer(bufnum), 0, -1, false, lines); err != nil {
			log.Fatal(err)
		}
	}
}

func handleNewLine(response []byte, client *nvim.Nvim, bufnum int) {
	// If the message doesn't contain the role/user pattern, check for line modifications
	re := regexp.MustCompile(`^start=(\w+)end=(\w+)lines\[(.*)\]`)
	matches := re.FindStringSubmatch(string(response))

	if len(matches) > 0 {
		// startnumStr := matches[1]
		// endnumStr := matches[2]
		lines := matches[3]

		// Convert start and end numbers to integers
		// startnum, err := strconv.Atoi(startnumStr)
		// if err != nil {
		// 	log.Fatalf("Error converting start number: %v", err)
		// 	return
		// }
		//
		// endnum, err := strconv.Atoi(endnumStr)
		// if err != nil {
		// 	log.Fatalf("Error converting end number: %v", err)
		// 	return
		// }

		// Split the userLines string into individual lines
		arrayUserlines := strings.Split(lines, "Ef232wefeEFAwdEFF")

		var linesBye [][]byte

		// Convert each line into []byte and append to the lines array
		for _, userLine := range arrayUserlines {
			linesBye = append(linesBye, []byte(userLine))
		}

		// Update only the specified line(s)
		if err := client.SetBufferLines(nvim.Buffer(bufnum), 0, -1, false, linesBye); err != nil {
			log.Fatal(err)
		}
	}
}

func handleUpdateLines(response []byte, client *nvim.Nvim, bufnum int, conn *websocket.Conn) {
	if string(response) == "FSADnkj34sd1QQW3O" {
		linesByte, err := client.BufferLines(nvim.Buffer(bufnum), 0, -1, false)
		if err != nil {
			log.Fatal(err)
		}

		var lines []string

		for lineid := range linesByte {
			lines = append(lines, string(linesByte[lineid]))
		}

		message := fmt.Sprintf("wfeFJEWO23ASD12oilines[%s]", strings.Join(lines, "Ef232wefeEFAwdEFF"))

		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func WriteLines(start int, end int, line string, url string) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket server: %v", err)
		os.Exit(1)
	}

	message := fmt.Sprintf("start=%dend=%dlines[%s]", start, end, line)

	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
		os.Exit(1)
	}

	conn.Close()
}

func StartClient(client *nvim.Nvim, bufnum int, url string) {
	// WebSocket server URL
	// url := "ws://localhost:8080/ws"

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Fatalf("Error connecting to WebSocket server: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		err = conn.WriteMessage(websocket.TextMessage, []byte("EWFSDNASKDJNQQWEOrole"))
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
			os.Exit(1)
		}

		receiveLines(client, bufnum, conn)
	}()
}
