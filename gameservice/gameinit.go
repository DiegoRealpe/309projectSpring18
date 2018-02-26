package main

import(
  "fmt"
)

func initgame(c *[NUMPLAYERS]client){
  packetIN := make(chan PacketIn)
	//packetOUT := make(chan PacketOut)
  go reciever(c[0], packetIN)
  go reciever(c[1], packetIN)
  for rcvpacket := range packetIN {
    id := rcvpacket.data[0]
    switch id {
    case 120:
      fmt.Println("sending a 120")
    }
  }

}
