package main

import (
	"fmt"
)

//long lived struct that should always be its own goroutine, it is initialized as the entry point for new connections and, when
//pairing is successful launches a game controller as a goroutine and sets the player connection to send packets there instead.
type matchMakingModel struct {

	//key value is their connection id
	waitingPlayers map[int]waitingPlayer //should only be modified by pairing routine
}

type waitingPlayer struct {
	connection *playerConnection
}

func startMatchmakingModel() matchMakingModel {
	fmt.Println("starting match making model")

	mmm := matchMakingModel{}

	mmm.waitingPlayers = make(map[int]waitingPlayer, 50)

	return mmm
}

//currently called whenever a player enters the lobby
func (mmm *matchMakingModel) tryToPair() {

	fmt.Println("trying to pair")
	for len(mmm.waitingPlayers) >= NUMPLAYERS {

		//split waiting players into a slice of players joining the game and a slice of those not
		gamePlayers := mmm.separatePlayersIntoAGame()

		mmm.startGame(gamePlayers)
	}
}

func (mmm *matchMakingModel) separatePlayersIntoAGame() []waitingPlayer {
	game := make([]waitingPlayer, NUMPLAYERS)

	//copy NUMPLAYERS elements from map into a slice and remove them from the map
	i := 0
	for key, val := range mmm.waitingPlayers {
		if i >= NUMPLAYERS {
			break
		}
		i++

		delete(mmm.waitingPlayers, key)
		game[i-1] = val
	}

	return game
}

func (mmm *matchMakingModel) acceptPlayer(connection *playerConnection){

	waitingPlayer := connectionToWaitingPlayer(connection)

	fmt.Println("added player to matchmaking pool with connection number",connection.client.clientNum)

	mmm.waitingPlayers[connection.id] = waitingPlayer
	mmm.tryToPair()
}

func (mmm *matchMakingModel) disconnectPlayer(id int) {
	disconnectingPlayer := mmm.waitingPlayers[id]

	fmt.Println("Player", disconnectingPlayer.connection.id, "has left matchmaking")
	disconnectingPlayer.connection.disconnect()

	delete(mmm.waitingPlayers, disconnectingPlayer.connection.id)
}

func (mmm *matchMakingModel) startGame(players []waitingPlayer) {
	fmt.Println("starting mock game with players:")
	for _, val := range players {
		fmt.Println("***", val.connection.client.clientNum)
	}

	gameOpts := GameOptions{numPlayers: NUMPLAYERS}
	for i, p := range players {
		gameOpts.players[i] = p.connection
	}

	gameOpts.connectionIDToPlayerNumberMap = makePlayerNumberMap(players)

	packetOutChannel := startPacketOutDispersionWithPlayers(players)
	packetInChannel := makePacketInChannelForAllPlayers(players)

	go runGameController(gameOpts, packetInChannel, packetOutChannel)
	send122PacketsToPlayers(players)
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
