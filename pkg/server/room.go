package server

import (
	"net"
)

// Client ...
type Client struct {
	Username string
	Conn     net.Conn
}

// CreateClient ...
func CreateClient(username string, conn net.Conn) *Client {
	return &Client{
		Username: username,
		Conn:     conn,
	}
}

// CreateRoom ...
func CreateRoom(name string) *Room {
	return &Room{
		Name:     name,
		Register: make(chan net.Conn),
		Clients:  make(map[string]*Client),
		Messages: make(chan string),
	}
}

// Room ...
type Room struct {
	Name     string
	Register chan net.Conn
	Clients  map[string]*Client
	Messages chan string
}

// Run ...
func (r *Room) Run() {
	for {
		select {
		case msg := <-r.Messages:
			r.broadcast(msg)
		}
	}
}

// Broadcast ...
func (r *Room) broadcast(msg string) {
	for _, client := range r.Clients {
		_, err := client.Conn.Write([]byte(msg))
		if err != nil {
			r.RemoveClient(client.Username)
		}
	}
}

// ClientCount ...
func (r *Room) ClientCount() int {
	return len(r.Clients)
}

// RemoveClient ...
func (r *Room) RemoveClient(username string) {
	if client, ok := r.Clients[username]; ok {
		delete(r.Clients, client.Username)
	}
}
