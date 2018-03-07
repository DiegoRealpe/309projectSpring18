package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

type client struct {
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	clientNum  int
}

var ports []int

const NUMPLAYERS = 2

func main() {
	ports = []int{5543, 9078}

	portHttpController := generatePortHttpController()

	//start listening for http
	startHttpServer(portHttpController)

	listenForConnections(portHttpController.connPasser)
}

func listenForConnections(connPasser chan<- net.Conn) {
	clients := new([NUMPLAYERS]client)
	i := 0
	for conn := range connPasser {
		clients[i].connection = conn
		clients[i].reader = bufio.NewReader(conn)
		clients[i].writer = bufio.NewWriter(conn)
		clients[i].clientNum = i
		i++
		if i == NUMPLAYERS-1 {
			go func(group *[NUMPLAYERS]client) {
				initgame(group)
			}(clients)
			i = 0
		}
	}
}

func generatePortHttpController() portHttpController {
	connPasser := make(chan net.Conn)
	return portHttpController{connPasser: connPasser}
}

func startHttpServer(portHttpController portHttpController) {
	http.HandleFunc("/", portHttpController.handlePortRequested)
	err := http.ListenAndServe(":6000", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
