package main

import (
	"fmt"

	"chat/internal/client"
)

var (
	host = "127.0.0.1"
	port = "8000"
)

func main() {
	c, err := client.CreateClient(host, port)
	defer c.Disconnect()

	if err != nil {
		fmt.Print(err)
		return
	}

	go c.WriteMessagesHandler()
	c.GetMessagesHandler()
}
