package main

import (
	"testing"
	"reflect"
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

	res := <- packetOut

	expected := PacketOut{22, []byte{121,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}
	if !reflect.DeepEqual(res,expected) {
		t.Fail()
	}

	close(packetIn)
	close(packetOut)
}


