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
  lc.dispersionMap[connecter.id] = connecter.packetOut

  for i := 0; i < NUMPLAYERS; i++{
    if lc.l.members[i] == nil{
      lc.l.members[i] == connecter
      lc.l.connectionIDToPlayerNumberMap[connecter.id] = i
      break
    }
  }
}

func (lc *LobbyController) removePlayer(id int){
  playernum := lc.packetRouterMap[id]
  lc.l.connectionIDToPlayerNumberMap[id] = 0
  lc.l.isReady[playernum] = 0
  lc.packetRouterMap[id] = 0
  lc.dispersionMap[id].close()
  lc.l.members[playernum] = nil
  lc.l.numMembers--
}

func (lc *LobbyController) allConnectionIDsBut(id int) []int {

	slice := make([]int,NUMPLAYERS-1)

	if debug {fmt.Println("sending to",slice)}

	i := 0
	for key, _ := range lc.connectionIDToPlayerNumberMap {
		if key != id{
			slice[i] = key
			i++
		}
	}

	if debug {fmt.Println("sending to",slice)}

	return slice
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
  go lc.listenAndDispersePackets(lc.dispersionMap, lc.packetOut)
}

func (lc *LobbyController) respondTo200(in *PacketIn, out chan<- PacketOut){//player is ready
  playernum := connectionIDToPlayerNumberMap[in.connectionId]
  isReady[playernum] = 1
  outpacket := PacketOut{
    size: 2,
    data: [2]byte,
    targetIds: lc.allConnectionIDsBut(in.connectionId)
  }
  outpacket.data[0] = byte(204)
  outpacket.data[1] = byte(playernum)
  out <- outpacket
}

func (lc *LobbyController) respondTo201(in *PacketIn, out chan<- PacketOut){//player is no longer ready
  playernum := connectionIDToPlayerNumberMap[in.connectionId]
  isReady[playernum] = 0
  outpacket := PacketOut{
    size: 2,
    data: [2]byte,
    targetIds: lc.allConnectionIDsBut(in.connectionId)
  }
  outpacket.data[0] = byte(205)
  outpacket.data[1] = byte(playernum)
  out <- outpacket
}

func (lc *LobbyController) respondTo202(in *PacketIn, out chan<- PacketOut){//player send a message in chat

}

func (lc *LobbyController) respondTo125(in *PacketIn, out chan<- PacketOut){//player leaves the lobby
  fmt.Println("recieved 125 packet...")
	disconnectingPlayer := lc.connectionIDToPlayerNumberMap[in.id]

	fmt.Println("Player", disconnectingPlayer, "has disconnected")
	lc.l.members[disconnectingPlayer].disconnect()

	packet126 := PacketOut{
    size : 2,
		data : []byte{126,disconnectingPlayer},
		targetIds: lc.allConnectionIDsBut(in.connectionId),
	}

	out <- packet126
}
