package server

import (
	"fmt"
	"log"
	"net"
	"os"
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

	// Get room's name
	roomLen, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %v quit.\n", conn.RemoteAddr())
		conn.Close()
		return
	}

	fmt.Println(string(data))

	roomName := string(data[:roomLen])
	room, exists := c.Rooms[roomName]
	if !exists {
		c.Rooms[roomName] = CreateRoom()
		room = c.Rooms[roomName]
	}

	fmt.Println(len(room.Clients))

	// Get client's nickname
	nameLength, err := conn.Read(data)
	if err != nil {
		log.Printf("Client %v quit.\n", conn.RemoteAddr())
		conn.Close()
		return
	}
	name := string(data[:nameLength])
	fmt.Println(name)

	go room.Run()

	if _, exists := room.Clients[name]; !exists {
		client := CreateClient(name, conn)
		room.Clients[name] = client
		joinedMsg := name + " joined to the room"
		room.Messages <- joinedMsg
		log.Println(joinedMsg)
	} else {
		in, err := conn.Write([]byte("This nickname is already exists in room :("))
		if err != nil {
			fmt.Printf("Error when send to client: %d\n", in)
			os.Exit(0)
		}
		return
	}

	for {
		msgLen, err := conn.Read(data)
		if err != nil {
			log.Printf("Client %s quit.\n", conn.RemoteAddr())
			conn.Close()
			return
		}

		msg := fmt.Sprintf("[%s]: %s", name, data[:msgLen])
		room.Messages <- msg
	}
}
