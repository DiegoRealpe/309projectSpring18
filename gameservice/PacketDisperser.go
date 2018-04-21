package main


/*
	to be used to disperse packets to a channel based on a map of ints to channels to send on,
	not naturally thread safe
 */

type packetDisperser struct{
	connections map[int]chan<- PacketOut
}

func (pd *packetDisperser) send(packet PacketOut){
	for _, id := range packet.targetIds{
		pd.connections[id] <- packet
	}
}
