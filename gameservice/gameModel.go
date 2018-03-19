package main

import (
	"fmt"
)

type GameOptions struct {
	numPlayers int
	connectionIDToPlayerNumberMap map[int]byte
}

type Game struct {
	numPlayers int
	connectionIDToPlayerNumberMap map[int]byte
}


func (gOpts GameOptions) buildGame() (g Game) {
	g.numPlayers = gOpts.numPlayers
	g.connectionIDToPlayerNumberMap = gOpts.connectionIDToPlayerNumberMap
	//ext...

	return
}

func (g *Game) respondTo120(in *PacketIn, out chan<- PacketOut) {
	fmt.Println("recieved 120 packet")

	playerNumber := g.connectionIDToPlayerNumberMap[in.connectionId]

	packet120 := ParseBytesTo120(in.data)

	packet121 := packet121{
		serverPlayerState: 121,
		playernumber:      playerNumber,
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

