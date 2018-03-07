package main

import "testing"

func TestMatchmakingSimplePairing(t *testing.T) {
	mmc := makeMatchmakingController()

	client0 := client{clientNum:0}
	connec0 := MakePlayerConnection(client0,nil)

	client1 := client{clientNum:1}
	connec1 := MakePlayerConnection(client1,nil)

	mmc.addConnectionToPool(connec0)
	mmc.addConnectionToPool(connec1)
}
