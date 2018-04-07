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

type packet206 struct {
	playerNumber int
	username string
}

type packet202 struct {//client sends a message in chat.
	message					string
}

type packet203 struct {//server broadcasts a message to the rest of the clients
	playerNumber			uint8
	message					string
}

type packet204 struct {//server informs clients that one client is ready.
	numReady				uint8
}

type packet205 struct {//server informs clients that one clinet is no longer ready.
	numUnready				uint8
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

func (packet *packet206) toBytes() []byte {
	rawData := make([]byte, 82)
	rawData[0] = 206
	rawData[1] = byte(packet.playerNumber)

	messageBytes := stringToUtf8Slice(packet.username,80)
	copy(rawData[2:81],messageBytes)

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

func (p packet204) toBytes() []byte{
	return []byte{204,p.numReady}
}

func (p packet205) toBytes() []byte{
	return []byte{205,p.numUnready}
}

func ParseBytesTo202(rawData []byte) packet202 {
	if len(rawData) != 401 {
		panic(rawData)
	}

	return packet202{
		message: utf8toString(rawData[1:401]),
	}
}

func (p packet203) toBytes() []byte{
	rawData := make([]byte, 402)
	rawData[0] = 203
	rawData[1] = p.playerNumber

	messageBytes := stringToUtf8Slice(p.message,400)
	copy(rawData[2:402],messageBytes)

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

func stringToUtf8Slice(s string, length int) []byte{
	if length < len(s){
		return []byte{}
	}

	b := make([]byte,length)
	for i:= 0 ; i < len(s); i++ {
		b[i] = s[i]
	}
	return b
}

func utf8toString(slice []byte) string{
	str := string(slice)


	for i, v := range str{
		if v == 0 {
			str = string(str[:i])
			break
		}
	}

	return str
}
