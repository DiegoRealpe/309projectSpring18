package main

import (
	"bytes"
	"encoding/binary"
	"math"
)

//packet120 with game state to be read by the server
type packet120 struct {
	clientPlayerState uint8 //Packet number 1
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
}

//ServerPacket with game states to be sent to clients
type packet121 struct {
	serverPlayerState uint8 //Packet number 2
	playernumber      uint8
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
	timestamp         float32
}

//ParseBytesTo120 Takes array of bytes and parses to a clientpacket struct
func ParseBytesTo120(rawData []byte) packet120 {
	if len(rawData) != 17 {
		panic(rawData)
	}
	Statebyte := rawData[:1] //Slicing the byte array into its subcomponents
	xPosByte := rawData[1:5]
	yPosByte := rawData[5:9]
	xVelByte := rawData[9:13]
	yVelByte := rawData[13:17]
	resultPacket := packet120{
		clientPlayerState: uint8(Statebyte[0]),
		xPosition:         BytestoFloat32(xPosByte),
		yPosition:         BytestoFloat32(yPosByte),
		xVelocity:         BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
	}
	return resultPacket
}

//Parse121 takes a server packet and readies it to be sent as a byte slice
func Parse121ToBytes(packet packet121) []byte {
	rawData := make([]byte, 22)

	rawData[0] = byte(121)
	rawData[1] = byte(packet.playernumber)

	xpos := Float32toBytes(packet.xPosition)
	copy(rawData[2:5],xpos)

	ypos := Float32toBytes(packet.yPosition)
	copy(rawData[6:9],ypos)

	xvel := Float32toBytes(packet.xVelocity)
	copy(rawData[10:13],xvel)

	yvel := Float32toBytes(packet.yVelocity)
	copy(rawData[14:17],yvel)

	time := Float32toBytes(packet.timestamp)
	copy(rawData[18:21], time)

	return rawData
}

//BytestoFloat32 Turns only a 4 byte slice into a float32 primitive
func BytestoFloat32(input []byte) float32 {
	if len(input) != 4 {
		panic(input)
	}
	bits := binary.LittleEndian.Uint32(input)
	float := math.Float32frombits(bits)
	return float
}

//Float32toBytes self explanatory
func Float32toBytes(input float32) []byte {
	var bytebuffer bytes.Buffer
	binary.Write(&bytebuffer, binary.LittleEndian, input)

	return bytebuffer.Bytes()
}
