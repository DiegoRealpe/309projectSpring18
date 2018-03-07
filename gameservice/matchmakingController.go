package main

import "fmt"

type matchMakingController struct {
	waitingPlayers []waitingPlayer
}

type waitingPlayer struct {
	connection *playerConnection
}

func makeMatchmakingController() matchMakingController {
	mmc := matchMakingController{}
	mmc.waitingPlayers = []waitingPlayer{}

	go mmc.startPairingRoutine()

	return mmc
}

func (mmc *matchMakingController) addConnectionToPool(connection *playerConnection) {
	waitingPlayer := connectionToWaitingPlayer(connection)
	mmc.waitingPlayers = append(mmc.waitingPlayers,waitingPlayer)

	fmt.Println("added player to matchmaking pool with connection number",connection.client.clientNum)
}

func connectionToWaitingPlayer(connection *playerConnection) waitingPlayer{
	rtn := waitingPlayer{}

	rtn.connection = connection

	return rtn
}

//currently called whenever a player enters the lobby
func (mmc *matchMakingController) tryToPair(){
}

func (mmc *matchMakingController) startPairingRoutine() {

}


//todo:
//add listeners to accept and send data to player connections, but
//for the time being, we can ignore all communications on this screen