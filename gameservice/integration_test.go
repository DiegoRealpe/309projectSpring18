package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestConnectClients(t *testing.T) { //used globally for tests

	t.Skip()

	go main()

	rw0 := connectUser()
	slice := make([]uint8, 2)
	rw0.Read(slice)
	fmt.Println(slice)

	rw1 := connectUser()
	slice = make([]uint8, 2)
	rw1.Read(slice)
	fmt.Println(slice)

	time.Sleep(100 * time.Millisecond)

	rw0.Write([]uint8{120, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	rw0.Flush()

	time.Sleep(100 * time.Millisecond)

	slice = make([]uint8, 50)
	rw0.Read(slice)
	fmt.Println("read from player 0:", slice)

	slice = make([]uint8, 50)
	rw1.Read(slice)
	fmt.Println("read from player 1:", slice)

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
