package server

import (
	"net"
)

// Client is used to get access to connection
type Client struct {
	Username string
	Conn     net.Conn
}

// CreateClient returns an instance of Client
func CreateClient(username string, conn net.Conn) *Client {
	return &Client{
		Username: username,
		Conn:     conn,
	}
}

// CreateRoom returns an instance of Room
func CreateRoom(name string) *Room {
	return &Room{
		Name:       name,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Messages:   make(chan string),
	}
}

// Room is manager of actions of clients
type Room struct {
	Name       string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Messages   chan string
}

// Run starts a room for message exchange
func (r *Room) Run() {
	for {
		select {
		case msg := <-r.Messages:
			r.Broadcast(msg)
		case client := <-r.Register:
			r.Clients[client.Username] = client
			r.Broadcast(client.Username + " joined to the room")
		case client := <-r.Unregister:
			client.Conn.Close()
			delete(r.Clients, client.Username)
			r.Broadcast(client.Username + "left the room")
		}
	}
}

// Broadcast sends messages for all clients in the room
func (r *Room) Broadcast(msg string) {
	for _, client := range r.Clients {
		_, err := client.Conn.Write([]byte(msg))
		if err != nil {
			r.Unregister <- client
		}
	}
}

// ClientCount returns count of clients
func (r *Room) ClientCount() int {
	return len(r.Clients)
}
