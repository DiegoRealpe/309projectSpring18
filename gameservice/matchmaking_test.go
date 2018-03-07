package main

import (
	"testing"
	"time"
)

func TestMatchmakingSimplePairing(t *testing.T) {
	t.Skip()

	mmc := makeMatchmakingController()

	client0 := client{clientNum:0}
	connec0 := MakePlayerConnection(client0,nil)

	client1 := client{clientNum:1}
	connec1 := MakePlayerConnection(client1,nil)

	client2 := client{clientNum:2}
	connec2 := MakePlayerConnection(client2,nil)

	client3 := client{clientNum:3}
	connec3 := MakePlayerConnection(client3,nil)

	mmc.addConnectionToPool(connec0)
	mmc.addConnectionToPool(connec1)
	mmc.addConnectionToPool(connec2)
	mmc.addConnectionToPool(connec3)

	time.Sleep(100 * time.Millisecond)
}
