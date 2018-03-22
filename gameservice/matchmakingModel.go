package main

import (
	"fmt"
)

//long lived struct that should always be its own goroutine, it is initialized as the entry point for new connections and, when
//pairing is successful launches a game controller as a goroutine and sets the player connection to send packets there instead.
type matchMakingModel struct {
	waitingPlayers []waitingPlayer //should only be modified by pairing routine
	waitingPlayerChan chan waitingPlayer
}

type waitingPlayer struct {
	connection *playerConnection
}

func startMatchmakingModel() matchMakingModel{
	fmt.Println("starting match making model")

	mmm := matchMakingModel{}

	mmm.waitingPlayers = getWaitingPlayerSLice()
	mmm.waitingPlayerChan = make(chan waitingPlayer, 10)

	//start listening on waiting player channel
	go mmm.startPairingRoutine()

	return mmm
}

//make slice with underlying array size 50 to hold players and decrease tha ammount of re-allocation
func getWaitingPlayerSLice() []waitingPlayer {
	var underlyingArray [50]waitingPlayer
	return underlyingArray[0:0]
}

//currently called whenever a player enters the lobby
func (mmm *matchMakingModel) tryToPair(){

	fmt.Println("trying to pair")
	for len(mmm.waitingPlayers) >= NUMPLAYERS {

		//split waiting players into a slice of players joining the game and a slice of those not
		gamePlayers := mmm.waitingPlayers[0:NUMPLAYERS]
		mmm.waitingPlayers = mmm.waitingPlayers[NUMPLAYERS:]

		mmm.startGame(gamePlayers)
	}
}

func (mmm *matchMakingModel) startPairingRoutine() {
	for player := range mmm.waitingPlayerChan {

		mmm.waitingPlayers = append(mmm.waitingPlayers,player)

		mmm.tryToPair()
	}
}

func (mmm *matchMakingModel) startGame(players []waitingPlayer){
	fmt.Println("starting mock game with players:")
	for _, val := range players{
		fmt.Println("***",val.connection.client.clientNum)
	}

	gameOpts := GameOptions{numPlayers:NUMPLAYERS}
	gameOpts.connectionIDToPlayerNumberMap = makePlayerNumberMap(players)

	packetOutChannel := startPacketOutDispersionWithPlayers(players)
	packetInChannel := makePacketInChannelForAllPlayers(players)

	go runGameController(gameOpts,packetInChannel,packetOutChannel)
	send122PacketsToPlayers(players)
}

//assigned sequentially
func makePlayerNumberMap(players []waitingPlayer) map[int]byte {
	m := make(map[int]byte )

	for i, val := range players{
		m[val.connection.id] = byte(i)
	}

	fmt.Println("map is",m)

	return m
}


//builds a single channel which sends to all clients
func startPacketOutDispersionWithPlayers(players []waitingPlayer) chan<- PacketOut{

	idDispersionMap := makeIDDispersionMap(players)

	toDisperse := make(chan PacketOut,50)
	go listenAndDispersePackets(idDispersionMap,toDisperse)

	return toDisperse
}

func makeIDDispersionMap(players []waitingPlayer) map[int]chan<- PacketOut{
	m := make(map[int]chan<- PacketOut)

	for _, val := range players {
		m[val.connection.id] = val.connection.packetOut
	}

	return m
}

func makePacketInChannelForAllPlayers(players []waitingPlayer) <-chan PacketIn{
	packetInChannel := make(chan PacketIn, 50)

	for _, player :=  range players{
		player.connection.SetNewPacketInChannel(packetInChannel)
	}

	return packetInChannel
}

//should be in line with player numbers because both were assigned sequentially from the same slice
func send122PacketsToPlayers(players []waitingPlayer){
	for num, player := range players{
		send122PacketToPlayer(player,num)
	}
}

func send122PacketToPlayer(player waitingPlayer,playerNum int){
	packet := PacketOut{}
	packet.size = 2
	packet.data = []byte{122,byte(playerNum)}

	player.connection.packetOut <- packet
}
