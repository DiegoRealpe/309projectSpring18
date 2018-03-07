package main

import (
	"fmt"
)

func reciever(c client, packetIN chan<- PacketIn) {
	for {
		fmt.Println("reading from connection", c.clientNum)
		readbyte, _ := c.reader.Peek(1)
		id := int(readbyte[0])
		fmt.Println("got it fam")
		var bytemessage []byte
		var rcvpacket PacketIn
		switch id {
		case 120:
			bytemessage, _ = c.reader.Peek(17)
			c.reader.Read(bytemessage)
			rcvpacket.size = 17
			rcvpacket.data = bytemessage
			packetIN <- rcvpacket
		default:
			c.reader.Discard(1)
		}
	}
}
