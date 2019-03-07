package main

import (
	"chat/internal/server"
	"log"
)

var (
	host = "127.0.0.1"
	port = "8000"
)

func main() {
	log.Println("Initiating server... (Ctrl-C to stop)")

	chat := server.CreateChat()
	go chat.Run()
	go chat.CleanChat(60)

	connector := server.CreateConnector(host, port)
	connector.AcceptConn(chat)
}
