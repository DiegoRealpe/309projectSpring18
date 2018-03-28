package main

//to be started as a goroutine
func listenAndDispersePackets(dispersionMap map[int]chan<- PacketOut, toDisperse <-chan PacketOut) {
	for packet := range toDisperse {
		disperseSinglePacket(packet, dispersionMap)
	}
}
func disperseSinglePacket(packet PacketOut, dispersionMap map[int]chan<- PacketOut) {
	for _, val := range packet.targetIds {
		connectionChan := dispersionMap[val]
		connectionChan <- packet
	}
}
