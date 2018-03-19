package main

import (
	"fmt"
)

type GameOptions struct {
	numPlayers int
	//ext...
}

type Game struct {
	numPlayers int
}


func (gOpts GameOptions) buildGame() (g Game) {
	g.numPlayers = gOpts.numPlayers
	//ext...
	return
}

func (g *Game) respondTo120(in *PacketIn, out chan<- PacketOut) {
	fmt.Println("recieved 120 packet")

	packet120 := ParseBytesTo120(in.data)

	fmt.Println(packet120)

	packet121 := packet121{
		serverPlayerState: 121,
		playernumber:      in.playerNum,
		xPosition:         packet120.xPosition,
		yPosition:         packet120.yPosition,
		xVelocity:         packet120.xVelocity,
		yVelocity:         packet120.yVelocity,
		timestamp:         0,
	}

	fmt.Println(packet121)

	packetOut := PacketOut{
		size: 22,
		data: packet121.toBytes(),
	}

	out <- packetOut
}

