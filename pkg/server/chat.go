package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// CreateChat returns an instance of Chat
func CreateChat() *Chat {
	return &Chat{
		Register: make(chan net.Conn),
		Rooms:    make(map[string]*Room),
	}
}

// Chat is the main struct of chat
type Chat struct {
	Register chan net.Conn
	Rooms    map[string]*Room
	sync.Mutex
}

// Run handles incoming connections
func (c *Chat) Run() {
	for {
		conn := <-c.Register
		log.Printf("New connection: [%v]", conn.RemoteAddr())

		go c.ProcessConn(conn)
	}
}

// ProcessConn initialises clients's connection
func (c *Chat) ProcessConn(conn net.Conn) {
	data := make([]byte, 254)

	// Getting client's nickname
	usernameLen, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %s has not been connected", conn.RemoteAddr())
		conn.Close()
		return
	}
	username := string(data[:usernameLen])

	// Getting room's name
	roomLen, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %s has not been connected", conn.RemoteAddr())
		conn.Close()
		return
	}
	roomname := string(data[:roomLen])

	c.ProcessRoom(roomname)

	// Joining the room
	c.Lock()
	defer c.Unlock()

	room := c.Rooms[roomname]
	var client *Client
	if c.IsUsernameUniq(username, roomname) {
		client = CreateClient(username, conn)
		room.Register <- client
	} else {
		_, err := conn.Write([]byte("This nickname is already exists in room"))
		if err != nil {
			log.Println("Error when send to client")
		}
		conn.Close()
		return
	}

	go c.ListenClient(client, room)
}

// ProcessRoom creates room if room does not exists
func (c *Chat) ProcessRoom(roomname string) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.Rooms[roomname]; !exists {
		c.Rooms[roomname] = CreateRoom(roomname)
		room := c.Rooms[roomname]
		go room.Run()
	}
}

// IsUsernameUniq checks if username is uniq in room
func (c *Chat) IsUsernameUniq(username, roomname string) bool {
	_, exists := c.Rooms[roomname].Clients[username]
	return !exists
}

// ListenClient gets client's messages
func (c *Chat) ListenClient(client *Client, room *Room) {
	data := make([]byte, 254)

	for {
		msgLen, err := client.Conn.Read(data)
		if err != nil {
			room.Unregister <- client
			log.Printf("Client %s quit", client.Conn.RemoteAddr())
			client.Conn.Close()
			return
		}

		rawMessage := string(data[:msgLen])
		switch rawMessage {
		case "/quit":
			room.Unregister <- client
			log.Printf("Client %s quit", client.Conn.RemoteAddr())
			client.Conn.Close()
			return
		case "/change_room":
			room.Unregister <- client
			c.ProcessConn(client.Conn)
			return
		default:
			room.Messages <- fmt.Sprintf("(%s) %s: %s", room.Name, client.Username, rawMessage)
		}
	}
}

// CleanChat removes rooms with zero clients once a minute
func (c *Chat) CleanChat() {
	for {
		time.Sleep(1 * time.Minute)
		c.Lock()
		for _, room := range c.Rooms {
			if room.ClientCount() == 0 {
				c.RemoveRoom(room)
			}
		}
		c.Unlock()
	}
}

// RemoveRoom removes the room
func (c *Chat) RemoveRoom(room *Room) {
	if room != nil {
		if room.ClientCount() == 0 {
			log.Printf("Room [%s] has been destroyed", room.Name)
			delete(c.Rooms, room.Name)
		}
	}
}
