package main

import "fmt"

type matchMakingController struct {
	matchMakingModel matchMakingModel
}

func makeMatchmakingController() matchMakingController {
	mmc := matchMakingController{}

	mmm := startMatchmakingModel()

	mmc.matchMakingModel = mmm
	return mmc
}

func (mmc *matchMakingController) addConnectionToPool(connection *playerConnection) {
	waitingPlayer := connectionToWaitingPlayer(connection)

	fmt.Println("added player to matchmaking pool with connection number",connection.client.clientNum)
	mmc.matchMakingModel.waitingPlayerChan <- waitingPlayer
}

func connectionToWaitingPlayer(connection *playerConnection) waitingPlayer{
	rtn := waitingPlayer{}

	rtn.connection = connection

	return rtn
}


//todo:
//add listeners to accept and send data to player connections, but
//for the time being, we can ignore all communications on this screen