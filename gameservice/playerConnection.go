package main

import (
	"sync"
	"fmt"
	"log"
	"io"
)

const maxReadBufferSize = 100

type playerConnection struct {
	client        client
	packetInMutex packetInMutex
	packetOut     chan PacketOut

	packetLength int //if no packet is being read, packetLength should be 0
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

func (pConn *playerConnection) startReading() {

	fmt.Println("player reading")
	for {

		pConn.tryToPeekAndSetNewPacketLength()
		pConn.tryToReadPacket()
	}
}

func (pConn *playerConnection) sendToPacketIn(data []byte){
	packet := PacketIn{size:len(data),data:data}

	fmt.Println("got",packet)

	//sending over packetIn must be mutexed because the channel mighn have been changed by another thread
	pConn.packetInMutex.mut.Lock()
	pConn.packetInMutex.packetIn <- packet
	pConn.packetInMutex.mut.Unlock()
}

func (pConn *playerConnection) startTransmitting() {
	for packet := range pConn.packetOut{
		fmt.Println("sending packet",packet)
		pConn.transmitPacket(packet)
	}
}

func (pConn *playerConnection) SetNewPacketInChannel(packetIn chan<- PacketIn) {
	pConn.packetInMutex.mut.Lock()

	pConn.packetInMutex.packetIn = packetIn

	pConn.packetInMutex.mut.Unlock()
}

func (pConn *playerConnection) transmitPacket(out PacketOut) {
	pConn.client.writer.Write(out.data)
	pConn.client.writer.Flush()
}

//to be started as a goroutine
func listenAndDispersePackets(connectionChannels []chan<- PacketOut, toDisperse <-chan PacketOut) {
	for packet := range toDisperse {
		for _, connectionChannel := range connectionChannels {
			connectionChannel <- packet
		}
	}
}

func (pConn *playerConnection) peekLengthIfNoCommand(){
	if(pConn.packetLength != 0) {
		pConn.tryToPeekAndSetNewPacketLength()
	}
}

func (pConn *playerConnection) tryToPeekAndSetNewPacketLength(){
	peeked, _ := pConn.client.reader.Peek(1)

	if len(peeked) == 1{

		packetLength := packetLengths[peeked[0]]

		if packetLength != 0{
			pConn.packetLength = packetLength
			fmt.Println("setting packet length",pConn.packetLength)
		}else{
			log.Fatal("Invalid Packet sent by client, got byte ",peeked[0])
		}
	}
}

func (pConn *playerConnection) tryToReadPacket(){
	peeked, error := pConn.client.reader.Peek(pConn.packetLength)

	if len(peeked) != pConn.packetLength || pConn.packetLength == 0{
		//if connection did not have full packet
		if error == io.EOF || error == nil{
			return
		}
		fmt.Println("tried to peek bytes and got non EOF error",error)
	}

	fmt.Println("full packet is",peeked)
	pConn.client.reader.Discard(pConn.packetLength)
	pConn.packetLength = 0

	pConn.sendToPacketIn(peeked)
}



var packetLengths = map[byte]int{
	8 : 2,
	120 : 17,
}