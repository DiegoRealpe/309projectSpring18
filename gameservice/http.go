package main

import (
	"io"
	"net"
	"net/http"
	"strconv"
)

type clientConnection struct {
	connection net.Conn
	port int
}

type portHttpController struct {
	connPasser chan clientConnection
}


func makePortHttpController() portHttpController {
	connPasser := make(chan clientConnection)
	return portHttpController{connPasser: connPasser}
}


//when an http request is sent, send the requester a port and start listening on that port
func (portHttpController *portHttpController) handlePortRequested(w http.ResponseWriter, r *http.Request) {

	if numPortsAvailable() <= 0 {
		io.WriteString(w, "no ports avaliable, sorry fam")
		return
	}

	usedport := requestPort()

	stringport := strconv.Itoa(usedport)

	io.WriteString(w, stringport)

	go func() { //accept the first attempted connection on the port
		ln, _ := net.Listen("tcp", ":"+stringport)

		conn, _ := ln.Accept()

		connClient := clientConnection{
			connection: conn,
			port:				usedport,
		}

		portHttpController.connPasser <- connClient
	}()

}

