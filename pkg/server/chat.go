package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
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
		log.Printf("New connection: [%s]", conn.RemoteAddr())

		go c.ProcessConn(conn)
	}
}

// ProcessConn initialises clients's connection
func (c *Chat) ProcessConn(conn net.Conn) {
	msgToClient := color.BlueString("Enter your name: ")
	username, err := DealWithClient(msgToClient, conn, true)
	if err != nil {
		return
	}
	msgToClient = color.BlueString("Enter room name you want to join: ")
	roomname, err := DealWithClient(msgToClient, conn, true)
	if err != nil {
		return
	}

	// Create room if room does not exist
	c.ProcessRoom(roomname)

	if c.IsUsernameUniq(username, roomname) {
		client := CreateClient(username, conn)
		c.Lock()
		room := c.Rooms[roomname]
		room.Register <- client
		c.Unlock()
		go c.ListenClient(client, room)
	} else {
		msgToClient = color.HiRedString("This nickname is already exists in this room.\n")
		if _, err = DealWithClient(msgToClient, conn, false); err != nil {
			return
		}

		go c.ProcessConn(conn)
		return
	}
}

// DealWithClient sends message to client, then reads message from client and return it if flag "withReading" is true
func DealWithClient(requestMsg string, conn net.Conn, withReading bool) (result string, err error) {
	data := make([]byte, 254)

	if _, err = conn.Write([]byte(requestMsg)); err != nil {
		conn.Close()
		return
	}

	if withReading {
		length, err := conn.Read(data)
		if err != nil {
			log.Printf("Client %s has not been connected", conn.RemoteAddr())
			conn.Close()
			return "", err
		}
		result = string(data[:length])
	}

	return
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
	c.Lock()
	defer c.Unlock()

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
		case "/list":
			room.SendClientList(client)
		default:
			room.Messages <- fmt.Sprintf("(%s) %s: %s", room.Name, client.Username, rawMessage)
		}
	}
}

// CleanChat removes rooms with zero clients once an n seconds
func (c *Chat) CleanChat(n int) {
	for {
		time.Sleep(time.Duration(n) * time.Second)
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
