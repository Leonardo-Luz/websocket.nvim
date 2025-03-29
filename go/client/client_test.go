package client

import "testing"

var client Client

// TODO: implement tests for client

func setup() {
	client = *Newclient("localhost", "8080",
		func(client_addr string, client_message []byte) {},
		func(client_addr string, client_message []byte) {},
		func(client_addr string, client_message []byte) {},
		func(client_addr string, client_message []byte) {},
		func(client_addr string, client_message []byte) {},
		func(client_addr string, client_message []byte) {},
	)
}

func testOnMessage(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}

func testOnConnect(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}

func testOnDisconnect(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}

func testOnOpen(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}

func testOnClose(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}

func testOnError(t *testing.T) {
	setup()

	if 2+2 != 4 {
		t.Errorf("Math Error!")
	}
}
