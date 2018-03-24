package main

import (
	"io"
	"net"
	"net/http"
	"strconv"
)

type portHttpController struct {
	connPasser chan net.Conn
}


func makePortHttpController() portHttpController {
	connPasser := make(chan net.Conn)
	return portHttpController{connPasser: connPasser}
}


//when an http request is sent, send the requester a port and start listening on that port
func (portHttpController *portHttpController) handlePortRequested(w http.ResponseWriter, r *http.Request) {

	if len(ports) < 1 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := ports[len(ports)-1]

	ports = ports[:len(ports)-1]

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)

	go func() { //accept the first attempted connection on the port
		ln, _ := net.Listen("tcp", ":"+stringport)

		conn, _ := ln.Accept()

		portHttpController.connPasser <- conn
	}()

}

