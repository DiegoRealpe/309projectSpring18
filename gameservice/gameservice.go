package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

struct client{
	net.Conn conn
	bufio.Reader reader
	bufio.Writer writer
	clientno int
}

var ports []int

var connPasser net.Conn

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

func main(){
	ports = {5543, 9078}
	connPasser = make(chan net.Conn)
	startHttpServer()
	group := new(*[]client)
	for _, conn := range connPasser {

	}
}

func startHttpServer() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
