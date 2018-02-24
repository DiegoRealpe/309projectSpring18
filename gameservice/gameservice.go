package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

struct client{
	net.Conn connection
	bufio.Reader reader
	bufio.Writer writer
	clientno int
}

var ports []int

const NUMPLAYERS int = 2

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
	group := new([NUMPLAYERS]client)
	i := 0
	for _, conn := range connPasser {
		*group[i] = client{
			connection: conn
			reader: bufio.NewReader(conn)
			writer: bufio.NewWriter(conn)
			clientno: i
		}
		i++
		if i == NUMPLAYERS - 1{
			go func(group interface{}){
				init(group)
			}(group)
			i = 0
		}
	}
}

func startHttpServer() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
