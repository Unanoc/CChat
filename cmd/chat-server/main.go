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

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	if err != nil {
		log.Printf("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		c.Register <- conn
	}
}
