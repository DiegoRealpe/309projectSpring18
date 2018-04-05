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
}

type waitingPlayer struct {
	connection *playerConnection
}

func startMatchmakingModel() matchMakingModel {
	fmt.Println("starting match making model")

	mmm := matchMakingModel{openSpaces:0}
	mmm.playerChan = make(chan *waitingPlayer, 50)

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

func (mmm *matchMakingModel) disconnectPlayer(id int) {

	//disconnectingPlayer, playerExisted := mmm.waitingPlayers[id]
	//
	////because we are disconnecting using 125 and connection being closed this sometimes is called twice per player
	//if playerExisted {
	//	fmt.Println("Player", disconnectingPlayer.connection.id, "has left matchmaking")
	//	disconnectingPlayer.connection.disconnect()
	//
	//	delete(mmm.waitingPlayers, disconnectingPlayer.connection.id)
	//}
}

//assigned sequentially
func makePlayerNumberMap(players []waitingPlayer) map[int]byte {
	m := make(map[int]byte)

	for i, val := range players {
		m[val.connection.id] = byte(i)
	}

	fmt.Println("map is", m)

	return m
}

//builds a single channel which sends to all clients
func startPacketOutDispersionWithPlayers(players []waitingPlayer) chan<- PacketOut {

	idDispersionMap := makeIDDispersionMap(players)

	toDisperse := make(chan PacketOut, 50)
	go listenAndDispersePackets(idDispersionMap, toDisperse)

	return toDisperse
}

func makeIDDispersionMap(players []waitingPlayer) map[int]chan<- PacketOut {
	m := make(map[int]chan<- PacketOut)

	for _, val := range players {
		m[val.connection.id] = val.connection.packetOut
	}

	return m
}

func makePacketInChannelForAllPlayers(players []waitingPlayer) <-chan PacketIn {
	packetInChannel := make(chan PacketIn, 50)

	for _, player := range players {
		player.connection.SetNewPacketInChannel(packetInChannel)
	}

	return packetInChannel
}

//should be in line with player numbers because both were assigned sequentially from the same slice
func send122PacketsToPlayers(players []waitingPlayer) {
	for num, player := range players {
		send122PacketToPlayer(player, num)
	}
}

func send122PacketToPlayer(player waitingPlayer, playerNum int) {
	packet := PacketOut{}
	packet.size = 2
	packet.data = []byte{122, byte(playerNum)}

	player.connection.packetOut <- packet
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
	mmm.openSpacesMut.Lock()
	mmm.openSpaces -= 1
	mmm.openSpacesMut.Unlock()
}
