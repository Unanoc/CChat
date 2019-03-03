package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

// CreateClient returns an instance of Client
func CreateClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

// Client keeps the connection of user
type Client struct {
	Conn net.Conn
}

// ProcessJoin organizes the connection process
func (c *Client) ProcessJoin() error {
	writeStr := make([]byte, 254)
	readStr := make([]byte, 254)

	// Get the name of client
	length, err := c.Conn.Read(readStr)
	if err != nil {
		return fmt.Errorf(color.RedString("Error when recieve from server"))
	}
	fmt.Printf("%s", readStr[:length])
	fmt.Scanf("%s", &writeStr)
	if _, err := c.Conn.Write([]byte(writeStr)); err != nil {
		return fmt.Errorf(color.RedString("Error when send to server"))
	}

	// Get the room name
	length, err = c.Conn.Read(readStr)
	if err != nil {
		return fmt.Errorf(color.RedString("Error when recieve from server"))
	}
	fmt.Printf("%s", readStr[:length])
	fmt.Scanf("%s", &writeStr)
	if _, err := c.Conn.Write([]byte(writeStr)); err != nil {
		return fmt.Errorf(color.RedString("Error when send to server"))
	}

	color.HiGreen("You have been successfully connected")
	return nil
}

// GetMessagesHandler gets messages for client
func (c *Client) GetMessagesHandler() {
	readStr := make([]byte, 254)
	for {
		length, err := c.Conn.Read(readStr)
		if err != nil {
			color.Red("Connection is closed")
			return
		}
		fmt.Printf("%s", readStr[:length])
	}
}

// WriteMessagesHandler sends client's messages
func (c *Client) WriteMessagesHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		writeStr, _, _ := reader.ReadLine()
		_, err := c.Conn.Write([]byte(writeStr))
		if err != nil {
			color.Red("Error when send to server")
			return
		}
	}
}
