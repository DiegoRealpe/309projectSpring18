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
	players [NUMPLAYERS]gamePlayer
}

type gamePlayer struct{
	connection *playerConnection
	isConnected bool
	isHost bool
	username string
	emoji string
}

func (gOpts GameOptions) buildGame() (g Game) {
	g.numPlayers = gOpts.numPlayers
	g.connectionIDToPlayerNumberMap = gOpts.connectionIDToPlayerNumberMap

	for i := 0; i < NUMPLAYERS; i++ {
		g.players[i].connection = gOpts.players[i]
		g.players[i].isConnected = true
	}
	//ext...

	return
}

func (g *Game) respondTo120(in *PacketIn, sendOut func(PacketOut) ) {
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

	sendOut(packetOut)
}

func (g *Game) respondTo123(in *PacketIn, sendOut func(PacketOut)) {
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

	sendOut(packetOut)
}

func (g *Game) respondTo125(in *PacketIn, sendOut func(PacketOut)){
	fmt.Println("recieved 125 packet...")
	disconnectingPlayer := g.connectionIDToPlayerNumberMap[in.connectionId]

	fmt.Println("Player", disconnectingPlayer, "has disconnected",)
	g.players[disconnectingPlayer].connection.disconnect()
	g.players[disconnectingPlayer].isConnected = false
	if(g.players[disconnectingPlayer].isHost){
		g.players[disconnectingPlayer].isHost = false
		fmt.Println("reassigning host")
		g.send127ToFirstAvaliablePlayer(sendOut);
	}

	packet126 := PacketOut{
		size : 2,
		data : []byte{126,disconnectingPlayer},
		targetIds: g.allConnectionIDsBut(in.connectionId),
	}

	sendOut(packet126)
}

func (g *Game) send122ToEveryone(sendOut func(PacketOut)) {
	packet122 := PacketOut{
		size: 1,
		data: []byte{122},
		targetIds: g.allConnectionIds(),
	}
	sendOut(packet122)
}

func (g *Game) send127ToFirstAvaliablePlayer(sendOut func(PacketOut)){

	/*
		algorithm: assign next host sequentially by iterating through players and choosing first active player
	 */
	 fmt.Print("reassigning hosts!!!!!!")
	for i, p := range g.players {
		if(p.isConnected){
			out := PacketOut{
				size: 1,
				data: []byte{127},
				targetIds: []int{p.connection.id},
			}
			fmt.Println("new host is",i)

			sendOut(out)
			g.players[i].isHost = true;
			break;
		}
	}
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

func (g *Game) allConnectionIds() []int{
	rtn := make([]int, g.numPlayers)

	for i := 0; i < g.numPlayers ; i+= 1 {
		rtn[i] = g.players[i].connection.id
	}
	return rtn
}
