package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/fatih/color"
)

var (
	host   = "127.0.0.1"
	port   = "8000"
	remote = host + ":" + port
)

func readAndWrite(msg []byte, conn net.Conn) error {
	writeStr, readStr := make([]byte, 254), make([]byte, 254)
	// Getting from server
	_, err := conn.Read(readStr)
	if err != nil {
		return fmt.Errorf("Unexpected error: %s", err)
	}

	// Writing to server
	writeStr = []byte(msg)
	if _, err := conn.Write([]byte(writeStr)); err != nil {
		return fmt.Errorf("Unexpected error: %s", err)
	}
	return nil
}

func TestInputNameAndRoom(t *testing.T) {
	conn, err := net.Dial("tcp", remote)
	defer conn.Close()

	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}

	for i := 0; i < 2; i++ {
		if err = readAndWrite([]byte("test"), conn); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	readStr := make([]byte, 254)
	// Getting message about success connection
	length, err := conn.Read(readStr)
	if err != nil {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}
	if string(readStr[:length]) != color.HiBlackString("test joined to the room [test]\n") {
		t.Error("Unexpected error: ", err.Error(), "\n")
		t.FailNow()
	}
}
