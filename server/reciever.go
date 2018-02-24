package main

func reciever(c client, packetIN chan<- Packet){
  for {
    fmt.Println("reading from connection", connNumber)
    readbyte, err := c.reader.Peek(1)
    id := int(readbyte[0])
    fmt.Println("got it fam")
    var bytemessage []byte
    switch id {
    case 120:
      bytemessage, _ = c.reader.Peek(17)
      packetIN <- bytemessage
    default:
      c.reader.Discard(1)
    }
  }
}
