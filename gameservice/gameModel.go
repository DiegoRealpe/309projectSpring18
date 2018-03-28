package main

import (
	"fmt"
)

type GameOptions struct {
	numPlayers int
	connectionIDToPlayerNumberMap map[int]byte
	players [NUMPLAYERS]*playerConnection
}

type Game struct {
	numPlayers int
	connectionIDToPlayerNumberMap map[int]byte
	players [NUMPLAYERS]*playerConnection
}


func (gOpts GameOptions) buildGame() (g Game) {
	g.numPlayers = gOpts.numPlayers
	g.connectionIDToPlayerNumberMap = gOpts.connectionIDToPlayerNumberMap
	for i := 0; i < NUMPLAYERS; i++ {
		g.players[i] = gOpts.players[i]
	}
	//ext...

	return
}

func (g *Game) respondTo120(in *PacketIn, out chan<- PacketOut) {
	if debug {fmt.Println("recieved 120 packet")}

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

	packetOut := PacketOut{
		size: 22,
		data: packet121.toBytes(),
		targetIds: g.allConnectionIDsBut(in.connectionId),
	}

	out <- packetOut
}

func (g *Game) respondTo123(in *PacketIn, out chan<- PacketOut) {
	if debug {fmt.Println("recieved 123 packet")}


	packet123 := ParseBytesTo123(in.data)

	packet124 := packet124{
		xPosition:         packet123.xPosition,
		yPosition:         packet123.yPosition,
		xVelocity:         packet123.xVelocity,
		yVelocity:         packet123.yVelocity,
		timestamp:         0,
	}

	packetOut := PacketOut{
		size: 21,
		data: packet124.toBytes(),
		targetIds: g.allConnectionIDsBut(in.connectionId),
	}

	out <- packetOut
}

func (g *Game) respondTo125(in *PacketIn, out chan<- PacketOut){
	fmt.Println("recieved 125 packet...")
	packet125 := ParseBytesTo125(in.data)
	disconnectingPlayer := packet125.playerNumber

	freePort(g.players[disconnectingPlayer].portNumber)
	g.players[disconnectingPlayer].isActive = 0
	fmt.Println("Player", disconnectingPlayer, "has disconnected")
	g.players[disconnectingPlayer].disconnect()
}

func (g *Game) allConnectionIDsBut(id int) []int {

	slice := make([]int,NUMPLAYERS-1)

	if debug {fmt.Println("sending to",slice)}

	i := 0
	for key, _ := range g.connectionIDToPlayerNumberMap {
		if key != id{
			slice[i] = key
			i++
		}
	}

	if debug {fmt.Println("sending to",slice)}

	return slice
}
