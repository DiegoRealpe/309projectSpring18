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

//from the client to update the ball packets
type packet123 struct {
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
}

type packet124 struct {
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
	timestamp         float32
}

type packet125 struct {
	playerNumber			uint8
}


//ParseBytesTo120 Takes array of bytes and parses to a clientpacket struct
func ParseBytesTo120(rawData []byte) packet120 {
	if len(rawData) != 17 {
		panic(rawData)
	}
	xPosByte := rawData[1:5]
	yPosByte := rawData[5:9]
	xVelByte := rawData[9:13]
	yVelByte := rawData[13:17]
	resultPacket := packet120{
		clientPlayerState: 1,
		xPosition:         BytestoFloat32(xPosByte),
		yPosition:         BytestoFloat32(yPosByte),
		xVelocity:         BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
	}
	return resultPacket
}

//Parse121 takes a server packet and readies it to be sent as a byte slice
func (packet *packet121) toBytes() []byte {
	rawData := make([]byte, 22)

	rawData[0] = byte(121)
	rawData[1] = byte(packet.playernumber)

	xpos := Float32toBytes(packet.xPosition)
	copy(rawData[2:6], xpos)

	ypos := Float32toBytes(packet.yPosition)
	copy(rawData[6:10], ypos)

	xvel := Float32toBytes(packet.xVelocity)
	copy(rawData[10:14], xvel)

	yvel := Float32toBytes(packet.yVelocity)
	copy(rawData[14:18], yvel)

	time := Float32toBytes(packet.timestamp)
	copy(rawData[18:22], time)

	return rawData
}

func (packet *packet124) toBytes() []byte {
	rawData := make([]byte, 21)

	rawData[0] = byte(124)

	xpos := Float32toBytes(packet.xPosition)
	copy(rawData[1:5], xpos)

	ypos := Float32toBytes(packet.yPosition)
	copy(rawData[5:9], ypos)

	xvel := Float32toBytes(packet.xVelocity)
	copy(rawData[9:13], xvel)

	yvel := Float32toBytes(packet.yVelocity)
	copy(rawData[13:17], yvel)

	time := Float32toBytes(packet.timestamp)
	copy(rawData[17:21], time)

	return rawData
}

func ParseBytesTo123(rawData []byte) packet123 {
	if len(rawData) != 17 {
		panic(rawData)
	}

	xPosByte := rawData[1:5]
	yPosByte := rawData[5:9]
	xVelByte := rawData[9:13]
	yVelByte := rawData[13:17]
	resultPacket := packet123{
		xPosition:         BytestoFloat32(xPosByte),
		yPosition:         BytestoFloat32(yPosByte),
		xVelocity:         BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
	}
	return resultPacket
}

func ParseBytesTo125(rawData []byte) packet125{
	if len(rawData) != 2 {
		panic(rawData)
	}

	playerNumberByte := rawData[1]

	resultPacket := packet125{
		playerNumber:			uint8(playerNumberByte),
	}
	return resultPacket

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
