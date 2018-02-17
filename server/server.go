package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

var ports []int

var connPasser chan net.Conn

//Game struct with 2 connections, readers and writers to give to the players for them to communicate
type Game struct {
	connections [2]net.Conn
	reader      [2]*bufio.Reader
	writer      [2]*bufio.Writer
}

//when an http request is sent, send the requester a port and start listening on that port
func handler(w http.ResponseWriter, r *http.Request) {

	if len(ports) < 1 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := ports[len(ports)-1]

	ports = ports[:len(ports)-1]

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)

	go func() {//accept the first attempted connection on the port
		ln, _ := net.Listen("tcp", ":"+stringport)

		conn, _ := ln.Accept()

		connPasser <- conn
	}()

}

func main() {

	var g Game

	ports = []int{5543, 9078}

	connPasser = make(chan net.Conn)

	go func() {//handle http requests
		http.HandleFunc("/", handler)
		err := http.ListenAndServe(":80", nil)

		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	i := 0

	for connected := range connPasser {//when a player connects, initialize their readers and writers
		g.connections[i] = connected
		g.reader[i] = bufio.NewReader(g.connections[i])
		g.writer[i] = bufio.NewWriter(g.connections[i])
		hellobyte := []byte{byte(122), byte(i)}
		g.writer[i].Write(hellobyte)
		g.writer[i].Flush()
		
		time.Sleep(2*time.Second)
		testpacket := ServerPacket{//testing; to be removed later
			serverPlayerState: 121,
			playernumber: uint8(i),
			xPosition: 0,
			yPosition: 0,
			xVelocity: 0,
			yVelocity: 0,
			timestamp: 0,
		}
		
		testbytes := ParseServerPacket(testpacket)
		
		g.writer[i].Write(testbytes)
		g.writer[i].Flush()
		
		go func() {
			ListenAndSend(g, i)
		}()
		
		i++
	}
}


