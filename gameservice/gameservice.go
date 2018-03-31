package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type client struct {
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	clientNum  int
	port int
}

const debug = false

const NUMPLAYERS = 2

func main() {

	fmt.Println("starting game service!")
	initPortService()

	portHttpController := makePortHttpController()

	matchMakingFunction := startMatchmakingController()

	go listenForConnections(portHttpController.connPasser,matchMakingFunction)

	//start listening for http
	startHttpServer(portHttpController)
}

func listenForConnections(connPasser <-chan clientConnection, matchMakingFunction func(connection *playerConnection)) {
	fmt.Println("listening for connections")

	currentClientNumber := 0

	for conn := range connPasser {
		fmt.Println("starting handling for a connection")

		client := client{}
		client.connection = conn.connection
		client.reader = bufio.NewReader(conn.connection)
		client.writer = bufio.NewWriter(conn.connection)
		client.clientNum = currentClientNumber
		client.port = conn.port

		currentClientNumber++

		playerConnection := MakePlayerConnection(client,nil)

		matchMakingFunction(playerConnection)
	}
}

func startHttpServer(portHttpController portHttpController) {
	http.HandleFunc("/", portHttpController.handlePortRequested)
	err := http.ListenAndServe(":6000", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
