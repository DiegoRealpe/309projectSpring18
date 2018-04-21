package main

import (
	"fmt"
	)

type gameController struct {
	g               Game

	disperser gameDisperser

	packetRouterMap map[byte]func(*PacketIn, func(PacketOut))
}

type gameDisperser struct{
	connections map[int]chan<- PacketOut
}

//should be a goroutine, but not start new goroutines
func runGameController(gameOptions GameOptions, in <-chan PacketIn) {
	fmt.Println("starting game controller!")

	controller := gameController{}

	controller.g = gameOptions.buildGame()

	controller.buildPacketMap()

	controller.configureDisperser()

	controller.g.send122ToEveryone(controller.disperser.send)
	controller.g.send127ToFirstAvaliablePlayer(controller.disperser.send)

	for p := range in {
		controller.respondToSinglePacket(&p)

		//TODO we need to make sure goroutine ends
	}
}

func (gc *gameController) configureDisperser(){
	gc.disperser.connections = make(map[int]chan<- PacketOut)

	//populate map to disperse
	for _, v := range gc.g.players {
		gc.disperser.connections[v.connection.id] = v.connection.packetOut
	}
}

func (controller *gameController) respondToSinglePacket(in *PacketIn) {
	packetType := in.parseType()

	controller.callHandlerFor(packetType, in)
}

//builds a map of packet types to handler functions
func (controller *gameController) buildPacketMap() {
	packetMap := map[byte]func(*PacketIn, func(PacketOut) ){}

	packetMap[120] = controller.g.respondTo120
	packetMap[123] = controller.g.respondTo123
	packetMap[125] = controller.g.respondTo125
	packetMap[130] = controller.g.respondTo130
	packetMap[133] = controller.g.respondTo133

	controller.packetRouterMap = packetMap
}

//route packet correctly
func (controller *gameController) callHandlerFor(packetType byte, in *PacketIn) {

	handlerFunc := controller.packetRouterMap[packetType]

	if handlerFunc == nil {
		fmt.Println("ERROR : Packet did not have recognized type. type was", packetType)
	} else {
		handlerFunc(in, controller.disperser.send)
	}
}

func (d gameDisperser) send(packet PacketOut){
	if debug {fmt.Println("dispersing packet",packet)}

	for _, id := range packet.targetIds{
		d.connections[id] <- packet
	}
}
