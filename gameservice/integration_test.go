package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"testing"
	"os"
	"time"
)


func TestMain(m *testing.M) {

	go main()

	m.Run()

	os.Exit(0)
}

func TestConnectClients(t *testing.T) { //used globally for test

	player0 := connectUser()
	player1 := connectUser()

	time.Sleep(100 * time.Millisecond)

	player1.Write([]byte{120,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68})
	player1.Flush()

	player0.Write([]byte{120,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68})
	player0.Flush()

	time.Sleep(100 * time.Millisecond)

}

func connectUser() *bufio.ReadWriter {
	resp, _ := http.Get("http://localhost:6000/tcpport")
	responseSlice := make([]byte, 4)
	resp.Body.Read(responseSlice)
	portstring := string(responseSlice)
	fmt.Println("Connecting to port", portstring)
	conn, _ := net.Dial("tcp", ":"+portstring)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
	return rw
}
