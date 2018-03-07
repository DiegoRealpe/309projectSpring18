package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestConnectClients(t *testing.T) { //used globally for tests

	t.Skip()

	go main()


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
