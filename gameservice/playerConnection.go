package main

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type playerConnection struct {
	client        client
	packetInMutex packetInMutex
	packetOut     chan PacketOut
	portNumber    int
	id            int
	isActive      bool

	portIsOpen    bool
	disconnectionMut sync.Mutex

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

	fmt.Println("making PlayerConnection for client number", client.clientNum)

	playerConnection := playerConnection{}

	playerConnection.client = client
	playerConnection.packetInMutex = packetInMutex{packetIn: packetIn}
	playerConnection.packetOut = make(chan PacketOut, 100)
	playerConnection.portNumber = client.port
	playerConnection.assignId()
	playerConnection.isActive = true
	playerConnection.portIsOpen = true

	go playerConnection.startReading()
	go playerConnection.startTransmitting()

	return &playerConnection
}

func (pConn *playerConnection) startReading() {

	if debug {
		fmt.Println("player reading")
	}
	for pConn.isActive {
		pConn.tryToPeekAndSetNewPacketLength()
		if !pConn.isActive {
			break
		}
		pConn.tryToReadPacket()
	}
	fmt.Println("No longer reading from player with id", pConn.id)
}

func (pConn *playerConnection) sendToPacketIn(data []byte) {
	packet := PacketIn{
		connectionId: pConn.id,
		size:         len(data),
		data:         data,
	}

	if debug {
		fmt.Println("got", packet)
	}

	//sending over packetIn must be mutexed because the channel mighn have been changed by another thread
	pConn.packetInMutex.mut.Lock()
	pConn.packetInMutex.packetIn <- packet
	pConn.packetInMutex.mut.Unlock()
}

//assign next available id to player
func (pConn *playerConnection) assignId() {
	idMutex.Lock()
	pConn.id = nextID
	nextID++
	idMutex.Unlock()
}

func (pConn *playerConnection) startTransmitting() {
	for packet := range pConn.packetOut {
		if debug {
			fmt.Println("sending packet", packet)
		}
		pConn.transmitPacket(packet)
		if !pConn.isActive {
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

func (pConn *playerConnection) tryToPeekAndSetNewPacketLength() {
	peeked, err := pConn.client.reader.Peek(1)
	if err == io.EOF {
		fmt.Println("!!!!!!!!EOF!!!!!!!!")
		pConn.isActive = false
		fmt.Println("sending zombie 125")
		pConn.sendZombie125()
	}

	if len(peeked) == 1 {

		packetLength := inputPacketLengths[peeked[0]]

		if packetLength != 0 {
			pConn.packetLength = packetLength
			if debug {
				fmt.Println("setting packet length", pConn.packetLength)
			}
		} else {
			log.Fatal("Invalid Packet sent by client, got byte ", peeked[0])
		}
	}
}

func (pConn *playerConnection) tryToReadPacket() {
	peeked, error := pConn.client.reader.Peek(pConn.packetLength)

	if len(peeked) != pConn.packetLength || pConn.packetLength == 0 {
		//if connection did not have full packet
		if error == io.EOF || error == nil {
			fmt.Println("!!!!!!EOF error!!!!")
			pConn.isActive = false
			fmt.Println("sending zombie 125")
			pConn.sendZombie125()
			return
		}
		fmt.Println("tried to peek bytes and got non EOF error", error)
	}

	if debug {
		fmt.Println("full packet is", peeked)
	}
	pConn.client.reader.Discard(pConn.packetLength)
	pConn.packetLength = 0

	pConn.sendToPacketIn(peeked)
}

//this method is thread safe via blocking
func (pConn *playerConnection) disconnect() {
	fmt.Println("disconnect called")

	pConn.disconnectionMut.Lock()

	if pConn.portIsOpen {
		pConn.client.connection.Close()
		freePort(pConn.portNumber)
		pConn.isActive = false
	}

	pConn.portIsOpen = false
	pConn.disconnectionMut.Unlock()
}

func (pConn *playerConnection) sendZombie125(){
	pConn.sendToPacketIn([]byte{125})
}

var inputPacketLengths = map[byte]int{
	120 : 17,
	123 : 17,
	125 : 1,
	200 : 1,
	201 : 1,
	202 : 401,
	208 : 25,
}
