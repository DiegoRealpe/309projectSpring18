package main

import(
  "fmt"
)

type LobbyController struct{
  Lobby l
  packetIn chan PacketIn
  packetOut chan PacketOut

  packetRouterMap map[int]func(*PacketIn, chan<- PacketOut)
  dispersionMap map[int]chan<- PacketOut
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
      lc.l.connectionIDToPlayerNumberMap[connecter.id] = i
      break
    }
  }
}

func (lc *LobbyController) LobbyReceiveAndRespond(){
  for in := range lc.packetIn{
    type := in.parseType()
    if lc.packetMap[type] == nil{
      fmt.Println("Invalid packet type, type was", type)
    } else {
      lc.packetMap[type](&in, lc.packetOut)
    }
  }
}

func (lc *LobbyController) buildLobbyPacketMap(){
  packetMap := map[byte]func(*PacketIn, chan<- PacketOut){}

  packetMap[200] = lc.l.respondTo200
  packetMap[201] = lc.l.respondTo201
  packetMap[202] = lc.l.respondTo202
  packetMap[125] = lc.l.respondTo125

  lc.packetRouterMap = packetMap
}

func startLobby(in chan PacketIn, out chan PacketOut){
  lc := LobbyController{}
  lc.l = buildLobby
  lc.buildLobbyPacketMap()
  lc.packetIn = in
  lc.packetOut = out
  go lc.LobbyReceiveAndRespond()
}
