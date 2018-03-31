package main

import(
  "fmt"
)

type Lobby struct{

  members [NUMPLAYERS]*playerConnection
  isReady [NUMPLAYERS]int
  numMembers int
  packetRouterMap map[int]func(*PacketIn, chan<- PacketOut)
}

func (l *Lobby) buildLobbyPacketMap(){
  packetMap := map[byte]func(*PacketIn, chan<- PacketOut){}

  packetMap[200] = respondTo200()
}

func (l *Lobby) respondTo200(*PacketIn, chan<- PacketOut){//player is ready

}

func (l *Lobby) respondTo201(*PacketIn, chan<- PacketOut){//player is no longer ready

}

func (l *Lobby) respondTo202(*PacketIn, chan<- PacketOut){//player send a message in chat

}

func (l *Lobby) respondTo125(*PacketIn, chan<- PacketOut){//player leaves the lobby

}
