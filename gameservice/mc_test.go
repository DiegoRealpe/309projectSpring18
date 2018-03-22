package main

import (
	"reflect"
	"testing"
	"fmt"
)

func Test120Response(t *testing.T) {

	var packetIn = make(chan PacketIn, 10)
	var packetOut = make(chan PacketOut, 10)

	packetIn <- PacketIn{0, 18, []byte{120,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68}}

	go runGameController(GameOptions{}, packetIn, packetOut)

	res := <-packetOut

	fmt.Println("got back packet",res)

	//PacketOut{22, []byte{121,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68, 0, 0, 0, 0}, []int{0,1}}

	//assertEquals(expected, res, t, "121 did not match 120")

	close(packetIn)
	close(packetOut)
}

func Test123Response(t *testing.T) {

	var packetIn = make(chan PacketIn, 10)
	var packetOut = make(chan PacketOut, 10)

	packetIn <- PacketIn{0, 18, []byte{123,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68}}

	go runGameController(GameOptions{}, packetIn, packetOut)

	res := <-packetOut

	fmt.Println("got back packet",res)

	//PacketOut{22, []byte{121,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68, 0, 0, 0, 0}, []int{0,1}}

	//assertEquals(expected, res, t, "121 did not match 120")

	close(packetIn)
	close(packetOut)
}


//message is nullable
func assertEquals(expected, actual interface{}, t *testing.T, message string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Log(message)
		t.Log("expected:", expected)
		t.Log("actual  :", actual)
		t.Fail()
	}
}
