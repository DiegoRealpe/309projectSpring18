package main

import (
	"sync"
)

type playerConnection struct {
	client        client
	packetInMutex packetInMutex
	packetOut     <-chan PacketOut
}

type packetInMutex struct {
	packetIn chan<- PacketIn
	mut      sync.Mutex
}

type packetOutMutex struct {
	packetOut chan<- PacketOut
	mut       sync.Mutex
}

func MakePlayerConnection(client client, packetIn chan<- PacketIn) *playerConnection {
	playerConnection := playerConnection{}

	playerConnection.client = client
	playerConnection.packetInMutex = packetInMutex{packetIn: packetIn}
	playerConnection.packetOut = make(chan PacketOut, 100)

	return &playerConnection
}

func (pconn *playerConnection) startReading() {

}

func (pconn *playerConnection) startTransmitting() {

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

//dispersion stops once the returned channel is closed
func startSendToAllRoutine(connectionOutChannels []<-chan PacketOut) chan<- PacketOut {
	returnChannel := make(chan PacketOut, 100)

	go listenAndDispersePackets(connectionOutChannels, returnChannel)

	return returnChannel
}

func listenAndDispersePackets(connectionChannels []<-chan PacketOut, toDisperse <-chan PacketOut) {
	for packet := range toDisperse {
		for _, connectionChannel := range connectionChannels {
			connectionChannel <- packet
		}
	}
}
