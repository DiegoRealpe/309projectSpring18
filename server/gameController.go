package main

import "fmt"

//should be a gorouting, but not start new goroutines
func packetInRouter(gameOptions GameOptions, in <-chan PacketIn, out chan<- PacketOut){
	fmt.Println("game controller   :::  ","starting router")
	g := gameOptions.buildGame()

	for p := range in{
		respondToSinglePacket(&p,out,&g)
	}
}

func respondToSinglePacket(in *PacketIn,out chan<- PacketOut, g *Game){
	packetType := parsePacketType(in)
	switch packetType{
	case 0 :
		g.respondTo0(in,out)
	case 1 :
		g.respondTo1(in,out)
	default:
		fmt.Println("Error: packet with type",packetType,"DNE")
	}
}

func parsePacketType(in *PacketIn) byte{
	return in.data[0]
}
