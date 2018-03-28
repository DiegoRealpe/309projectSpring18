package main

import (
	"sync"
	"fmt"
	"log"
	"io"
)

type playerConnection struct {
	client 			client
	packetInMutex 	packetInMutex
	packetOut 		chan PacketOut
	portNumber		int
	id				int
	isActive 		int

	packetLength int //if no packet is being read, packetLength should be 0
}

type packetInMutex struct {
	packetIn chan<- PacketIn
	mut      sync.Mutex
}

//needs to be a tread safe type
var nextID = 1
var idMutex = sync.Mutex{}

func MakePlayerConnection(client client, packetIn chan<- PacketIn) *playerConnection {

	fmt.Println("making PlayerConnection for client number",client.clientNum)

	playerConnection := playerConnection{}

	playerConnection.client = client
	playerConnection.packetInMutex = packetInMutex{packetIn: packetIn}
	playerConnection.packetOut = make(chan PacketOut, 100)
	playerConnection.portNumber = client.port
	playerConnection.assignId()
	playerConnection.isActive = 1

	go playerConnection.startReading()
	go playerConnection.startTransmitting()

	return &playerConnection
}

func (pConn *playerConnection) startReading() {

	if debug {fmt.Println("player reading")}
	for pConn.isActive == 1 {
		pConn.tryToPeekAndSetNewPacketLength()
		pConn.tryToReadPacket()
	}
	fmt.Println("No longer reading from player with id", pConn.id)
}

func (pConn *playerConnection) sendToPacketIn(data []byte){
	packet := PacketIn{
		connectionId: pConn.id,
		size:len(data),
		data:data,
	}

	if debug {fmt.Println("got",packet)}

	//sending over packetIn must be mutexed because the channel mighn have been changed by another thread
	pConn.packetInMutex.mut.Lock()
	pConn.packetInMutex.packetIn <- packet
	pConn.packetInMutex.mut.Unlock()
}

//assign next available id to player
func (pConn *playerConnection) assignId(){
	idMutex.Lock()
	pConn.id = nextID
	nextID++
	idMutex.Unlock()
}

func (pConn *playerConnection) startTransmitting() {
	for packet := range pConn.packetOut{
		if debug {fmt.Println("sending packet",packet)}
		pConn.transmitPacket(packet)
		if pConn.isActive == 0 {
			break
		}
	}
	fmt.Println("No longer transmitting to player with id", pConn.id)
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


func (pConn *playerConnection) peekLengthIfNoCommand(){
	if(pConn.packetLength != 0) {
		pConn.tryToPeekAndSetNewPacketLength()
	}
}

func (pConn *playerConnection) tryToPeekAndSetNewPacketLength(){
	peeked, _ := pConn.client.reader.Peek(1)

	if len(peeked) == 1{

		packetLength := inputPacketLengths[peeked[0]]

		if packetLength != 0{
			pConn.packetLength = packetLength
			if debug {fmt.Println("setting packet length",pConn.packetLength)}
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

	if debug {fmt.Println("full packet is",peeked)}
	pConn.client.reader.Discard(pConn.packetLength)
	pConn.packetLength = 0

	pConn.sendToPacketIn(peeked)
}

func (pConn *playerConnection) disconnect(){
	pConn.client.connection.Close()
}

var inputPacketLengths = map[byte]int{
	120 : 17,
	123 : 17,
	125 : 2,
}
