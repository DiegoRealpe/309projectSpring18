package main

import (
	"bytes"
	"encoding/binary"
	"math"
)

//ClientPacket with game state to be read by the server
type ClientPacket struct {
	clientPlayerState uint8 //Packet number 1
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
}

//ServerPacket with game states to be sent to clients
type ServerPacket struct {
	serverPlayerState uint8 //Packet number 2
	playernumber      uint8
	xPosition         float32
	yPosition         float32
	xVelocity         float32
	yVelocity         float32
	timestamp         float32
}

//ParseBytes Takes array of bytes and parses to a clientpacket struct
func ParseBytes(rawData []byte) ClientPacket {
	if len(rawData) != 17 {
		panic(rawData)
	}
	Statebyte := rawData[:1] //Slicing the byte array into its subcomponents
	xPosByte := rawData[1:5]
	yPosByte := rawData[5:9]
	xVelByte := rawData[9:13]
	yVelByte := rawData[13:17]
	resultPacket := ClientPacket{
		clientPlayerState: uint8(Statebyte[0]),
		xPosition:         BytestoFloat32(xPosByte),
		yPosition:         BytestoFloat32(yPosByte),
		xVelocity:         BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
	}
	return resultPacket
}

//ParseServerPacket takes a server packet and readies it to be sent as a byte slice
func ParseServerPacket(packet ServerPacket) []byte {
	rawData := make([]byte, 22)

	i := 0

	rawData[0] = byte(121)
	rawData[1] = byte(packet.playernumber)
	xpos := Float32toBytes(packet.xPosition)
	for i = 2; i <= 5; i++ {
		rawData[i] = xpos[i - 2]
	}
	ypos := Float32toBytes(packet.yPosition)
	for i = 6; i <= 9; i++ {
		rawData[i] = ypos[i - 6]
	}
	xvel := Float32toBytes(packet.xVelocity)
	for i = 10; i <= 13; i++ {
		rawData[i] = xvel[i - 10]
	}
	yvel := Float32toBytes(packet.yVelocity)
	for i = 14; i <= 17; i++ {
		rawData[i] = yvel[i - 14]
	}
	time := Float32toBytes(packet.timestamp)
	for i = 18; i <= 21; i++ {
		rawData[i] = time[i - 18]
	}

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
