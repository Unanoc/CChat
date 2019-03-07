package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"chat/internal/client"
	"chat/internal/server"

	"github.com/fatih/color"
)

var (
	host      = "127.0.0.1"
	port      = "8000"
	chat      = server.CreateChat()
	connector = server.CreateConnector(host, port)
)

// start server
func init() {
	go chat.Run()
	go chat.CleanChat(1)
	go connector.AcceptConn(chat)
}

func ReadAndWrite(testCase []string, conn net.Conn) error {
	readStr := make([]byte, 254)

	for _, input := range testCase {
		// Getting from server
		if _, err := conn.Read(readStr); err != nil {
			return fmt.Errorf("Unexpected error: %s", err)
		}

		// Writing to server
		if _, err := conn.Write([]byte(input)); err != nil {
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}

	return nil
}

func TestConnectingToRoom(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	// start client
	c, err := client.CreateClient(host, port)
	defer c.Disconnect()
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}

	testCaseConnectToRoom := []string{"client", "test_room1"}
	if err = ReadAndWrite(testCaseConnectToRoom, c.Conn); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Getting message about success connection
	readStr := make([]byte, 254)
	length, err := c.Conn.Read(readStr)
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}
	if string(readStr[:length]) != color.HiBlackString(
		"%s joined to the room [%s]\n",
		testCaseConnectToRoom[0],
		testCaseConnectToRoom[1],
	) {
		t.FailNow()
	}
}

func TestClientsInRoom(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testCase := [][]string{
		[]string{"client_1", "test_room2"},
		[]string{"client_2", "test_room2"},
		[]string{"client_3", "other_test_room"},
	}

	// start clients
	for _, clientCase := range testCase {
		client, err := client.CreateClient(host, port)
		defer client.Disconnect()
		if err != nil {
			t.Error("Unexpected error: ", err.Error(), "\n")
			t.FailNow()
		}

		if err = ReadAndWrite(clientCase, client.Conn); err != nil {
			t.Error(err)
			t.FailNow()
		}
		go client.GetMessagesHandler()
	}

	// checking
	testRoom := chat.Rooms["test_room2"]
	otherTestRoom := chat.Rooms["other_test_room"]

	if testRoom.ClientCount() != 2 && otherTestRoom.ClientCount() != 1 {
		t.FailNow()
	}
}

func TestMsgLength(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client, err := client.CreateClient(host, port)
	defer client.Disconnect()
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}

	testCase := []string{"client", "test_room3"}
	if err = ReadAndWrite(testCase, client.Conn); err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Miss message about success connection
	readStr := make([]byte, 2*254)
	_, err = client.Conn.Read(readStr)
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}

	// msg's length is 276 bytes
	testMsg := "Ever man are put down his very. " +
		"And marry may table him avoid. " +
		"Hard sell it were into it upon. " +
		"He forbade affixed parties of assured " +
		"to me windows. Happiness him nor she " +
		"disposing provision. Add astonished " +
		"principles precaution yet friendship stimulated."

	// Writing to server
	if _, err := client.Conn.Write([]byte(testMsg)); err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}
	length, err := client.Conn.Read(readStr)
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}

	if string(readStr[21:length]) != testMsg[:254]+"\n" {
		t.FailNow()
	}
}

func TestCleanRoom(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testCase := [][]string{
		[]string{"client_1", "test_room_for_clean1"},
		[]string{"client_2", "test_room_for_clean2"},
		[]string{"client_3", "test_room_for_clean3"},
		[]string{"client_4", "test_room_for_clean4"},
		[]string{"client_5", "test_room_for_clean5"},
	}

	// start clients
	clients := make([]*client.Client, 0)
	for _, clientCase := range testCase {
		client, err := client.CreateClient(host, port)
		if err != nil {
			t.Error("Unexpected error: ", err.Error(), "\n")
			t.FailNow()
		}
		clients = append(clients, client)

		if err = ReadAndWrite(clientCase, client.Conn); err != nil {
			t.Error(err)
			t.FailNow()
		}
		go client.GetMessagesHandler()
	}

	currentCountRoom := len(chat.Rooms)
	if currentCountRoom < len(testCase) {
		t.FailNow()
	}

	for _, client := range clients {
		client.Disconnect()
	}

	time.Sleep(3 * time.Second)
	currentCountRoom = len(chat.Rooms)
	if currentCountRoom != 0 {
		t.FailNow()
	}
}
