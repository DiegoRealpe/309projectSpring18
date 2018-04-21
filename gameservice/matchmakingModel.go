package main

import (
	"fmt"
	"sync"
)

//long lived struct that should always be its own goroutine, it is initialized as the entry point for new connections and, when
//pairing is successful launches a game controller as a goroutine and sets the player connection to send packets there instead.
type matchMakingModel struct {
	playerChan chan *waitingPlayer

	openSpaces int
	openSpacesMut sync.Mutex

	disconnected map[int]bool
	disconnectedMut sync.Mutex
}

type waitingPlayer struct {
	connection *playerConnection
}

func startMatchmakingModel() matchMakingModel {
	fmt.Println("starting match making model")

	mmm := matchMakingModel{openSpaces:0}
	mmm.playerChan = make(chan *waitingPlayer, 50)
	mmm.disconnected = make(map[int]bool)

	return mmm
}


func (mmm *matchMakingModel) acceptPlayer(connection *playerConnection){

	waitingPlayer := connectionToWaitingPlayer(connection)

	fmt.Println("added player to matchmaking pool with connection number",connection.client.clientNum)

	mmm.playerChan <- &waitingPlayer

	mmm.openSpacesMut.Lock()
	if mmm.openSpaces == 0{
		go mmm.runNewLobby()
		mmm.openSpaces = NUMPLAYERS
	}
	mmm.openSpacesMut.Unlock()
}

func (mmm *matchMakingModel) runNewLobby(){
	go startLobby(mmm)
}

//method may be called twice due to concurrency setup, should have no functional difference
//between one call and 2
func (mmm *matchMakingModel) disconnectPlayer(id int) {
	if debug {fmt.Println("connection with id",id,"quit from matchmaking")}

	mmm.disconnectedMut.Lock()
	mmm.disconnected[id] = true
	mmm.disconnectedMut.Unlock()
}

func (mmm *matchMakingModel) connectionIdHasDisconnected(id int) bool{
	mmm.disconnectedMut.Lock()
	_ , existed := mmm.disconnected[id]
	mmm.disconnectedMut.Unlock()
	return existed
}


func (mmm *matchMakingModel) respondTo125(in *PacketIn) {
	fmt.Println("recieved 125 packet...")

	mmm.disconnectPlayer(in.connectionId)
}

func connectionToWaitingPlayer(connection *playerConnection) waitingPlayer {
	rtn := waitingPlayer{}

	rtn.connection = connection

	return rtn
}

func (mmm *matchMakingModel) decrementOpenSpaces() {
	mmm.decrementOpenSpacesBy(1)
}

func (mmm *matchMakingModel) decrementOpenSpacesBy(by int) {
	mmm.openSpacesMut.Lock()
	mmm.openSpaces -= by
	mmm.openSpacesMut.Unlock()
}