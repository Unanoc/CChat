package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

// CreateClient ...
func CreateClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

// Client ...
type Client struct {
	Conn net.Conn
}

// ProcessJoin ...
func (c *Client) ProcessJoin() error {
	writeStr := make([]byte, 254)

	fmt.Printf(color.BlueString("Enter your name: "))
	fmt.Scanf("%s", &writeStr)
	if _, err := c.Conn.Write([]byte(writeStr)); err != nil {
		return fmt.Errorf(color.RedString("Error when send to server\n"))
	}

	fmt.Printf(color.BlueString("Enter room name you want to join: "))
	fmt.Scanf("%s", &writeStr)
	if _, err := c.Conn.Write([]byte(writeStr)); err != nil {
		return fmt.Errorf(color.RedString("Error when send to server\n"))
	}

	fmt.Println(color.GreenString("You have been successfully connected."))
	return nil
}

// GetMessagesHandler ...
func (c *Client) GetMessagesHandler() {
	readStr := make([]byte, 254)
	for {
		length, err := c.Conn.Read(readStr)
		if err != nil {
			fmt.Printf(color.RedString("Error when read from server. Error:%s\n", err))
			return
		}
		fmt.Println(string(readStr[:length]))
	}
}

// WriteMessagesHandler ...
func (c *Client) WriteMessagesHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		writeStr, _, _ := reader.ReadLine()
		if string(writeStr) == "quit" {
			fmt.Print(color.RedString("Communication terminated.\n"))
			return
		}

		in, err := c.Conn.Write([]byte(writeStr))
		if err != nil {
			fmt.Printf(color.RedString("Error when send to server: %d\n", in))
			return
		}
	}
}
