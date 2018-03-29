package main

import "fmt"

type matchMakingController struct {
	packetIn        chan PacketIn
	newConnectionIn chan *playerConnection

	matchMakingModel matchMakingModel
}

//returns function for adding player to pool
func startMatchmakingController() func(connection *playerConnection) {
	fmt.Println("starting matchmaking controller")

	mmc := matchMakingController{}
	mmm := startMatchmakingModel()
	mmc.matchMakingModel = mmm
	mmc.packetIn = make(chan PacketIn, 50)
	mmc.newConnectionIn = make(chan *playerConnection, 50)

	go mmc.runWaitForPacketsAndPlayers()

	return mmc.addConnectionToPool
}

func (mmc *matchMakingController) runWaitForPacketsAndPlayers() {
	for {
		select {
		case packet := <-mmc.packetIn:
			mmc.handleSinglePacketIn(&packet)
		case connection := <-mmc.newConnectionIn:
			fmt.Println("accepting player")
			mmc.matchMakingModel.acceptPlayer(connection)
		}
	}
}

func (mmc *matchMakingController) handleSinglePacketIn(in *PacketIn) {
	if in.data[0] == 125 && in.size == 1 {
		mmc.matchMakingModel.respondTo125(in)
	}
}

func (mmc *matchMakingController) addConnectionToPool(connection *playerConnection) {
	connection.SetNewPacketInChannel(mmc.packetIn)
	mmc.newConnectionIn <- connection
}
