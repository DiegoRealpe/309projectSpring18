package main

//Object to execute parsing functions
type byteParser struct {
}

type clientPacket struct {
	clientPlayerState uint8 //Packet number 1
	xPosition         float32
	yPosition         float32
	xVelolcity        float32
	yVelocity         float32
}

type serverPacket struct {
	serverPlayerState uint8 //Packet number 2
	playernumber      uint8
	xPosition         float32
	yPosition         float32
	xVelolcity        float32
	yVelocity         float32
}

//ParseBytes: Takes array of bytes and parses to a clientpacket struct
func (b *byteParser) ParseBytes(rawData []byte) clientPacket {
	resultPacket := new(clientPacket) //Creating empty clientpack instance
	Statebyte := rawData[:1]          //Slicing the byte array into its subcomponents
	//xPosByte := rawData[1:5]
	//yPosByte := rawData[5:9]
	//xVelByte := rawData[9:13]
	//yVelByte := rawData[13:17]

	resultPacket.clientPlayerState = uint8(Statebyte[0])

	return *resultPacket
}
