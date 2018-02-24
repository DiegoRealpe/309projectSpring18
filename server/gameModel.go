package main

import "fmt"

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


func (g *Game) respondTo0(in *Packet, out chan<- Packet){
	fmt.Println("game model   :::  ","logic for packet:",in)

	data := make([]byte,2)
	pack1 := Packet{2,data}

	out <- pack1
}

func (g *Game) respondTo1(in *Packet, out chan<- Packet){
	fmt.Println("game model   :::  ","doing whatever packet 1 would do")
}

