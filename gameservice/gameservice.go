package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

type client struct{
	connection net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	clientno int
}

var ports []int

const NUMPLAYERS int = 2

var connPasser chan net.Conn

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
	ports = []int{5543, 9078}
	connPasser = make(chan net.Conn)
	startHttpServer()
	group := new([NUMPLAYERS]client)
	i := 0
	for conn := range connPasser {
		group[i].connection = conn
		group[i].reader = bufio.NewReader(conn)
		group[i].writer = bufio.NewWriter(conn)
		group[i].clientno = i
		i++
		if i == NUMPLAYERS - 1{
			go func(group *[NUMPLAYERS]client){
				initgame(group)
			}(group)
			i = 0
		}
	}
}

func startHttpServer() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":6000", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
