package main

import (
	"encoding/binary"
	"math"
	"bytes"
)

type clientPacket struct {
	clientPlayerState uint8 //Packet number 1
	xPosition         float32
	yPosition         float32
	xVelocity        float32
	yVelocity         float32
}

type serverPacket struct {
	serverPlayerState uint8 //Packet number 2
	playernumber      uint8
	xPosition         float32
	yPosition         float32
	xVelocity        float32
	yVelocity         float32
	timestamp         float32
}

//ParseBytes: Takes array of bytes and parses to a clientpacket struct
func ParseBytes(rawData []byte) clientPacket {
	if len(rawData) != 17 {
		panic(rawData)
	}
	Statebyte := rawData[:1] //Slicing the byte array into its subcomponents
	xPosByte := rawData[1:5]
	yPosByte := rawData[5:9]
	xVelByte := rawData[9:13]
	yVelByte := rawData[13:17]
	resultPacket := clientPacket{
		clientPlayerState: uint8(Statebyte[0]),
		xPosition:         BytestoFloat32(xPosByte),
		yPosition:         BytestoFloat32(yPosByte),
		xVelocity:        BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
	}
	return resultPacket
}

func ParseServerPacket(packet serverPacket) []byte {
	var rawData []byte

	i := 0
	
	rawData = append(rawData, byte(121))
	rawData = append(rawData, byte(packet.playernumber))
	xpos := Float32toBytes(packet.xPosition)
	for i = 2; i <= 5; i++ {
		rawData = append(rawData, xpos[i])
	}
	ypos := Float32toBytes(packet.yPosition)
	for i = 6; i <= 9; i++{
		rawData = append(rawData, ypos[i])
	}
	xvel := Float32toBytes(packet.xVelocity)
	for i = 10; i <= 13; i++ {
		rawData = append(rawData, xvel[i])
	}
	yvel := Float32toBytes(packet.yVelocity)
	for i = 14; i <= 17; i++ {
		rawData = append(rawData, yvel[i])
	}
	time := Float32toBytes(packet.timestamp)
	for i = 18; i <= 21; i++{
		rawData = append(rawData, time[i])
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
/*
//DecodePackage: Takes array of bytes and parses to a clientpacket struct
func (b *byteParser) DecodePackage(serverPacket) []Byte {
	//magic
}
*/
