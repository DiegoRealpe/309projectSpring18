package main

type PacketIn struct {
	connectionId	int
	size      		int
	data      		[]byte
}

type PacketOut struct {
	size 		int
	data 		[]byte
	targetIds 	[]int
}

func (p *PacketIn) parseType() byte {
	return p.data[0]
}
