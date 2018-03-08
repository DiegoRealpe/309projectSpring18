package main

import (
	"sync"
	"fmt"
	"time"
)

const maxReadBufferSize = 100

type playerConnection struct {
	client        client
	packetInMutex packetInMutex
	packetOut     chan PacketOut

	readBuffer [maxReadBufferSize]byte
	readBufferSize int
	packetSize int
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
		time.Sleep(time.Second)
		//todo: respond correctly to read error

		//allData, err := ioutil.ReadAll(pConn.client.reader)
		//fmt.Println(err)
		allData := []byte{}

		fmt.Println("allData is",allData)

		if len(allData) > 0 {
			fmt.Println(allData)
			pConn.dataWasSent(allData)
		}

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

func (pConn *playerConnection) dataWasSent(bytes []byte) {
	for pConn.readBufferSize + len(bytes) > maxReadBufferSize{
		bytes = pConn.partialUnload(bytes);
	}

	copy(pConn.readBuffer[pConn.readBufferSize:], bytes)
	pConn.readBufferSize += len(bytes)

	fmt.Println("read is", pConn.readBuffer)
	pConn.unloadReadBuffer()
}

func (pConn *playerConnection) unloadReadBuffer() {
	for pConn.readBufferSize >= pConn.packetSize {
		pConn.packetSize = pConn.getPacketSize()
		fmt.Println("got full command", pConn.readBuffer[0:pConn.readBufferSize])

		pConn.readBufferSize -= pConn.packetSize
		copy(
			pConn.readBuffer[0:],
			pConn.readBuffer[pConn.packetSize:])
	}
}
func (pConn *playerConnection) partialUnload(bytes []byte) []byte {
	unloadAmmount := (maxReadBufferSize - pConn.readBufferSize)

	dst := pConn.readBuffer[pConn.readBufferSize:]
	src := bytes[0:unloadAmmount]

	copy(dst,src)
	pConn.readBufferSize += unloadAmmount

	fmt.Println("read buffer is", pConn.readBuffer)

	pConn.unloadReadBuffer()

	return bytes[unloadAmmount+1:]
}

func (pConn *playerConnection) getPacketSize() int{
	if pConn.readBufferSize == 0 {
		return -1
	}

	packetSize := packetSizeMap[pConn.readBuffer[0]]

	if packetSize == 0{
		fmt.Println("packet", pConn.readBuffer[0],"does not exist, handling is not complete yet")
	}

	return packetSize
}

var packetSizeMap = map[byte]int{
	120 : 18,
}