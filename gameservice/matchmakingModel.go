package main

import (
	"fmt"
)

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
		mmm.waitingPlayers = mmm.waitingPlayers[2:]

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
	packetOutChannel := startPacketOutDispersionWithPlayers(players)
	packetInChannel := makePacketInChannelForAllPlayers(players)

	go runGameController(gameOpts,packetInChannel,packetOutChannel)
	send122PacketsToPlayers(players)
}


//builds a single channel which sends to all clients
func startPacketOutDispersionWithPlayers(players []waitingPlayer) chan<- PacketOut{
	chanSlice := make([]chan<- PacketOut,NUMPLAYERS)
	for i, player := range players{
		chanSlice[i] = player.connection.packetOut
	}

	toDisperse := make(chan PacketOut,50)
	go listenAndDispersePackets(chanSlice,toDisperse)

	return toDisperse
}

func makePacketInChannelForAllPlayers(players []waitingPlayer) <-chan PacketIn{
	packetInChannel := make(chan PacketIn, 50)

	for _, player :=  range players{
		player.connection.SetNewPacketInChannel(packetInChannel)
	}

	return packetInChannel
}

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

func startPlayersReading(players []waitingPlayer){
	for _, val := range players{
		go val.connection.startReading()
	}
}