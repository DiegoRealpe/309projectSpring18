package main

/*
	LobbyController.go contains all the higher-level and communication-related functionality for
	the lobby. this includes grabbing players from the chan of waiting players and routing packets
	to the correct handler.
 */


import (
	"fmt"
)

type LobbyController struct {
	l Lobby

	mmm *matchMakingModel
	packetIn chan PacketIn

	disperser packetDisperser

	packetRouterMap map[byte]func(*PacketIn,func(PacketOut))
}

func startLobby(mmm *matchMakingModel) {
	fmt.Println("startingn new lobby")

	lc := LobbyController{}
	lc.mmm = mmm
	lc.packetIn = make(chan PacketIn, 50)

	lc.l = lc.makeLobby()
	lc.buildPacketMap()

	lc.disperser.connections = make(map[int]chan<- PacketOut)

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

	for packet := range lc.packetIn { //just listen for packets
		lc.handleSinglePacket(packet)

		if packet.data[0] == 125 {
			break
		}else if lc.l.readyToMoveToGameScene{
			lc.startGCFromLobby()
			return
		}
	}

	fmt.Println("Lobby closing")
}

func (lc *LobbyController) addSinglePlayer(newPlayer *waitingPlayer) {

	if lc.mmm.connectionIdHasDisconnected(newPlayer.connection.id){
		newPlayer.connection.disconnect()
		return
	}

	lc.mmm.decrementOpenSpaces()
	fmt.Println("lobby added player with id", newPlayer.connection.id, ", current size is",lc.l.size+1)

	newPlayer.connection.SetNewPacketInChannel(lc.packetIn)

	lc.startDispersingToConnection(newPlayer.connection)

	lc.l.addPlayer(newPlayer,lc.disperser.send)
}

func (lc *LobbyController) handleSinglePacket(packet PacketIn){
	if debug {fmt.Println("lc got a packet", packet.data)}

	packetByte := packet.data[0]

	//call handler function based on value from map
	handlerFunction := lc.packetRouterMap[packetByte]
	handlerFunction(&packet,lc.disperser.send)
}


func (lc *LobbyController) startDispersingToConnection(connection *playerConnection){
	lc.disperser.connections[connection.id] = connection.packetOut
}

func (lc *LobbyController) makeLobby() (l Lobby) {
	l.size = 0
	l.messages = []chatMessage{}
	l.players = [NUMPLAYERS]lobbyPlayer{}

	return
}

//builds a map of packet types to handler functions
func (lc *LobbyController) buildPacketMap() {
	packetMap := map[byte]func(*PacketIn, func(PacketOut)){}

	packetMap[200] = lc.l.respondTo200
	packetMap[201] = lc.l.respondTo201
	packetMap[202] = lc.l.respondTo202
	packetMap[125] = lc.l.respondTo125
	packetMap[208] = lc.l.respondTo208

	lc.packetRouterMap = packetMap
}

func (lc *LobbyController) startGCFromLobby() {
	options := GameOptions{
		numPlayers: NUMPLAYERS,
		connectionIDToPlayerNumberMap: lc.getConnIDToPlayerNumberMap(),
	}
	for i, p := range lc.l.players {
		options.players[i] = p.connection
	}

	go runGameController(options, lc.packetIn)
}

func (lc *LobbyController) getConnIDToPlayerNumberMap() map[int]byte {
	idmap := make(map[int]byte)
	for i, p := range lc.l.players {
		idmap[p.connection.id] = byte(i)
	}
	return idmap
}

func (lc *LobbyController) startReadyTimer() bool {
	return true
}