package main

import "fmt"

//should be a gorouting, but not start new goroutines
func packetInRouter(gameOptions GameOptions, in <-chan Packet, out chan<- Packet){
	fmt.Println("game controller   :::  ","starting router")
	g := gameOptions.buildGame()

	for p := range in{
		respondToSinglePacket(&p,out,&g)
	}
}

func respondToSinglePacket(in *Packet,out chan<- Packet, g *Game){
	packetType := in.parseType()
	switch packetType{
	case 0 :
		g.respondTo0(in,out)
	case 1 :
		g.respondTo1(in,out)
	default:
		fmt.Println("Error: packet with type",packetType,"DNE")
	}
}

func buildPacketMap(g *Game){
	//packetMap := map[byte](func(*PacketIn,chan PacketOut))


}