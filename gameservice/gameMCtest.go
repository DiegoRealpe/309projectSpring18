package main

import (
	"testing"
	"time"
	"fmt"
)


func TestSum(t *testing.T) {

	packetOut := make(chan PacketOut,100)
	packetIn := make(chan PacketIn,100)

	data := make([]byte,2)
	pack1 := PacketIn{2,2,data}

	g := GameOptions{}

	go mockSocketTransmitter(packetOut)
	go runGameController(g,packetIn,packetOut)

	packetIn <- pack1

	time.Sleep(5 * time.Second)
	//closing packets
	close(packetIn)
	close(packetOut)

}


func mockSocketTransmitter(packetOut <-chan PacketOut){
	for p := range packetOut{
		fmt.Println("sending out packet:",p)
	}

}

