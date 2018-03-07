package main

import (
	"sync"
	"fmt"
)

type playerConnection struct {
	client        client
	packetInMutex packetInMutex
	packetOut     chan PacketOut
}

type packetInMutex struct {
	packetIn chan<- PacketIn
	mut      sync.Mutex
}

func MakePlayerConnection(client client, packetIn chan<- PacketIn) *playerConnection {

	fmt.Println("making PlayerConnection for client number",client.clientNum)

	playerConnection := playerConnection{}

	playerConnection.client = client
	playerConnection.packetInMutex = packetInMutex{packetIn: packetIn}
	playerConnection.packetOut = make(chan PacketOut, 100)

	go playerConnection.startReading()
	go playerConnection.startTransmitting()

	return &playerConnection
}

func (pconn *playerConnection) startReading() {

}

func (pconn *playerConnection) startTransmitting() {
	for packet := range pconn.packetOut{
		fmt.Println("sending packet",packet)
		pconn.transmitPacket(packet)
	}
}

func (pconn *playerConnection) send(packetIn PacketIn) {
	mutex := &pconn.packetInMutex.mut
	packetChannel := &pconn.packetInMutex.packetIn

	mutex.Lock()
	*packetChannel <- packetIn
	mutex.Unlock()
}

func (pconn *playerConnection) SetNewPacketInChannel(packetIn chan<- PacketIn) {
	pconn.packetInMutex.mut.Lock()

	pconn.packetInMutex.packetIn = packetIn

	pconn.packetInMutex.mut.Unlock()
}

func (pconn *playerConnection) transmitPacket(out PacketOut) {
	pconn.client.writer.Write(out.data)
	pconn.client.writer.Flush()
}

//to be started as a goroutine
func listenAndDispersePackets(connectionChannels []chan<- PacketOut, toDisperse <-chan PacketOut) {
	for packet := range toDisperse {
		for _, connectionChannel := range connectionChannels {
			connectionChannel <- packet
		}
	}
}

