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

	matchMakingController := makeMatchmakingController()

	listenForConnections(portHttpController.connPasser,matchMakingController)
}

func listenForConnections(connPasser chan<- net.Conn, matchMakingController matchMakingController) {

	currentClientNumber := 1

	for conn := range connPasser {
		client := client{}
		client.connection = conn
		client.reader = bufio.NewReader(conn)
		client.writer = bufio.NewWriter(conn)
		client.clientNum = currentClientNumber

		currentClientNumber++

		playerConnection := MakePlayerConnection(client,nil)

		matchMakingController.addConnectionToPool(playerConnection)
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
