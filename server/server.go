package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

var ports []int

var connPasser chan net.Conn

//Game struct with 2 connections, readers and writers to give to the players for them to communicate
type Game struct {
	connections [2]net.Conn
	reader      [2]*bufio.Reader
	writer      [2]*bufio.Writer
}

func handler(w http.ResponseWriter, r *http.Request) {

	if len(ports) < 1 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := ports[len(ports)-1]

	ports = ports[:len(ports)-1]

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)

	go func() {
		ln, _ := net.Listen("tcp", ":"+stringport)

		conn, _ := ln.Accept()

		writer := bufio.NewWriter(conn)
		writer.WriteString("welcome to port " + stringport + " stranger\nAre you ready to play some FUTBOL?????")
		writer.Flush()

		connPasser <- conn
	}()

}

func main() {

	var g Game

	ports = []int{5543, 9078}

	connPasser = make(chan net.Conn)

	go func() {
		http.HandleFunc("/", handler)
		err := http.ListenAndServe(":80", nil)

		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	i := 0

	for i < 2 {
		for connected := range connPasser {
			g.connections[i] = connected
			i++
		}
	}

	for i = 0; i < len(g.connections); i++ { //create readers and writers for the connections
		g.reader[i] = bufio.NewReader(g.connections[i])
		g.writer[i] = bufio.NewWriter(g.connections[i])
	}

	for i = 0; i < len(g.connections); i++ {
		go func() {
			ListenAndSend(g, i)
		}()
	}
}

//restful api
//crud api
//create, read, update, delete
