package main

import (
	"testing"
	"fmt"
	"time"
)

func TestSum(t *testing.T) {

	t.Skip()
	//packetOut := make(chan PacketOut,100)
	//packetIn := make(chan PacketIn,100)
	//
	//data := make([]byte,2)
	//pack1 := PacketIn{2,2,data}
	//
	//g := GameOptions{}
	//
	//go mockSocketTransmitter(packetOut)
	//go runGameController(g,packetIn,packetOut)
	//
	//packetIn <- pack1
	//
	//time.Sleep(5 * time.Second)
	////closing packets
	//close(packetIn)
	//close(packetOut)

}

func Test120Response(t *testing.T){//used globally for tests
	var packetIn = make(chan PacketIn,10)
	var packetOut = make(chan PacketOut,10)

	packetIn <- PacketIn{0,18,[]byte{120,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} }

	go runGameController(GameOptions{}, packetIn, packetOut)

	time.Sleep(100 * time.Millisecond)

	res := <- packetOut
	fmt.Println("got result",res)

	

	close(packetIn)
	close(packetOut)
}


func mockSocketTransmitter(packetOut <-chan PacketOut){
	for p := range packetOut{
		fmt.Println("sending out packet:",p)
	}

}

