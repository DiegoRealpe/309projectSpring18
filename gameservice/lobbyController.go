package main

import(
  "fmt"
)

type LobbyController struct{
  Lobby l
  packetIn chan PacketIn
  packetOut chan PacketOut

  packetRouterMap map[int]func(*PacketIn, chan<- PacketOut)
}

func (lc *LobbyController) addPlayer(connecter *playerConnection){
  if(lc.l.numMembers == NUMPLAYERS){
    fmt.Println("Error, tried to add a player to a full lobby")
    return
  }
  connecter.SetNewPacketInChannel(lc.packetIn)

  for i := 0; i < NUMPLAYERS; i++{
    if lc.l.members[i] == nil{
      lc.l.members[i] == connecter
    }
  }
}

func (lc *LobbyController) LobbyReceiveAndRespond(){
  for in := range lc.packetIn{
    type := in.parseType()
    
  }
}

func (lc *LobbyController) buildLobbyPacketMap(){
  packetMap := map[byte]func(*PacketIn, chan<- PacketOut){}

  packetMap[200] = respondTo200()
  packetMap[201] = respondTo201()
  packetMap[202] = respondTo202()
  packetMap[125] = respondTo125()

  lc.packetRouterMap = packetMap
}

func startLobby(in chan PacketIn, out chan PacketOut){
  lc := LobbyController{}
  lc.l = buildLobby
  lc.buildLobbyPacketMap
  lc.packetIn = in
  lc.packetOut = out
}
