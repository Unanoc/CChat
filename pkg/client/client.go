package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

// CreateClient returns an instance of Client
func CreateClient(host, port string) (*Client, error) {
	client := &Client{
		Host:   host,
		Port:   port,
		Remote: host + ":" + port,
	}

	conn, err := net.Dial("tcp", client.Remote)
	if err != nil {
		return nil, fmt.Errorf(color.RedString("Server not found"))
	}
	client.Conn = conn

	return client, nil
}

// Client keeps the connection of user
type Client struct {
	Conn   net.Conn
	Host   string
	Port   string
	Remote string
}

// Disconnect closes client's connection
func (c *Client) Disconnect() {
	c.Conn.Close()
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
