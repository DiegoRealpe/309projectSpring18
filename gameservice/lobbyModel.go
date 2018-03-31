package main

import(
  "fmt"
)

type Lobby struct{

  members [NUMPLAYERS]*playerConnection
  isReady [NUMPLAYERS]int
  numMembers int
  connectionIDToPlayerNumberMap map[int]int
}

func (l *Lobby) buildLobbyPacketMap(){
  packetMap := map[byte]func(*PacketIn, chan<- PacketOut){}

  packetMap[200] = respondTo200()
}

func (l *Lobby) respondTo200(in *PacketIn, out chan<- PacketOut){//player is ready
  playernum := connectionIDToPlayerNumberMap[in.connectionId]
  isReady[playernum] = 1
  var outpacket [2]byte
  outpacket[0] = byte(204)
  outpacket[1] = byte(playernum)

}

func (l *Lobby) respondTo201(in *PacketIn, out chan<- PacketOut){//player is no longer ready

}

func (l *Lobby) respondTo202(in *PacketIn, out chan<- PacketOut){//player send a message in chat

}

func (l *Lobby) respondTo125(in *PacketIn, out chan<- PacketOut){//player leaves the lobby

}
