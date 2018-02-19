package main

import (
	"fmt"
)

//ListenAndSend function which:
//Listens for bytes and parses into clientpackets
//Sends bytes back to all clients with serverpackets
func ListenAndSend(g Game, connNumber int) {
	go func(){
		for {
			fmt.Println("reading from connection", connNumber)
			if g.reader[connNumber].Buffered() > 0 {
				readbyte, _ := g.reader[connNumber].Peek(1)
				id := int(readbyte[0])
				fmt.Println("fuck")
				if id == 120 {
					bytemessage, err := g.reader[connNumber].Peek(17)
					if err != nil {
						fmt.Println("packet error: " + err.Error())
						g.reader[connNumber].ReadByte()
					} else {
						//parse the byte message
						rcvpacket := ParseBytes(bytemessage)
						fmt.Println(rcvpacket.clientPlayerState)
						//construct a message to broadcast to the clients
						sendpacket := ServerPacket{
							serverPlayerState: 121,
							playernumber:      uint8(connNumber),
							xPosition:         rcvpacket.xPosition,
							yPosition:         rcvpacket.yPosition,
							xVelocity:         rcvpacket.xVelocity,
							yVelocity:         rcvpacket.yVelocity,
							timestamp:         0,
						}
						
						sendbytes := ParseServerPacket(sendpacket)
						
						//broadcast that message to the clients
						for _, reciever := range g.writer {
							reciever.Write(sendbytes)
						}
						
					}
					g.reader[connNumber].Read(bytemessage)
				} else {
					g.reader[connNumber].Discard(1)
				}
			}
		}
	}()
}
