package main


type PacketIn struct{
	playerNum byte
	size int
	data []byte
}

type PacketOut struct{
	size int
	data []byte
}

func (p *PacketIn) parseType() byte{
	return p.data[0]
}
