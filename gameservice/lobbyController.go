package main

import(
  "fmt"
)

type LobbyController struct{
  l Lobby
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
  lc.dispersionMap[connecter.id] = connecter.packetOut

  for i := 0; i < NUMPLAYERS; i++{
    if lc.l.members[i] == nil{
      lc.l.members[i] = connecter
      lc.l.connectionIDToPlayerNumberMap[connecter.id] = i
      break
    }
  }
}

func (lc *LobbyController) removePlayer(id int){
  playernum := lc.l.connectionIDToPlayerNumberMap[id]
  lc.l.connectionIDToPlayerNumberMap[id] = 0
  lc.l.isReady[playernum] = 0
  lc.l.connectionIDToPlayerNumberMap[id] = 0
  lc.l.members[playernum] = nil
  lc.l.numMembers--
}

func (lc *LobbyController) allConnectionIDsBut(id int) []int {

	slice := make([]int,NUMPLAYERS-1)

	if debug {fmt.Println("sending to",slice)}

	i := 0
	for key, _ := range lc.l.connectionIDToPlayerNumberMap {
		if key != id{
			slice[i] = key
			i++
		}
	}

	if debug {fmt.Println("sending to",slice)}

	return slice
}

func (lc *LobbyController) LobbyReceiveAndRespond(){
  for in := range lc.packetIn {
    packetType := int((&in).parseType())
    if lc.packetRouterMap[packetType] == nil{
      fmt.Println("Invalid packet type, type was", packetType)
    } else {
      lc.packetRouterMap[packetType](&in, lc.packetOut)
    }
  }
}

func (lc *LobbyController) buildLobbyPacketMap(){
  packetMap := map[int]func(*PacketIn, chan<- PacketOut){}

  lc.packetRouterMap[200] = lc.respondTo200
  lc.packetRouterMap[201] = lc.respondTo201
  lc.packetRouterMap[202] = lc.respondTo202
  lc.packetRouterMap[125] = lc.respondTo125

  lc.packetRouterMap = packetMap
}

func startLobby(in chan PacketIn, out chan PacketOut){
  lc := LobbyController{}
  lc.l = Lobby{}
  lc.buildLobbyPacketMap()
  lc.packetIn = in
  lc.packetOut = out
  go lc.LobbyReceiveAndRespond()
  go listenAndDispersePackets(lc.dispersionMap, lc.packetOut)
}

func (lc *LobbyController) respondTo200(in *PacketIn, out chan<- PacketOut){//player is ready
  playernum := lc.l.connectionIDToPlayerNumberMap[in.connectionId]
  lc.l.isReady[playernum] = 1
  outpacket := PacketOut{
    size: 2,
    data: []byte{0, 0},
    targetIds: lc.allConnectionIDsBut(in.connectionId),
  }
  outpacket.data[0] = byte(204)
  outpacket.data[1] = byte(playernum)
  out <- outpacket
}

func (lc *LobbyController) respondTo201(in *PacketIn, out chan<- PacketOut){//player is no longer ready
  playernum := lc.l.connectionIDToPlayerNumberMap[in.connectionId]
  lc.l.isReady[playernum] = 0
  outpacket := PacketOut{
    size: 2,
    data: []byte{0, 0},
    targetIds: lc.allConnectionIDsBut(in.connectionId),
  }
  outpacket.data[0] = byte(205)
  outpacket.data[1] = byte(playernum)
  out <- outpacket
}

func (lc *LobbyController) respondTo202(in *PacketIn, out chan<- PacketOut){//player send a message in chat

}

func (lc *LobbyController) respondTo125(in *PacketIn, out chan<- PacketOut){//player leaves the lobby
  fmt.Println("recieved 125 packet...")
	disconnectingPlayer := lc.l.connectionIDToPlayerNumberMap[in.connectionId]

	fmt.Println("Player", disconnectingPlayer, "has disconnected")
	lc.l.members[disconnectingPlayer].disconnect()

	packet126 := PacketOut{
    size : 2,
		data : []byte{126, byte(disconnectingPlayer)},
		targetIds: lc.allConnectionIDsBut(in.connectionId),
	}

	out <- packet126
}
