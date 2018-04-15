package main

import (
	"fmt"
	"sync")

type gameController struct {
	out             chan PacketOut
	g               Game

	disperser gameDisperser

	packetRouterMap map[byte]func(*PacketIn, chan<- PacketOut)
}

type gameDisperser struct{
	mut sync.Mutex
	connections map[int]chan<- PacketOut
}

//should be a gorouting, but not start new goroutines
func runGameController(gameOptions GameOptions, in <-chan PacketIn, out chan PacketOut) {
	fmt.Println("starting game controller")

	controller := gameController{}
	controller.g = gameOptions.buildGame()
	controller.buildPacketMap()
	controller.out = out

	controller.disperser.connections = make(map[int]chan<- PacketOut)

	for _, player := range gameOptions.players {
		controller.disperser.connections[player.id] = player.packetOut
	}

	go controller.runGameDispersion()

	for p := range in {
		controller.respondToSinglePacket(&p)

		//TODO we need to make sure goroutine ends
	}
}

func (gc *gameController) runGameDispersion() {
	for packet := range gc.out {
		gc.disperser.mut.Lock()

		for _, id := range packet.targetIds{
			gc.disperser.connections[id] <- packet
		}
		gc.disperser.mut.Unlock()
	}
}

func (controller *gameController) respondToSinglePacket(in *PacketIn) {
	packetType := in.parseType()

	controller.callHandlerFor(packetType, in)
}

//builds a map of packet types to handler functions
func (controller *gameController) buildPacketMap() {
	packetMap := map[byte]func(*PacketIn, chan<- PacketOut){}

	packetMap[120] = controller.g.respondTo120
	packetMap[123] = controller.g.respondTo123
	packetMap[125] = controller.g.respondTo125

	controller.packetRouterMap = packetMap
}

//route packet correctly
func (controller *gameController) callHandlerFor(packetType byte, in *PacketIn) {

	handlerFunc := controller.packetRouterMap[packetType]

	if handlerFunc == nil {
		fmt.Println("ERROR : Packet did not have recognized type. type was", packetType)
	} else {
		handlerFunc(in, controller.out)
	}
}
//autolayout
