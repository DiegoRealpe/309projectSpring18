package main


type Packet struct{
	size int
	data []byte
}

func (p *Packet) parseType() byte{
	return p.data[0]
}