package main

import (
	"encoding/binary"
	"math"
)

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
		xVelolcity:        BytestoFloat32(xVelByte),
		yVelocity:         BytestoFloat32(yVelByte),
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

/*//Float32toBytes self explanatory
func Float32toBytes(input float32) []Byte {
	bits := binary.LittleEndian.Uint32(input)
	float := math.Float32frombits(bits)
	return float
}

//DecodePackage: Takes array of bytes and parses to a clientpacket struct
func (b *byteParser) DecodePackage(serverPacket) []Byte {
	//magic
}
*/
