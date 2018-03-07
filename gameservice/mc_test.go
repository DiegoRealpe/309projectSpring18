package main

import (
	"testing"
	"reflect"
)


func Test120Response(t *testing.T){
	var packetIn = make(chan PacketIn,10)
	var packetOut = make(chan PacketOut,10)

	packetIn <- PacketIn{0,18,[]byte{120,0,1,2,3,4,5,6,7,8,9,10,0,0,0,0,0,0} }

	go runGameController(GameOptions{}, packetIn, packetOut)

	res := <- packetOut

	expected := PacketOut{22, []byte{121,0,0,1,2,3,4,5,6,7,8,9,10,0,0,0,0,0,0,0,0,0}}

	assertEquals(expected,res,t,"121 did not match 120")

	close(packetIn)
	close(packetOut)
}

//message is nullable
func assertEquals(expected, actual interface{}, t *testing.T, message string){
	if !reflect.DeepEqual(actual,expected) {
		t.Log(message)
		t.Log("expected:",expected)
		t.Log("actual  :",actual)
		t.Fail()
	}
}


