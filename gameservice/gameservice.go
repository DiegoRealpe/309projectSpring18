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
}

var ports []int

const NUMPLAYERS = 2

func main() {

	fmt.Println("starting game service!")

	ports = []int{6001, 6002, 6003, 6004, 6005, 6006} //todo: make a staic function with static variables for this

	portHttpController := makePortHttpController()

	matchMakingController := makeMatchmakingController()

	go listenForConnections(portHttpController.connPasser,matchMakingController)

	//start listening for http
	startHttpServer(portHttpController)
}

func listenForConnections(connPasser <-chan net.Conn, matchMakingController matchMakingController) {
	fmt.Println("listening for connections")

	currentClientNumber := 0

	for conn := range connPasser {
		fmt.Println("starting handling for a connection")

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

func startHttpServer(portHttpController portHttpController) {
	http.HandleFunc("/", portHttpController.handlePortRequested)
	err := http.ListenAndServe(":6000", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
