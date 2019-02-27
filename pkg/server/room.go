package server

import (
	"fmt"
	"net"
	"os"

	uuid "github.com/satori/go.uuid"
)

// Client ...
type Client struct {
	ID       string
	Username string
	Conn     net.Conn
}

// CreateClient ...
func CreateClient(username string, conn net.Conn) *Client {
	id := uuid.NewV4().String()
	return &Client{
		ID:       id,
		Username: username,
		Conn:     conn,
	}
}

// CreateRoom ...
func CreateRoom() *Room {
	id := uuid.NewV4().String()
	return &Room{
		ID:       id,
		Register: make(chan net.Conn),
		Clients:  make(map[string]*Client),
		Messages: make(chan string),
	}
}

// Room ...
type Room struct {
	ID       string
	Register chan net.Conn
	Clients  map[string]*Client
	Messages chan string
}

// Run ...
func (r *Room) Run() {
	for {
		select {
		case msg := <-r.Messages:
			fmt.Println(msg, " is recieved")
			r.broadcast(msg)
		}
	}
}

// Broadcast ...
func (r *Room) broadcast(msg string) {
	for _, client := range r.Clients {
		in, err := client.Conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("Error when send to client: %d\n", in)
			os.Exit(0)
		}
	}
}
