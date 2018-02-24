package main

import (
	"fmt"
)

type GameOptions struct {
	numPlayers int
	//ext...
}

func (gOpts GameOptions) buildGame() (g Game) {
	g.numPlayers = gOpts.numPlayers
	//ext...
	return
}

type Game struct {
	numPlayers int
}


func (g *Game) respondTo0(in *PacketIn, out chan<- PacketOut){
	fmt.Println("game model   :::  ","logic for packet:",in)

	data := make([]byte,2)
	pack1 := PacketOut{2,data}

	out <- pack1
}

func (g *Game) respondTo1(in *PacketIn, out chan<- PacketOut){
	fmt.Println("game model   :::  ","doing whatever packet 1 would do")
}


func (g *Game) respondTo120(in *PacketIn, out chan<- PacketOut){
	fmt.Println("recieved 120 with data", *in)

	//packet120 := ParseBytesTo120(in.data)
	//
	//packet121 := packet121{
	//	serverPlayerState : 121,
	//	playernumber:
	//	xPosition : packet120.xPosition,
	//}


}

//playernumber      uint8
//xPosition         float32
//yPosition         float32
//xVelocity         float32
//yVelocity         float32
//timestamp         float32
