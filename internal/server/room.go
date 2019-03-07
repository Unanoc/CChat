package server

import (
	"chat/pkg/queue"
	"net"

	"github.com/fatih/color"
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
		Storage:    queue.CreateQueue(128),
	}
}

// Room is manager of actions of clients
type Room struct {
	Name       string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Messages   chan string
	Storage    *queue.Queue
}

// Run starts a room for message exchange
func (r *Room) Run() {
	var msg string

	for {
		select {
		case msg = <-r.Messages:
			msg = msg + "\n"
		case client := <-r.Register:
			r.Clients[client.Username] = client
			if history := r.Storage.FromHeadToTail(); history != nil {
				r.SendHistory(history, client)
			}
			msg = color.HiBlackString("%s joined to the room [%s]\n", client.Username, r.Name)
		case client := <-r.Unregister:
			delete(r.Clients, client.Username)
			msg = color.HiBlackString("%s left the room [%s]\n", client.Username, r.Name)
		}

		r.Broadcast(msg)
		r.Storage.Push(msg)
	}
}

// Broadcast sends messages for all clients in the room
func (r *Room) Broadcast(msg string) {
	for _, client := range r.Clients {
		r.SendToClient(msg, client)
	}
}

// SendToClient sends message to client
func (r *Room) SendToClient(msg string, client *Client) {
	_, err := client.Conn.Write([]byte(msg))
	if err != nil {
		r.Unregister <- client
	}
}

// SendHistory sends last 128 messages of the room to client
func (r *Room) SendHistory(storage []string, client *Client) {
	for _, msg := range storage {
		r.SendToClient(msg, client)
	}
}

// SendClientList sends list of clients of the room
func (r *Room) SendClientList(reciever *Client) {
	for _, client := range r.Clients {
		r.SendToClient(color.HiYellowString("%s\n", client.Username), reciever)
	}
}

// ClientCount returns count of clients
func (r *Room) ClientCount() int {
	return len(r.Clients)
}
