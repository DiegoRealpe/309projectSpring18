package main

import "fmt"

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


}