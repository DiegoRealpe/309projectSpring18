package main

import "fmt"

type gameController struct{
	in <-chan Packet
	out chan<- Packet
	g Game
	packetRouterMap map[byte]func(*Packet,chan<- Packet)
}

//should be a gorouting, but not start new goroutines
func runGameController(gameOptions GameOptions, in <-chan Packet, out chan<- Packet){
	fmt.Println("game controller   :::  ","starting router")

	controller := gameController{}
	controller.g = gameOptions.buildGame()
	controller.buildPacketMap()

	for p := range in{
		controller.respondToSinglePacket(&p)
	}
}

func (controller *gameController) respondToSinglePacket(in *Packet){
	packetType := in.parseType()

	controller.callHandlerFor(packetType,in)
}

//builds a map of packet types to handler functions
func (controller *gameController) buildPacketMap() {
	packetMap := map[byte](func(*Packet,chan<- Packet)){}

	packetMap[0] = controller.g.respondTo0
	packetMap[1] = controller.g.respondTo1

	controller.packetRouterMap = packetMap
}

//route packet correctly
func (controller *gameController) callHandlerFor(packetType byte,in *Packet) {

	handlerFunc := controller.packetRouterMap[packetType]

	if handlerFunc == nil{
		fmt.Println("ERROR : Packet did not have recognized type. type was",packetType)
	}else{
		handlerFunc(in,controller.out)
	}
}