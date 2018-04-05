package main

import "fmt"

type gameController struct {
	out             chan<- PacketOut
	g               Game
	packetRouterMap map[byte]func(*PacketIn, chan<- PacketOut)
}

//should be a gorouting, but not start new goroutines
func runGameController(gameOptions GameOptions, in <-chan PacketIn, out chan<- PacketOut) {
	fmt.Println("starting game controller")

	controller := gameController{}
	controller.g = gameOptions.buildGame()
	controller.buildPacketMap()
	controller.out = out

	for p := range in {
		controller.respondToSinglePacket(&p)

		//TODO we need to make sure goroutine ends
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
