package main

import (
	"chat/pkg/server"
	"log"
	"net"
	"os"
)

var (
	host   = "127.0.0.1"
	port   = "8000"
	remote = host + ":" + port
)

func main() {
	log.Println("Initiating server... (Ctrl-C to stop)")

	c := server.CreateChat()
	go c.Run()
	go c.CleanChat()

	listener, err := net.Listen("tcp", remote)
	defer listener.Close()

	if err != nil {
		log.Printf("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting client: %s\n", err.Error())
			os.Exit(0)
		}

		c.Register <- conn
	}
}
