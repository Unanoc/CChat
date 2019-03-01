package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// CreateChat ...
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

// Run ...
func (c *Chat) Run() {
	for {
		conn := <-c.Register
		log.Printf("New connection: [%v]", conn.RemoteAddr())

		go c.ProcessConn(conn)
	}
}

// ProcessConn ...
func (c *Chat) ProcessConn(conn net.Conn) {
	data := make([]byte, 254)

	// Getting client's nickname
	usernameLen, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %v quit.\n", conn.RemoteAddr())
		conn.Close()
		return
	}
	username := string(data[:usernameLen])

	// Getting room's name
	roomLen, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %v quit.\n", conn.RemoteAddr())
		conn.Close()
		return
	}
	roomname := string(data[:roomLen])

	c.ProcessRoom(roomname)

	// Joining the room
	c.Lock()
	defer c.Unlock()

	room := c.Rooms[roomname]
	if c.IsUsernameUniq(username, roomname) {
		client := CreateClient(username, conn)
		room.Clients[username] = client
		room.Messages <- fmt.Sprint(username, " joined to the room")
	} else {
		in, err := conn.Write([]byte("This nickname is already exists in room :("))
		if err != nil {
			fmt.Printf("Error when send to client: %d\n", in)
		}
		return
	}

	go c.ListenClient(conn, username, room)
}

// ProcessRoom ...
func (c *Chat) ProcessRoom(roomname string) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.Rooms[roomname]; !exists {
		c.Rooms[roomname] = CreateRoom(roomname)
		room := c.Rooms[roomname]
		go room.Run()
	}
}

// IsUsernameUniq ...
func (c *Chat) IsUsernameUniq(username, roomname string) bool {
	_, exists := c.Rooms[roomname].Clients[username]
	return !exists
}

// ListenClient ...
func (c *Chat) ListenClient(conn net.Conn, username string, room *Room) {
	data := make([]byte, 254)

	for {
		msgLen, err := conn.Read(data)
		if err != nil {
			log.Printf("Client %s quit.\n", conn.RemoteAddr())
			room.Messages <- fmt.Sprintf("%s left the room", username)
			room.RemoveClient(username)
			conn.Close()
			return
		}

		msg := fmt.Sprintf("(%s) %s: %s", room.Name, username, data[:msgLen])
		room.Messages <- msg
	}
}

// CleanChat ...
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

// RemoveRoom ...
func (c *Chat) RemoveRoom(room *Room) {
	if room != nil {
		if room.ClientCount() == 0 {
			log.Printf("Room [%s] has been destroyed", room.Name)
			delete(c.Rooms, room.Name)
		}
	}
}
