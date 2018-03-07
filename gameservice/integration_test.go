package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"testing"
	"os"
)


func TestMain(m *testing.M) {

	go main()

	m.Run()

	os.Exit(0)
}

func TestConnectClients(t *testing.T) { //used globally for test

	connectUser()
	connectUser()


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
