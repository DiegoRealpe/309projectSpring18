package main

import (
	"fmt"
	"sync"
)

type LobbyController struct {
	l Lobby

	mmm *matchMakingModel
	packetIn chan PacketIn
	packetOut chan PacketOut

	disperser lobbyDisperser
}

type lobbyDisperser struct{
	mut sync.Mutex
	connections [] <-chan PacketOut
}

func startLobby(mmm *matchMakingModel) {
	fmt.Println("startingn new lobby")

	lc := LobbyController{}
	lc.mmm = mmm
	lc.packetIn = make(chan PacketIn, 50)
	lc.packetOut = make(chan PacketOut, 50)

	lc.l = lc.makeLobby()

	go lc.runLobbyDispersion()

	lobbyFull := false
	for !lobbyFull { //listen for players and packets
		select {
		case newPlayer := <-lc.mmm.playerChan:
			lc.addSinglePlayer(newPlayer)
			lobbyFull = lc.l.size == NUMPLAYERS
		case packet := <-lc.packetIn:
			lc.handleSinglePacket(packet)
		}
	}

	fmt.Println("lobby is full")

	for packet := range lc.packetIn { //just listen for players
		lc.handleSinglePacket(packet)

		//TODO make logic to break loop to end lobby and GORoutine if all players have disconnected
	}

	fmt.Println("Lobby closing")
}

func (lc *LobbyController) addSinglePlayer(newPlayer *waitingPlayer) {
	lc.mmm.decrementOpenSpaces()
	lc.l.size += 1
	fmt.Println("lobby added player with id", newPlayer.connection.id, ", current size is",lc.l.size)
	newPlayer.connection.SetNewPacketInChannel(lc.packetIn)
}

func (lc *LobbyController) handleSinglePacket(packet PacketIn){
	fmt.Println("lc got a packet", packet.data)
}

func (lc *LobbyController) runLobbyDispersion() {
	
}

func (lc *LobbyController) makeLobby() (l Lobby) {

	l.size = 0

	return
}
