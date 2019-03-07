package server

import (
	"log"
	"net"
)

// Connector is network layer
type Connector struct {
	Host   string
	Port   string
	Remote string
}

// CreateConnector returns an instance of Connector
func CreateConnector(host, port string) *Connector {
	return &Connector{
		Host:   host,
		Port:   port,
		Remote: host + ":" + port,
	}
}

// AcceptConn creates listener and then listen to new connections
func (c *Connector) AcceptConn(chat *Chat) {
	listener, err := net.Listen("tcp", c.Remote)
	defer listener.Close()

	if err != nil {
		log.Printf("Error when listen: %s, Err: %s\n", c.Remote, err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting client: %s\n", err.Error())
			return
		}

		chat.Register <- conn
	}
}
