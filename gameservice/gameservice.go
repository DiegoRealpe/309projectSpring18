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

const NUMPLAYERS = 2

func main() {

	fmt.Println("starting game service!")
	initPortService()

	//ports = []int{6001, 6002, 6003, 6004, 6005, 6006, 6007, 6008, 6009, 6010, 6011, 6012} //todo: make a staic function with static variables for this

	portHttpController := makePortHttpController()

	matchMakingController := makeMatchmakingController()

	go listenForConnections(portHttpController.connPasser,matchMakingController)

	//start listening for http
	startHttpServer(portHttpController)
}

func listenForConnections(connPasser <-chan clientConnection, matchMakingController matchMakingController) {
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

		matchMakingController.addConnectionToPool(playerConnection)
	}
}

func startHttpServer(portHttpController portHttpController) {
	http.HandleFunc("/", portHttpController.handlePortRequested)
	err := http.ListenAndServe(":6000", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
