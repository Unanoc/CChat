package server

import (
	"log"
	"net"
)
type Connector struct {
	Host string
	Port string
	Remote string
}

func CreateConnector(host, port string) *Connector{
	return &Connector {
		Host: host,
		Port: port,
		Remote: host + ":" + port,
	}
}

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