package main

import (
	"fmt"
	"net"

	"chat/pkg/client"

	"github.com/fatih/color"
)

var (
	host   = "127.0.0.1"
	port   = "8000"
	remote = host + ":" + port
)

func main() {
	conn, err := net.Dial("tcp", remote)
	defer conn.Close()

	if err != nil {
		color.Red("Server not found.")
		return
	}

	c := client.CreateClient(conn)
	if err := c.ProcessJoin(); err != nil {
		fmt.Print(err)
		return
	}

	go c.WriteMessagesHandler()
	c.GetMessagesHandler()
}
